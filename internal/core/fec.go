package core

import (
	"fmt"
	"sync"

	"github.com/klauspost/reedsolomon"
)

// encoderCache caches Reed-Solomon encoders to avoid repeated creation overhead
var encoderCache = struct {
	cache map[string]reedsolomon.Encoder
	mutex sync.RWMutex
}{
	cache: make(map[string]reedsolomon.Encoder),
}

// shardPool is a pool of byte slices used for FEC shards to reduce memory allocations
var shardPool = sync.Pool{
	New: func() interface{} {
		// Initial capacity of 1024 bytes, will grow as needed
		return make([]byte, 1024)
	},
}

// dataPool is a pool of byte slices used for decoded data to reduce memory allocations
var dataPool = sync.Pool{
	New: func() interface{} {
		// Initial capacity of 1024 bytes, will grow as needed
		return make([]byte, 1024)
	},
}

// getEncoderFromCache retrieves a Reed-Solomon encoder from the cache or creates a new one
func getEncoderFromCache(k, m int) (reedsolomon.Encoder, error) {
	key := fmt.Sprintf("%d:%d", k, m)
	
	// Try to get from cache first
	encoderCache.mutex.RLock()
	enc, exists := encoderCache.cache[key]
	encoderCache.mutex.RUnlock()
	
	if exists {
		return enc, nil
	}
	
	// Create new encoder if not in cache
	enc, err := reedsolomon.New(k, m)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	encoderCache.mutex.Lock()
	encoderCache.cache[key] = enc
	encoderCache.mutex.Unlock()
	
	return enc, nil
}

// FEC represents a forward error correction encoder/decoder.
type FEC struct {
	enc        reedsolomon.Encoder
	k          int // number of data shards
	m          int // number of parity shards
	lastDataSize int // size of the last encoded data
}

// ReturnShards returns the shards to the memory pool for reuse.
// This should be called after processing the shards to avoid memory leaks.
func (f *FEC) ReturnShards(shards [][]byte) {
	for _, shard := range shards {
		if shard != nil {
			shardPool.Put(shard[:cap(shard)])
		}
	}
}

// ReturnData returns the decoded data to the memory pool for reuse.
// This should be called after processing the data to avoid memory leaks.
func (f *FEC) ReturnData(data []byte) {
	if data != nil {
		dataPool.Put(data[:cap(data)])
	}
}

// NewFEC creates a new FEC with the specified number of data and parity shards.
func NewFEC(k, m int) (*FEC, error) {
	enc, err := getEncoderFromCache(k, m)
	if err != nil {
		return nil, fmt.Errorf("failed to get reedsolomon encoder: %w", err)
	}
	return &FEC{
		enc: enc,
		k:   k,
		m:   m,
	}, nil
}

// Encode encodes the data into shards, including parity shards.
func (f *FEC) Encode(data []byte) ([][]byte, error) {
	// Store the original data size
	f.lastDataSize = len(data)
	
	// Calculate shard size
	shardSize := (len(data) + f.k - 1) / f.k
	
	// Create shards using memory pool
	shards := make([][]byte, f.k+f.m)
	for i := range shards {
		// Get a byte slice from the pool
		shard := shardPool.Get().([]byte)
		
		// Ensure it's large enough
		if cap(shard) < shardSize {
			// If not large enough, allocate a new one
			shard = make([]byte, shardSize)
		} else {
			// If large enough, resize it
			shard = shard[:shardSize]
		}
		
		shards[i] = shard
	}
	
	// Copy data to data shards (first k shards) row-wise
	for i := 0; i < f.k; i++ {
		start := i * shardSize
		end := start + shardSize
		if end > len(data) {
			end = len(data)
		}
		if start < len(data) {
			copy(shards[i], data[start:end])
		}
	}
	
	// Encode parity shards
	if err := f.enc.Encode(shards); err != nil {
		// Return shards to pool before returning error
		for _, shard := range shards {
			shardPool.Put(shard[:cap(shard)])
		}
		return nil, fmt.Errorf("failed to encode shards: %w", err)
	}
	
	return shards, nil
}

// Decode decodes the shards back into data.
func (f *FEC) Decode(shards [][]byte) ([]byte, error) {
	// Verify shard count
	if len(shards) != f.k+f.m {
		return nil, fmt.Errorf("invalid shard count: expected %d, got %d", f.k+f.m, len(shards))
	}
	
	// Verify shard sizes
	if len(shards) > 0 {
		// Find the first non-nil shard to determine expected size
		var shardSize int
		var found bool
		for _, shard := range shards {
			if shard != nil {
				shardSize = len(shard)
				found = true
				break
			}
		}
		
		// If all shards are nil, that's an error
		if !found {
			return nil, fmt.Errorf("all shards are nil")
		}
		
		// Verify that all non-nil shards have the same size
		for _, shard := range shards {
			if shard != nil && len(shard) != shardSize {
				return nil, fmt.Errorf("shards have different sizes")
			}
		}
		
		// Reconstruct missing shards if needed
		if err := f.enc.Reconstruct(shards); err != nil {
			return nil, fmt.Errorf("failed to reconstruct shards: %w", err)
		}
		
		// Calculate the original data size
		// Use the stored original data size if available, otherwise calculate it
		dataSize := f.lastDataSize
		if dataSize == 0 {
			// Fallback to the old method if lastDataSize is not set
			dataSize = f.k * shardSize
		}
		
		// Get a byte slice from the pool for the decoded data
		data := dataPool.Get().([]byte)
		
		// Ensure it's large enough
		if cap(data) < dataSize {
			// If not large enough, allocate a new one
			data = make([]byte, dataSize)
		} else {
			// If large enough, resize it
			data = data[:dataSize]
		}
		
		// Copy data shards (first k shards) into the data slice
		offset := 0
		for i := 0; i < f.k; i++ {
			shardLen := len(shards[i])
			if offset+shardLen > dataSize {
				shardLen = dataSize - offset
			}
			if shardLen > 0 {
				copy(data[offset:], shards[i][:shardLen])
				offset += shardLen
			}
		}
		
		// Return shards to pool
		for _, shard := range shards {
			if shard != nil {
				shardPool.Put(shard[:cap(shard)])
			}
		}
		
		return data, nil
	}
	
	return nil, fmt.Errorf("no shards provided")
}
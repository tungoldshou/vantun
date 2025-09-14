package core

import (
	"fmt"

	"github.com/klauspost/reedsolomon"
)

// FEC represents a forward error correction encoder/decoder.
type FEC struct {
	enc reedsolomon.Encoder
	k   int // number of data shards
	m   int // number of parity shards
}

// NewFEC creates a new FEC with the specified number of data and parity shards.
func NewFEC(k, m int) (*FEC, error) {
	enc, err := reedsolomon.New(k, m)
	if err != nil {
		return nil, fmt.Errorf("failed to create reedsolomon encoder: %w", err)
	}
	return &FEC{
		enc: enc,
		k:   k,
		m:   m,
	}, nil
}

// Encode encodes the data into shards, including parity shards.
func (f *FEC) Encode(data []byte) ([][]byte, error) {
	// Calculate shard size
	shardSize := (len(data) + f.k - 1) / f.k
	
	// Create shards
	shards := make([][]byte, f.k+f.m)
	for i := range shards {
		shards[i] = make([]byte, shardSize)
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
		// In a real implementation, we would store the original data size
		// For now, we'll reconstruct the data by concatenating all data shards
		dataSize := f.k * shardSize
		data := make([]byte, 0, dataSize)
		
		// Concatenate all data shards (first k shards)
		for i := 0; i < f.k; i++ {
			data = append(data, shards[i]...)
		}
		
		return data, nil
	}
	
	return nil, fmt.Errorf("no shards provided")
}
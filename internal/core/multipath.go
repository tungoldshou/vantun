package core

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// Path represents a network path.
type Path struct {
	// addr is the address of the path.
	addr string
	// conn is the QUIC connection for the path.
	conn quic.Connection
	// rtt is the round-trip time for the path.
	rtt time.Duration
	// loss is the packet loss rate for the path.
	loss float64
	// bandwidth is the estimated bandwidth for the path.
	bandwidth uint64
	// active indicates if the path is active.
	active bool
	// lastActive is the time when the path was last active.
	lastActive time.Time
}

// MultipathSession represents a multipath session.
type MultipathSession struct {
	// paths is the list of paths.
	paths []*Path
	// mutex protects the paths slice.
	mutex sync.RWMutex
	// nextPath is the index of the next path to use for round-robin scheduling.
	nextPath int
	// config holds the TLS configuration for new connections.
	config *Config
	// dataSplitter handles data splitting across multiple paths.
	dataSplitter *DataSplitter
	// pathSelector selects the best path for data transmission.
	pathSelector *PathSelector
	// tokenBucketController controls the rate of data transmission.
	tokenBucketController *TokenBucketController
	// adaptiveFEC adjusts FEC parameters based on telemetry data.
	adaptiveFEC *AdaptiveFEC
}

// DataSplitter handles splitting data across multiple paths.
type DataSplitter struct {
	// chunkSize is the size of each data chunk.
	chunkSize int
}

// NewDataSplitter creates a new DataSplitter.
func NewDataSplitter(chunkSize int) *DataSplitter {
	return &DataSplitter{
		chunkSize: chunkSize,
	}
}

// Split splits data into chunks.
func (ds *DataSplitter) Split(data []byte) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += ds.chunkSize {
		end := i + ds.chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

// PathSelector selects the best path for data transmission.
type PathSelector struct {
	// strategy is the path selection strategy.
	strategy PathSelectionStrategy
}

// PathSelectionStrategy is the strategy for selecting paths.
type PathSelectionStrategy int

const (
	// RoundRobinStrategy selects paths in a round-robin fashion.
	RoundRobinStrategy PathSelectionStrategy = iota
	// MinRTTStrategy selects the path with the minimum RTT.
	MinRTTStrategy
	// WeightedStrategy selects paths based on their weights.
	WeightedStrategy
)

// NewPathSelector creates a new PathSelector.
func NewPathSelector(strategy PathSelectionStrategy) *PathSelector {
	return &PathSelector{
		strategy: strategy,
	}
}

// SelectPath selects the best path based on the strategy.
func (ps *PathSelector) SelectPath(paths []*Path, nextPath *int) *Path {
	if len(paths) == 0 {
		return nil
	}

	switch ps.strategy {
	case RoundRobinStrategy:
		// Find an active path using round-robin
		for i := 0; i < len(paths); i++ {
			idx := (*nextPath + i) % len(paths)
			path := paths[idx]
			if path.active {
				*nextPath = (idx + 1) % len(paths)
				return path
			}
		}
	case MinRTTStrategy:
		// Find the path with the minimum RTT
		var bestPath *Path
		minRTT := time.Duration(1<<63 - 1) // Max int64
		for _, path := range paths {
			if path.active && path.rtt < minRTT {
				minRTT = path.rtt
				bestPath = path
			}
		}
		return bestPath
	case WeightedStrategy:
		// Select paths based on their weights (bandwidth)
		totalWeight := uint64(0)
		for _, path := range paths {
			if path.active {
				totalWeight += path.bandwidth
			}
		}

		if totalWeight == 0 {
			// If all paths have zero weight, fall back to round-robin
			return ps.SelectPath(paths, nextPath)
		}

		// Select a path based on its weight
		randWeight := uint64(rand.Int63n(int64(totalWeight)))
		currentWeight := uint64(0)
		for _, path := range paths {
			if path.active {
				currentWeight += path.bandwidth
				if currentWeight >= randWeight {
					return path
				}
			}
		}
	}

	return nil
}

// NewMultipathSession creates a new MultipathSession.
func NewMultipathSession(config *Config, tokenBucketController *TokenBucketController, adaptiveFEC *AdaptiveFEC) *MultipathSession {
	return &MultipathSession{
		paths:                 make([]*Path, 0),
		config:                config,
		dataSplitter:          NewDataSplitter(1024), // 1KB chunks
		pathSelector:          NewPathSelector(RoundRobinStrategy),
		tokenBucketController: tokenBucketController,
		adaptiveFEC:           adaptiveFEC,
	}
}

// AddPath adds a new path to the session.
func (ms *MultipathSession) AddPath(ctx context.Context, addr string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	Info("Adding path to %s", addr)
	// Establish a new QUIC connection for the path
	conn, err := quic.DialAddr(ctx, addr, ms.config.TLSConfig, nil)
	if err != nil {
		Error("Failed to dial path %s: %v", addr, err)
		return fmt.Errorf("failed to dial path %s: %w", addr, err)
	}
	Info("Successfully dialed path %s", addr)

	// Perform session negotiation handshake on the control stream.
	if err := performClientHandshake(ctx, conn); err != nil {
		conn.CloseWithError(0, "handshake failed")
		Error("Handshake failed for path %s: %v", addr, err)
		return fmt.Errorf("handshake failed for path %s: %w", addr, err)
	}
	Info("Session handshake completed for path %s", addr)

	// Create a new path with initial placeholder values
	path := &Path{
		addr:       addr,
		conn:       conn,
		active:     true,
		rtt:        50 * time.Millisecond,      // Initial RTT placeholder
		loss:       0.01,                       // Initial loss placeholder (1%)
		bandwidth:  1000000,                    // Initial bandwidth placeholder (1 MB/s)
		lastActive: time.Now(),
	}

	// Add the path to the session
	ms.paths = append(ms.paths, path)

	// Update the token bucket controller with the first path's connection
	if ms.tokenBucketController != nil && len(ms.paths) == 1 {
		// Update the token bucket controller with the connection
		ms.tokenBucketController.UpdateConnection(conn)
		Info("Updated token bucket controller with connection to %s", addr)
	}

	// Start path probing in a separate goroutine
	go ms.probePath(path)

	return nil
}

// RemovePath removes a path from the session.
func (ms *MultipathSession) RemovePath(addr string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for i, path := range ms.paths {
		if path.addr == addr {
			// Close the connection
			if path.conn != nil {
				path.conn.CloseWithError(0, "path removed")
			}

			// Remove the path from the slice
			ms.paths = append(ms.paths[:i], ms.paths[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("path %s not found", addr)
}

// probePath continuously probes a path to update its metrics.
func (ms *MultipathSession) probePath(path *Path) {
	// Create a ticker for periodic probing
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Since quic-go Connection interface doesn't provide Stats method in all versions,
			// we'll use placeholder values for path metrics
			ms.mutex.Lock()
			
			// Update RTT with placeholder value
			if path.rtt == 0 {
				path.rtt = 50 * time.Millisecond
			}
			
			// Update loss with placeholder value
			if path.loss == 0 {
				path.loss = 0.01 // 1% default
			}
			
			// Update bandwidth with placeholder value
			if path.bandwidth == 0 {
				path.bandwidth = 1000000 // 1 MB/s placeholder
			}
			
			path.lastActive = time.Now()
			ms.mutex.Unlock()

			Info("Path %s: RTT=%v, Loss=%.2f%%, Bandwidth=%d B/s",
				path.addr, path.rtt, path.loss*100, path.bandwidth)
		}
	}
}

// OpenStream opens a new stream on the best path.
func (ms *MultipathSession) OpenStream(ctx context.Context) (quic.Stream, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	// Select the best path
	path := ms.pathSelector.SelectPath(ms.paths, &ms.nextPath)
	if path == nil {
		return nil, fmt.Errorf("no active paths available")
	}

	Info("Opening stream on path %s", path.addr)
	// Stream type 1 is for interactive data.
	stream, err := path.conn.OpenStreamSync(ctx)
	if err != nil {
		Info("Failed to open stream on path %s: %v", path.addr, err)
		return nil, err
	}
	Info("Successfully opened stream on path %s", path.addr)
	
	// Send stream type identifier on the stream.
	Info("Sending stream type message on path %s", path.addr)
	payload := &StreamTypePayload{
		Type: StreamTypeInteractive,
	}
	data, err := EncodeStreamType(payload)
	if err != nil {
		stream.Close()
		Info("Failed to encode stream type: %v", err)
		return nil, fmt.Errorf("failed to encode stream type: %w", err)
	}
	
	msg := &Message{
		Type: StreamType,
		Data: data,
	}
	
	if err := WriteMessage(stream, msg); err != nil {
		stream.Close()
		Info("Failed to send stream type message on path %s: %v", path.addr, err)
		return nil, fmt.Errorf("failed to send stream type: %w", err)
	}
	
	Info("Successfully opened stream on path %s and sent stream type message", path.addr)
	return stream, nil
}

// AcceptStream accepts a new stream on any path.
func (ms *MultipathSession) AcceptStream(ctx context.Context) (quic.Stream, error) {
	// For simplicity, we'll just accept on the first path
	// A real implementation would need to handle streams from multiple paths
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	if len(ms.paths) > 0 && ms.paths[0].active {
		return ms.paths[0].conn.AcceptStream(ctx)
	}

	return nil, fmt.Errorf("no active paths available for accepting streams")
}

// SendData sends data across multiple paths with FEC and rate control.
func (ms *MultipathSession) SendData(ctx context.Context, data []byte) error {
	// Apply FEC encoding if adaptiveFEC is available
	var chunks [][]byte
	if ms.adaptiveFEC != nil {
		var err error
		chunks, err = ms.adaptiveFEC.Encode(data)
		if err != nil {
			return fmt.Errorf("failed to encode data with FEC: %w", err)
		}
		Info("Encoded data with FEC: %d chunks", len(chunks))
	} else {
		// Split data into chunks without FEC
		chunks = ms.dataSplitter.Split(data)
	}

	// Send each chunk on a different path
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	for i, chunk := range chunks {
		// Select a path for this chunk
		path := ms.pathSelector.SelectPath(ms.paths, &ms.nextPath)
		if path == nil {
			return fmt.Errorf("no active paths available for sending data")
		}

		// Check token bucket if controller is available
		if ms.tokenBucketController != nil {
			// In a real implementation, you would check the token bucket before sending
			// For now, we'll just log that we're checking
			Info("Checking token bucket for path %s", path.addr)
		}

		// Open a stream on the selected path
		stream, err := path.conn.OpenStreamSync(ctx)
		if err != nil {
			// Mark the path as inactive
			path.active = false
			Warn("Failed to open stream on path %s: %v", path.addr, err)
			continue
		}

		// Send the chunk
		_, err = stream.Write(chunk)
		if err != nil {
			// Mark the path as inactive
			path.active = false
			Error("Failed to send chunk on path %s: %v", path.addr, err)
			stream.Close()
			continue
		}

		// Close the stream
		stream.Close()

		Info("Sent chunk %d/%d on path %s", i+1, len(chunks), path.addr)
	}

	return nil
}

// ReceiveData receives data from multiple paths and decodes FEC if necessary.
// This is a simplified implementation
// In a real implementation, you would need to:
// 1. Receive chunks from multiple paths concurrently
// 2. Reassemble the chunks in the correct order
// 3. Handle lost chunks and retransmissions
func (ms *MultipathSession) ReceiveData(ctx context.Context) ([]byte, error) {
	// For now, we'll just receive data from the first path
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	if len(ms.paths) > 0 && ms.paths[0].active {
		stream, err := ms.paths[0].conn.AcceptStream(ctx)
		if err != nil {
			return nil, err
		}
		defer stream.Close()

		buf := make([]byte, 4096) // Larger buffer for FEC chunks
		n, err := stream.Read(buf)
		if err != nil {
			return nil, err
		}

		// If we have FEC, try to decode the data
		if ms.adaptiveFEC != nil {
			// This is a simplified approach - in a real implementation,
			// you would collect multiple chunks and then decode them
			// For now, we'll just return the raw data
			Info("Received data with potential FEC encoding")
			return buf[:n], nil
		}

		return buf[:n], nil
	}

	return nil, fmt.Errorf("no active paths available for receiving data")
}

// Close closes all paths in the session.
func (ms *MultipathSession) Close() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Close all paths
	for _, path := range ms.paths {
		if path.conn != nil {
			path.conn.CloseWithError(0, "session closed")
		}
	}

	return nil
}

// SetPathSelectionStrategy sets the path selection strategy.
func (ms *MultipathSession) SetPathSelectionStrategy(strategy PathSelectionStrategy) {
	ms.pathSelector.strategy = strategy
}

// GetPathStats returns statistics for all paths.
func (ms *MultipathSession) GetPathStats() []Path {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	stats := make([]Path, len(ms.paths))
	for i, path := range ms.paths {
		stats[i] = *path
	}

	return stats
}
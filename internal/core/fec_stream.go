package core

import (
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

// FECStream wraps a quic.Stream to add FEC capabilities.
type FECStream struct {
	stream quic.Stream
	fec    *FEC
	k      int // number of data shards
	m      int // number of parity shards
}

// NewFECStream creates a new FECStream.
func NewFECStream(stream quic.Stream, k, m int) (*FECStream, error) {
	fec, err := NewFEC(k, m)
	if err != nil {
		return nil, fmt.Errorf("failed to create FEC: %w", err)
	}
	return &FECStream{
		stream: stream,
		fec:    fec,
		k:      k,
		m:      m,
	}, nil
}

// Write writes data to the stream with FEC encoding.
func (f *FECStream) Write(p []byte) (n int, err error) {
	// Encode data into shards
	shards, err := f.fec.Encode(p)
	if err != nil {
		return 0, fmt.Errorf("failed to encode data: %w", err)
	}
	
	// Send each shard with a header indicating shard index and total count
	for i, shard := range shards {
		// Create header: [shard_index][total_shards][data_length][shard_data]
		header := make([]byte, 8)
		header[0] = byte(i)
		header[1] = byte(f.k + f.m)
		// Note: This is a simplified header format. In a real implementation,
		// you would need to handle data length and use a more robust framing.
		
		// Write header and shard data
		if _, err := f.stream.Write(header); err != nil {
			return 0, fmt.Errorf("failed to write header: %w", err)
		}
		if _, err := f.stream.Write(shard); err != nil {
			return 0, fmt.Errorf("failed to write shard: %w", err)
		}
	}
	
	return len(p), nil
}

// Read reads data from the stream and decodes it with FEC.
func (f *FECStream) Read(p []byte) (n int, err error) {
	// This is a simplified implementation. In a real implementation,
	// you would need to buffer incoming shards and reconstruct the data
	// when enough shards are received.
	
	// For now, we'll just read from the underlying stream directly.
	return f.stream.Read(p)
}

// Close closes the underlying stream.
func (f *FECStream) Close() error {
	return f.stream.Close()
}

// StreamID returns the ID of the underlying stream.
func (f *FECStream) StreamID() quic.StreamID {
	return f.stream.StreamID()
}

// SetDeadline sets the read and write deadlines for the stream.
func (f *FECStream) SetDeadline(t time.Time) error {
	return f.stream.SetDeadline(t)
}

// SetReadDeadline sets the read deadline for the stream.
func (f *FECStream) SetReadDeadline(t time.Time) error {
	return f.stream.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline for the stream.
func (f *FECStream) SetWriteDeadline(t time.Time) error {
	return f.stream.SetWriteDeadline(t)
}
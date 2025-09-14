package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// Obfuscator obfuscates data to make it look like HTTP/3 traffic.
type Obfuscator struct {
	// enabled indicates if obfuscation is enabled
	enabled bool
	// http3Obfuscator is the HTTP/3 obfuscator instance
	http3Obfuscator *HTTP3Obfuscator
}

// ObfuscatorConfig holds the configuration for the Obfuscator.
type ObfuscatorConfig struct {
	// Enabled indicates if obfuscation is enabled
	Enabled bool
	// FrameTypes is a list of HTTP/3 frame types to use for obfuscation
	FrameTypes []byte
	// MinPadding is the minimum padding size
	MinPadding int
	// MaxPadding is the maximum padding size
	MaxPadding int
}

// NewObfuscator creates a new Obfuscator.
func NewObfuscator(config ObfuscatorConfig) *Obfuscator {
	// If no frame types are specified, use a default set
	frameTypes := config.FrameTypes
	if len(frameTypes) == 0 {
		frameTypes = []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7}
	}

	// Ensure minPadding and maxPadding have reasonable defaults
	minPadding := config.MinPadding
	if minPadding < 0 {
		minPadding = 0
	}

	maxPadding := config.MaxPadding
	if maxPadding < minPadding {
		maxPadding = minPadding + 100
	}

	// Create HTTP3 obfuscator if obfuscation is enabled
	var http3Obfs *HTTP3Obfuscator
	if config.Enabled {
		http3Obfs = NewHTTP3Obfuscator(frameTypes, minPadding, maxPadding)
	}

	return &Obfuscator{
		enabled:         config.Enabled,
		http3Obfuscator: http3Obfs,
	}
}

// Obfuscate obfuscates the data to make it look like HTTP/3 traffic.
// If obfuscation is not enabled, it returns the data as-is.
func (o *Obfuscator) Obfuscate(data []byte) ([]byte, error) {
	if !o.enabled {
		// If obfuscation is not enabled, return the data as-is
		return data, nil
	}

	// Use the HTTP/3 obfuscator for better obfuscation
	return o.http3Obfuscator.Obfuscate(data)
}

// Deobfuscate deobfuscates the data.
// If obfuscation is not enabled, it returns the data as-is.
func (o *Obfuscator) Deobfuscate(data []byte) ([]byte, error) {
	if !o.enabled {
		// If obfuscation is not enabled, return the data as-is
		return data, nil
	}

	// Use the HTTP/3 obfuscator for better obfuscation
	return o.http3Obfuscator.Deobfuscate(data)
}

// ObfuscatorStream wraps a QUIC stream to obfuscate data.
type ObfuscatorStream struct {
	quic.Stream
	obfuscator *Obfuscator
	// readBuf is a buffer for reading obfuscated data
	readBuf *bytes.Buffer
	// writeBuf is a buffer for writing obfuscated data
	writeBuf *bytes.Buffer
	// mutex protects the buffers
	mutex sync.Mutex
}

// NewObfuscatorStream creates a new ObfuscatorStream.
func NewObfuscatorStream(stream quic.Stream, obfuscator *Obfuscator) *ObfuscatorStream {
	return &ObfuscatorStream{
		Stream:     stream,
		obfuscator: obfuscator,
		readBuf:    &bytes.Buffer{},
		writeBuf:   &bytes.Buffer{},
	}
}

// Write writes data to the stream, obfuscating it first if obfuscation is enabled.
func (os *ObfuscatorStream) Write(p []byte) (n int, err error) {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	obfuscatedData, err := os.obfuscator.Obfuscate(p)
	if err != nil {
		return 0, err
	}

	// Write the obfuscated data to the underlying stream
	_, err = os.Stream.Write(obfuscatedData)
	if err != nil {
		return 0, err
	}

	// Return the original data length, not the obfuscated length
	return len(p), nil
}

// Read reads data from the stream, deobfuscating it if obfuscation is enabled.
func (os *ObfuscatorStream) Read(p []byte) (n int, err error) {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	// If we have data in the read buffer, use it first
	if os.readBuf.Len() > 0 {
		return os.readBuf.Read(p)
	}

	// For simplicity in this implementation, we'll read a fixed amount of data
	// In a production implementation, you would want to handle partial reads better
	tempBuf := make([]byte, 4096)
	n, err = os.Stream.Read(tempBuf)
	if err != nil {
		return 0, err
	}

	// Deobfuscate the data
	deobfuscatedData, err := os.obfuscator.Deobfuscate(tempBuf[:n])
	if err != nil {
		return 0, err
	}

	// Copy the deobfuscated data to the read buffer
	os.readBuf.Write(deobfuscatedData)

	// Now read from the read buffer to the output buffer
	return os.readBuf.Read(p)
}

// HTTP3Obfuscator obfuscates data to look like real HTTP/3 traffic.
type HTTP3Obfuscator struct {
	// frameTypes is a list of HTTP/3 frame types to use for obfuscation
	frameTypes []byte
	// minPadding is the minimum padding size
	minPadding int
	// maxPadding is the maximum padding size
	maxPadding int
	// rng is a random number generator
	rng *rand.Rand
}

// NewHTTP3Obfuscator creates a new HTTP3Obfuscator.
func NewHTTP3Obfuscator(frameTypes []byte, minPadding, maxPadding int) *HTTP3Obfuscator {
	return &HTTP3Obfuscator{
		frameTypes: frameTypes,
		minPadding: minPadding,
		maxPadding: maxPadding,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Obfuscate obfuscates the data to look like HTTP/3 traffic.
// This creates fake HTTP/3 frames.
func (h *HTTP3Obfuscator) Obfuscate(data []byte) ([]byte, error) {
	// HTTP/3 frame format:
	// - Type (1 byte)
	// - Length (1-8 bytes, varint)
	// - Payload (variable)

	// We'll create a series of fake HTTP/3 frames with our data embedded in them.

	var buf bytes.Buffer

	// Split data into chunks to fit into multiple frames
	chunkSize := 1024 // Arbitrary chunk size
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]

		// Always use DATA frame type (0x0) for actual data to ensure it can be extracted
		frameType := byte(0x0)

		// Use the actual chunk data
		frameData := chunk

		frameLength := len(frameData)

		// Write frame header
		buf.WriteByte(frameType)

		// Write length as proper HTTP/3 varint
		if err := writeVarInt(&buf, uint64(frameLength)); err != nil {
			return nil, err
		}

		// Write frame payload
		buf.Write(frameData)

		// Add some fake padding frames with random data
		// The probability of adding a padding frame increases with the obfuscation level
		if h.rng.Float32() < 0.4 { // 40% chance to add a padding frame
			paddingType := byte(0x1) // PADDING frame type
			paddingLength := h.minPadding + h.rng.Intn(h.maxPadding-h.minPadding+1)

			// Write padding frame header
			buf.WriteByte(paddingType)

			// Write padding length as proper HTTP/3 varint
			if err := writeVarInt(&buf, uint64(paddingLength)); err != nil {
				return nil, err
			}

			// Write random padding data
			for j := 0; j < paddingLength; j++ {
				buf.WriteByte(byte(h.rng.Intn(256)))
			}
		}
	}

	return buf.Bytes(), nil
}

// Deobfuscate deobfuscates the data from HTTP/3-like format.
func (h *HTTP3Obfuscator) Deobfuscate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	reader := bytes.NewReader(data)

	for reader.Len() > 0 {
		// Read frame type
		frameType, err := reader.ReadByte()
		if err != nil {
			// If we've reached the end of the data, break
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("incomplete frame header: %v", err)
		}

		// Read frame length as varint
		frameLength, err := readVarInt(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read frame length: %v", err)
		}

		// Check if we have enough data for the frame payload
		if uint64(reader.Len()) < frameLength {
			return nil, fmt.Errorf("incomplete frame payload: need %d bytes, have %d", frameLength, reader.Len())
		}

		// Process frame based on type
		if frameType == 0x0 { // DATA frame
			frameData := make([]byte, frameLength)
			_, err := io.ReadFull(reader, frameData)
			if err != nil {
				return nil, fmt.Errorf("failed to read frame data: %v", err)
			}
			buf.Write(frameData)
		} else {
			// Skip other frame types (e.g., padding 0x1)
			_, err := io.CopyN(io.Discard, reader, int64(frameLength))
			if err != nil {
				return nil, fmt.Errorf("failed to skip frame data: %v", err)
			}
		}
	}

	return buf.Bytes(), nil
}

// writeVarInt writes an HTTP/3 variable-length integer to the buffer
func writeVarInt(buf *bytes.Buffer, value uint64) error {
	switch {
	case value <= 0x3F: // 1 byte
		buf.WriteByte(byte(value))
	case value <= 0x3FFF: // 2 bytes
		buf.WriteByte(byte(0x40 | (value >> 8)))
		buf.WriteByte(byte(value & 0xFF))
	case value <= 0x3FFFFFFF: // 4 bytes
		buf.WriteByte(byte(0x80 | (value >> 24)))
		buf.WriteByte(byte((value >> 16) & 0xFF))
		buf.WriteByte(byte((value >> 8) & 0xFF))
		buf.WriteByte(byte(value & 0xFF))
	case value <= 0x3FFFFFFFFFFFFFFF: // 8 bytes
		buf.WriteByte(byte(0xC0 | (value >> 56)))
		buf.WriteByte(byte((value >> 48) & 0xFF))
		buf.WriteByte(byte((value >> 40) & 0xFF))
		buf.WriteByte(byte((value >> 32) & 0xFF))
		buf.WriteByte(byte((value >> 24) & 0xFF))
		buf.WriteByte(byte((value >> 16) & 0xFF))
		buf.WriteByte(byte((value >> 8) & 0xFF))
		buf.WriteByte(byte(value & 0xFF))
	default:
		return fmt.Errorf("value too large for varint encoding: %d", value)
	}
	return nil
}

// readVarInt reads an HTTP/3 variable-length integer from the reader
func readVarInt(reader *bytes.Reader) (uint64, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	switch {
	case firstByte&0xC0 == 0x00: // 1 byte
		return uint64(firstByte & 0x3F), nil
	case firstByte&0xC0 == 0x40: // 2 bytes
		secondByte, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		return (uint64(firstByte&0x3F) << 8) | uint64(secondByte), nil
	case firstByte&0xE0 == 0x80: // 4 bytes
		buf := make([]byte, 3)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			return 0, err
		}
		return (uint64(firstByte&0x1F) << 24) | (uint64(buf[0]) << 16) | (uint64(buf[1]) << 8) | uint64(buf[2]), nil
	case firstByte&0xF0 == 0xC0: // 8 bytes
		buf := make([]byte, 7)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			return 0, err
		}
		return (uint64(firstByte&0x0F) << 56) | (uint64(buf[0]) << 48) | (uint64(buf[1]) << 40) | (uint64(buf[2]) << 32) |
			(uint64(buf[3]) << 24) | (uint64(buf[4]) << 16) | (uint64(buf[5]) << 8) | uint64(buf[6]), nil
	default:
		return 0, fmt.Errorf("invalid varint prefix: 0x%02x", firstByte)
	}
}

// ObfuscatorSession wraps a Session to provide obfuscated streams.
type ObfuscatorSession struct {
	*Session
	obfuscator *Obfuscator
}

// NewObfuscatorSession creates a new ObfuscatorSession.
func NewObfuscatorSession(session *Session, obfuscator *Obfuscator) *ObfuscatorSession {
	return &ObfuscatorSession{
		Session:    session,
		obfuscator: obfuscator,
	}
}

// OpenInteractiveStream opens a new interactive stream with obfuscation.
func (os *ObfuscatorSession) OpenInteractiveStream(ctx context.Context) (quic.Stream, error) {
	stream, err := os.Session.OpenInteractiveStream(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap the stream with obfuscation
	return NewObfuscatorStream(stream, os.obfuscator), nil
}

// AcceptInteractiveStream accepts a new interactive stream with obfuscation.
func (os *ObfuscatorSession) AcceptInteractiveStream(ctx context.Context) (quic.Stream, error) {
	stream, err := os.Session.AcceptInteractiveStream(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap the stream with obfuscation
	return NewObfuscatorStream(stream, os.obfuscator), nil
}

// OpenBulkStream opens a new bulk stream with obfuscation.
func (os *ObfuscatorSession) OpenBulkStream(ctx context.Context) (quic.Stream, error) {
	stream, err := os.Session.OpenBulkStream(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap the stream with obfuscation
	return NewObfuscatorStream(stream, os.obfuscator), nil
}

// AcceptBulkStream accepts a new bulk stream with obfuscation.
func (os *ObfuscatorSession) AcceptBulkStream(ctx context.Context) (quic.Stream, error) {
	stream, err := os.Session.AcceptBulkStream(ctx)
	if err != nil {
		return nil, err
	}

	// Wrap the stream with obfuscation
	return NewObfuscatorStream(stream, os.obfuscator), nil
}
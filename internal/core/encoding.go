package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/fxamacker/cbor/v2"
)

// lengthBufPool is a pool of byte slices used for length prefixes to reduce memory allocations
var lengthBufPool = sync.Pool{
	New: func() interface{} {
		// Length buffers are fixed size of 4 bytes
		return make([]byte, 4)
	},
}

// EncodeSessionInit encodes a SessionInitPayload into a CBOR byte slice.
func EncodeSessionInit(payload *SessionInitPayload) ([]byte, error) {
	data, err := cbor.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SessionInit payload: %w", err)
	}
	return data, nil
}

// DecodeSessionInit decodes a CBOR byte slice into a SessionInitPayload.
func DecodeSessionInit(data []byte) (*SessionInitPayload, error) {
	var payload SessionInitPayload
	if err := cbor.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SessionInit payload: %w", err)
	}
	return &payload, nil
}

// EncodeSessionAccept encodes a SessionAcceptPayload into a CBOR byte slice.
func EncodeSessionAccept(payload *SessionAcceptPayload) ([]byte, error) {
	data, err := cbor.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SessionAccept payload: %w", err)
	}
	return data, nil
}

// DecodeSessionAccept decodes a CBOR byte slice into a SessionAcceptPayload.
func DecodeSessionAccept(data []byte) (*SessionAcceptPayload, error) {
	var payload SessionAcceptPayload
	if err := cbor.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SessionAccept payload: %w", err)
	}
	return &payload, nil
}

// EncodeStreamType encodes a StreamTypePayload into a CBOR byte slice.
func EncodeStreamType(payload *StreamTypePayload) ([]byte, error) {
	data, err := cbor.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal StreamType payload: %w", err)
	}
	return data, nil
}

// DecodeStreamType decodes a CBOR byte slice into a StreamTypePayload.
func DecodeStreamType(data []byte) (*StreamTypePayload, error) {
	var payload StreamTypePayload
	if err := cbor.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal StreamType payload: %w", err)
	}
	return &payload, nil
}

// WriteMessage writes a message with a length prefix
func WriteMessage(stream io.Writer, msg *Message) error {
	data, err := cbor.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	// Get length buffer from pool
	lengthBuf := lengthBufPool.Get().([]byte)
	defer lengthBufPool.Put(lengthBuf)
	
	// Write length prefix (4 bytes, big endian)
	length := uint32(len(data))
	binary.BigEndian.PutUint32(lengthBuf, length)
	
	if _, err := stream.Write(lengthBuf); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}
	
	// Write message data
	if _, err := stream.Write(data); err != nil {
		return fmt.Errorf("failed to write message data: %w", err)
	}
	
	return nil
}

// ReadMessage reads a message with a length prefix
func ReadMessage(stream io.Reader) (*Message, error) {
	// Get length buffer from pool
	lengthBuf := lengthBufPool.Get().([]byte)
	defer lengthBufPool.Put(lengthBuf)
	
	// Read length prefix (4 bytes)
	if _, err := io.ReadFull(stream, lengthBuf); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}
	
	length := binary.BigEndian.Uint32(lengthBuf)
	
	// Validate length to prevent excessive memory allocation
	if length > 1024*1024 { // 1MB limit
		return nil, fmt.Errorf("message too large: %d bytes", length)
	}
	
	// Read message data
	data := make([]byte, length)
	if _, err := io.ReadFull(stream, data); err != nil {
		return nil, fmt.Errorf("failed to read message data: %w", err)
	}
	
	// Unmarshal message
	var msg Message
	if err := cbor.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	
	return &msg, nil
}
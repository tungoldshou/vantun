package core

import (
	"context"
	"fmt"

	"github.com/quic-go/quic-go"
)

// Stream types
const (
	StreamTypeInteractive = 1
	StreamTypeBulk        = 2
	StreamTypeTelemetry   = 3
)

// OpenInteractiveStream opens a new interactive stream.
func (s *Session) OpenInteractiveStream(ctx context.Context) (quic.Stream, error) {
	// Stream type 1 is for interactive data.
	stream, err := s.conn.OpenStreamSync(ctx)
	if err != nil {
		Error("Failed to open interactive stream: %v", err)
		return nil, err
	}
	
	// Send stream type identifier on the stream.
	payload := &StreamTypePayload{
		Type: StreamTypeInteractive,
	}
	data, err := EncodeStreamType(payload)
	if err != nil {
		stream.Close()
		Error("Failed to encode stream type: %v", err)
		return nil, fmt.Errorf("failed to encode stream type: %w", err)
	}
	
	msg := &Message{
		Type: StreamType,
		Data: data,
	}
	
	if err := WriteMessage(stream, msg); err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to send stream type: %w", err)
	}
	
	return stream, nil
}

// AcceptInteractiveStream accepts a new interactive stream.
func (s *Session) AcceptInteractiveStream(ctx context.Context) (quic.Stream, error) {
	stream, err := s.conn.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	
	// Read and verify stream type identifier from the stream.
	msg, err := ReadMessage(stream)
	if err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to read stream type: %w", err)
	}
	
	if msg.Type != StreamType {
		stream.Close()
		return nil, fmt.Errorf("expected StreamType message, got %d", msg.Type)
	}
	
	payload, err := DecodeStreamType(msg.Data)
	if err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to decode stream type payload: %w", err)
	}
	
	if payload.Type != StreamTypeInteractive {
		stream.Close()
		return nil, fmt.Errorf("expected interactive stream type, got %d", payload.Type)
	}
	
	return stream, nil
}

// OpenBulkStream opens a new bulk stream.
func (s *Session) OpenBulkStream(ctx context.Context) (quic.Stream, error) {
	// Stream type 2 is for bulk data.
	stream, err := s.conn.OpenStreamSync(ctx)
	if err != nil {
		Error("Failed to open bulk stream: %v", err)
		return nil, err
	}
	
	// Send stream type identifier on the stream.
	payload := &StreamTypePayload{
		Type: StreamTypeBulk,
	}
	data, err := EncodeStreamType(payload)
	if err != nil {
		stream.Close()
		Error("Failed to encode stream type: %v", err)
		return nil, fmt.Errorf("failed to encode stream type: %w", err)
	}
	
	msg := &Message{
		Type: StreamType,
		Data: data,
	}
	
	if err := WriteMessage(stream, msg); err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to send stream type: %w", err)
	}
	
	return stream, nil
}

// AcceptBulkStream accepts a new bulk stream.
func (s *Session) AcceptBulkStream(ctx context.Context) (quic.Stream, error) {
	stream, err := s.conn.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	
	// Read and verify stream type identifier from the stream.
	msg, err := ReadMessage(stream)
	if err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to read stream type: %w", err)
	}
	
	if msg.Type != StreamType {
		stream.Close()
		return nil, fmt.Errorf("expected StreamType message, got %d", msg.Type)
	}
	
	payload, err := DecodeStreamType(msg.Data)
	if err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to decode stream type payload: %w", err)
	}
	
	if payload.Type != StreamTypeBulk {
		stream.Close()
		return nil, fmt.Errorf("expected bulk stream type, got %d", payload.Type)
	}
	
	return stream, nil
}

// OpenTelemetryStream opens a new telemetry stream.
func (s *Session) OpenTelemetryStream(ctx context.Context) (quic.Stream, error) {
	// Stream type 3 is for telemetry data.
	stream, err := s.conn.OpenStreamSync(ctx)
	if err != nil {
		Error("Failed to open telemetry stream: %v", err)
		return nil, err
	}
	
	// Send stream type identifier on the stream.
	payload := &StreamTypePayload{
		Type: StreamTypeTelemetry,
	}
	data, err := EncodeStreamType(payload)
	if err != nil {
		stream.Close()
		Error("Failed to encode stream type: %v", err)
		return nil, fmt.Errorf("failed to encode stream type: %w", err)
	}
	
	msg := &Message{
		Type: StreamType,
		Data: data,
	}
	
	if err := WriteMessage(stream, msg); err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to send stream type: %w", err)
	}
	
	return stream, nil
}

// AcceptTelemetryStream accepts a new telemetry stream.
func (s *Session) AcceptTelemetryStream(ctx context.Context) (quic.Stream, error) {
	stream, err := s.conn.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	
	// Read and verify stream type identifier from the stream.
	msg, err := ReadMessage(stream)
	if err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to read stream type: %w", err)
	}
	
	if msg.Type != StreamType {
		stream.Close()
		return nil, fmt.Errorf("expected StreamType message, got %d", msg.Type)
	}
	
	payload, err := DecodeStreamType(msg.Data)
	if err != nil {
		stream.Close()
		return nil, fmt.Errorf("failed to decode stream type payload: %w", err)
	}
	
	if payload.Type != StreamTypeTelemetry {
		stream.Close()
		return nil, fmt.Errorf("expected telemetry stream type, got %d", payload.Type)
	}
	
	return stream, nil
}
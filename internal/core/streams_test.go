package core

import (
	"context"
	"testing"
)

func TestSessionOpenInteractiveStream(t *testing.T) {
	// Test case 1: Opening an interactive stream
	ctx := context.Background()
	conn := &MockQUICConnection{}
	session := &Session{conn: conn}

	stream, err := session.OpenInteractiveStream(ctx)
	if err != nil {
		t.Errorf("Failed to open interactive stream: %v", err)
	}

	if stream == nil {
		t.Error("OpenInteractiveStream returned nil stream")
	}
}

func TestSessionAcceptInteractiveStream(t *testing.T) {
	// Test case 2: Accepting an interactive stream
	ctx := context.Background()
	mockStream := &MockQUICStream{}
	conn := &MockQUICConnection{}
	conn.AddStream(mockStream)
	session := &Session{conn: conn}

	stream, err := session.AcceptInteractiveStream(ctx)
	if err != nil {
		t.Errorf("Failed to accept interactive stream: %v", err)
	}

	if stream == nil {
		t.Error("AcceptInteractiveStream returned nil stream")
	}
}

func TestSessionOpenBulkStream(t *testing.T) {
	// Test case 3: Opening a bulk stream
	ctx := context.Background()
	conn := &MockQUICConnection{}
	session := &Session{conn: conn}

	stream, err := session.OpenBulkStream(ctx)
	if err != nil {
		t.Errorf("Failed to open bulk stream: %v", err)
	}

	if stream == nil {
		t.Error("OpenBulkStream returned nil stream")
	}
}

func TestSessionAcceptBulkStream(t *testing.T) {
	// Test case 4: Accepting a bulk stream
	ctx := context.Background()
	mockStream := &MockQUICStream{}
	conn := &MockQUICConnection{}
	conn.AddStream(mockStream)
	session := &Session{conn: conn}

	stream, err := session.AcceptBulkStream(ctx)
	if err != nil {
		t.Errorf("Failed to accept bulk stream: %v", err)
	}

	if stream == nil {
		t.Error("AcceptBulkStream returned nil stream")
	}
}

func TestStreamDataExchange(t *testing.T) {
	// Test case 5: Data exchange through streams
	ctx := context.Background()
	mockStream := &MockQUICStream{}
	conn := &MockQUICConnection{}
	conn.AddStream(mockStream)
	session := &Session{conn: conn}

	// Accept a stream
	stream, err := session.AcceptInteractiveStream(ctx)
	if err != nil {
		t.Fatalf("Failed to accept interactive stream: %v", err)
	}

	// Write data to the stream
	testData := []byte("Hello, VANTUN!")
	n, err := stream.Write(testData)
	if err != nil {
		t.Errorf("Failed to write to stream: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	// Verify data was written to mock stream
	if string(mockStream.writeData) != string(testData) {
		t.Errorf("Expected written data %s, got %s", string(testData), string(mockStream.writeData))
	}

	// Set up mock stream to return data for reading
	mockStream.readData = testData

	// Read data from the stream
	buf := make([]byte, len(testData))
	n, err = stream.Read(buf)
	if err != nil {
		t.Errorf("Failed to read from stream: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to read %d bytes, read %d", len(testData), n)
	}

	// Verify data was read correctly
	if string(buf) != string(testData) {
		t.Errorf("Expected read data %s, got %s", string(testData), string(buf))
	}
}
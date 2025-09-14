package core

import (
	"testing"
)

func TestFECSimple(t *testing.T) {
	// Test with a simple case
	k, m := 3, 2
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create test data
	data := []byte("Hello, World! This is a test message.")
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Verify shard count
	expectedShards := k + m
	if len(shards) != expectedShards {
		t.Errorf("Expected %d shards, got %d", expectedShards, len(shards))
	}

	// Decode the shards back to data
	decodedData, err := fec.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData) != string(data) {
		t.Errorf("Decoded data does not match original. Got %s, expected %s", 
			string(decodedData), string(data))
	}
}
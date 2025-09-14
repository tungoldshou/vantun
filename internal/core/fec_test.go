package core

import (
	"testing"
)

func TestFECEncodeDecode(t *testing.T) {
	// Test case 1: Basic encoding and decoding
	k, m := 10, 3
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create test data
	data := []byte("This is a test message for FEC encoding and decoding.")
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Verify shard count
	expectedShards := k + m
	if len(shards) != expectedShards {
		t.Errorf("Expected %d shards, got %d", expectedShards, len(shards))
	}

	// Verify shard sizes
	shardSize := len(shards[0])
	for i, shard := range shards {
		if len(shard) != shardSize {
			t.Errorf("Shard %d has size %d, expected %d", i, len(shard), shardSize)
		}
	}

	// Decode the shards back to data
	decodedData, err := fec.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData[:len(data)]) != string(data) {
		t.Errorf("Decoded data does not match original. Got %s, expected %s", 
			string(decodedData[:len(data)]), string(data))
	}
}

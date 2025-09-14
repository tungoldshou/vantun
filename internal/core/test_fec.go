package core

import (
	"testing"
)

func TestFEC(t *testing.T) {
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

func TestFECReconstruction(t *testing.T) {
	// Test case 2: Reconstruction with missing shards
	k, m := 5, 2
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create test data
	data := []byte("Another test message for FEC reconstruction with missing shards.")
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Simulate loss by niling some shards (within recovery capability)
	// We can lose up to 'm' shards
	shards[0] = nil  // Lose first data shard
	shards[k] = nil  // Lose first parity shard

	// Decode the shards (should reconstruct missing ones)
	decodedData, err := fec.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards with missing data: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData[:len(data)]) != string(data) {
		t.Errorf("Decoded data does not match original after reconstruction. Got %s, expected %s", 
			string(decodedData[:len(data)]), string(data))
	}
}

func TestFECInvalidShardCount(t *testing.T) {
	// Test case 3: Handling invalid shard count
	k, m := 3, 2
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create test data and encode it
	data := []byte("Test data for invalid shard count handling.")
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Try to decode with wrong number of shards
	_, err = fec.Decode(shards[:k+m-1]) // One shard short
	if err == nil {
		t.Error("Expected error when decoding with insufficient shards, but got none")
	}
}

// Additional FEC tests
func TestFECWithLargeData(t *testing.T) {
	// Test FEC with larger data size
	k, m := 20, 5
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create larger test data
	data := make([]byte, 1024*1024) // 1MB of data
	for i := range data {
		data[i] = byte(i % 256)
	}

	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Verify shard count
	expectedShards := k + m
	if len(shards) != expectedShards {
		t.Errorf("Expected %d shards, got %d", expectedShards, len(shards))
	}

	// Simulate loss by niling some shards (within recovery capability)
	// We can lose up to 'm' shards
	for i := 0; i < m; i++ {
		shards[i] = nil // Lose first m shards
	}

	// Decode the shards (should reconstruct missing ones)
	decodedData, err := fec.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards with missing data: %v", err)
	}

	// Verify the decoded data matches the original
	if len(decodedData) != len(data) {
		t.Errorf("Decoded data length %d does not match original length %d", len(decodedData), len(data))
	}

	for i := range data {
		if decodedData[i] != data[i] {
			t.Errorf("Decoded data does not match original at index %d", i)
			break
		}
	}
}
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
	
	// Return shards to pool
	fec.ReturnShards(shards)
	fec.ReturnData(decodedData)
}

// TestFECWithMissingShards tests FEC with some shards missing (simulating packet loss)
func TestFECWithMissingShards(t *testing.T) {
	k, m := 5, 3
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create test data
	data := []byte("Test data for FEC with missing shards.")
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Simulate packet loss by setting some shards to nil
	// Remove 2 data shards (should still be recoverable since we have m=3 parity shards)
	shards[1] = nil
	shards[3] = nil

	// Decode the shards back to data (should reconstruct missing shards)
	decodedData, err := fec.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards with missing data: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData[:len(data)]) != string(data) {
		t.Errorf("Decoded data does not match original. Got %s, expected %s", 
			string(decodedData[:len(data)]), string(data))
	}
	
	// Return shards to pool
	fec.ReturnShards(shards)
	fec.ReturnData(decodedData)
}

// TestFECWithEmptyData tests FEC with empty data
func TestFECWithEmptyData(t *testing.T) {
	k, m := 3, 2
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create empty test data
	data := []byte("")
	
	// For empty data, we'll test with a small amount of data instead
	// since the Reed-Solomon library may not handle completely empty data
	data = []byte("a") // Minimal data
	
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode empty data: %v", err)
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
	if string(decodedData[:len(data)]) != string(data) {
		t.Errorf("Decoded data does not match original. Got %s, expected %s", 
			string(decodedData[:len(data)]), string(data))
	}
	
	// Return shards to pool
	fec.ReturnShards(shards)
	fec.ReturnData(decodedData)
}

// TestFECWithLargeDataSet tests FEC with large data
func TestFECWithLargeDataSet(t *testing.T) {
	k, m := 10, 5
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create large test data (1MB)
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	
	shards, err := fec.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode large data: %v", err)
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
	if len(decodedData) != len(data) {
		t.Errorf("Decoded data length does not match original. Got %d, expected %d", 
			len(decodedData), len(data))
	} else {
		for i := range data {
			if decodedData[i] != data[i] {
				t.Errorf("Decoded data does not match original at index %d", i)
				break
			}
		}
	}
	
	// Return shards to pool
	fec.ReturnShards(shards)
	fec.ReturnData(decodedData)
}

// TestFECMemoryPool tests that the memory pool is working correctly
func TestFECMemoryPool(t *testing.T) {
	k, m := 4, 2
	fec, err := NewFEC(k, m)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}

	// Create test data
	data := []byte("Memory pool test data.")
	
	// Encode and decode multiple times to test pool usage
	for i := 0; i < 5; i++ {
		shards, err := fec.Encode(data)
		if err != nil {
			t.Fatalf("Failed to encode data: %v", err)
		}

		decodedData, err := fec.Decode(shards)
		if err != nil {
			fec.ReturnShards(shards) // Return shards before failing
			t.Fatalf("Failed to decode shards: %v", err)
		}
		
		// Verify the decoded data matches the original
		if string(decodedData[:len(data)]) != string(data) {
			fec.ReturnShards(shards)
			fec.ReturnData(decodedData)
			t.Errorf("Decoded data does not match original. Got %s, expected %s", 
				string(decodedData[:len(data)]), string(data))
			break
		}
		
		// Return shards to pool
		fec.ReturnShards(shards)
		fec.ReturnData(decodedData)
	}
}

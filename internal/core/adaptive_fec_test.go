package core

import (
	"testing"
	"time"
)

func TestAdaptiveFECAdjust(t *testing.T) {
	// Test case 1: Basic adjustment with high loss
	k, m := 10, 3
	minParity, maxParity := 1, 10
	adaptiveFEC, err := NewAdaptiveFEC(k, m, minParity, maxParity)
	if err != nil {
		t.Fatalf("Failed to create AdaptiveFEC: %v", err)
	}

	// Test with high loss data (should increase parity shards)
	highLossData := &TelemetryData{
		RTT:          50 * time.Millisecond,
		Loss:         0.15, // 15% loss
		Bandwidth:    1000000,
		DeliveryRate: 1000000, // Match bandwidth for 100% efficiency
		Timestamp:    time.Now(),
	}

	previousM := adaptiveFEC.m
	err = adaptiveFEC.Adjust(highLossData)
	if err != nil {
		t.Fatalf("Failed to adjust FEC: %v", err)
	}

	// Should have increased parity shards
	if adaptiveFEC.m <= previousM {
		t.Errorf("Expected parity shards to increase with high loss, but went from %d to %d", previousM, adaptiveFEC.m)
	}

	// Record the new value after high loss adjustment
	highLossM := adaptiveFEC.m

	// Test with low loss data (should decrease parity shards)
	lowLossData := &TelemetryData{
		RTT:          50 * time.Millisecond,
		Loss:         0.005, // 0.5% loss
		Bandwidth:    1000000,
		DeliveryRate: 1000000, // Match bandwidth for 100% efficiency
		Timestamp:    time.Now(),
	}

	err = adaptiveFEC.Adjust(lowLossData)
	if err != nil {
		t.Fatalf("Failed to adjust FEC: %v", err)
	}

	// Should have decreased parity shards from the high loss value
	if adaptiveFEC.m >= highLossM {
		t.Errorf("Expected parity shards to decrease with low loss, but went from %d to %d", highLossM, adaptiveFEC.m)
	}
}

func TestAdaptiveFECBounds(t *testing.T) {
	// Test case 2: Verify FEC stays within bounds
	k, m := 5, 3
	minParity, maxParity := 2, 5
	adaptiveFEC, err := NewAdaptiveFEC(k, m, minParity, maxParity)
	if err != nil {
		t.Fatalf("Failed to create AdaptiveFEC: %v", err)
	}

	// Test with very high loss that would exceed maxParity
	veryHighLossData := &TelemetryData{
		RTT:          50 * time.Millisecond,
		Loss:         0.5, // 50% loss
		Bandwidth:    1000000,
		DeliveryRate: 1000000, // Match bandwidth for 100% efficiency
		Timestamp:    time.Now(),
	}

	// Apply adjustment multiple times to try to exceed bounds
	for i := 0; i < 10; i++ {
		err = adaptiveFEC.Adjust(veryHighLossData)
		if err != nil {
			t.Fatalf("Failed to adjust FEC: %v", err)
		}
	}

	// Verify m doesn't exceed maxParity
	if adaptiveFEC.m > maxParity {
		t.Errorf("Parity shards %d exceeded maximum %d", adaptiveFEC.m, maxParity)
	}

	// Test with very low loss that would go below minParity
	veryLowLossData := &TelemetryData{
		RTT:          50 * time.Millisecond,
		Loss:         0.001, // 0.1% loss
		Bandwidth:    1000000,
		DeliveryRate: 1000000, // Match bandwidth for 100% efficiency
		Timestamp:    time.Now(),
	}

	// Apply adjustment multiple times to try to go below bounds
	for i := 0; i < 10; i++ {
		err = adaptiveFEC.Adjust(veryLowLossData)
		if err != nil {
			t.Fatalf("Failed to adjust FEC: %v", err)
		}
	}

	// Verify m doesn't go below minParity
	if adaptiveFEC.m < minParity {
		t.Errorf("Parity shards %d below minimum %d", adaptiveFEC.m, minParity)
	}
}

func TestAdaptiveFECEncodeDecode(t *testing.T) {
	// Test case 3: Verify encoding/decoding works with adaptive FEC
	k, m := 5, 2
	minParity, maxParity := 1, 5
	adaptiveFEC, err := NewAdaptiveFEC(k, m, minParity, maxParity)
	if err != nil {
		t.Fatalf("Failed to create AdaptiveFEC: %v", err)
	}

	// Create test data
	data := []byte("Test data for adaptive FEC encoding and decoding.")

	// Encode with initial FEC settings
	shards, err := adaptiveFEC.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Verify shard count
	expectedShards := k + adaptiveFEC.m
	if len(shards) != expectedShards {
		t.Errorf("Expected %d shards, got %d", expectedShards, len(shards))
	}

	// Decode with the same FEC instance (should work)
	decodedData, err := adaptiveFEC.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData) != string(data) {
		t.Errorf("Decoded data does not match original. Got %s, expected %s", 
			string(decodedData), string(data))
	}

	// Adjust FEC parameters
	highLossData := &TelemetryData{
		RTT:          50 * time.Millisecond,
		Loss:         0.2, // 20% loss
		Bandwidth:    1000000,
		DeliveryRate: 1000000, // Match bandwidth for 100% efficiency
		Timestamp:    time.Now(),
	}

	err = adaptiveFEC.Adjust(highLossData)
	if err != nil {
		t.Fatalf("Failed to adjust FEC: %v", err)
	}

	// Encode with adjusted FEC settings
	shards2, err := adaptiveFEC.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data with adjusted FEC: %v", err)
	}

	// Verify shard count changed according to new settings
	expectedShards2 := k + adaptiveFEC.m
	if len(shards2) != expectedShards2 {
		t.Errorf("Expected %d shards after adjustment, got %d", expectedShards2, len(shards2))
	}

	// Decode the new shards with the adjusted FEC (should work)
	decodedData2, err := adaptiveFEC.Decode(shards2)
	if err != nil {
		t.Fatalf("Failed to decode adjusted shards: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData2) != string(data) {
		t.Errorf("Decoded data from adjusted FEC does not match original. Got %s, expected %s", 
			string(decodedData2), string(data))
	}
}
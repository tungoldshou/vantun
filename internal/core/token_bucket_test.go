package core

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	// Create a token bucket with 100 tokens per second rate and 200 token capacity
	bucket := NewTokenBucket(100, 200)

	// Test initial state - bucket should be full
	if !bucket.Consume(200) {
		t.Error("Expected to be able to consume full capacity initially")
	}

	// Try to consume more tokens than available
	if bucket.Consume(1) {
		t.Error("Expected consume to fail when bucket is empty")
	}

	// Wait for some tokens to be added
	time.Sleep(100 * time.Millisecond)

	// Try to consume tokens again
	// We should have about 10 tokens (100 tokens/sec * 0.1 sec)
	if !bucket.Consume(10) {
		t.Error("Expected to be able to consume tokens after waiting")
	}

	// Test rate limiting
	// Wait for 1 second to accumulate tokens
	time.Sleep(1 * time.Second)

	// We should have about 100 tokens (up to capacity of 200)
	if !bucket.Consume(100) {
		t.Error("Expected to be able to consume 100 tokens after 1 second")
	}
}

func TestTokenBucketRate(t *testing.T) {
	// Create a token bucket with 50 tokens per second rate and 100 token capacity
	bucket := NewTokenBucket(50, 100)

	// Test GetRate
	if bucket.GetRate() != 50 {
		t.Errorf("Expected rate 50, got %f", bucket.GetRate())
	}

	// Test SetRate
	bucket.SetRate(75)
	if bucket.GetRate() != 75 {
		t.Errorf("Expected rate 75 after setting, got %f", bucket.GetRate())
	}
}

func TestAdaptiveFEC(t *testing.T) {
	// Create an adaptive FEC with 10 data shards and 3 parity shards
	adaptiveFEC, err := NewAdaptiveFEC(10, 3, 1, 5)
	if err != nil {
		t.Fatalf("Failed to create adaptive FEC: %v", err)
	}

	// Test basic encode/decode
	data := []byte("This is a test message for adaptive FEC.")
	shards, err := adaptiveFEC.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}

	// Verify shard count
	expectedShards := 10 + 3
	if len(shards) != expectedShards {
		t.Errorf("Expected %d shards, got %d", expectedShards, len(shards))
	}

	// Decode the shards back to data
	decodedData, err := adaptiveFEC.Decode(shards)
	if err != nil {
		t.Fatalf("Failed to decode shards: %v", err)
	}

	// Verify the decoded data matches the original
	if string(decodedData[:len(data)]) != string(data) {
		t.Errorf("Decoded data does not match original. Got %s, expected %s", 
			string(decodedData[:len(data)]), string(data))
	}

	// Test adjustment with high loss
	telemetryHighLoss := &TelemetryData{
		Loss: 0.15, // 15% loss
	}

	err = adaptiveFEC.Adjust(telemetryHighLoss)
	if err != nil {
		t.Fatalf("Failed to adjust FEC: %v", err)
	}

	// Should have increased parity shards
	if adaptiveFEC.m <= 3 {
		t.Errorf("Expected parity shards to increase with high loss, got %d", adaptiveFEC.m)
	}

	// Test adjustment with low loss
	telemetryLowLoss := &TelemetryData{
		Loss: 0.005, // 0.5% loss
	}

	err = adaptiveFEC.Adjust(telemetryLowLoss)
	if err != nil {
		t.Fatalf("Failed to adjust FEC: %v", err)
	}

	// Should have decreased parity shards (but not below min)
	if adaptiveFEC.m < 1 {
		t.Errorf("Expected parity shards to not go below minimum, got %d", adaptiveFEC.m)
	}
}
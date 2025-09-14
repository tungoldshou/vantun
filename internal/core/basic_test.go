package core

import (
	"testing"
)

func TestBasicFunctionality(t *testing.T) {
	// Test that basic types and functions are available
	t.Log("Basic functionality test")

	// Test FEC
	fec, err := NewFEC(10, 3)
	if err != nil {
		t.Fatalf("Failed to create FEC: %v", err)
	}
	if fec == nil {
		t.Error("Expected non-nil FEC")
	}

	// Test AdaptiveFEC
	adaptiveFEC, err := NewAdaptiveFEC(10, 3, 1, 5)
	if err != nil {
		t.Fatalf("Failed to create AdaptiveFEC: %v", err)
	}
	if adaptiveFEC == nil {
		t.Error("Expected non-nil AdaptiveFEC")
	}

	// Test TokenBucket
	tokenBucket := NewTokenBucket(1000, 5000)
	if tokenBucket == nil {
		t.Error("Expected non-nil TokenBucket")
	}

	// Test DataSplitter
	dataSplitter := NewDataSplitter(1024)
	if dataSplitter == nil {
		t.Error("Expected non-nil DataSplitter")
	}

	// Test PathSelector
	pathSelector := NewPathSelector(RoundRobinStrategy)
	if pathSelector == nil {
		t.Error("Expected non-nil PathSelector")
	}

	t.Log("All basic functionality tests passed")
}
package core

import (
	"bytes"
	"math/rand"
	"testing"
)

// NewTestRand creates a new random number generator with a fixed seed for testing
func NewTestRand() *rand.Rand {
	return rand.New(rand.NewSource(42)) // Fixed seed for reproducible tests
}

func TestHTTP3ObfuscatorRoundTripSimple(t *testing.T) {
	// Create HTTP3 obfuscator with fixed seed for reproducible results
	frameTypes := []byte{0x0, 0x1, 0x2}
	obfuscator := &HTTP3Obfuscator{
		frameTypes: frameTypes,
		minPadding: 10,
		maxPadding: 100,
		// Use a fixed seed for testing
		rng: NewTestRand(),
	}

	// Test data
	testData := []byte("Hello, VANTUN!")

	// Obfuscate data
	obfuscatedData, err := obfuscator.Obfuscate(testData)
	if err != nil {
		t.Fatalf("Failed to obfuscate data: %v", err)
	}

	// Verify obfuscated data is different from original
	if bytes.Equal(obfuscatedData, testData) {
		t.Error("Obfuscated data is the same as original data")
	}

	// Verify obfuscated data is not empty
	if len(obfuscatedData) == 0 {
		t.Error("Obfuscated data is empty")
	}

	// Deobfuscate data using the same obfuscator instance to ensure consistency
	deobfuscatedData, err := obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Fatalf("Failed to deobfuscate data: %v", err)
	}

	// Verify round-trip preserves original data
	if !bytes.Equal(deobfuscatedData, testData) {
		t.Errorf("Round-trip failed: got %s, expected %s", 
			string(deobfuscatedData), string(testData))
	}
}
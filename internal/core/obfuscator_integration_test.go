package core

import (
	"bytes"
	"testing"
)

func TestObfuscatorIntegration(t *testing.T) {
	// Create obfuscator
	config := ObfuscatorConfig{
		Enabled:    true,
		MinPadding: 10,
		MaxPadding: 100,
	}
	obfuscator := NewObfuscator(config)

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

	// Deobfuscate data
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
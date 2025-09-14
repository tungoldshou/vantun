package core

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHTTP3ObfuscatorDebug(t *testing.T) {
	// Create HTTP3 obfuscator
	frameTypes := []byte{0x0, 0x1, 0x2}
	obfuscator := NewHTTP3Obfuscator(frameTypes, 10, 100)

	// Test data
	testData := []byte("Hello, VANTUN!")

	// Obfuscate data
	obfuscatedData, err := obfuscator.Obfuscate(testData)
	if err != nil {
		t.Fatalf("Failed to obfuscate data: %v", err)
	}

	// Print obfuscated data for debugging
	fmt.Printf("Original data: %s\n", testData)
	fmt.Printf("Obfuscated data length: %d\n", len(obfuscatedData))
	fmt.Printf("Obfuscated data: %v\n", obfuscatedData)

	// Deobfuscate data
	deobfuscatedData, err := obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Fatalf("Failed to deobfuscate data: %v", err)
	}

	// Print deobfuscated data for debugging
	fmt.Printf("Deobfuscated data: %s\n", deobfuscatedData)

	// Verify round-trip preserves original data
	if !bytes.Equal(deobfuscatedData, testData) {
		t.Errorf("Round-trip failed: got %s, expected %s", 
			string(deobfuscatedData), string(testData))
	}
}
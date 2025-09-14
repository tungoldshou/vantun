package core

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestObfuscatorObfuscate(t *testing.T) {
	// Test case 1: Basic obfuscation
	config := ObfuscatorConfig{
		Enabled: true,
	}
	obfuscator := NewObfuscator(config)
	data := []byte("This is test data for obfuscation.")

	obfuscatedData, err := obfuscator.Obfuscate(data)
	if err != nil {
		t.Errorf("Obfuscator.Obfuscate failed: %v", err)
	}

	// Verify obfuscated data is different from original (due to obfuscation)
	if bytes.Equal(obfuscatedData, data) {
		t.Errorf("Expected obfuscated data to be different from original")
	}
}

func TestObfuscatorDeobfuscate(t *testing.T) {
	// Test case 2: Basic deobfuscation
	config := ObfuscatorConfig{
		Enabled: true,
	}
	obfuscator := NewObfuscator(config)
	originalData := []byte("This is test data for deobfuscation.")

	// First obfuscate the data
	obfuscatedData, err := obfuscator.Obfuscate(originalData)
	if err != nil {
		t.Fatalf("Obfuscator.Obfuscate failed: %v", err)
	}

	// Then deobfuscate it
	deobfuscatedData, err := obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("Obfuscator.Deobfuscate failed: %v", err)
	}

	// Verify the round-trip preserves the original data
	if !bytes.Equal(deobfuscatedData, originalData) {
		t.Errorf("Round-trip failed: got %s, expected %s", 
			string(deobfuscatedData), string(originalData))
	}
}

func TestHTTP3ObfuscatorObfuscate(t *testing.T) {
	// Test case 3: HTTP/3 obfuscation
	http3Obfuscator := NewHTTP3Obfuscator([]byte{0x0, 0x1, 0x2}, 10, 100)
	data := []byte("This is test data for HTTP/3 obfuscation.")

	obfuscatedData, err := http3Obfuscator.Obfuscate(data)
	if err != nil {
		t.Errorf("HTTP3Obfuscator.Obfuscate failed: %v", err)
	}

	// Verify obfuscated data exists
	if len(obfuscatedData) == 0 {
		t.Error("HTTP3Obfuscator.Obfuscate returned empty data")
	}

	// Verify obfuscated data is different from original
	if bytes.Equal(obfuscatedData, data) {
		t.Error("HTTP3Obfuscator.Obfuscate did not modify the data")
	}
}

func TestHTTP3ObfuscatorDeobfuscate(t *testing.T) {
	// Test case 4: HTTP/3 deobfuscation
	http3Obfuscator := NewHTTP3Obfuscator([]byte{0x0, 0x1, 0x2}, 10, 100)
	originalData := []byte("This is test data for HTTP/3 deobfuscation.")

	// First obfuscate the data
	obfuscatedData, err := http3Obfuscator.Obfuscate(originalData)
	if err != nil {
		t.Fatalf("HTTP3Obfuscator.Obfuscate failed: %v", err)
	}

	// Then deobfuscate it
	deobfuscatedData, err := http3Obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("HTTP3Obfuscator.Deobfuscate failed: %v", err)
	}

	// Verify the round-trip preserves the original data
	if !bytes.Equal(deobfuscatedData, originalData) {
		t.Errorf("Round-trip failed: got %s, expected %s", 
			string(deobfuscatedData), string(originalData))
	}
}

func TestHTTP3ObfuscatorRoundTrip(t *testing.T) {
	// Test case 5: Complete round-trip with various data sizes
	http3Obfuscator := NewHTTP3Obfuscator([]byte{0x0, 0x1, 0x2}, 10, 100)
	
	// Test with different data sizes
	testCases := []struct {
		name string
		data []byte
	}{
		{"Small data", []byte("Hello")},
		{"Medium data", []byte("This is a medium-sized message for testing purposes.")},
		{"Large data", make([]byte, 2048)}, // 2KB of zero bytes
	}
	
	// Fill the large data with random values
	for i := range testCases[2].data {
		testCases[2].data[i] = byte(rand.Intn(256))
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Obfuscate
			obfuscatedData, err := http3Obfuscator.Obfuscate(tc.data)
			if err != nil {
				t.Fatalf("HTTP3Obfuscator.Obfuscate failed: %v", err)
			}
			
			// Deobfuscate
			deobfuscatedData, err := http3Obfuscator.Deobfuscate(obfuscatedData)
			if err != nil {
				t.Errorf("HTTP3Obfuscator.Deobfuscate failed: %v", err)
			}
			
			// Verify the round-trip preserves the original data
			if !bytes.Equal(deobfuscatedData, tc.data) {
				t.Errorf("Round-trip failed for %s: got %d bytes, expected %d bytes", 
					tc.name, len(deobfuscatedData), len(tc.data))
			}
		})
	}
}

func TestHTTP3ObfuscatorEdgeCases(t *testing.T) {
	// Test case 6: Edge cases
	http3Obfuscator := NewHTTP3Obfuscator([]byte{0x0, 0x1, 0x2}, 10, 100)
	
	// Test with empty data
	emptyData := []byte{}
	obfuscatedData, err := http3Obfuscator.Obfuscate(emptyData)
	if err != nil {
		t.Errorf("HTTP3Obfuscator.Obfuscate failed with empty data: %v", err)
	}
	
	deobfuscatedData, err := http3Obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("HTTP3Obfuscator.Deobfuscate failed with empty data: %v", err)
	}
	
	if len(deobfuscatedData) != 0 {
		t.Error("Expected empty data after round-trip of empty data")
	}
	
	// Test with single byte
	singleByte := []byte{42}
	obfuscatedData, err = http3Obfuscator.Obfuscate(singleByte)
	if err != nil {
		t.Errorf("HTTP3Obfuscator.Obfuscate failed with single byte: %v", err)
	}
	
	deobfuscatedData, err = http3Obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("HTTP3Obfuscator.Deobfuscate failed with single byte: %v", err)
	}
	
	if !bytes.Equal(deobfuscatedData, singleByte) {
		t.Errorf("Round-trip failed for single byte: got %v, expected %v", 
			deobfuscatedData, singleByte)
	}
}
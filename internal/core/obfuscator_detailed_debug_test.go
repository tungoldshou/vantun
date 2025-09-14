package core

import (
	"bytes"
	"fmt"
	"testing"
)

func TestObfuscatorDetailedDebug(t *testing.T) {
	// Create HTTP3 obfuscator directly
	frameTypes := []byte{0x0, 0x1, 0x2}
	http3Obfs := NewHTTP3Obfuscator(frameTypes, 10, 100)

	// Create obfuscator through wrapper
	config := ObfuscatorConfig{
		Enabled:    true,
		MinPadding: 10,
		MaxPadding: 100,
	}
	obfuscator := NewObfuscator(config)

	// Test data
	testData := []byte("Hello, VANTUN!")

	// Test HTTP3 obfuscator directly
	fmt.Println("=== Testing HTTP3Obfuscator directly ===")
	http3Obfuscated, err := http3Obfs.Obfuscate(testData)
	if err != nil {
		t.Fatalf("Failed to obfuscate data with HTTP3Obfuscator: %v", err)
	}

	fmt.Printf("Original data: %s\n", testData)
	fmt.Printf("HTTP3 obfuscated length: %d\n", len(http3Obfuscated))
	fmt.Printf("HTTP3 obfuscated data: %v\n", http3Obfuscated)

	http3Deobfuscated, err := http3Obfs.Deobfuscate(http3Obfuscated)
	if err != nil {
		t.Fatalf("Failed to deobfuscate data with HTTP3Obfuscator: %v", err)
	}

	fmt.Printf("HTTP3 deobfuscated data: %s\n", http3Deobfuscated)
	if !bytes.Equal(http3Deobfuscated, testData) {
		t.Errorf("HTTP3 obfuscator round-trip failed: got %s, expected %s", 
			string(http3Deobfuscated), string(testData))
	}

	// Test obfuscator wrapper
	fmt.Println("\n=== Testing Obfuscator wrapper ===")
	obfsObfuscated, err := obfuscator.Obfuscate(testData)
	if err != nil {
		t.Fatalf("Failed to obfuscate data with Obfuscator: %v", err)
	}

	fmt.Printf("Obfuscator obfuscated length: %d\n", len(obfsObfuscated))
	fmt.Printf("Obfuscator obfuscated data: %v\n", obfsObfuscated)

	obfsDeobfuscated, err := obfuscator.Deobfuscate(obfsObfuscated)
	if err != nil {
		t.Fatalf("Failed to deobfuscate data with Obfuscator: %v", err)
	}

	fmt.Printf("Obfuscator deobfuscated data: %s\n", obfsDeobfuscated)
	if !bytes.Equal(obfsDeobfuscated, testData) {
		t.Errorf("Obfuscator round-trip failed: got %s, expected %s", 
			string(obfsDeobfuscated), string(testData))
	}

	// Compare the obfuscated data
	if !bytes.Equal(http3Obfuscated, obfsObfuscated) {
		fmt.Println("Obfuscated data from HTTP3Obfuscator and Obfuscator are different")
	} else {
		fmt.Println("Obfuscated data from HTTP3Obfuscator and Obfuscator are the same")
	}
}
//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
	"vantun/internal/core"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	
	// Test the varint encoding/decoding directly
	testValue := uint64(1024) // This is our chunk size
	fmt.Printf("Testing varint encoding/decoding for value: %d\n", testValue)
	
	var buf bytes.Buffer
	err := core.WriteVarInt(&buf, testValue)
	if err != nil {
		fmt.Printf("Error writing varint: %v\n", err)
		return
	}
	
	fmt.Printf("Encoded bytes: %x\n", buf.Bytes())
	
	reader := bytes.NewReader(buf.Bytes())
	decodedValue, err := core.ReadVarInt(reader)
	if err != nil {
		fmt.Printf("Error reading varint: %v\n", err)
		return
	}
	
	fmt.Printf("Decoded value: %d\n", decodedValue)
	fmt.Printf("Match: %t\n", testValue == decodedValue)
}

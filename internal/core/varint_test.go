package core

import (
	"bytes"
	"fmt"
	"testing"
)

func TestVarIntEncoding(t *testing.T) {
	// Test values
	testValues := []uint64{0, 1, 14, 63, 64, 16383, 16384, 1073741823, 1073741824}

	for _, value := range testValues {
		t.Run(fmt.Sprintf("Value%d", value), func(t *testing.T) {
			var buf bytes.Buffer
			
			// Encode
			err := writeVarInt(&buf, value)
			if err != nil {
				t.Fatalf("Failed to encode value %d: %v", value, err)
			}
			
			encoded := buf.Bytes()
			fmt.Printf("Value %d encoded as %v (length %d)\n", value, encoded, len(encoded))
			
			// Decode
			reader := bytes.NewReader(encoded)
			decoded, err := readVarInt(reader)
			if err != nil {
				t.Fatalf("Failed to decode value %d: %v", value, err)
			}
			
			if decoded != value {
				t.Errorf("Decoded value %d does not match original %d", decoded, value)
			}
			
			// Check that all bytes were consumed
			if reader.Len() != 0 {
				t.Errorf("Not all bytes were consumed during decoding")
			}
		})
	}
}
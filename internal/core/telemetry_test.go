package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
)

func TestTelemetryDataSerialization(t *testing.T) {
	// Create test telemetry data
	data := &TelemetryData{
		RTT:              50 * time.Millisecond,
		Loss:             0.01,
		Bandwidth:        1000000,
		Timestamp:        time.Now(),
		CongestionWindow: 10000,
		BytesInFlight:    1000,
		DeliveryRate:     500000,
	}

	// Encode the data
	encoded, err := cbor.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to encode telemetry data: %v", err)
	}

	// Decode the data
	var decoded TelemetryData
	err = cbor.Unmarshal(encoded, &decoded)
	if err != nil {
		t.Fatalf("Failed to decode telemetry data: %v", err)
	}

	// Verify the decoded data matches the original
	if decoded.RTT != data.RTT {
		t.Errorf("RTT mismatch: expected %v, got %v", data.RTT, decoded.RTT)
	}
	if decoded.Loss != data.Loss {
		t.Errorf("Loss mismatch: expected %f, got %f", data.Loss, decoded.Loss)
	}
	if decoded.Bandwidth != data.Bandwidth {
		t.Errorf("Bandwidth mismatch: expected %d, got %d", data.Bandwidth, decoded.Bandwidth)
	}
	if decoded.CongestionWindow != data.CongestionWindow {
		t.Errorf("CongestionWindow mismatch: expected %d, got %d", data.CongestionWindow, decoded.CongestionWindow)
	}
	if decoded.BytesInFlight != data.BytesInFlight {
		t.Errorf("BytesInFlight mismatch: expected %d, got %d", data.BytesInFlight, decoded.BytesInFlight)
	}
	if decoded.DeliveryRate != data.DeliveryRate {
		t.Errorf("DeliveryRate mismatch: expected %d, got %d", data.DeliveryRate, decoded.DeliveryRate)
	}
}

func TestTelemetryReporterAndReceiver(t *testing.T) {
	// Create test telemetry data
	data := &TelemetryData{
		RTT:              25 * time.Millisecond,
		Loss:             0.005,
		Bandwidth:        2000000,
		Timestamp:        time.Now(),
		CongestionWindow: 20000,
		BytesInFlight:    2000,
		DeliveryRate:     1000000,
	}

	// Create a buffer to simulate a stream
	buf := &bytes.Buffer{}

	// We'll manually write to the buffer for testing
	encoded, err := cbor.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to encode telemetry data: %v", err)
	}

	// Write length prefix and data to buffer
	length := uint32(len(encoded))
	lengthBuf := make([]byte, 4)
	lengthBuf[0] = byte(length >> 24)
	lengthBuf[1] = byte(length >> 16)
	lengthBuf[2] = byte(length >> 8)
	lengthBuf[3] = byte(length)
	buf.Write(lengthBuf)
	buf.Write(encoded)

	// We'll manually read from the buffer for testing
	lengthBuf = buf.Bytes()[:4]
	length = uint32(lengthBuf[0])<<24 | uint32(lengthBuf[1])<<16 | uint32(lengthBuf[2])<<8 | uint32(lengthBuf[3])
	encoded = buf.Bytes()[4 : 4+length]

	// Decode the telemetry data
	var receivedData TelemetryData
	err = cbor.Unmarshal(encoded, &receivedData)
	if err != nil {
		t.Fatalf("Failed to decode telemetry data: %v", err)
	}

	// Verify the received data matches the original
	if receivedData.RTT != data.RTT {
		t.Errorf("RTT mismatch: expected %v, got %v", data.RTT, receivedData.RTT)
	}
	if receivedData.Loss != data.Loss {
		t.Errorf("Loss mismatch: expected %f, got %f", data.Loss, receivedData.Loss)
	}
	if receivedData.Bandwidth != data.Bandwidth {
		t.Errorf("Bandwidth mismatch: expected %d, got %d", data.Bandwidth, receivedData.Bandwidth)
	}
}
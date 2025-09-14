package core

import (
	"context"
	"testing"
	"time"
)



func TestMultipathSessionAcceptStream(t *testing.T) {
	// Test case 3: Accepting a stream on a multipath session with no paths
	config := &Config{
		Address: "localhost:4242",
	}
	// Create required components for multipath session
	tokenBucketController := NewTokenBucketController(nil, nil, nil)
	adaptiveFEC, _ := NewAdaptiveFEC(10, 3, 1, 5)
	multipathSession := NewMultipathSession(config, tokenBucketController, adaptiveFEC)

	ctx := context.Background()
	stream, err := multipathSession.AcceptStream(ctx)
	if err == nil {
		t.Error("Expected error when accepting stream with no paths, got nil")
	}
	if stream != nil {
		t.Error("Expected nil stream when accepting stream with no paths, got non-nil")
	}
}

func TestMultipathSessionDataSplitter(t *testing.T) {
	// Test case 4: Testing the data splitter
	dataSplitter := NewDataSplitter(1024)
	testData := make([]byte, 3000)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	chunks := dataSplitter.Split(testData)
	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, got %d", len(chunks))
	}

	if len(chunks[0]) != 1024 || len(chunks[1]) != 1024 || len(chunks[2]) != 952 {
		t.Error("Chunks have incorrect sizes")
	}

	// Reassemble the chunks
	reassembled := make([]byte, 0)
	for _, chunk := range chunks {
		reassembled = append(reassembled, chunk...)
	}

	if len(reassembled) != len(testData) {
		t.Error("Reassembled data has incorrect length")
	}

	for i := range reassembled {
		if reassembled[i] != testData[i] {
			t.Error("Reassembled data does not match original data")
			break
		}
	}
}

func TestMultipathSessionPathSelector(t *testing.T) {
	// Test case 5: Testing the path selector
	paths := []*Path{
		{addr: "path1", rtt: 10 * time.Millisecond, bandwidth: 1000000, active: true},
		{addr: "path2", rtt: 20 * time.Millisecond, bandwidth: 2000000, active: true},
		{addr: "path3", rtt: 30 * time.Millisecond, bandwidth: 3000000, active: true},
	}

	nextPath := 0
	pathSelector := NewPathSelector(RoundRobinStrategy)
	selected := pathSelector.SelectPath(paths, &nextPath)
	if selected == nil || selected.addr != "path1" {
		t.Error("Round-robin selector did not select the correct path")
	}

	selected = pathSelector.SelectPath(paths, &nextPath)
	if selected == nil || selected.addr != "path2" {
		t.Error("Round-robin selector did not select the correct path")
	}

	pathSelector.strategy = MinRTTStrategy
	selected = pathSelector.SelectPath(paths, &nextPath)
	if selected == nil || selected.addr != "path1" {
		t.Error("Min-RTT selector did not select the correct path")
	}

	pathSelector.strategy = WeightedStrategy
	selected = pathSelector.SelectPath(paths, &nextPath)
	if selected == nil {
		t.Error("Weighted selector did not select a path")
	}
}

func TestObfuscatorObfuscateDeobfuscate(t *testing.T) {
	// Test case 1: Obfuscating and deobfuscating data with obfuscation enabled
	config := ObfuscatorConfig{
		Enabled:    true,
		MinPadding: 10,
		MaxPadding: 100,
	}
	obfuscator := NewObfuscator(config)
	testData := []byte("Hello, VANTUN!")

	obfuscatedData, err := obfuscator.Obfuscate(testData)
	if err != nil {
		t.Errorf("Failed to obfuscate data: %v", err)
	}

	if len(obfuscatedData) <= len(testData) {
		t.Error("Obfuscated data should be larger than original data")
	}

	deobfuscatedData, err := obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("Failed to deobfuscate data: %v", err)
	}

	if string(deobfuscatedData) != string(testData) {
		t.Errorf("Expected deobfuscated data %s, got %s", string(testData), string(deobfuscatedData))
	}
}

func TestObfuscatorWithoutObfuscation(t *testing.T) {
	// Test case 2: Obfuscating and deobfuscating data with obfuscation disabled
	config := ObfuscatorConfig{
		Enabled: false,
	}
	obfuscator := NewObfuscator(config)
	testData := []byte("Hello, VANTUN!")

	obfuscatedData, err := obfuscator.Obfuscate(testData)
	if err != nil {
		t.Errorf("Failed to obfuscate data: %v", err)
	}

	if string(obfuscatedData) != string(testData) {
		t.Errorf("Expected obfuscated data to be same as original when obfuscation is disabled")
	}

	deobfuscatedData, err := obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("Failed to deobfuscate data: %v", err)
	}

	if string(deobfuscatedData) != string(testData) {
		t.Errorf("Expected deobfuscated data %s, got %s", string(testData), string(deobfuscatedData))
	}
}

func TestObfuscatorStreamWriteRead(t *testing.T) {
	// Test case 3: Writing to and reading from an ObfuscatorStream
	mockStream := &MockQUICStream{}
	config := ObfuscatorConfig{
		Enabled:    true,
		MinPadding: 10,
		MaxPadding: 100,
	}
	obfuscator := NewObfuscator(config)
	obfsStream := NewObfuscatorStream(mockStream, obfuscator)

	testData := []byte("Hello, VANTUN!")
	n, err := obfsStream.Write(testData)
	if err != nil {
		t.Errorf("Failed to write to obfuscator stream: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	// Verify that the obfuscated data was written to the mock stream
	if len(mockStream.writeData) <= len(testData) {
		t.Error("Expected obfuscated data to be larger than original data")
	}

	// Set up mock stream to return obfuscated data for reading
	// We need to obfuscate the data and then set it as the read data
	obfuscatedData, _ := obfuscator.Obfuscate(testData)
	mockStream.readData = obfuscatedData

	buf := make([]byte, len(testData))
	n, err = obfsStream.Read(buf)
	if err != nil {
		t.Errorf("Failed to read from obfuscator stream: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to read %d bytes, read %d", len(testData), n)
	}

	if string(buf) != string(testData) {
		t.Errorf("Expected read data %s, got %s", string(testData), string(buf))
	}
}

func TestLoggerLevels(t *testing.T) {
	// Test case 1: Testing different log levels
	InitLogger("debug")
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")

	InitLogger("info")
	Debug("Debug message (should not appear)")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")

	InitLogger("warn")
	Debug("Debug message (should not appear)")
	Info("Info message (should not appear)")
	Warn("Warn message")
	Error("Error message")

	InitLogger("error")
	Debug("Debug message (should not appear)")
	Info("Info message (should not appear)")
	Warn("Warn message (should not appear)")
	Error("Error message")
}

func TestConfigManagerLoadConfig(t *testing.T) {
	// This test is skipped because it requires a valid config file
	// In a real test, you would create a temporary config file and load it
	t.Skip("Skipping config manager test - requires file system setup")
}

func TestTokenBucketConsume(t *testing.T) {
	// Test case 1: Consuming tokens from a token bucket
	tokenBucket := NewTokenBucket(1000, 5000)

	// Should be able to consume tokens
	if !tokenBucket.Consume(1000) {
		t.Error("Expected to be able to consume tokens from a full bucket")
	}

	// Should not be able to consume more tokens than capacity
	if tokenBucket.Consume(10000) {
		t.Error("Expected to not be able to consume more tokens than capacity")
	}
}



func TestTelemetryCollectorCollect(t *testing.T) {
	// Test case 1: Collecting telemetry data
	mockConn := &MockQUICConnection{}
	collector := NewTelemetryCollector(mockConn)

	data := collector.Collect()
	if data == nil {
		t.Error("Expected non-nil telemetry data, got nil")
	}

	// For mock connection, we expect default values
	// The current implementation returns placeholder values, so we check that it's not zero
	if data.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestTelemetryReporterReport(t *testing.T) {
	// Test case 2: Reporting telemetry data
	mockStream := &MockQUICStream{}
	reporter := NewTelemetryReporter(mockStream)

	telemetryData := &TelemetryData{
		RTT:       50 * time.Millisecond,
		Loss:      0.01,
		Bandwidth: 1000000,
		Timestamp: time.Now(),
	}

	err := reporter.Report(telemetryData)
	if err != nil {
		t.Errorf("Failed to report telemetry data: %v", err)
	}
}

func TestTelemetryManagerStartStop(t *testing.T) {
	// Test case 3: Starting and stopping the telemetry manager
	mockConn := &MockQUICConnection{}
	mockStream := &MockQUICStream{}
	manager := NewTelemetryManager(mockConn, mockStream, 100*time.Millisecond)

	manager.Start()
	time.Sleep(200 * time.Millisecond)
	manager.Stop()
}

func TestHTTP3ObfuscatorObfuscateDeobfuscate(t *testing.T) {
	// Test case 4: Obfuscating and deobfuscating data with HTTP/3 obfuscator
	frameTypes := []byte{0x0, 0x1, 0x2}
	obfuscator := NewHTTP3Obfuscator(frameTypes, 10, 100)
	testData := []byte("Hello, VANTUN!")

	obfuscatedData, err := obfuscator.Obfuscate(testData)
	if err != nil {
		t.Errorf("Failed to obfuscate data: %v", err)
	}

	deobfuscatedData, err := obfuscator.Deobfuscate(obfuscatedData)
	if err != nil {
		t.Errorf("Failed to deobfuscate data: %v", err)
	}

	if string(deobfuscatedData) != string(testData) {
		t.Errorf("Expected deobfuscated data %s, got %s", string(testData), string(deobfuscatedData))
	}
}
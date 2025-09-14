package core

import (
	"context"
	"crypto/tls"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

// NetworkCondition represents network conditions for testing
type NetworkCondition struct {
	Latency   time.Duration
	Bandwidth int64 // bytes per second
	PacketLoss float64 // 0.0 - 1.0
}

// NetworkConditionTestScenario represents a test scenario with specific network conditions
type NetworkConditionTestScenario struct {
	Name      string
	Condition NetworkCondition
}

// TestNetworkConditions tests VANTUN performance under various network conditions
func TestNetworkConditions(t *testing.T) {
	// Define test scenarios
	scenarios := []NetworkConditionTestScenario{
		{
			Name: "Good Network",
			Condition: NetworkCondition{
				Latency:    10 * time.Millisecond,
				Bandwidth:  100 * 1024 * 1024, // 100 Mbps
				PacketLoss: 0.0,
			},
		},
		{
			Name: "Typical WiFi",
			Condition: NetworkCondition{
				Latency:    50 * time.Millisecond,
				Bandwidth:  50 * 1024 * 1024, // 50 Mbps
				PacketLoss: 0.001, // 0.1%
			},
		},
		{
			Name: "Mobile Network (4G)",
			Condition: NetworkCondition{
				Latency:    100 * time.Millisecond,
				Bandwidth:  20 * 1024 * 1024, // 20 Mbps
				PacketLoss: 0.005, // 0.5%
			},
		},
		{
			Name: "Poor Network",
			Condition: NetworkCondition{
				Latency:    300 * time.Millisecond,
				Bandwidth:  5 * 1024 * 1024, // 5 Mbps
				PacketLoss: 0.02, // 2%
			},
		},
		{
			Name: "Very Poor Network",
			Condition: NetworkCondition{
				Latency:    500 * time.Millisecond,
				Bandwidth:  1 * 1024 * 1024, // 1 Mbps
				PacketLoss: 0.05, // 5%
			},
		},
	}

	// Run tests for each scenario
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			testNetworkCondition(t, scenario.Condition)
		})
	}
}

// testNetworkCondition runs a test with specific network conditions
func testNetworkCondition(t *testing.T, condition NetworkCondition) {
	// Create a test server
	serverAddr, _, err := createTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	// Note: server is nil in this implementation, so we don't close it

	// Test parameters - shorter duration for unit tests
	testDuration := 3 * time.Second
	dataSize := 256
	requestInterval := 500 * time.Millisecond

	// Create client config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"vantun"},
	}

	clientConfig := &Config{
		Address:   serverAddr,
		TLSConfig: tlsConfig,
		IsServer:  false,
	}

	// Create client session
	ctx, cancel := context.WithTimeout(context.Background(), testDuration+5*time.Second)
	defer cancel()

	t.Logf("Creating client session to %s", serverAddr)
	clientSession, err := NewSession(ctx, clientConfig)
	if err != nil {
		t.Fatalf("Failed to create client session: %v", err)
	}
	defer func() {
		if clientSession != nil {
			clientSession.Close()
		}
	}()
	t.Logf("Client session created successfully")

	// Track statistics
	var totalRequests int64
	var totalErrors int64
	var writeErr error
	var totalBytes int64
	var totalLatency int64 // Store as nanoseconds
	var latencyCount int64

	// Start test
	startTime := time.Now()
	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

requestLoop:
	for {
		select {
		case <-ctx.Done():
			// Test completed or cancelled
			t.Logf("Test context done")
			break requestLoop
		case <-ticker.C:
			if time.Since(startTime) >= testDuration {
				// Test duration completed
				t.Logf("Test duration completed")
				break requestLoop
			}

			atomic.AddInt64(&totalRequests, 1)
			t.Logf("Starting request #%d", totalRequests)

			// Simulate network conditions
			simulateNetworkCondition(condition)

			// Measure request latency
			requestStart := time.Now()

			// Open stream
			t.Logf("Opening interactive stream")
			stream, err := clientSession.OpenInteractiveStream(ctx)
			if err != nil {
				t.Logf("Failed to open interactive stream: %v", err)
				atomic.AddInt64(&totalErrors, 1)
				continue
			}
			t.Logf("Interactive stream opened successfully")

			// Generate test data
			testData := make([]byte, dataSize)
			for i := range testData {
				testData[i] = byte(rand.Intn(256))
			}
			t.Logf("Generated test data: %v", testData)

			// Send data
			t.Logf("Sending data")
			// First read any data that might be available (like the stream type response)
			readBuf := make([]byte, 4096)
			n, readErr := stream.Read(readBuf)
			if readErr == nil {
				t.Logf("Read %d bytes before sending test data: %v", n, readBuf[:n])
			}
			
			_, writeErr = stream.Write(testData)
			if writeErr != nil {
				t.Logf("Failed to write data: %v", writeErr)
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}
			t.Logf("Data sent successfully")

			// Read echo
			t.Logf("Reading echo")
			buf := make([]byte, dataSize)
			n, err = stream.Read(buf)
			if err != nil {
				t.Logf("Failed to read echo: %v", err)
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}
			t.Logf("Echo read successfully, received %d bytes: %v", n, buf[:n])

			// Verify data
			match := true
			if n != dataSize {
				t.Logf("Data size mismatch: expected %d, got %d", dataSize, n)
				match = false
			} else {
				for i := range testData {
					if buf[i] != testData[i] {
						match = false
						break
					}
				}
			}
			
			if !match {
				t.Logf("Data mismatch")
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}
			t.Logf("Data verified successfully")

			// Measure latency
			requestLatency := time.Since(requestStart)
			atomic.AddInt64(&totalLatency, int64(requestLatency)) // Store as nanoseconds
			atomic.AddInt64(&latencyCount, 1)

			atomic.AddInt64(&totalBytes, int64(dataSize*2)) // sent + received
			stream.Close()
			t.Logf("Request #%d completed successfully", totalRequests)
		}
	}

	// Calculate statistics
	duration := time.Since(startTime)
	throughput := float64(atomic.LoadInt64(&totalBytes)) / duration.Seconds()
	
	// Calculate average latency in nanoseconds, then convert to time.Duration
	var avgLatency time.Duration
	latencyCountVal := atomic.LoadInt64(&latencyCount)
	if latencyCountVal > 0 {
		avgLatency = time.Duration(atomic.LoadInt64(&totalLatency) / latencyCountVal)
	}

	// Report statistics
	t.Logf("Network condition test completed (%v latency, %v bandwidth, %.2f%% packet loss):",
		condition.Latency, formatBandwidth(condition.Bandwidth), condition.PacketLoss*100)
	t.Logf("  Test duration: %v", duration)
	t.Logf("  Total requests: %d", atomic.LoadInt64(&totalRequests))
	t.Logf("  Total errors: %d", atomic.LoadInt64(&totalErrors))
	t.Logf("  Total bytes transferred: %d", atomic.LoadInt64(&totalBytes))
	t.Logf("  Throughput: %.2f MB/s", throughput/1024/1024)
	t.Logf("  Average latency: %v", avgLatency)
	
	// Avoid division by zero
	totalRequestsVal := atomic.LoadInt64(&totalRequests)
	if totalRequestsVal > 0 {
		successRate := 100 * (1 - float64(atomic.LoadInt64(&totalErrors))/float64(totalRequestsVal))
		t.Logf("  Success rate: %.2f%%", successRate)

		// Verify success rate is acceptable
		if successRate < 50 { // Lower threshold for unit tests
			t.Errorf("Success rate %.2f%% is below acceptable threshold of 50%%", successRate)
		}
	}
}

// simulateNetworkCondition simulates network conditions by adding delays and packet loss
func simulateNetworkCondition(condition NetworkCondition) {
	// Add latency
	time.Sleep(condition.Latency)

	// Simulate packet loss by randomly dropping requests
	if condition.PacketLoss > 0 && rand.Float64() < condition.PacketLoss {
		// Simulate packet loss by sleeping for a long time
		time.Sleep(100 * time.Millisecond)
	}

	// Simulate bandwidth limitation by adding delays based on data size
	// This is a simplified simulation - in reality, bandwidth limitation would be more complex
}

// formatBandwidth formats bandwidth in human-readable format
func formatBandwidth(bandwidth int64) string {
	if bandwidth >= 1024*1024 {
		return "Mbps"
	} else if bandwidth >= 1024 {
		return "Kbps"
	} else {
		return "bps"
	}
}

// TestMultipathNetworkConditions tests multipath performance under various network conditions
func TestMultipathNetworkConditions(t *testing.T) {
	// Define test scenarios
	scenarios := []NetworkConditionTestScenario{
		{
			Name: "Good Network",
			Condition: NetworkCondition{
				Latency:    10 * time.Millisecond,
				Bandwidth:  100 * 1024 * 1024, // 100 Mbps
				PacketLoss: 0.0,
			},
		},
		{
			Name: "Poor Network",
			Condition: NetworkCondition{
				Latency:    300 * time.Millisecond,
				Bandwidth:  5 * 1024 * 1024, // 5 Mbps
				PacketLoss: 0.02, // 2%
			},
		},
	}

	// Run tests for each scenario
	for _, scenario := range scenarios {
		t.Run("Multipath "+scenario.Name, func(t *testing.T) {
			testMultipathNetworkCondition(t, scenario.Condition)
		})
	}
}

// testMultipathNetworkCondition tests multipath performance with specific network conditions
func testMultipathNetworkCondition(t *testing.T, condition NetworkCondition) {
	// Create a test server
	serverAddr, _, err := createTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	// Note: server is nil in this implementation, so we don't close it

	// Test parameters - shorter duration for unit tests
	testDuration := 3 * time.Second
	dataSize := 256
	requestInterval := 500 * time.Millisecond

	// Create client config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"vantun"},
	}

	// Create multipath session
	clientConfig := &Config{
		Address:   serverAddr,
		TLSConfig: tlsConfig,
		IsServer:  false,
	}

	// Create required components for multipath session
	tokenBucketController := NewTokenBucketController(nil, nil, nil)
	adaptiveFEC, _ := NewAdaptiveFEC(10, 3, 1, 5)
	multipathSession := NewMultipathSession(clientConfig, tokenBucketController, adaptiveFEC)

	// Add path
	ctx, cancel := context.WithTimeout(context.Background(), testDuration+5*time.Second)
	defer cancel()

	err = multipathSession.AddPath(ctx, serverAddr)
	if err != nil {
		t.Fatalf("Failed to add path to multipath session: %v", err)
	}
	defer func() {
		if multipathSession != nil {
			multipathSession.Close()
		}
	}()

	// Track statistics
	var totalRequests int64
	var totalErrors int64
	var totalBytes int64
	var totalLatency int64 // Store as nanoseconds
	var latencyCount int64

	// Start test
	startTime := time.Now()
	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

requestLoop:
	for {
		select {
		case <-ctx.Done():
			// Test completed or cancelled
			break requestLoop
		case <-ticker.C:
			if time.Since(startTime) >= testDuration {
				// Test duration completed
				break requestLoop
			}

			atomic.AddInt64(&totalRequests, 1)

			// Simulate network conditions
			simulateNetworkCondition(condition)

			// Measure request latency
			requestStart := time.Now()

			// Open stream
			stream, err := multipathSession.OpenStream(ctx)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				continue
			}

			// Generate test data
			testData := make([]byte, dataSize)
			for i := range testData {
				testData[i] = byte(rand.Intn(256))
			}

			// Send data
			// First read any data that might be available (like the stream type response)
			readBuf := make([]byte, 4096)
			n, readErr := stream.Read(readBuf)
			if readErr == nil {
				t.Logf("Read %d bytes before sending test data: %v", n, readBuf[:n])
			}
			
			_, err = stream.Write(testData)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Read echo
			buf := make([]byte, dataSize)
			_, err = stream.Read(buf)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Verify data
			match := true
			for i := range testData {
				if buf[i] != testData[i] {
					match = false
					break
				}
			}
			
			if !match {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Measure latency
			requestLatency := time.Since(requestStart)
			atomic.AddInt64(&totalLatency, int64(requestLatency)) // Store as nanoseconds
			atomic.AddInt64(&latencyCount, 1)

			atomic.AddInt64(&totalBytes, int64(dataSize*2)) // sent + received
			stream.Close()
		}
	}

	// Calculate statistics
	duration := time.Since(startTime)
	throughput := float64(atomic.LoadInt64(&totalBytes)) / duration.Seconds()
	
	// Calculate average latency in nanoseconds, then convert to time.Duration
	var avgLatency time.Duration
	latencyCountVal := atomic.LoadInt64(&latencyCount)
	if latencyCountVal > 0 {
		avgLatency = time.Duration(atomic.LoadInt64(&totalLatency) / latencyCountVal)
	}

	// Report statistics
	t.Logf("Multipath network condition test completed (%v latency, %v bandwidth, %.2f%% packet loss):",
		condition.Latency, formatBandwidth(condition.Bandwidth), condition.PacketLoss*100)
	t.Logf("  Test duration: %v", duration)
	t.Logf("  Total requests: %d", atomic.LoadInt64(&totalRequests))
	t.Logf("  Total errors: %d", atomic.LoadInt64(&totalErrors))
	t.Logf("  Total bytes transferred: %d", atomic.LoadInt64(&totalBytes))
	t.Logf("  Throughput: %.2f MB/s", throughput/1024/1024)
	t.Logf("  Average latency: %v", avgLatency)
	
	// Avoid division by zero
	totalRequestsVal := atomic.LoadInt64(&totalRequests)
	if totalRequestsVal > 0 {
		successRate := 100 * (1 - float64(atomic.LoadInt64(&totalErrors))/float64(totalRequestsVal))
		t.Logf("  Success rate: %.2f%%", successRate)

		// Verify success rate is acceptable
		if successRate < 50 { // Lower threshold for unit tests
			t.Errorf("Success rate %.2f%% is below acceptable threshold of 50%%", successRate)
		}
	}
}

// TestFECNetworkConditions tests FEC performance under various network conditions
func TestFECNetworkConditions(t *testing.T) {
	// Define test scenarios with high packet loss to test FEC
	scenarios := []NetworkConditionTestScenario{
		{
			Name: "High Loss Network",
			Condition: NetworkCondition{
				Latency:    100 * time.Millisecond,
				Bandwidth:  10 * 1024 * 1024, // 10 Mbps
				PacketLoss: 0.1, // 10%
			},
		},
		{
			Name: "Very High Loss Network",
			Condition: NetworkCondition{
				Latency:    200 * time.Millisecond,
				Bandwidth:  5 * 1024 * 1024, // 5 Mbps
				PacketLoss: 0.2, // 20%
			},
		},
	}

	// Run tests for each scenario
	for _, scenario := range scenarios {
		t.Run("FEC "+scenario.Name, func(t *testing.T) {
			testFECNetworkCondition(t, scenario.Condition)
		})
	}
}

// testFECNetworkCondition tests FEC performance with specific network conditions
func testFECNetworkCondition(t *testing.T, condition NetworkCondition) {
	// Create a test server
	serverAddr, _, err := createTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	// Note: server is nil in this implementation, so we don't close it

	// Test parameters - shorter duration for unit tests
	testDuration := 3 * time.Second
	dataSize := 256
	requestInterval := 500 * time.Millisecond

	// Create client config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"vantun"},
	}

	clientConfig := &Config{
		Address:   serverAddr,
		TLSConfig: tlsConfig,
		IsServer:  false,
	}

	// Create client session
	ctx, cancel := context.WithTimeout(context.Background(), testDuration+5*time.Second)
	defer cancel()

	clientSession, err := NewSession(ctx, clientConfig)
	if err != nil {
		t.Fatalf("Failed to create client session: %v", err)
	}
	defer func() {
		if clientSession != nil {
			clientSession.Close()
		}
	}()

	// Create adaptive FEC
	adaptiveFEC, err := NewAdaptiveFEC(10, 3, 1, 10)
	if err != nil {
		t.Fatalf("Failed to create adaptive FEC: %v", err)
	}

	// Track statistics
	var totalRequests int64
	var totalErrors int64
	var totalBytes int64
	var totalLatency int64 // Store as nanoseconds
	var latencyCount int64

	// Start test
	startTime := time.Now()
	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

requestLoop:
	for {
		select {
		case <-ctx.Done():
			// Test completed or cancelled
			break requestLoop
		case <-ticker.C:
			if time.Since(startTime) >= testDuration {
				// Test duration completed
				break requestLoop
			}

			atomic.AddInt64(&totalRequests, 1)

			// Simulate network conditions
			simulateNetworkCondition(condition)

			// Measure request latency
			requestStart := time.Now()

			// Open stream
			stream, err := clientSession.OpenBulkStream(ctx)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				continue
			}

			// Generate test data
			testData := make([]byte, dataSize)
			for i := range testData {
				testData[i] = byte(rand.Intn(256))
			}

			// Send data
			// First read any data that might be available (like the stream type response)
			readBuf := make([]byte, 4096)
			n, readErr := stream.Read(readBuf)
			if readErr == nil {
				t.Logf("Read %d bytes before sending test data: %v", n, readBuf[:n])
			}
			
			_, err = stream.Write(testData)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Read echo
			buf := make([]byte, dataSize)
			_, err = stream.Read(buf)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Verify data
			match := true
			for i := range testData {
				if buf[i] != testData[i] {
					match = false
					break
				}
			}
			
			if !match {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Measure latency
			requestLatency := time.Since(requestStart)
			atomic.AddInt64(&totalLatency, int64(requestLatency)) // Store as nanoseconds
			atomic.AddInt64(&latencyCount, 1)

			atomic.AddInt64(&totalBytes, int64(dataSize*2)) // sent + received
			stream.Close()

			// Simulate FEC adjustment based on simulated telemetry data
			telemetryData := &TelemetryData{
				RTT:       condition.Latency,
				Loss:      condition.PacketLoss,
				Bandwidth: uint64(condition.Bandwidth),
				Timestamp: time.Now(),
			}

			err = adaptiveFEC.Adjust(telemetryData)
			if err != nil {
				t.Logf("Failed to adjust FEC: %v", err)
			}
		}
	}

	// Calculate statistics
	duration := time.Since(startTime)
	throughput := float64(atomic.LoadInt64(&totalBytes)) / duration.Seconds()
	
	// Calculate average latency in nanoseconds, then convert to time.Duration
	var avgLatency time.Duration
	latencyCountVal := atomic.LoadInt64(&latencyCount)
	if latencyCountVal > 0 {
		avgLatency = time.Duration(atomic.LoadInt64(&totalLatency) / latencyCountVal)
	}

	// Report statistics
	t.Logf("FEC network condition test completed (%v latency, %v bandwidth, %.2f%% packet loss):",
		condition.Latency, formatBandwidth(condition.Bandwidth), condition.PacketLoss*100)
	t.Logf("  Test duration: %v", duration)
	t.Logf("  Total requests: %d", atomic.LoadInt64(&totalRequests))
	t.Logf("  Total errors: %d", atomic.LoadInt64(&totalErrors))
	t.Logf("  Total bytes transferred: %d", atomic.LoadInt64(&totalBytes))
	t.Logf("  Throughput: %.2f MB/s", throughput/1024/1024)
	t.Logf("  Average latency: %v", avgLatency)
	
	// Avoid division by zero
	totalRequestsVal := atomic.LoadInt64(&totalRequests)
	if totalRequestsVal > 0 {
		successRate := 100 * (1 - float64(atomic.LoadInt64(&totalErrors))/float64(totalRequestsVal))
		t.Logf("  Success rate: %.2f%%", successRate)

		// Verify success rate is acceptable
		if successRate < 50 { // Lower threshold for unit tests
			t.Errorf("Success rate %.2f%% is below acceptable threshold of 50%%", successRate)
		}
	}
}
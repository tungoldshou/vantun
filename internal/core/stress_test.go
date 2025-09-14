package core

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestStressConcurrentConnections tests concurrent connections
func TestStressConcurrentConnections(t *testing.T) {
	// Create a test server
	serverAddr, server, err := createTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer server.Close()

	// Test parameters
	numClients := 10
	requestsPerClient := 100
	dataSize := 1024

	// Create client config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"vantun"},
	}

	// Track statistics
	var totalRequests int64
	var totalErrors int64
	var totalBytes int64

	// Create clients concurrently
	var wg sync.WaitGroup
	clientErrChan := make(chan error, numClients*requestsPerClient)

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			clientConfig := &Config{
				Address:   serverAddr,
				TLSConfig: tlsConfig,
				IsServer:  false,
			}

			// Create client session
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			clientSession, err := NewSession(ctx, clientConfig)
			if err != nil {
				clientErrChan <- err
				return
			}
			defer clientSession.Close()

			// Send multiple requests
			for j := 0; j < requestsPerClient; j++ {
				atomic.AddInt64(&totalRequests, 1)

				// Open stream
				stream, err := clientSession.OpenInteractiveStream(ctx)
				if err != nil {
					atomic.AddInt64(&totalErrors, 1)
					clientErrChan <- err
					continue
				}

				// Generate test data
				testData := make([]byte, dataSize)
				_, err = rand.Read(testData)
				if err != nil {
					atomic.AddInt64(&totalErrors, 1)
					stream.Close()
					clientErrChan <- err
					continue
				}

				// Send data
				_, err = stream.Write(testData)
				if err != nil {
					atomic.AddInt64(&totalErrors, 1)
					stream.Close()
					clientErrChan <- err
					continue
				}

				// Read echo
				buf := make([]byte, dataSize)
				_, err = stream.Read(buf)
				if err != nil {
					atomic.AddInt64(&totalErrors, 1)
					stream.Close()
					clientErrChan <- err
					continue
				}

				// Verify data
				for k := range testData {
					if buf[k] != testData[k] {
						atomic.AddInt64(&totalErrors, 1)
						stream.Close()
						clientErrChan <- err
						continue
					}
				}

				atomic.AddInt64(&totalBytes, int64(dataSize*2)) // sent + received
				stream.Close()
			}
		}(i)
	}

	// Wait for all clients to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All clients finished
	case <-time.After(60 * time.Second):
		t.Fatal("Test timed out")
	}

	close(clientErrChan)

	// Report statistics
	t.Logf("Stress test completed:")
	t.Logf("  Total requests: %d", atomic.LoadInt64(&totalRequests))
	t.Logf("  Total errors: %d", atomic.LoadInt64(&totalErrors))
	t.Logf("  Total bytes transferred: %d", atomic.LoadInt64(&totalBytes))
	t.Logf("  Success rate: %.2f%%", 100*(1-float64(atomic.LoadInt64(&totalErrors))/float64(atomic.LoadInt64(&totalRequests))))

	// Check for errors
	if atomic.LoadInt64(&totalErrors) > 0 {
		t.Errorf("Found %d errors during stress test", atomic.LoadInt64(&totalErrors))
	}

	// Verify success rate is acceptable
	successRate := 1 - float64(atomic.LoadInt64(&totalErrors))/float64(atomic.LoadInt64(&totalRequests))
	if successRate < 0.95 {
		t.Errorf("Success rate %.2f%% is below acceptable threshold of 95%%", successRate*100)
	}
}

// TestStressHighThroughput tests high throughput with bulk data
func TestStressHighThroughput(t *testing.T) {
	// Create a test server
	serverAddr, server, err := createTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer server.Close()

	// Test parameters
	dataSize := 10 * 1024 * 1024 // 10MB
	numStreams := 5

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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientSession, err := NewSession(ctx, clientConfig)
	if err != nil {
		t.Fatalf("Failed to create client session: %v", err)
	}
	defer clientSession.Close()

	// Track statistics
	var totalBytes int64
	var totalErrors int64
	startTime := time.Now()

	// Open multiple streams and transfer data concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, numStreams)

	for i := 0; i < numStreams; i++ {
		wg.Add(1)
		go func(streamID int) {
			defer wg.Done()

			// Open bulk stream
			stream, err := clientSession.OpenBulkStream(ctx)
			if err != nil {
				errChan <- err
				atomic.AddInt64(&totalErrors, 1)
				return
			}
			defer stream.Close()

			// Generate test data
			testData := make([]byte, dataSize)
			_, err = rand.Read(testData)
			if err != nil {
				errChan <- err
				atomic.AddInt64(&totalErrors, 1)
				return
			}

			// Write data
			n, err := stream.Write(testData)
			if err != nil {
				errChan <- err
				atomic.AddInt64(&totalErrors, 1)
				return
			}
			atomic.AddInt64(&totalBytes, int64(n))

			// Read echo
			buf := make([]byte, dataSize)
			totalRead := 0
			for totalRead < dataSize {
				n, err := stream.Read(buf[totalRead:])
				if err != nil {
					errChan <- err
					atomic.AddInt64(&totalErrors, 1)
					return
				}
				totalRead += n
			}
			atomic.AddInt64(&totalBytes, int64(totalRead))

			// Verify data
			for j := range testData {
				if buf[j] != testData[j] {
					errChan <- err
					atomic.AddInt64(&totalErrors, 1)
					return
				}
			}
		}(i)
	}

	// Wait for all streams to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All streams finished
	case <-time.After(120 * time.Second):
		t.Fatal("Test timed out")
	}

	close(errChan)

	// Calculate throughput
	duration := time.Since(startTime)
	throughput := float64(atomic.LoadInt64(&totalBytes)) / duration.Seconds()

	// Report statistics
	t.Logf("High throughput test completed:")
	t.Logf("  Data transferred: %d bytes", atomic.LoadInt64(&totalBytes))
	t.Logf("  Duration: %v", duration)
	t.Logf("  Throughput: %.2f MB/s", throughput/1024/1024)
	t.Logf("  Total errors: %d", atomic.LoadInt64(&totalErrors))

	// Check for errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Stream error: %v", err)
		}
	}

	// Verify no errors occurred
	if atomic.LoadInt64(&totalErrors) > 0 {
		t.Errorf("Found %d errors during high throughput test", atomic.LoadInt64(&totalErrors))
	}
}

// TestStressLongRunningConnection tests long-running connection stability
func TestStressLongRunningConnection(t *testing.T) {
	// Create a test server
	serverAddr, server, err := createTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer server.Close()

	// Test parameters
	testDuration := 30 * time.Second // Shortened for testing
	requestInterval := 100 * time.Millisecond
	dataSize := 256

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
	ctx, cancel := context.WithTimeout(context.Background(), testDuration+10*time.Second)
	defer cancel()

	clientSession, err := NewSession(ctx, clientConfig)
	if err != nil {
		t.Fatalf("Failed to create client session: %v", err)
	}
	defer clientSession.Close()

	// Track statistics
	var totalRequests int64
	var totalErrors int64
	var totalBytes int64

	// Start test
	startTime := time.Now()
	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Test completed or cancelled
			goto report
		case <-ticker.C:
			if time.Since(startTime) >= testDuration {
				// Test duration completed
				goto report
			}

			atomic.AddInt64(&totalRequests, 1)

			// Open stream
			stream, err := clientSession.OpenInteractiveStream(ctx)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				continue
			}

			// Generate test data
			testData := make([]byte, dataSize)
			_, err = rand.Read(testData)
			if err != nil {
				atomic.AddInt64(&totalErrors, 1)
				stream.Close()
				continue
			}

			// Send data
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
			for i := range testData {
				if buf[i] != testData[i] {
					atomic.AddInt64(&totalErrors, 1)
					stream.Close()
					continue
				}
			}

			atomic.AddInt64(&totalBytes, int64(dataSize*2)) // sent + received
			stream.Close()
		}
	}

report:
	// Report statistics
	t.Logf("Long-running connection test completed:")
	t.Logf("  Test duration: %v", time.Since(startTime))
	t.Logf("  Total requests: %d", atomic.LoadInt64(&totalRequests))
	t.Logf("  Total errors: %d", atomic.LoadInt64(&totalErrors))
	t.Logf("  Total bytes transferred: %d", atomic.LoadInt64(&totalBytes))
	t.Logf("  Requests per second: %.2f", float64(atomic.LoadInt64(&totalRequests))/time.Since(startTime).Seconds())
	t.Logf("  Success rate: %.2f%%", 100*(1-float64(atomic.LoadInt64(&totalErrors))/float64(atomic.LoadInt64(&totalRequests))))

	// Verify success rate is acceptable
	successRate := 1 - float64(atomic.LoadInt64(&totalErrors))/float64(atomic.LoadInt64(&totalRequests))
	if successRate < 0.99 {
		t.Errorf("Success rate %.2f%% is below acceptable threshold of 99%%", successRate*100)
	}
}
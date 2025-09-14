//go:build ignore

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"time"

	"vantun/internal/cli"
	"vantun/internal/core"
)

var (
	serverAddr = flag.String("server-addr", "localhost:4242", "VANTUN server address")
	testDuration = flag.Duration("duration", 24*time.Hour, "Test duration")
	dataSize = flag.Int("data-size", 1024, "Size of data to send in each request (bytes)")
	requestInterval = flag.Duration("interval", 1*time.Second, "Interval between requests")
)

func main() {
	flag.Parse()

	// Load client configuration
	config, err := cli.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Override server address
	config.Address = *serverAddr

	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"vantun"},
	}

	// Create core configuration
	coreConfig := &core.Config{
		Address:   config.Address,
		TLSConfig: tlsConfig,
		IsServer:  config.Server,
	}

	// Create token bucket
	tokenBucket := core.NewTokenBucket(config.TokenBucketRate, config.TokenBucketCapacity)
	
	// Create adaptive FEC
	adaptiveFEC, err := core.NewAdaptiveFEC(config.FECData, config.FECParity, 1, 10)
	if err != nil {
		log.Fatalf("Failed to create adaptive FEC: %v", err)
	}
	
	// Create session
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	session, err := core.NewSession(ctx, coreConfig)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create token bucket controller
	controller := core.NewTokenBucketController(tokenBucket, adaptiveFEC, session.conn)
	controller.Start()
	defer controller.Stop()

	log.Printf("Starting %v stability test", *testDuration)
	
	// Create data to send
	data := make([]byte, *dataSize)
	for i := range data {
		data[i] = byte(i % 256)
	}
	
	// Test start time
	startTime := time.Now()
	
	// Run test for specified duration
	ticker := time.NewTicker(*requestInterval)
	defer ticker.Stop()
	
	requestCount := 0
	errorCount := 0
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Test cancelled")
			return
		case <-ticker.C:
			// Check if test duration has elapsed
			if time.Since(startTime) >= *testDuration {
				log.Println("Test duration completed")
				break
			}
			
			// Open interactive stream
			stream, err := session.OpenInteractiveStream(ctx)
			if err != nil {
				log.Printf("Failed to open interactive stream: %v", err)
				errorCount++
				continue
			}
			
			// Send data
			if _, err := stream.Write(data); err != nil {
				log.Printf("Failed to write data: %v", err)
				stream.Close()
				errorCount++
				continue
			}
			
			// Read response
			buf := make([]byte, len(data))
			if _, err := stream.Read(buf); err != nil {
				log.Printf("Failed to read response: %v", err)
				stream.Close()
				errorCount++
				continue
			}
			
			// Close stream
			stream.Close()
			
			requestCount++
			
			// Log progress every 100 requests
			if requestCount%100 == 0 {
				log.Printf("Sent %d requests, %d errors", requestCount, errorCount)
			}
		}
	}
	
	log.Printf("Test completed. Sent %d requests, %d errors", requestCount, errorCount)
}

// generateTLSConfig generates a simple self-signed TLS config for testing.
func generateTLSConfig() *tls.Config {
	// For this test, we'll return nil to use the default config.
	// In practice, you would need to ensure the client and server have matching TLS configurations.
	// For now, we'll rely on the insecure skip verify option in the server's TLS config.
	return nil
}

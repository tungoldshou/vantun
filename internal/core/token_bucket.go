package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// TokenBucket controls the rate of data transmission.
type TokenBucket struct {
	// rate is the rate at which tokens are added to the bucket (tokens per second).
	rate float64
	// capacity is the maximum number of tokens the bucket can hold.
	capacity float64
	// tokens is the current number of tokens in the bucket.
	tokens float64
	// lastUpdate is the time when the token count was last updated.
	lastUpdate time.Time
	// mutex protects the token bucket fields.
	mutex sync.Mutex
}

// NewTokenBucket creates a new TokenBucket with the specified rate and capacity.
func NewTokenBucket(rate, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // Start with a full bucket
		lastUpdate: time.Now(),
	}
}

// updateTokens updates the token count based on the elapsed time.
func (tb *TokenBucket) updateTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tokensToAdd := elapsed * tb.rate
	
	tb.tokens += tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	
	tb.lastUpdate = now
}

// Consume attempts to consume the specified number of tokens.
// It returns true if there were enough tokens, false otherwise.
func (tb *TokenBucket) Consume(tokens float64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	
	tb.updateTokens()
	
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	
	return false
}

// GetRate returns the current rate of the token bucket.
func (tb *TokenBucket) GetRate() float64 {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	return tb.rate
}

// SetRate sets the rate of the token bucket.
func (tb *TokenBucket) SetRate(rate float64) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.rate = rate
}

// TokenBucketController controls the rate of data transmission based on telemetry data.
type TokenBucketController struct {
	// bucket is the token bucket being controlled.
	bucket *TokenBucket
	// collector is the telemetry collector.
	collector *TelemetryCollector
	// reporter is the telemetry reporter.
	reporter *TelemetryReporter
	// adaptiveFEC is the adaptive FEC controller.
	adaptiveFEC *AdaptiveFEC
	// ctx is the context for the controller.
	ctx context.Context
	// cancel is the cancel function for the controller.
	cancel context.CancelFunc
}

// NewTokenBucketController creates a new TokenBucketController.
func NewTokenBucketController(bucket *TokenBucket, adaptiveFEC *AdaptiveFEC, conn quic.Connection) *TokenBucketController {
	ctx, cancel := context.WithCancel(context.Background())
	return &TokenBucketController{
		bucket:      bucket,
		collector:   NewTelemetryCollector(conn),
		reporter:    NewTelemetryReporter(nil), // Stream will be set later
		adaptiveFEC: adaptiveFEC,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// SetTelemetryStream sets the telemetry stream for reporting.
func (tbc *TokenBucketController) SetTelemetryStream(stream quic.Stream) {
	tbc.reporter = NewTelemetryReporter(stream)
}

// UpdateConnection updates the connection used for telemetry collection.
func (tbc *TokenBucketController) UpdateConnection(conn quic.Connection) {
	tbc.collector = NewTelemetryCollector(conn)
}

// Start starts the token bucket controller.
// It periodically collects telemetry data and adjusts the token bucket rate.
func (tbc *TokenBucketController) Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Collect telemetry every second
		defer ticker.Stop()
		
		for {
			select {
			case <-tbc.ctx.Done():
				return
			case <-ticker.C:
				// Collect telemetry data
				data := tbc.collector.Collect()
				
				// Report telemetry data
				if err := tbc.reporter.Report(data); err != nil {
					fmt.Printf("Failed to report telemetry: %v\n", err)
				}
				
				// Adjust token bucket rate based on telemetry data
				// This is a simple example. In a real implementation, you would
				// use a more sophisticated algorithm.
				if data.Loss > 0.05 { // If loss is greater than 5%
					// Reduce rate by 10%
					newRate := tbc.bucket.GetRate() * 0.9
					tbc.bucket.SetRate(newRate)
					fmt.Printf("Reducing rate to %f due to high loss\n", newRate)
				} else if data.Loss < 0.01 { // If loss is less than 1%
					// Increase rate by 10%
					newRate := tbc.bucket.GetRate() * 1.1
					// Ensure rate doesn't exceed a reasonable maximum
					// In a real implementation, this would be based on bandwidth estimation
					if newRate < 10000000 { // 10 MB/s limit
						tbc.bucket.SetRate(newRate)
						fmt.Printf("Increasing rate to %f due to low loss\n", newRate)
					}
				}
				
				// Adjust FEC parameters based on telemetry data
				if tbc.adaptiveFEC != nil {
					if err := tbc.adaptiveFEC.Adjust(data); err != nil {
						fmt.Printf("Failed to adjust FEC: %v\n", err)
					}
				}
			}
		}
	}()
}

// Stop stops the token bucket controller.
func (tbc *TokenBucketController) Stop() {
	tbc.cancel()
}
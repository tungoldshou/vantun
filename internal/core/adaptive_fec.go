package core

import (
	"fmt"
	"time"
)

// AdaptiveFEC adjusts FEC parameters based on telemetry data.
type AdaptiveFEC struct {
	// fec is the underlying FEC encoder/decoder.
	fec *FEC
	// k is the number of data shards.
	k int
	// m is the number of parity shards.
	m int
	// minParity is the minimum number of parity shards.
	minParity int
	// maxParity is the maximum number of parity shards.
	maxParity int
}

// NewAdaptiveFEC creates a new AdaptiveFEC.
func NewAdaptiveFEC(k, m, minParity, maxParity int) (*AdaptiveFEC, error) {
	fec, err := NewFEC(k, m)
	if err != nil {
		return nil, fmt.Errorf("failed to create FEC: %w", err)
	}
	
	return &AdaptiveFEC{
		fec:       fec,
		k:         k,
		m:         m,
		minParity: minParity,
		maxParity: maxParity,
	}, nil
}

// Adjust adjusts the FEC parameters based on telemetry data using a more sophisticated algorithm.
func (af *AdaptiveFEC) Adjust(data *TelemetryData) error {
	// Calculate a new number of parity shards based on multiple factors:
	// 1. Packet loss rate
	// 2. RTT (higher RTT means more expensive retransmissions)
	// 3. Bandwidth (higher bandwidth can afford more redundancy)
	// 4. Delivery rate relative to bandwidth (network efficiency)
	
	// Base adjustment based on loss rate
	lossFactor := af.calculateLossFactor(data.Loss)
	
	// RTT factor (higher RTT increases the cost of retransmissions)
	rttFactor := af.calculateRTTFactor(data.RTT)
	
	// Bandwidth factor (higher bandwidth can afford more redundancy)
	bandwidthFactor := af.calculateBandwidthFactor(data.Bandwidth)
	
	// Efficiency factor (delivery rate relative to bandwidth)
	efficiencyFactor := af.calculateEfficiencyFactor(data.DeliveryRate, data.Bandwidth)
	
	// Combine all factors to determine the adjustment
	combinedFactor := lossFactor * rttFactor * bandwidthFactor * efficiencyFactor
	
	// Calculate new number of parity shards
	newM := int(float64(af.m) * combinedFactor)
	
	// Ensure newM is within bounds
	if newM < af.minParity {
		newM = af.minParity
	}
	if newM > af.maxParity {
		newM = af.maxParity
	}
	
	// Apply smoothing to avoid rapid changes
	newM = af.smoothAdjustment(newM)
	
	// If the number of parity shards has changed, create a new FEC
	if newM != af.m {
		fmt.Printf("Adjusting FEC: changing parity shards from %d to %d (lossFactor=%.2f, rttFactor=%.2f, bandwidthFactor=%.2f, efficiencyFactor=%.2f)\n", 
			af.m, newM, lossFactor, rttFactor, bandwidthFactor, efficiencyFactor)
		fec, err := NewFEC(af.k, newM)
		if err != nil {
			return fmt.Errorf("failed to create new FEC: %w", err)
		}
		af.fec = fec
		af.m = newM
	}
	
	return nil
}

// calculateLossFactor calculates a factor based on packet loss rate
func (af *AdaptiveFEC) calculateLossFactor(loss float64) float64 {
	// More conservative mapping for loss factor:
	// 0% loss -> 1.0 (no change)
	// 0.1% loss -> 0.65 (decrease)
	// 0.5% loss -> 0.75 (decrease)
	// 1% loss -> 1.0 (no change)
	// 2% loss -> 1.05 (slight increase)
	// 5% loss -> 1.15 (moderate increase)
	// 10% loss -> 1.3 (significant increase)
	// 20% loss -> 1.6 (large increase)
	
	if loss <= 0.0 {
		return 1.0
	}
	
	// Threshold below which we decrease redundancy
	if loss < 0.01 { // Less than 1% loss
		// Decrease redundancy for very low loss rates, more aggressively
		// Scale from 0.65 to 1.0 as loss goes from 0% to 1%
		return 0.65 + (0.35 * loss / 0.01)
	}
	
	// For loss rates between 1% and higher, increase redundancy more conservatively
	if loss <= 0.02 { // 1-2% loss
		return 1.0 + ((loss - 0.01) * 5.0) // Scale from 1.0 to 1.05
	}
	
	if loss <= 0.05 { // 2-5% loss
		return 1.05 + ((loss - 0.02) * 5.0) // Scale from 1.05 to 1.15
	}
	
	if loss <= 0.1 { // 5-10% loss
		return 1.15 + ((loss - 0.05) * 3.0) // Scale from 1.15 to 1.3
	}
	
	// For higher loss rates, increase more aggressively
	return 1.3 + ((loss - 0.1) * 3.0) // Scale from 1.3 upward
}

// calculateRTTFactor calculates a factor based on RTT
func (af *AdaptiveFEC) calculateRTTFactor(rtt time.Duration) float64 {
	// RTT factor increases with RTT, but more conservatively:
	// 10ms -> 1.0 (baseline)
	// 50ms -> 1.08
	// 100ms -> 1.12
	// 200ms -> 1.16
	
	rttMs := float64(rtt.Milliseconds())
	
	// Linear mapping for RTT factor, more conservative
	if rttMs <= 10.0 {
		return 1.0
	}
	
	return 1.0 + (rttMs-10.0)/300.0
}

// calculateBandwidthFactor calculates a factor based on available bandwidth
func (af *AdaptiveFEC) calculateBandwidthFactor(bandwidth uint64) float64 {
	// Bandwidth factor decreases with higher bandwidth (can afford more redundancy):
	// 1 Mbps -> 1.05
	// 10 Mbps -> 1.025
	// 100 Mbps -> 1.0
	// 1000 Mbps -> 0.975
	
	bandwidthMbps := float64(bandwidth) / (1024.0 * 1024.0)
	
	// Logarithmic mapping for bandwidth factor, more conservative
	if bandwidthMbps <= 0.0 {
		return 1.05
	}
	
	// Decreasing factor with increasing bandwidth, but more conservative
	factor := 1.05 - (0.075 * (bandwidthMbps / 100.0))
	if factor < 0.95 {
		factor = 0.95
	}
	
	return factor
}

// calculateEfficiencyFactor calculates a factor based on network efficiency
func (af *AdaptiveFEC) calculateEfficiencyFactor(deliveryRate, bandwidth uint64) float64 {
	// Efficiency factor increases when delivery rate is low relative to bandwidth:
	// 100% efficiency -> 1.0 (no change)
	// 50% efficiency -> 1.3 (increase redundancy)
	// 10% efficiency -> 2.0 (significant increase)
	
	if bandwidth <= 0 {
		return 1.0
	}
	
	efficiency := float64(deliveryRate) / float64(bandwidth)
	
	// Inverse relationship - lower efficiency means we need more redundancy
	if efficiency >= 1.0 {
		return 1.0
	}
	
	return 2.0 - efficiency
}

// smoothAdjustment applies smoothing to avoid rapid changes in FEC parameters
func (af *AdaptiveFEC) smoothAdjustment(newM int) int {
	// Limit the change to at most 2 shards per adjustment
	maxChange := 2
	
	if newM > af.m+maxChange {
		return af.m + maxChange
	}
	
	if newM < af.m-maxChange {
		return af.m - maxChange
	}
	
	return newM
}

// Encode encodes the data using the current FEC parameters.
func (af *AdaptiveFEC) Encode(data []byte) ([][]byte, error) {
	return af.fec.Encode(data)
}

// Decode decodes the shards using the current FEC parameters.
func (af *AdaptiveFEC) Decode(shards [][]byte) ([]byte, error) {
	return af.fec.Decode(shards)
}
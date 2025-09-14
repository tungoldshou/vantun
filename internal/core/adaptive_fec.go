package core

import (
	"fmt"
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

// Adjust adjusts the FEC parameters based on telemetry data.
func (af *AdaptiveFEC) Adjust(data *TelemetryData) error {
	// This is a simple example. In a real implementation, you would
	// use a more sophisticated algorithm to determine the optimal
	// number of parity shards.
	
	// Calculate a new number of parity shards based on loss rate
	// This is a very simplified algorithm for demonstration purposes.
	var newM int
	if data.Loss > 0.1 { // High loss (10%+)
		// Increase redundancy
		newM = af.m + 2
	} else if data.Loss > 0.05 { // Medium loss (5-10%)
		// Slightly increase redundancy
		newM = af.m + 1
	} else if data.Loss < 0.01 { // Low loss (<1%)
		// Decrease redundancy to improve throughput
		newM = af.m - 1
	} else {
		// Keep current redundancy for loss between 1-5%
		newM = af.m
	}
	
	// Ensure newM is within bounds
	if newM < af.minParity {
		newM = af.minParity
	}
	if newM > af.maxParity {
		newM = af.maxParity
	}
	
	// If the number of parity shards has changed, create a new FEC
	if newM != af.m {
		fmt.Printf("Adjusting FEC: changing parity shards from %d to %d\n", af.m, newM)
		fec, err := NewFEC(af.k, newM)
		if err != nil {
			return fmt.Errorf("failed to create new FEC: %w", err)
		}
		af.fec = fec
		af.m = newM
	}
	
	return nil
}

// Encode encodes the data using the current FEC parameters.
func (af *AdaptiveFEC) Encode(data []byte) ([][]byte, error) {
	return af.fec.Encode(data)
}

// Decode decodes the shards using the current FEC parameters.
func (af *AdaptiveFEC) Decode(shards [][]byte) ([]byte, error) {
	return af.fec.Decode(shards)
}
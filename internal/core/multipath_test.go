package core

import (
	"testing"
	"time"
)

func TestPathSelector(t *testing.T) {
	// Create test paths
	paths := []*Path{
		{
			addr:      "path1",
			rtt:       10 * time.Millisecond,
			loss:      0.01,
			bandwidth: 1000000, // 1 MB/s
			active:    true,
		},
		{
			addr:      "path2",
			rtt:       20 * time.Millisecond,
			loss:      0.02,
			bandwidth: 2000000, // 2 MB/s
			active:    true,
		},
		{
			addr:      "path3",
			rtt:       5 * time.Millisecond,
			loss:      0.05,
			bandwidth: 500000, // 0.5 MB/s
			active:    true,
		},
	}

	// Test RoundRobinStrategy
	selector := NewPathSelector(RoundRobinStrategy)
	nextPath := 0

	// Select paths multiple times
	selected := make(map[string]int)
	for i := 0; i < 6; i++ {
		path := selector.SelectPath(paths, &nextPath)
		if path == nil {
			t.Fatal("Expected a path to be selected")
		}
		selected[path.addr]++
	}

	// Each path should be selected twice with round-robin
	if selected["path1"] != 2 || selected["path2"] != 2 || selected["path3"] != 2 {
		t.Errorf("Round-robin selection not working correctly. Got: %v", selected)
	}

	// Test MinRTTStrategy
	selector = NewPathSelector(MinRTTStrategy)
	nextPath = 0

	// Select path multiple times - should always select path3 (lowest RTT)
	for i := 0; i < 3; i++ {
		path := selector.SelectPath(paths, &nextPath)
		if path == nil {
			t.Fatal("Expected a path to be selected")
		}
		if path.addr != "path3" {
			t.Errorf("Expected path3 (lowest RTT) to be selected, got %s", path.addr)
		}
	}

	// Test WeightedStrategy
	selector = NewPathSelector(WeightedStrategy)
	nextPath = 0

	// Select paths multiple times - higher bandwidth paths should be selected more often
	selected = make(map[string]int)
	for i := 0; i < 100; i++ {
		path := selector.SelectPath(paths, &nextPath)
		if path == nil {
			t.Fatal("Expected a path to be selected")
		}
		selected[path.addr]++
	}

	// Path2 should be selected most often (highest bandwidth)
	if selected["path2"] <= selected["path1"] || selected["path2"] <= selected["path3"] {
		t.Errorf("Weighted selection not working correctly. Got: %v", selected)
	}
}

func TestDataSplitter(t *testing.T) {
	// Create a data splitter with 10-byte chunks
	splitter := NewDataSplitter(10)

	// Test with data smaller than chunk size
	smallData := []byte("hello")
	chunks := splitter.Split(smallData)
	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk, got %d", len(chunks))
	}
	if string(chunks[0]) != "hello" {
		t.Errorf("Expected 'hello', got %s", string(chunks[0]))
	}

	// Test with data exactly matching chunk size
	exactData := []byte("0123456789")
	chunks = splitter.Split(exactData)
	if len(chunks) != 1 {
		t.Errorf("Expected 1 chunk, got %d", len(chunks))
	}
	if string(chunks[0]) != "0123456789" {
		t.Errorf("Expected '0123456789', got %s", string(chunks[0]))
	}

	// Test with data larger than chunk size
	largeData := []byte("0123456789abcdefghij")
	chunks = splitter.Split(largeData)
	if len(chunks) != 2 {
		t.Errorf("Expected 2 chunks, got %d", len(chunks))
	}
	if string(chunks[0]) != "0123456789" {
		t.Errorf("Expected '0123456789', got %s", string(chunks[0]))
	}
	if string(chunks[1]) != "abcdefghij" {
		t.Errorf("Expected 'abcdefghij', got %s", string(chunks[1]))
	}

	// Test with data much larger than chunk size
	veryLargeData := []byte("0123456789abcdefghij0123456789abcdefghij0123456789abcdefghij")
	chunks = splitter.Split(veryLargeData)
	if len(chunks) != 6 {
		t.Errorf("Expected 6 chunks, got %d", len(chunks))
	}
}
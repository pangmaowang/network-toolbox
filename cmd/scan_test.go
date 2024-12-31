package cmd

import (
	"net"
	"testing"
	"time"
)

func TestScanPorts(t *testing.T) {
	target := "127.0.0.1"

	// Using higher port numbers that are less likely to be in use
	testPort := 50080
	listener, err := net.Listen("tcp", ":50080")
	if err != nil {
		t.Fatalf("Error starting listener: %v", err)
	}
	defer listener.Close()

	go func() {
		_, _ = listener.Accept()
	}()

	time.Sleep(100 * time.Millisecond)

	scanPort := testPort - 1 // 50079
	endPort := testPort + 1  // 50081

	timeout := 500 * time.Millisecond

	results := scanPorts(target, scanPort, endPort, timeout)

	expected := []scanResult{
		{Port: 50079, Open: false},
		{Port: 50080, Open: true},
		{Port: 50081, Open: false},
	}

	// Convert results to a map for easier comparison
	resultMap := make(map[int]bool)
	for _, r := range results {
		resultMap[r.Port] = r.Open
	}

	// Check if we got the expected number of results
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	// Check each expected result
	for _, exp := range expected {
		open, exists := resultMap[exp.Port]
		if !exists {
			t.Errorf("Port %d was not scanned", exp.Port)
			continue
		}
		if open != exp.Open {
			t.Errorf("Port %d: expected open=%v, got open=%v", exp.Port, exp.Open, open)
		}
	}
}

func TestInvalidTarget(t *testing.T) {
	results := scanPorts("invalid-target", 1, 10, time.Second)
	if len(results) > 0 {
		t.Errorf("对于无效目标，应返回空结果，但获得了 %d 个结果", len(results))
	}
}

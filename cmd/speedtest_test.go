package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadSpeed(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send 1MB of test data
		data := make([]byte, 1024*1024)
		for i := range data {
			data[i] = 'A'
		}
		w.Write(data)
	}))
	defer ts.Close()

	// Test valid URL
	speed, err := testDownloadSpeed(ts.URL)
	if err != nil {
		t.Errorf("testDownloadSpeed failed: %v", err)
	}
	if speed <= 0 {
		t.Errorf("Expected positive speed, got %f", speed)
	}

	// Test invalid URL
	_, err = testDownloadSpeed("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL but got nil")
	}

	// Test server error
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts2.Close()

	_, err = testDownloadSpeed(ts2.URL)
	if err == nil {
		t.Error("Expected error for server error but got nil")
	}
}

// TODO: find a way to mock FTP server and test it

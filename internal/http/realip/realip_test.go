package realip

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransport_RoundTrip(t *testing.T) {
	// Create a test server to capture requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			t.Error("Expected X-Real-IP header to be set")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create transport with real IP middleware
	transport := New(http.DefaultTransport)

	// Create a request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send request through transport
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetLocalIP(t *testing.T) {
	ip := getLocalIP()
	if ip == "" {
		t.Log("Warning: Could not get local IP address")
	} else {
		t.Logf("Local IP: %s", ip)
	}
}

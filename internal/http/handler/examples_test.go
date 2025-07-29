package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func Example_updateGaugeMetric() {
	client := &http.Client{Timeout: 5 * time.Second}

	url := "http://localhost:8080/update/gauge/cpu_usage/85.5"
	req, _ := http.NewRequest("POST", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_updateCounterMetric() {
	client := &http.Client{Timeout: 5 * time.Second}

	url := "http://localhost:8080/update/counter/requests_count/100"
	req, _ := http.NewRequest("POST", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_updateMetricJSON() {
	client := &http.Client{Timeout: 5 * time.Second}

	payload := map[string]interface{}{
		"id":    "memory_usage",
		"type":  "gauge",
		"value": 75.2,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8080/update/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_batchUpdateMetrics() {
	client := &http.Client{Timeout: 5 * time.Second}

	metrics := []map[string]interface{}{
		{"id": "cpu_usage", "type": "gauge", "value": 85.5},
		{"id": "memory_usage", "type": "gauge", "value": 67.8},
		{"id": "requests_count", "type": "counter", "delta": 100},
	}

	jsonData, _ := json.Marshal(metrics)
	req, _ := http.NewRequest("POST", "http://localhost:8080/updates/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_getMetric() {
	client := &http.Client{Timeout: 5 * time.Second}

	url := "http://localhost:8080/value/gauge/cpu_usage"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_getMetricJSON() {
	client := &http.Client{Timeout: 5 * time.Second}

	payload := map[string]interface{}{
		"id":   "cpu_usage",
		"type": "gauge",
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8080/value/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_pingServer() {
	client := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequest("GET", "http://localhost:8080/ping", nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

func Example_getAllMetrics() {
	client := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()
}

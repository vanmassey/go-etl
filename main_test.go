package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleProcessMicroservice(t *testing.T) {
	// 1. Arrange a mock dataset array simulating inbound API traffic
	mockInput := []UserRecord{
		{ID: 1, Name: "Alice", Age: 25, State: "CA"},
		{ID: 2, Name: "Bob", Age: 19, State: "NY"},
		{ID: 3, Name: "Charlie", Age: 30, State: "CA"},
	}

	payload, err := json.Marshal(mockInput)
	if err != nil {
		t.Fatalf("Failed to marshal mock input data: %v", err)
	}

	// 2. Construct a mock HTTP request hitting our processing endpoint
	req, err := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create mock HTTP request: %v", err)
	}

	// 3. Initialize an HTTP response recorder to capture the microservice reply
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleProcess)

	// 4. Act: Serve the request over the mock server environment
	handler.ServeHTTP(rr, req)

	// 5. Assert: Verify the response code is 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got: %d", rr.Code)
	}

	// 6. Decode and validate the filtered streaming response payload
	var results []UserRecord
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response payload: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected exactly 2 matched records, but received: %d", len(results))
	}

	for _, user := range results {
		if user.State != "CA" || user.Age <= 21 {
			t.Errorf("Error: Invalid record leaked through pipeline filter constraints: %+v", user)
		}
	}
}

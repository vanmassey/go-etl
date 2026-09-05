package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandleProcessAPI verifies standard concurrent JSON streaming traffic
func TestHandleProcessAPI(t *testing.T) {
	mockInput := []UserRecord{
		{ID: 1, Name: "Alice", Age: 25, State: "CA"},
		{ID: 2, Name: "Bob", Age: 19, State: "NY"},
		{ID: 3, Name: "Charlie", Age: 30, State: "CA"},
	}

	payload, err := json.Marshal(mockInput)
	if err != nil {
		t.Fatalf("Failed to marshal mock input: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/process", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleProcess)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got: %d", rr.Code)
	}

	var results []UserRecord
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 matches, got: %d", len(results))
	}
}

// TestDataFramePipelineDirect isolates and validates core channel filtering
func TestDataFramePipelineDirect(t *testing.T) {
	df := NewDataFrame[UserRecord](5)

	go func() {
		df.inputStream <- UserRecord{ID: 1, Name: "Tony", Age: 45, State: "CA"}
		df.inputStream <- UserRecord{ID: 2, Name: "Steve", Age: 20, State: "NY"}
		close(df.inputStream)
	}()

	processedDf := df.Filter(2, func(u UserRecord) bool {
		return u.Age > 21 && u.State == "CA"
	})

	count := 0
	for row := range processedDf.inputStream {
		count++
		if row.Name != "Tony" {
			t.Errorf("Unexpected record passed through: %s", row.Name)
		}
	}

	if count != 1 {
		t.Errorf("Expected 1 pipeline match, got: %d", count)
	}
}

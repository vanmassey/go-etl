package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// UserRecord represents a typed row in our dataset
type UserRecord struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	State string `json:"state"`
}

// DataFrame represents our data pipeline stream
type DataFrame[T any] struct {
	inputStream chan T
}

// NewDataFrame initializes our high-speed data stream
func NewDataFrame[T any](bufferSize int) *DataFrame[T] {
	return &DataFrame[T]{
		inputStream: make(chan T, bufferSize),
	}
}

// Filter runs a concurrent filtering predicate across multiple worker threads
func (df *DataFrame[T]) Filter(numWorkers int, predicate func(T) bool) *DataFrame[T] {
	outputStream := make(chan T, cap(df.inputStream))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range df.inputStream {
				if predicate(row) {
					outputStream <- row
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputStream)
	}()

	return &DataFrame[T]{inputStream: outputStream}
}

// handleProcess Ingests JSON batches over HTTP and streams them through the dataframe pipeline
func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var records []UserRecord
	if err := json.NewDecoder(r.Body).Decode(&records); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Initialize our streaming pipeline
	df := NewDataFrame[UserRecord](100)

	// Stream the parsed JSON records into the channel concurrently
	go func() {
		for _, record := range records {
			df.inputStream <- record
		}
		close(df.inputStream)
	}()

	// Execute concurrent filtering logic across 4 goroutines
	processedDf := df.Filter(4, func(u UserRecord) bool {
		return u.Age > 21 && u.State == "CA"
	})

	// Collect the filtered results to send back in the HTTP response
	var results []UserRecord
	for row := range processedDf.inputStream {
		results = append(results)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/process", handleProcess)

	fmt.Println("Starting go-etl microservice engine on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

// UserRecord represents a typed row in our dataset
type UserRecord struct {
	ID    int
	Name  string
	Email string
	Age   int
	State string
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

func main() {
	start := time.Now()
	
	// 1. Initialize our DataFrame stream
	df := NewDataFrame[UserRecord](100)

	// 2. Concurrently read and stream data from data.csv line-by-line
	go func() {
		defer close(df.inputStream)

		file, err := os.Open("data.csv")
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		
		// Skip the header row (id, name, email, age, state)
		if _, err := reader.Read(); err != nil {
			return
		}

		// Read line-by-line streaming loop
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break // End of file reached smoothly
			}
			if err != nil {
				continue // Skip bad/corrupted rows
			}

			// Convert string fields into our type-safe fields
			id, _ := strconv.Atoi(record[0])
			age, _ := strconv.Atoi(record[3])

			// Stream the structured data straight into the pipeline
			df.inputStream <- UserRecord{
				ID:    id,
				Name:  record[1],
				Email: record[2],
				Age:   age,
				State: record[4],
			}
		}
	}()

	// 3. Filter data concurrently using 4 worker processes
	// Let's hunt for users over 21 living in California (CA)
	processedDf := df.Filter(4, func(u UserRecord) bool {
		return u.Age > 21 && u.State == "CA"
	})

	// 4. Consume and print the streamed outputs
	matchedCount := 0
	fmt.Println("Ingesting and processing data.csv concurrently...")
	fmt.Println("---------------------------------------------------------")
	
	for row := range processedDf.inputStream {
		matchedCount++
		fmt.Printf("[Matched User] ID: %d | %s | Age: %d | %s\n", row.ID, row.Name, row.Age, row.State)
	}

	fmt.Println("---------------------------------------------------------")
	fmt.Printf("Engine Completed in %v\n", time.Since(start))
	fmt.Printf("Total matches processed: %d\n", matchedCount)
}

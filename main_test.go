package main

import (
	"testing"
)

func TestDataFrameFilter(t *testing.T) {
	// 1. Initialize a test DataFrame with a buffer of 10 rows
	df := NewDataFrame[UserRecord](10)

	// 2. Insert mock test records into the pipeline channel
	go func() {
		df.inputStream <- UserRecord{ID: 1, Name: "Test User 1", Age: 25, State: "CA"}
		df.inputStream <- UserRecord{ID: 2, Name: "Test User 2", Age: 19, State: "NY"}
		df.inputStream <- UserRecord{ID: 3, Name: "Test User 3", Age: 30, State: "CA"}
		close(df.inputStream)
	}()

	// 3. Execute the filter logic using 2 worker routines
	processedDf := df.Filter(2, func(u UserRecord) bool {
		return u.Age > 21 && u.State == "CA"
	})

	// 4. Verify the results output correctly
	matchedCount := 0
	for row := range processedDf.inputStream {
		matchedCount++
		if row.State != "CA" || row.Age <= 21 {
			t.Errorf("Error: Row leaked through filter improperly: %+v", row)
		}
	}

	if matchedCount != 2 {
		t.Errorf("Error: Expected exactly 2 matches, but engine processed: %d", matchedCount)
	}
}

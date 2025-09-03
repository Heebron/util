package util

import (
	"testing"
	"time"
)

// TestNewTimeSeriesMap tests the creation of a new TimeSeries instance
func TestNewTimeSeriesMap(t *testing.T) {
	// Test case 1: Create with hourly precision
	hourly := NewTimeSeriesMap[int](time.Hour)
	if hourly == nil {
		t.Error("Expected non-nil TimeSeries instance")
	}
	if hourly.PrecisionMask != time.Hour {
		t.Errorf("Expected precision mask to be %v, got %v", time.Hour, hourly.PrecisionMask)
	}
	if hourly.Len() != 0 {
		t.Errorf("Expected new TimeSeries to be empty, got length %d", hourly.Len())
	}

	// Test case 2: Create with minute precision
	minutely := NewTimeSeriesMap[string](time.Minute)
	if minutely == nil {
		t.Error("Expected non-nil TimeSeries instance")
	}
	if minutely.PrecisionMask != time.Minute {
		t.Errorf("Expected precision mask to be %v, got %v", time.Minute, minutely.PrecisionMask)
	}
	if minutely.Len() != 0 {
		t.Errorf("Expected new TimeSeries to be empty, got length %d", minutely.Len())
	}

	// Test case 3: Create with custom precision
	customPrecision := 5 * time.Minute
	custom := NewTimeSeriesMap[float64](customPrecision)
	if custom == nil {
		t.Error("Expected non-nil TimeSeries instance")
	}
	if custom.PrecisionMask != customPrecision {
		t.Errorf("Expected precision mask to be %v, got %v", customPrecision, custom.PrecisionMask)
	}
	if custom.Len() != 0 {
		t.Errorf("Expected new TimeSeries to be empty, got length %d", custom.Len())
	}
}

// TestGet tests the Get method of TimeSeries
func TestGet(t *testing.T) {
	// Setup
	ts := NewTimeSeriesMap[int](time.Hour)
	now := time.Now()
	truncatedNow := now.Truncate(time.Hour)

	// Add a value
	value := 42
	ts.data[truncatedNow] = &value

	// Test case 1: Get existing value with exact truncated time
	result := ts.Get(truncatedNow)
	if result == nil {
		t.Error("Expected non-nil result for existing time entry")
	} else if *result != value {
		t.Errorf("Expected value %d, got %d", value, *result)
	}

	// Test case 2: Get with time that truncates to the same hour
	result = ts.Get(truncatedNow.Add(30 * time.Minute))
	if result == nil {
		t.Error("Expected non-nil result for time in same hour")
	} else if *result != value {
		t.Errorf("Expected value %d, got %d", value, *result)
	}

	// Test case 3: Get non-existent value
	result = ts.Get(truncatedNow.Add(2 * time.Hour))
	if result != nil {
		t.Errorf("Expected nil result for non-existent time entry, got %v", *result)
	}
}

// TestUpdate tests the Update method of TimeSeries
func TestUpdate(t *testing.T) {
	// Setup
	ts := NewTimeSeriesMap[int](time.Hour)
	now := time.Now()
	truncatedNow := now.Truncate(time.Hour)

	// Test case 1: Update non-existent entry (should create)
	initialValue := 10
	ts.Update(truncatedNow, func() *int {
		return &initialValue
	}, func(v *int) {
		// This shouldn't be called for a new entry
		t.Error("Update function called for new entry")
	})

	result := ts.Get(truncatedNow)
	if result == nil {
		t.Error("Expected non-nil result after Update")
	} else if *result != 10 {
		t.Errorf("Expected value %d, got %d", 10, *result)
	}

	// Test case 2: Update existing entry
	ts.Update(truncatedNow, func() *int {
		// This shouldn't be called for an existing entry
		t.Error("Init function called for existing entry")
		return new(int)
	}, func(v *int) {
		*v = *v + 5
	})

	result = ts.Get(truncatedNow)
	if result == nil {
		t.Error("Expected non-nil result after second Update")
	} else if *result != 15 {
		t.Errorf("Expected value %d, got %d", 15, *result)
	}

	// Test case 3: Update with time that truncates to the same hour
	ts.Update(truncatedNow.Add(30*time.Minute), func() *int {
		// This shouldn't be called for an existing entry
		t.Error("Init function called for existing entry")
		return new(int)
	}, func(v *int) {
		*v = *v + 7
	})

	result = ts.Get(truncatedNow)
	if result == nil {
		t.Error("Expected non-nil result after third Update")
	} else if *result != 22 {
		t.Errorf("Expected value %d, got %d", 22, *result)
	}
}

// TestLen tests the Len method of TimeSeries
func TestLen(t *testing.T) {
	// Setup
	ts := NewTimeSeriesMap[string](time.Hour)

	// Test case 1: Empty TimeSeries
	if ts.Len() != 0 {
		t.Errorf("Expected length 0 for empty TimeSeries, got %d", ts.Len())
	}

	// Test case 2: Add one entry
	now := time.Now()
	value := "test"
	ts.data[now.Truncate(time.Hour)] = &value

	if ts.Len() != 1 {
		t.Errorf("Expected length 1 after adding one entry, got %d", ts.Len())
	}

	// Test case 3: Add another entry for a different hour
	anotherValue := "another"
	ts.data[now.Add(2*time.Hour).Truncate(time.Hour)] = &anotherValue

	if ts.Len() != 2 {
		t.Errorf("Expected length 2 after adding second entry, got %d", ts.Len())
	}

	// Test case 4: Overwrite an existing entry (length shouldn't change)
	updatedValue := "updated"
	ts.data[now.Truncate(time.Hour)] = &updatedValue

	if ts.Len() != 2 {
		t.Errorf("Expected length to remain 2 after overwriting entry, got %d", ts.Len())
	}
}

// TestClear tests the Clear method of TimeSeries
func TestClear(t *testing.T) {
	// Setup
	ts := NewTimeSeriesMap[float64](time.Hour)
	now := time.Now()

	// Add some entries
	value1 := 1.1
	value2 := 2.2
	ts.data[now.Truncate(time.Hour)] = &value1
	ts.data[now.Add(time.Hour).Truncate(time.Hour)] = &value2

	// Verify initial state
	if ts.Len() != 2 {
		t.Errorf("Expected initial length 2, got %d", ts.Len())
	}

	// Test case 1: Clear the TimeSeries
	ts.Clear()

	// Verify it's empty
	if ts.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", ts.Len())
	}

	// Test case 2: Clear an already empty TimeSeries
	ts.Clear()

	// Verify it's still empty
	if ts.Len() != 0 {
		t.Errorf("Expected length 0 after clearing empty TimeSeries, got %d", ts.Len())
	}
}

// TestKeys tests the Keys method of TimeSeries
func TestKeys(t *testing.T) {
	// Setup
	ts := NewTimeSeriesMap[int](time.Hour)
	now := time.Now()
	baseTime := now.Truncate(time.Hour)

	// Test case 1: Empty TimeSeries
	keys := ts.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected empty keys slice for empty TimeSeries, got %d keys", len(keys))
	}

	// Add entries in non-chronological order
	value1 := 1
	value2 := 2
	value3 := 3

	// Add in reverse order to test sorting
	thirdTime := baseTime.Add(2 * time.Hour)
	secondTime := baseTime.Add(time.Hour)
	firstTime := baseTime

	ts.data[thirdTime] = &value3
	ts.data[firstTime] = &value1
	ts.data[secondTime] = &value2

	// Test case 2: TimeSeries with entries
	keys = ts.Keys()

	// Check length
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Check order (should be chronological)
	if len(keys) >= 3 {
		if !keys[0].Equal(firstTime) {
			t.Errorf("Expected first key to be %v, got %v", firstTime, keys[0])
		}
		if !keys[1].Equal(secondTime) {
			t.Errorf("Expected second key to be %v, got %v", secondTime, keys[1])
		}
		if !keys[2].Equal(thirdTime) {
			t.Errorf("Expected third key to be %v, got %v", thirdTime, keys[2])
		}
	}
}

// TestTimeSeriesWithStructs tests TimeSeries with a struct type
func TestTimeSeriesWithStructs(t *testing.T) {
	// Define a test struct
	type TestStruct struct {
		Name  string
		Count int
	}

	// Setup
	ts := NewTimeSeriesMap[TestStruct](time.Hour)
	now := time.Now()

	// Test case: Add and retrieve a struct
	initialStruct := TestStruct{Name: "Test", Count: 1}
	ts.Update(now, func() *TestStruct {
		return &initialStruct
	}, func(v *TestStruct) {
		// This shouldn't be called for a new entry
	})

	result := ts.Get(now)
	if result == nil {
		t.Error("Expected non-nil result after Update with struct")
	} else {
		if result.Name != initialStruct.Name {
			t.Errorf("Expected Name %s, got %s", initialStruct.Name, result.Name)
		}
		if result.Count != initialStruct.Count {
			t.Errorf("Expected Count %d, got %d", initialStruct.Count, result.Count)
		}
	}

	// Update the struct
	ts.Update(now, func() *TestStruct {
		return new(TestStruct)
	}, func(v *TestStruct) {
		v.Count++
	})

	result = ts.Get(now)
	if result == nil {
		t.Error("Expected non-nil result after Update with struct")
	} else {
		if result.Count != 2 {
			t.Errorf("Expected Count 2 after update, got %d", result.Count)
		}
	}
}

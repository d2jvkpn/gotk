package gox

import (
	// "fmt"
	"testing"
	"time"
)

func TestTimeCeil(t *testing.T) {
	// Set up a test time
	testTime, err := time.Parse("2006-01-02 15:04:05", "2022-12-11 12:34:56")
	if err != nil {
		t.Errorf("Error parsing test time: %v", err)
	}

	// Test the second time unit
	expected, err := time.Parse("2006-01-02 15:04:05", "2022-12-11 12:34:56")
	if err != nil {
		t.Errorf("Error parsing expected time: %v", err)
	}
	result, err := TimeCeil(testTime, "second")
	if err != nil {
		t.Errorf("Error calling TimeCeil: %v", err)
	}
	if !result.Equal(expected) {
		t.Errorf("Expected TimeCeil(%v, \"second\") to be %v, got %v", testTime, expected, result)
	}

	// Test the minute time unit
	expected, err = time.Parse("2006-01-02 15:04:05", "2022-12-11 12:34:00")
	if err != nil {
		t.Errorf("Error parsing expected time: %v", err)
	}
	result, err = TimeCeil(testTime, "minute")
	if err != nil {
		t.Errorf("Error calling TimeCeil: %v", err)
	}
	if !result.Equal(expected) {
		t.Errorf("Expected TimeCeil(%v, \"minute\") to be %v, got %v", testTime, expected, result)
	}

	// Test the hour time unit
	expected, err = time.Parse("2006-01-02 15:04:05", "2022-12-11 12:00:00")
	if err != nil {
		t.Errorf("Error parsing expected time: %v", err)
	}
	result, err = TimeCeil(testTime, "hour")
	if err != nil {
		t.Errorf("Error calling TimeCeil: %v", err)
	}
	if !result.Equal(expected) {
		t.Errorf("Expected TimeCeil(%v, \"hour\") to be %v, got %v", testTime, expected, result)
	}

	// Test the day time unit
	expected, err = time.Parse("2006-01-02 15:04:05", "2022-12-11 00:00:00")
	if err != nil {
		t.Errorf("Error parsing expected time: %v", err)
	}
	result, err = TimeCeil(testTime, "day")
	if err != nil {
		t.Errorf("Error calling TimeCeil: %v", err)
	}
}

func TestTimeFloor_second(t *testing.T) {
	// Set up test input and expected output
	at := time.Date(2020, 1, 1, 12, 30, 45, 0, time.UTC)
	expected := time.Date(2020, 1, 1, 12, 30, 45, 0, time.UTC)

	// Call the function and check the output
	out, err := TimeFloor(at, "second")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !out.Equal(expected) {
		t.Errorf("Unexpected output. Expected: %v, Got: %v", expected, out)
	}
}

func TestTimeFloor_minute(t *testing.T) {
	// Set up test input and expected output
	at := time.Date(2020, 1, 1, 12, 30, 45, 0, time.UTC)
	expected := time.Date(2020, 1, 1, 12, 30, 0, 0, time.UTC)

	// Call the function and check the output
	out, err := TimeFloor(at, "minute")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !out.Equal(expected) {
		t.Errorf("Unexpected output. Expected: %v, Got: %v", expected, out)
	}
}

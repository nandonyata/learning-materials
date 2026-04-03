package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	var testCases = []struct {
		name  string
		delay time.Duration
	}{
		{
			name:  "Delay 1 second",
			delay: 1 * time.Second,
		},
		{
			name:  "Delay half second",
			delay: 500 * time.Millisecond,
		},
	}

	// orderResult
	for _, e := range testCases {
		t.Run(e.name, func(t *testing.T) {
			// Reset the orderResult slice
			orderResult = []string{}

			// Set delays
			eatTime = e.delay
			sleepTime = e.delay
			thinkTime = e.delay

			// Run the dine function
			dine()

			// Check if the orderResult is correct
			if len(orderResult) != len(philosophers) {
				t.Errorf("Expected %d philosophers, got %d", len(philosophers), len(orderResult))
			}
		})
	}
}

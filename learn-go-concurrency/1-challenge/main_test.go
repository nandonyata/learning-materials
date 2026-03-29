package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s  string
		wg *sync.WaitGroup
	}{
		// TODO: Add test cases.
		{
			name: "Testing Update Message To Pewpew",
			s:    "pewpew",
			wg:   &sync.WaitGroup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wg.Add(1)
			updateMessage(tt.s, tt.wg)
			tt.wg.Wait()

			if msg != tt.s {
				t.Errorf("Expected pewpew, but got %s", msg)
			}
		})
	}
}

func Test_printMessage(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
		{
			name: "Testing printing message",
		},
	}

	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w
	main() // call the main to populate the "msg"
	w.Close()
	read, _ := io.ReadAll(r)
	output := string(read)
	os.Stdout = stdOut

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(output, "Hello, universe!") {
				t.Errorf("Expected Hello, universe!, but got %s", output)
			}
			if !strings.Contains(output, "Hello, cosmos!") {
				t.Errorf("Expected Hello, cosmos!, but got %s", output)
			}
			if !strings.Contains(output, "Hello, world!") {
				t.Errorf("Expected Hello, world!, but got %s", output)
			}
		})
	}
}

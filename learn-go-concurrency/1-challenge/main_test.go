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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdOut := os.Stdout

			r, w, _ := os.Pipe()
			os.Stdout = w

			printMessage()

			w.Close()

			read, _ := io.ReadAll(r)
			output := string(read)

			os.Stdout = stdOut

			if !strings.Contains(output, "pewpew") { // we got the value "pewpew" from "Test_updateMessage"
				t.Errorf("Expected pewpew, but got %s", output)

			}
		})
	}
}

package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printSomething(t *testing.T) {
	stdOut := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = writer

	var wg sync.WaitGroup
	wg.Add(1)
	go printSomething("pewpew", &wg)
	wg.Wait()

	writer.Close()

	result, _ := io.ReadAll(reader)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "pewpew") {
		t.Errorf("Expected: pewpew, got: %s", output)
	}
}

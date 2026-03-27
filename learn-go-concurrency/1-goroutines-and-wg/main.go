package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup

	words := []string{
		"Abe",
		"Cee",
		"Book",
		"Egg",
		"Lorem",
		"Ipsum",
		"Dolor",
	}

	for i, v := range words {
		wg.Add(1)
		go printSomething(fmt.Sprintf("%d: %s", i, v), &wg)
	}

	wg.Wait()
}

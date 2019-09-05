package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kjnsn/password-checker/pkg/filter"
)

func main() {
	words := make([][]byte, 0, 0)
	if err := readLines(&words); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Words read: %d\n", len(words))

	f := filter.NewFilter(uint(len(words)))

	// Test the filter for false positive rate.
	fpRate := f.EstimateFalsePositiveRate(uint(len(words)))
	fmt.Printf("False positive rate: %f\n", fpRate)

	// Add all the words.
	for _, word := range words {
		f.Add(word)
	}
}

func readLines(words *[][]byte) error {
	if words == nil {
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		*words = append(*words, scanner.Bytes())
	}
	return scanner.Err()
}

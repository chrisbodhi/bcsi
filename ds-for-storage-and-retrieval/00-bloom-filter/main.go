package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const wordsPath = "/usr/share/dict/words"

func loadWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	words, err := loadWords(wordsPath)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	var b bloomFilter = newTrivialBloomFilter()
	// We want the size to be a bit more than the number
	// of words in /usr/share/dict/words, which is 235,886
	var cb bloomFilter = newCbBloomFilter(0.01, 24000)

	wordsAdded := 0
	// Add every other word (even indices)
	for i := 0; i < len(words); i += 2 {
		b.add(words[i])
		cb.add(words[i])
		wordsAdded++
	}

	// Make sure there are no false negatives
	for i := 0; i < len(words); i += 2 {
		word := words[i]
		if !b.maybeContains(word) {
			log.Fatalf("false negative for word %q\n", word)
		}
		if !cb.maybeContains(word) {
			log.Fatalf("cb: false negative for word %q\n", word)
		}
	}

	falsePositives := 0
	cbFalsePositives := 0
	numChecked := 0

	// None of the words at odd indices were added, so whenever
	// maybeContains returns true, it's a false positive
	for i := 1; i < len(words); i += 2 {
		if b.maybeContains(words[i]) {
			falsePositives++
		}
		if cb.maybeContains(words[i]) {
			cbFalsePositives++
		}
		numChecked++
	}

	falsePositiveRate := float64(falsePositives) / float64(numChecked)
	cbFalsePositiveRate := float64(cbFalsePositives) / float64(numChecked)

	fmt.Printf("Elapsed time: %s\n", time.Since(start))
	fmt.Println("* * * * * *")
	fmt.Printf("Memory usage: %d bytes\n", b.memoryUsage())
	fmt.Printf("False positive rate: %0.2f%%\n", 100*falsePositiveRate)
	fmt.Printf("Words added: %d\n", wordsAdded)

	fmt.Println("* * * * * *")

	fmt.Printf("CB memory usage: %d bytes (%d kb)\n", cb.memoryUsage(), cb.memoryUsage() / 1000)
	fmt.Printf("CB false positive rate: %0.2f%%\n", 100*cbFalsePositiveRate)
}

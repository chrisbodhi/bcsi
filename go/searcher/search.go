package searcher

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chrisbodhi/bcsi/go/xkcd/indexer"
	"github.com/chrisbodhi/bcsi/go/xkcd/printer"
)

// Search searches the xkcd index
func Search(search string) {
	fmt.Printf("\nSearching for %s\n\n", search)

	f, err := os.Open(indexer.IndexPath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := csv.NewReader(f)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			fmt.Println("\nDone with your search.")
			break
		}
		lineMatches := strings.Contains(strings.ToLower(strings.Join(line, " ")), strings.ToLower(search))
		if lineMatches {
			printer.Print(line)
		}
	}
}

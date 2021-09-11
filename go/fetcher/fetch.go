package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chrisbodhi/bcsi/xkcd/indexer"
)

// FetchComic gets the URL's contents and returns them, ready for saving
func FetchComic(url string) indexer.Comic {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(url)
		fmt.Println(err)
		fmt.Println("fml")
	}

	defer resp.Body.Close()

	var comic indexer.Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		log.Fatal(err)
	}
	return comic
}

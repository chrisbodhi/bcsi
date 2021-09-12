package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Comic represents the minimal information we require, out of all of the fields provided by the xkcd JSON
type Comic struct {
	Num        int
	Transcript string
	Alt        string
}

// FetchComic gets the URL's contents and returns them, ready for saving
func FetchComic(url string) Comic {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(url)
		fmt.Println(err)
		fmt.Println("fml")
	}

	defer resp.Body.Close()

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		log.Fatal(err)
	}

	return comic
}

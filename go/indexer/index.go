package indexer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

var home, _ = os.UserHomeDir()

// IndexPath is where our xkcd index file is stored
var IndexPath = path.Join(home, ".xkcddb.csv")

const LatestXkcdUrl = "https://xkcd.com/info.0.json"

// Comic represents the minimal information we require, out of all of the fields provided by the xkcd JSON
type Comic struct {
	Num        int
	Transcript string
	Alt        string
}

func newComic(num int, transcript string, alt string) *Comic {
	c := Comic{num, transcript, alt}
	return &c
}

func exists() bool {
	_, err := os.Stat(IndexPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func updateIndex() {
	fmt.Println("booyah")
}

func structToRow(c Comic) string {
	return fmt.Sprintf("%d,\"%s\",\"%s\"\n", c.Num, c.Transcript, c.Alt)
}

func appendToIndex(row string) {
	index, err := os.OpenFile(IndexPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cannot open file")
		log.Fatal(err)
	}

	defer index.Close()

	if _, err := index.WriteString(row); err != nil {
		fmt.Println("cannot write string", row)
		log.Fatal(err)
	}
}

func populateIndex() {
	resp, err := http.Get(LatestXkcdUrl)
	if err != nil {
		fmt.Println("fml")
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	row := structToRow(comic)
	appendToIndex(row)
}

func makeIndex() {
	err := os.WriteFile(IndexPath, []byte("num,transcript,alt\n"), 0666)
	if err != nil {
		log.Print(err)
		log.Fatal("Could not create index.")
	}
	populateIndex()
}

// GetOrMake retrieves the index, if it exists, so that we may update it.
// If the index does not exist, we create the file and populate it with all of
// the data.
func GetOrMake() {
	if exists() {
		fmt.Println("Checking for updates...")
		updateIndex()
	} else {
		fmt.Println("Creating index...")
		makeIndex()
	}
}

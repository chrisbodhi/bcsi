package indexer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

var home, _ = os.UserHomeDir()

// IndexPath is where our xkcd index file is stored
var IndexPath = path.Join(home, ".xkcddb.csv")

const latestXkcdURL = "https://xkcd.com/info.0.json"

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

func indexExists() bool {
	_, err := os.Stat(IndexPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func updateIndex(startingIndex int) {
	latestComic := fetchComic(latestXkcdURL)
	indexLatestNum := getLatestNum()
	if latestComic.Num > indexLatestNum {
		// append latestComic to our index
		row := structToRow(latestComic)
		appendToIndex(row)
		// start to iterate
		start := latestComic.Num - 1 // since we have already added the latest comic, start with the comic that immediately precedes it
		for ; start > indexLatestNum; start-- {
			url := fmt.Sprint("https://xkcd.com/%d/info.0.json", start)
			addComicToIndex(url)
		}
	}
}

func structToRow(c Comic) string {
	return fmt.Sprintf("%d,\"%s\",\"%s\"\n", c.Num, c.Transcript, c.Alt)
}

func getIndexContents() [][]string {
	file, err := os.Open(IndexPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records[1:]
}

func getLatestNum() int {
	records := getIndexContents()
	var max int
	for _, record := range records {
		num, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		if num > max {
			max = num
		}
	}
	return max
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

func fetchComic(url string) Comic {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("fml")
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		log.Fatal(err)
	}
	return comic
}

func addComicToIndex(url string) {
	comic := fetchComic(url)
	row := structToRow(comic)
	appendToIndex(row)
}

func populateIndex(url string) {
	addComicToIndex(url)
}

func makeIndex(url string) {
	err := os.WriteFile(IndexPath, []byte("num,transcript,alt\n"), 0666)
	if err != nil {
		log.Print(err)
		log.Fatal("Could not create index.")
	}
	populateIndex(url)
}

// GetOrMake retrieves the index, if it exists, so that we may update it.
// If the index does not exist, we create the file and populate it with all of
// the data.
func GetOrMake() {
	if indexExists() {
		fmt.Println("Checking for updates...")
		latestFromIndex := getLatestNum()
		updateIndex(latestFromIndex)
	} else {
		fmt.Println("Creating index...")
		makeIndex(latestXkcdURL)
	}
}

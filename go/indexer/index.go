package indexer

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/chrisbodhi/bcsi/go/xkcd/fetcher"
)

var home, _ = os.UserHomeDir()

// IndexPath is where our xkcd index file is stored
var IndexPath = path.Join(home, ".xkcddb.csv")

const latestXkcdURL = "https://xkcd.com/info.0.json"

func newComic(num int, transcript string, alt string) *fetcher.Comic {
	c := fetcher.Comic{num, transcript, alt}
	return &c
}

func indexExists() bool {
	_, err := os.Stat(IndexPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func updateIndex(startingIndex int, latestComic fetcher.Comic) {
	if latestComic.Num > startingIndex {
		// append latestComic to our index
		row := structToRow(latestComic)
		appendToIndex(row)
		// start to iterate
		start := latestComic.Num - 1 // since we have already added the latest comic, start with the comic that immediately precedes it
		loopIt(start, startingIndex)
	}
}

// when we populate to start with, we pass in (latest from xkcd.com/info.0.json - 1) and zero
func loopIt(start int, startingIndex int) {
	for ; start > startingIndex; start-- {
		if start == 404 {
			fmt.Println("Just passing through...")
			continue
		}
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", start)
		addComicToIndex(url)
	}
}

func structToRow(c fetcher.Comic) []string {
	quoteRe := regexp.MustCompile(`"`)

	transcript := quoteRe.ReplaceAllString(c.Transcript, "'")
	alt := quoteRe.ReplaceAllString(c.Alt, "'")
	num := fmt.Sprintf("%d", c.Num)

	transcript = strings.ReplaceAll(transcript, "\n", " ")
	alt = strings.ReplaceAll(alt, "\n", " ")

	return []string{num, transcript, alt}
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
		fmt.Println("cannot read from file")
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

func appendToIndex(row []string) {
	f, err := os.OpenFile(IndexPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	w := csv.NewWriter(f)
	w.Write(row)

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	fmt.Print(".")
}

func addComicToIndex(url string) {
	comic := fetcher.FetchComic(url)
	row := structToRow(comic)
	appendToIndex(row)
}

func populateIndex(url string) {
	addComicToIndex(url)
	latest := getLatestNum()
	loopIt(latest-1, 0)
}

func makeIndex() {
	f, err := os.Create(IndexPath)
	if err != nil {
		log.Fatal(err)
	}
	records := [][]string{
		{"num", "transcript", "alt"},
	}

	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

// GetOrMake retrieves the index, if it exists, so that we may update it.
// If the index does not exist, we create the file and populate it with all of
// the data.
func GetOrMake() {
	if indexExists() {
		fmt.Println("Checking for updates...")
		latestFromIndex := getLatestNum()
		latestComic := fetcher.FetchComic(latestXkcdURL)
		updateIndex(latestFromIndex, latestComic)
		fmt.Println("All up to date.")
	} else {
		fmt.Println("Creating index...")
		makeIndex()
		populateIndex(latestXkcdURL)
		fmt.Println("Index created.")
	}
}

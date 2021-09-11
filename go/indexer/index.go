package indexer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
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

func structToRow(c Comic) string {
	quoteRe := regexp.MustCompile(`"`)
	transcript := quoteRe.ReplaceAllString(c.Transcript, "'")
	alt := quoteRe.ReplaceAllString(c.Alt, "'")
	transcript = strings.ReplaceAll(transcript, "\n", " ")
	alt = strings.ReplaceAll(alt, "\n", " ")
	return fmt.Sprintf("%d,\"%s\",\"%s\"\n", c.Num, transcript, alt)
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

// TODO replace with CSV writer
// https://pkg.go.dev/encoding/csv@go1.17#Writer.Write
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

func addComicToIndex(url string) {
	comic := fetchComic(url)
	row := structToRow(comic)
	appendToIndex(row)
}

func populateIndex(url string) {
	addComicToIndex(url)
	latest := getLatestNum()
	loopIt(latest-1, 0)
}

func makeIndex(url string) {
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
		fmt.Println("All up to date.")
	} else {
		fmt.Println("Creating index...")
		makeIndex(latestXkcdURL)
		fmt.Println("Index created.")
	}
}

package printer

import "fmt"

// Print puts the formatted line on stdout for the user's clicking pleasure
func Print(line []string) {
	url := fmt.Sprintf("https://xkcd.com/%s", line[0])
	transcript := line[1]
	alt := line[2]
	fmt.Println("URL:\t\t", url)
	fmt.Println("Alt:\t\t", alt)
	fmt.Println("Transcript:\t", transcript)
	fmt.Println("* * * * *")

}

package main

import (
	"bufio"
	"fmt"
	"os"
)

// S is So Srs, it'S uppercaSe.
const S = "🐢 "

func looper() ([]byte, error) {
	var s []byte
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(S)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("\nSee ya, slow poke. 💨")
			return s, err
		}
		if b == '\n' {
			break
		}
		s = append(s, b)
	}
	return s, nil
}

func main() {
	bs, err := looper()
	if err != nil {
		return
	}
	fmt.Println(string(bs))
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func eval(b []byte) {
	input := string(b)
	inputs := strings.Fields(input)
	userCmd := inputs[0]
	userArgs := inputs[1:]
	cmd := exec.Command(userCmd, userArgs...)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", output)
}

func handleInput() {
	bs, err := looper()
	if err != nil {
		return
	}
	eval(bs)
	handleInput()
}

func main() {
	handleInput()
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
	// TODO split on space
	// TODO create a command that is first elem
	// TODO pass that to exec.Command
	// https://pkg.go.dev/os/exec#Command
	cmd := exec.Command(input)
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

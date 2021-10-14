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

var state = map[string]bool{
	"continue": true,
}

func looper() ([]byte, error) {
	var s []byte
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(S)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			leave()
			return s, err
		}
		if b == '\n' {
			break
		}
		s = append(s, b)
	}
	return s, nil
}

// 'exit' command
func leave() {
	state["continue"] = false
	fmt.Println("\nSee ya, slow poke. 💨")
}

// 'cd' command
func changeDir(s string) {
	os.Chdir(s)
}

// default; calls to parent shell
func callOut(userCmd string, userArgs []string) {
	cmd := exec.Command(userCmd, userArgs...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("I'm just a simple turtle, I don't know how to", userCmd)
	}
	fmt.Printf("%s", output)
}

func eval(b []byte) {
	input := string(b)
	inputs := strings.Fields(input)
	if len(inputs) > 0 {
		userCmd := inputs[0]
		userArgs := inputs[1:]
		switch userCmd {
		case "exit":
			leave()
			return
		case "cd":
			var dest string
			if len(userArgs) == 0 {
				dest = "/"
			} else {
				dest = userArgs[0]
			}
			changeDir(dest)
			return
		default:
			callOut(userCmd, userArgs)
		}
	}
}

func handleInput() {
	bs, err := looper()
	if err != nil {
		return
	}
	eval(bs)
	if state["continue"] {
		handleInput()
	}
}

func main() {
	handleInput()
}

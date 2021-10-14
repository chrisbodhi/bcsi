package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const prompt = "🐢 "
const interrupt = "signal: interrupt"

var state = map[string]bool{
	"continue": true,
}

func looper() ([]byte, error) {
	var s []byte
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(prompt)
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
		if err.Error() != interrupt {
			fmt.Println("I'm just a simple turtle, I don't know how to", userCmd)
		}
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// TODO gets weird when I start the shell and then just type "exit"
	handleInput()
}

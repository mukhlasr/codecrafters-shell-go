package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Fatal(err)
			break
		}

		cmd := parseCmd(str[:len(str)-1])

		switch cmd.name {
		case "exit":
			handleExit(cmd)
		default:
			fmt.Printf("%s: command not found\n", cmd.name)
		}
	}
}

type cmd struct {
	name string
	args []string
}

func parseCmd(str string) cmd {
	// Split the string into words
	words := strings.Fields(str)

	// The first word is the command name
	name := words[0]

	// The rest of the words are the arguments
	args := words[1:]

	return cmd{name, args}
}

func handleExit(cmd cmd) {
	if cmd.name != "exit" {
		return
	}

	code := 0

	if len(cmd.args) == 1 {
		var err error
		code, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Printf("exit: %s: numeric argument required\n", cmd.args[0])
			code = 255
		}
	}

	os.Exit(code)
}

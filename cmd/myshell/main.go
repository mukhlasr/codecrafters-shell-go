package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	lastExitCode int
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd := parseCmd()
		if len(cmd) < 1 {
			continue
		}

		builtinName, builtinArgs := cmd[0], cmd[1:]
		builtinFunc, ok := getBuiltIns()[builtinName]
		if !ok {
			lastExitCode = execute(cmd)
			continue
		}

		builtinFunc(builtinArgs)
	}
}

func parseCmd() []string {
	r := bufio.NewReader(os.Stdin)

	var (
		res           []string
		onDoubleQuote bool
		onWhitespace  bool
		str           string
	)

loop:
	for {
		c, err := r.ReadByte()
		if errors.Is(err, io.EOF) {
			os.Exit(0)
		}

		if err != nil {
			log.Fatalln(err)
		}

		if isWhiteSpace(c) {
			if onDoubleQuote {
				str += string(c)
				continue
			}

			if onWhitespace {
				continue
			}

			res = append(res, str)
			str = ""
			continue
		}

		onWhitespace = false

		switch c {
		case '"':
			onDoubleQuote = !onDoubleQuote
		case '\n':
			if onDoubleQuote {
				str += string(c)
				fmt.Print("> ")
				continue
			}

			res = append(res, str)
			break loop
		default:
			str += string(c)
		}
	}

	return res
}

func handleExit(args []string) {
	code := 0

	if len(args) > 0 {
		var err error
		code, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("exit: %s: numeric argument required\n", args[0])
			code = 255
		}
	}

	os.Exit(code)
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleType(args []string) {
	cmdName := args[0]
	_, ok := getBuiltIns()[cmdName]
	if ok {
		fmt.Printf("%s is a shell builtin\n", cmdName)
		return
	}

	path, ok := isCmdExists(cmdName)
	if !ok {
		fmt.Printf("%s: not found\n", cmdName)
		return
	}

	fmt.Printf("%s is %s\n", cmdName, path)
}

func execute(cmd []string) (exitCode int) {
	if len(cmd) < 1 {
		return
	}

	path, exist := isCmdExists(cmd[0])

	if !exist {
		fmt.Printf("%s: command not found\n", cmd[0])
		return
	}

	execCmd := exec.Command(path, cmd[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()

	var execErr *exec.ExitError

	if errors.As(err, &execErr) {
		code := execErr.ExitCode()
		return code
	}

	if err != nil {
		fmt.Println(err)
		return 1
	}

	return
}

// isCmdExists returns value is the path of the command and return true if exists.
// If the command does not exist, it returns false.
func isCmdExists(cmd string) (string, bool) {
	path := os.Getenv("PATH")
	paths := strings.Split(path, ":")

	for _, p := range paths {
		cmdPath := fmt.Sprintf("%s/%s", p, cmd)
		_, err := os.Stat(cmdPath)
		if err == nil {
			return cmdPath, true
		}
	}

	return "", false
}

func getBuiltIns() map[string]func([]string) {
	return map[string]func([]string){
		"exit": handleExit,
		"echo": handleEcho,
		"type": handleType,
	}
}

// isWhiteSpace returns true if the given character is a white space(' ', '\t', '\r')
// \n is not considered as a white space.
func isWhiteSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

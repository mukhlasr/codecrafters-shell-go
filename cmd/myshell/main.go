package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	str, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}

	switch cmd := str[:len(str)-1]; cmd {
	default:
		fmt.Printf("%s: command not found\n", cmd)
	}

}

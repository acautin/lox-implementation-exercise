package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: tree-walk.exe [script]")
		os.Exit(65)
	} else if len(os.Args) == 2 {
		fmt.Println("Running file: " + os.Args[1] + "")
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("got from file:\n" + string(bytes))
	run(string(bytes))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func run(source string) {
	tokens := scanner.ScanTokens(source)
	for _, token := range tokens {
		fmt.Println(token)
	}
}

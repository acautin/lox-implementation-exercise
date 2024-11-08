package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/acautin/lox-implementation-exercise/tree-walk/interpreter"
	"github.com/acautin/lox-implementation-exercise/tree-walk/parser"
	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: tree [script]")
		os.Exit(65)
	} else if len(os.Args) == 2 {
		fmt.Println("Running file: " + os.Args[1])
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
	run(string(bytes))
}

func runPrompt() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !input.Scan() {
			break
		}
		line := input.Text()
		run(line)
	}
	if err := input.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func run(source string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "Error:", r)
		}
	}()

	// Step 1: Scan the source code into tokens
	tokens := scanner.ScanTokens(source)

	// Step 2: Parse the tokens into an AST
	expression, err := parser.Parse(tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Parse error:", err)
		return
	}

	// Step 3: Interpret the AST
	interp := interpreter.NewInterpreter()
	result, err := interp.Interpret(expression)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Interpretation error:", err)
		return
	}

	// Step 4: Print the result
	fmt.Println(result)
}

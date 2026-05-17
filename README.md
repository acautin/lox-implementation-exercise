# Lox Implementation Exercise

This repository contains a tree-walk interpreter for the **Lox programming language**, as described in the book [*Crafting Interpreters*](https://craftinginterpreters.com/) by Robert Nystrom. The implementation is written in **Go**.

## Project Structure

The interpreter is located in the `tree-walk/` directory and is divided into several modules:

- **`scanner/`**: A lexical analyzer that converts raw source code into tokens. It supports keywords, literals, and multi-line comments.
- **`parser/`**: A recursive descent parser that builds an Abstract Syntax Tree (AST) from the tokens.
- **`interpreter/`**: A visitor-based evaluator that executes the AST.
- **`main.go`**: The entry point providing both a REPL and script execution modes.

## Current Status

The project currently supports **expressions** and basic evaluation:

- [x] **Lexical Scanning**: Full support for Lox tokens, including nested comments.
- [x] **AST Generation**: Support for literals, unary, binary, and grouping expressions.
- [x] **Expression Evaluation**:
    - Arithmetic: `+`, `-`, `*`, `/` (including string concatenation).
    - Comparison: `>`, `>=`, `<`, `<=`
    - Equality: `==`, `!=`
    - Logic: `!`
- [ ] **Statements & State**: Variables, assignments, and block scopes (Pending).
- [ ] **Control Flow**: `if`, `while`, `for` (Pending).
- [ ] **Functions & Classes** (Pending).

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) 1.18 or higher.

### Building

To compile the interpreter into an executable from the root directory:

```bash
go build -C tree-walk -o ../golox .
```

The `-o` flag allows you to specify the name and path of the output binary (in this case, `golox`). 

### Running the Interpreter

From the root directory:

#### REPL Mode
Start an interactive session:
```bash
go run tree-walk/main.go
```

#### Script Mode
Execute a `.lox` file:
```bash
go run tree-walk/main.go path/to/script.lox
```

## Examples

You can currently evaluate expressions like these in the REPL:

```lox
> (5 - 3) * 2 == 4
true
> "Go" + " " + "Lox"
Go Lox
> !(5 > 10)
true
```

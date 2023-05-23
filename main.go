package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	e "github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/interpreter"
	"github.com/redundant4u/Golox/internal/parser"
	"github.com/redundant4u/Golox/internal/scanner"
)

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Colud not read file: %w", err)
	}

	_ = run(string(bytes))

	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Could not read line: %w", err)
		}

		err = run(string(line))
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func run(source string) error {
	sc := scanner.New(source)
	tokens, err := sc.ScanTokens()
	if err != nil {
		return fmt.Errorf("Failed to scan tokens")
	}

	parser := parser.New(tokens)
	statements := parser.Parse()
	if e.HadError {
		return fmt.Errorf("Failed to parse tokens")
	}

	interpreter := interpreter.New()
	interpreter.Interpret(statements)

	return nil
}

func main() {
	var err error

	if len(os.Args) > 2 {
		fmt.Println("Usage: go run main.go or golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err = runFile(os.Args[1])
		if e.HadError {
			os.Exit(65)
		} else if e.HadRuntimeError {
			os.Exit(70)
		}
	} else {
		err = runPrompt()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

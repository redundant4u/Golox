package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/redundant4u/Golox/internal/scanner"
)

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Cloud not read file: %w", err)
	}

	return run(string(bytes))
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadBytes('\n')
		errCheck(err)

		err = run(string(line))
		errCheck(err)
	}
}

func run(source string) error {
	sc := scanner.New(source)
	tokens, err := sc.ScanTokens()
	errCheck(err)

	for _, token := range tokens {
		fmt.Println(token)
	}

	return nil
}

func main() {
	var err error

	fmt.Println(len(os.Args))

	if len(os.Args) > 2 {
		fmt.Println("Usage: go run main.go OR golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err = runFile(os.Args[1])
	} else {
		err = runPrompt()
	}

	errCheck(err)
}

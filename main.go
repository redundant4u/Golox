package main

import (
	"Golox/lox"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	hadError        = false
	hadRuntimeError = false
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runFile(file string) {
	dat, err := ioutil.ReadFile(file)
	check(err)
	run(string(dat))

	if hadError {
		os.Exit(65)
	}

	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		dat, err := reader.ReadBytes('\n')
		check(err)
		run(string(dat))
		// hadError = false
	}
}

func run(src string) {
	scanner := lox.NewScanner(src)
	tokens := scanner.ScanTokens()
	parser := lox.NewParser(tokens)
	statements := parser.Parse()

	interpreter := lox.NewInterpreter()
	interpreter.Interprete(statements)

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }
}

func main() {
	file := flag.String("file", "", "the script file to execute")
	flag.Parse()

	args := flag.Args()

	if len(args) > 1 {
		fmt.Println("Usage: ./main [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(*file)
	} else {
		runPrompt()
	}
}

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		dat, err := reader.ReadBytes('\n')
		check(err)
		run(string(dat))
	}
}

func run(src string) {
	fmt.Println(src)
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: ./golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

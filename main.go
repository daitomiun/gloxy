package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/daitonium/gloxy/ast"
)

func main() {
	args := os.Args[1:]
	fmt.Printf("args -> %v\n", len(args))
	if len(args) > 1 {
		fmt.Println("Usage: gloxy [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runFile(path string) {
	fmt.Println("Run file")
	fmt.Printf("path: %s \n", path)
	bytes, err := os.ReadFile(path)
	check(err)
	fmt.Println(string(bytes))
	run(string(bytes))

	// Indicate an error in the exit code
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	fmt.Println("Run prompt")
	reader := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("> ")
		reader.Scan()
		line := reader.Text()
		if len(line) == 0 {
			continue
		}
		run(line)
		hadError = false

	}
}

func run(source string) {
	scanner := scanner{
		source:  source,
		tokens:  []ast.Token{},
		start:   0,
		current: 0,
		line:    1,
	}

	tokens := scanner.scanTokens()

	for i, token := range tokens {
		fmt.Printf("%d -> %v \n", i, token)
	}

}

var hadError bool

func codeError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[line %v ] Error %s: %s \n", line, where, message)
	hadError = true
}

package main

import (
	"bufio"
	"fmt"
	"lox/interpret"
	"lox/loxError"
	"lox/parse"
	"lox/scan"
	"lox/token"
	"os"
)

func main() {
	if len(os.Args) - 1 > 1 {
		fmt.Println("Usage: jlox [script]");
		os.Exit(64);
	} else if len(os.Args) - 1 == 1 {
		runFile(os.Args[1])
	} else {
		runPrompt();
	}
}

func runFile(path string) {
	code, err := os.ReadFile(path)
	if err == nil {
		run(string(code))
	}

	if (loxError.HadError) {os.Exit(65)}
	if (loxError.HadRuntimeError) {os.Exit(70)}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("> ")
		line, _, _ := reader.ReadLine()
		if (line == nil) {
			break;
		}

		run(string(line))
		loxError.HadError = false
	}
}

type Scanner struct {
	source string
	scanTokens func() []token.Token
}

func run(source string) {
	scanner := scan.NewScanner(source)
	tokens := scan.ScanTokens(scanner)

	statements := parse.Parse(tokens)
	if loxError.HadError {return}

	interpret.Resolve(interpret.GlobalEnv, statements)
	if loxError.HadError {return}

	interpret.Interpret(statements)
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-basic <filename.bas>")
		return
	}

	filename := os.Args[1]
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	lexer := NewLexer(string(code))
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.errors) > 0 {
		for _, msg := range parser.errors {
			fmt.Printf("Parser error: %s\n", msg)
		}
		return
	}

	evaluator := NewEvaluator(program)
	evaluator.Run()
}

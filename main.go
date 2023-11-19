package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ChukwuEmekaAjah/glox/interpreter"
)

var hadError bool

func main() {

	args := os.Args
	if len(args) > 2 {
		println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}

}

func runFile(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		println("There was a problem opening file. Reason is:", err.Error())
		return err
	}

	sourceCode, err := io.ReadAll(file)

	run(string(sourceCode))

	return err
}

func runPrompt() error {

	for {
		print("> ")
		expression, err := bufio.NewReader(os.Stdin).ReadString('\n')
		expression = expression[0 : len(expression)-1] // remove newline
		if err != nil {
			println("There was an error reading expression input", err.Error())
			return err
		}

		if expression == "" {
			return nil
		}
		run(string(expression))
		print("\n")
	}
}

func run(sourceCode string) {
	scanner := interpreter.NewScanner(sourceCode)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Printf("token is %v \n", token)
	}
	print("running code now")
}

func reportError(line uint, message string) {
	report(line, "", message)
}

func report(line uint, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	hadError = true
}

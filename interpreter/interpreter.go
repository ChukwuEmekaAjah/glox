package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// HadError specifies if an error was encountered
var HadError bool

// HadRuntimeError specifies errors generated at interpretation time
var HadRuntimeError bool

// Interpreter is an instance of the interpreter
type Interpreter struct {
	expression Expr
}

// Start starts the interperter to in REPL mode or file execution mode
func (interpreter *Interpreter) Start() {
	args := os.Args
	if len(args) > 2 {
		println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		err := interpreter.runFile(args[1])
		if err != nil {
			panic(err)
		}
	} else {
		err := interpreter.runPrompt()
		if err != nil {
			panic(err)
		}
	}
}

// Interpret evaluates the code
func (interpreter *Interpreter) Interpret() {
	value := evaluate(interpreter.expression)

	if len(Errors) > 0 {
		for _, err := range Errors {
			fmt.Println(err)
		}

		return
	}

	fmt.Print(value)
}

func (interpreter *Interpreter) runFile(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		println("There was a problem opening file. Reason is:", err.Error())
		return err
	}

	sourceCode, err := io.ReadAll(file)

	interpreter.run(string(sourceCode))

	if HadError {
		os.Exit(65)
	}

	if HadRuntimeError {
		os.Exit(70)
	}

	return err
}

func (interpreter *Interpreter) runPrompt() error {

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
		interpreter.run(string(expression))
		print("\n")
		HadError = false
	}
}

func (interpreter *Interpreter) run(sourceCode string) {
	scanner := NewScanner(sourceCode)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	expression := parser.Parse()
	println(expression.print())
	interpreter.expression = expression
	interpreter.Interpret()
	if HadError {
		return
	}
}

// ReportError tells us what line an error occurred and what the error message is
func ReportError(line uint, message string) {
	Report(line, "", message)
}

// Report tells us where an error occurred
func Report(line uint, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	HadError = true
}

// Error reports error found in the parser
func Error(token Token, message string) {
	if token.tokenType == "EOF" {
		Report(token.line, "at end", message)
	} else {
		Report(token.line, fmt.Sprintf("at '%v'", token.lexeme), message)
	}
}

// RuntimeError reports interpreter runtime error
func RuntimeError(token Token, message string) {
	fmt.Printf("%s\n[line %d]", message, token.line)
	HadRuntimeError = true
}

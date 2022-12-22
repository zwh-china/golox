package VM

import (
	"bufio"
	"fmt"
	"golox/lox/interpreter"
	"golox/lox/lexer"
	"golox/lox/parser"
	"golox/utils"
	"os"
)

type VM struct {
	hadError        bool
	hadRuntimeError bool
	vmLexer         *lexer.Lexer
	vmParser        *parser.Parser
	vmInterpreter   *interpreter.Interpreter
}

func (v *VM) RunFile(path string) {
	fileBytes, _ := os.ReadFile(path)
	v.run(string(fileBytes[:]))
	// Indicate an error in the exit code.
	if v.hadError {
		os.Exit(65)
	}
	if v.hadRuntimeError {
		os.Exit(70)
	}
}

func (v *VM) RunStr(code string) {
	v.run(code)
	// Indicate an error in the exit code.
	if v.hadError {
		os.Exit(65)
	}
}

func (v *VM) RunPrompt() {
	var line string
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("> ")
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			break
		} else {
			line = string(lineBytes[:])
			if line == "" {
				break
			}
		}
		v.run(line)
		v.hadError = false
	}
}

func (v *VM) run(source string) {
	v.vmLexer = lexer.NewLexer(source)
	tokens, lexerError := v.vmLexer.ScanTokens()
	if lexerError.HasError {
		utils.RaiseError(lexerError.Line, lexerError.Reason)
		v.hadError = true
	}

	v.vmParser = parser.NewParser(tokens)
	expression, parseError := v.vmParser.Parse()

	if parseError.HasError {
		v.hadError = true
	}

	if v.hadError {
		return
	}

	runtimeError := v.vmInterpreter.Interpret(expression)

	if runtimeError.HasError {
		v.hadRuntimeError = true
	}

}

func (v *VM) SetError(error bool) {
	v.hadError = error
}

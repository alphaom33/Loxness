package loxError

import (
	"fmt"
	"lox/token"
	"os"
)

var HadError = false
var HadRuntimeError = false

func Error(line int, message string) {
	Report(line, "", message)
}

func TokenError(tokeny token.Token, message string) {
	if tokeny.TokenType == token.EOF {
		Report(tokeny.Line, " at end", message)
	} else {
        Report(tokeny.Line, " at '" + tokeny.Lexeme + "'", message)
    }
}

func Report(line int, where string, message string) {
	os.Stderr.WriteString(fmt.Sprintf("[line %d] Error %s: %s\n", line, where, message))
	HadError = true
}

func ThrowRuntimeError(error RuntimeError) {
	os.Stderr.WriteString(fmt.Sprintf("%s, \n[line %d]", error.Error(), error.Token.Line))
	HadRuntimeError = true
}
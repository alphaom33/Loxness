package parse

import (
	"fmt"
	"lox/token"
)

type ParseError struct {
	Token token.Token
	Message string
}

func (e ParseError) Error() string {
  return fmt.Sprintf("%+v, %s", e.Token, e.Message)
}
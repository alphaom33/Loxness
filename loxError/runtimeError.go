package loxError

import (
	"fmt"
	"lox/token"
)

type RuntimeError struct {
  Token token.Token
  Message string
}

func (e RuntimeError) Error() string {
  return fmt.Sprintf("%+v, %s", e.Token, e.Message)
}
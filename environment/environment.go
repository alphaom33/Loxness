package environment

import (
	"lox/loxError"
	"lox/token"
)

type Environment struct {
  values map[string]any
}

func Define(e Environment, name string, value any) {
  e.values[name] = value 
}

func Get(e Environment, name token.Token) (any, error) {
 val, ok := e.values[name.Lexeme] 
  if ok {
    return val, nil
  }

  return nil, loxError.RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}
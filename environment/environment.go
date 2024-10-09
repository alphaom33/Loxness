package environment

import (
	"fmt"
	"lox/loxError"
	"lox/token"
)

type Environment struct {
  values map[string]any
}

func MakeEnvironment() Environment {
  return Environment{make(map[string]any)}
}

func Define(e *Environment, name string, value any) *Environment {
  e.values[name] = value 
  return e
}

func Get(e Environment, name token.Token) (any, error) {
 val, ok := e.values[name.Lexeme] 
  if ok {
    return val, nil
  }

  return nil, loxError.RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}

func Assign(e *Environment, name token.Token, value any) error {
  _, ok := e.values[name.Lexeme]
  if ok {
    e.values[name.Lexeme] = value
    return nil
  }

   return loxError.RuntimeError{name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)} 
}
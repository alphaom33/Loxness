package environment

import (
	"fmt"
	"lox/loxError"
	"lox/token"
)

type Environment struct {
  enclosing *Environment
  values map[string]any
  name string
}

func MakeEnvironment(parent *Environment, n string) Environment {
  a := Environment{parent, make(map[string]any), n}
  return a
}

func Define(e *Environment, name string, value any) error {
  env := e
  for env != nil {
    _, ok := env.values[name]
    if ok {
      return loxError.RuntimeError{token.Token{}, fmt.Sprintf("Variable '%s' is already defined in this scope", env.name)}
    }
    
    env = env.enclosing
  }

  e.values[name] = value 
  return nil
}

func Get(e *Environment, name token.Token) (any, error) {
 val, ok := e.values[name.Lexeme] 
  if ok {
    return val, nil
  }
  
  if e.enclosing != nil {
    return Get(e.enclosing, name)
  }

  return nil, loxError.RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}

func Assign(e *Environment, name token.Token, value any) error {
  _, ok := e.values[name.Lexeme]
  if ok {
    e.values[name.Lexeme] = value
    return nil
  }

  if e.enclosing != nil {
    return Assign(e.enclosing, name, value)
  }

   return loxError.RuntimeError{name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)} 
}
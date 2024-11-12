package environment

import (
	"fmt"
	"lox/loxError"
	"lox/token"
)

type Environment struct {
  Enclosing *Environment
  values map[string]any
  name string
}

func MakeEnvironment(parent *Environment, n string) Environment {
  a := Environment{parent, make(map[string]any), n}
  return a
}

func Define(e *Environment, name string, value any) error {
  e.values[name] = value
  return nil
}

func GetAt(e *Environment, distance int, name string) any {
  val, _ := ancestor(e, distance).values[name]
  return val
}

func AssignAt(e *Environment, distance int, name token.Token, value any) {
  ancestor(e, distance).values[name.Lexeme] = value
}

func ancestor(e *Environment, distance int) *Environment {
  for i := 0; i < distance; i++ {
    e = e.Enclosing
  }
  return e
}

func Get(e *Environment, name token.Token) (any, error) {
 val, ok := e.values[name.Lexeme] 
  if ok {
    return val, nil
  }
  
  if e.Enclosing != nil {
    return Get(e.Enclosing, name)
  }

  return nil, loxError.RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."}
}

func Assign(e *Environment, name token.Token, value any) error {
  _, ok := e.values[name.Lexeme]
  if ok {
    e.values[name.Lexeme] = value
    return nil
  }

  if e.Enclosing != nil {
    return Assign(e.Enclosing, name, value)
  }

   return loxError.RuntimeError{name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)} 
}

func (e Environment) GetName() string {
  return e.name
}
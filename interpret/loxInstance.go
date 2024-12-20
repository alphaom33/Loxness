package interpret

import (
	"lox/loxError"
	"lox/token"
	"reflect"
)

type LoxInstance struct {
  Class *LoxClass
  Fields map[string]any
}

func (e LoxInstance) String() string {
  return e.Class.Name.Lexeme + " Instance"
}

func (e LoxInstance) Get(name token.Token) (any, error) {
  if e.Class != nil {
    val, ok := e.Class.Getters[name.Lexeme]
    if ok {
      return val.Bind(e).Call(GlobalEnv, nil)
    }
  }
  
  val, ok := e.Fields[name.Lexeme]
  if ok {
    return val, nil
  }

  if e.Class != nil {
    method, _ := e.Class.FindMethod(name.Lexeme)
    if !reflect.DeepEqual(method, LoxFunction{}) {
      return method.Bind(e), nil
    }
  }

  return nil, loxError.RuntimeError{name, "Undefined Property '" + name.Lexeme + "'."}
}

func (e LoxInstance) Set(name token.Token, value any) {
  e.Fields[name.Lexeme] = value
}
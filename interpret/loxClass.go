package interpret

import (
	"lox/environment"
	"lox/token"
)

type LoxClass struct {
  LoxInstance
  Name token.Token
  Superclasses []LoxClass
  Methods map[string]LoxFunction
  Getters map[string] LoxFunction
}

func (e LoxClass) Call(env environment.Environment, arguments []any) (any, error) {
  instance := LoxInstance{&e, make(map[string]any)}
  initializer, err := e.FindMethod("init")
  if err == nil {
    initializer.Bind(instance).Call(env, arguments)
  }
  return instance, nil
}

func (e LoxClass) Arity() int {
  initializer, err := e.FindMethod("init")
  if err != nil {return 0}
  return initializer.Arity()
}

type MethodNotFoundError struct {
  Name string
}
func (e MethodNotFoundError) Error() string {
  return "Method" + e.Name + "not found"
}

func (e LoxClass) FindMethod(name string) (LoxFunction, error) {
  val, ok := e.Methods[name]
  if ok {
    return val, nil
  }

  for _, superclass := range e.Superclasses {
    val, err := superclass.FindMethod(name)
    if err == nil {
      return val, nil
    }
  }

  return LoxFunction{}, MethodNotFoundError{name}
}
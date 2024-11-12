package interpret

import (
	"fmt"
	"lox/environment"
)

type LoxFunction struct {
  Declaration Function
  Closure environment.Environment
  IsInitializer bool
}

func (e LoxFunction) Bind(instance LoxInstance) LoxFunction {
  env := environment.MakeEnvironment(&e.Closure, "")
  environment.Define(&env, "this", instance)
  return LoxFunction{e.Declaration, env, e.IsInitializer}
}

func (e LoxFunction) Call(_ environment.Environment, arguments []any) (any, error) {
  
  envy := environment.MakeEnvironment(&e.Closure, "func")
  for i := 0; i < len(e.Declaration.Params); i++ {
    environment.Define(&envy, e.Declaration.Params[i].Lexeme, arguments[i])
  }
  err := executeBlock(e.Declaration.Body, envy)
  rE, ok := err.(ReturnError)
  if ok {
    if e.IsInitializer {return environment.GetAt(&e.Closure, 0, "this"), nil}
    return rE.Value, nil
  }

  if e.IsInitializer {return environment.GetAt(&e.Closure, 0, "this"), nil}
  return nil, err
}

func (e LoxFunction) Arity() int {
  return len(e.Declaration.Params)
}

func (e LoxFunction) String() string {
 return fmt.Sprintf("<fn %s>", e.Declaration.Name.Lexeme) 
}
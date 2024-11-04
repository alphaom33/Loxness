package interpret

import (
	"lox/environment"
)

type LoxLambda struct {
  Declaration LLambda
  Closure environment.Environment
}

func (e LoxLambda) Call(env environment.Environment, arguments []any) (any, error) {
  envy := environment.MakeEnvironment(&e.Closure, "func")
  for i := 0; i < len(e.Declaration.Params); i++ {
    environment.Define(&envy, e.Declaration.Params[i].Lexeme, arguments[i])
  }

  err := executeBlock(e.Declaration.Body, envy)
  rE, ok := err.(ReturnError)
  if ok {
    return rE.Value, nil
  }
  return nil, err
}

func (e LoxLambda) Arity() int {
  return len(e.Declaration.Params)
}

func (e LoxLambda) String() string {
 return "lambda";
}
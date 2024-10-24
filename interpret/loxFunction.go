package interpret

import (
	"fmt"
	"lox/environment"
)

type LoxFunction struct {
  Declaration Function
}

func (e LoxFunction) Call(env environment.Environment, arguments []any) (any, error) {
  envy := environment.MakeEnvironment(&env, "func")
  for i := 0; i < len(e.Declaration.Params); i++ {
    environment.Define(&envy, e.Declaration.Params[i].Lexeme, arguments[i])
  }

  err := executeBlock(e.Declaration.Body, envy)
  return nil, err
}

func (e LoxFunction) Arity() int {
  return len(e.Declaration.Params)
}

func (e LoxFunction) String() string {
 return fmt.Sprintf("<fn %v>", e.Declaration.Name.Lexeme) 
}
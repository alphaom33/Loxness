package interpret

import "lox/environment"

type LoxCallable interface {
  Call(env environment.Environment, arguments []any) (any, error)
  Arity() int
}
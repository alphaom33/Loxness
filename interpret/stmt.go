package interpret

import (
	"lox/token"
  "lox/environment"
)

type Stmt interface {
  VisitStmt(environment.Environment) error
}

type Block struct {
  Statements []Stmt
}

type Expression struct {
  Expression Expr
}

type Print struct {
  Expression Expr
}

type Var struct {
  Name token.Token
  Initializer Expr
}

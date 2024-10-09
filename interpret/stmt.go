package interpret

import (
	"lox/token"
)

type Stmt interface {
  VisitStmt() error
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

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

type Class struct {
  Name token.Token
  superclass Variable
  Methods []Function
  StaticMethods []Function
  Getters []Function
}

type Expression struct {
  Expression Expr
}

type Function struct {
  Name token.Token
  Params []token.Token
  Body []Stmt
}

type If struct {
  Condition Expr
  ThenBranch Stmt
  ElseBranch Stmt
}

type Print struct {
  Expression Expr
}

type Return struct {
  Keyword token.Token
  Value Expr
}

type Var struct {
  Name token.Token
  Initializer Expr
}

type While struct {
  Condition Expr
  Body Stmt
}

type Break struct {
}

package  interpret

import (
  "lox/token"
  "lox/environment"
)

type Expr interface {
  AstPrint() string
  VisitExpr(environment.Environment) (any, error)
}

type Ternary struct {
  Condition Expr
  OnTrue Expr
  OnFalse Expr
}

type Binary struct {
  Left Expr
  Operator token.Token
  Right Expr
}

type Call struct {
  Callee Expr
  Paren token.Token
  Arguments []Expr
}

type Grouping struct {
  Expression Expr
}

type Literal struct {
  Value any
}

type Logical struct {
  Left Expr
  Operator token.Token
  Right Expr
}

type Unary struct {
  Operator token.Token
  Right Expr
}

type Variable struct {
  Name token.Token
}

type Assign struct {
  Name token.Token
  Value Expr
}
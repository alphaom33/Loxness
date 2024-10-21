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

type Grouping struct {
  Expression Expr
}

type Literal struct {
  Value any
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

type Undefined struct {
  
}
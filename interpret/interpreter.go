package interpret

import (
	"fmt"
	"lox/environment"
	"lox/loxError"
)

var env environment.Environment

func Interpret(statements []Stmt) {
  fmt.Println(statements)
  for _, statement := range statements {
    err := execute(statement)
    if err != nil {
      rE, _ := err.(loxError.RuntimeError)
      loxError.ThrowRuntimeError(rE)
      break
    }
  }
}

func (e Expression) VisitStmt() error {
  e.Expression.VisitExpr()
  return nil
}

func (e Print) VisitStmt() error {
  val, _ := e.Expression.VisitExpr()
  fmt.Println(Stringify(val))

  return nil
}

func (e Var) VisitStmt() error {
  var value any
  if e.Initializer != nil {
    var err error
    value, err = e.Initializer.VisitExpr()
    if err != nil {return err}
  }

  environment.Define(env, e.Name.Lexeme, value)
  return nil
}


func execute(stmt Stmt) error {
  fmt.Print(stmt)
  return stmt.VisitStmt()
}

func Stringify(object any) string {
  if object == nil {return "nil"}

  oD, ok := object.(float32)

  if ok {
    text := fmt.Sprintf("%f", oD)
    if text[len(text) - 1] == '0' {
      text = text[0:len(text) - 2]
    }
    return text
  }

  return fmt.Sprintf("%+v", object)
}
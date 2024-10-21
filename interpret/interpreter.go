package interpret

import (
	"fmt"
	"lox/environment"
	"lox/loxError"
)

var globalEnv environment.Environment = environment.MakeEnvironment(nil, "asdf")

func Interpret(statements []Stmt) {
  for _, statement := range statements {
    err := execute(statement, globalEnv)
    if err != nil {
      rE, _ := err.(loxError.RuntimeError)
      loxError.ThrowRuntimeError(rE)
      break
    }
  }
}

func (e Expression) VisitStmt(env environment.Environment) error {
  e.Expression.VisitExpr(env)
  return nil
}

func (e Print) VisitStmt(env environment.Environment) error {
  val, err := e.Expression.VisitExpr(env)
  if err != nil {return err}
  fmt.Println(Stringify(val))

  return nil
}

func (e Var) VisitStmt(env environment.Environment) error {
  var value any = Undefined{}
  if e.Initializer != nil {
    var err error
    value, err = e.Initializer.VisitExpr(env)
    if err != nil {return err}
  }

  return environment.Define(&env, e.Name.Lexeme, value)
}

func (e While) VisitStmt(env environment.Environment) error {
  val, err := evaluate(e.Condition, env)
  if err != nil {return err}
  for isTruthy(val) {
    err = execute(e.Body, env)
    if err != nil {return err}
    
    val, err = evaluate(e.Condition, env)
    if err != nil {return err}
  }

  return nil
}

func (e Block) VisitStmt(env environment.Environment) error {
  executeBlock(e.Statements, environment.MakeEnvironment(&env, ""))
  return nil
}

func (e If) VisitStmt(env environment.Environment) error {
  b, err := evaluate(e.Condition, env)
  if err != nil {return err}
  if isTruthy(b) {
    return execute(e.ThenBranch, env)
  } else if e.ElseBranch != nil {
    return execute(e.ThenBranch, env)
  }

  return nil
}

func execute(stmt Stmt, env environment.Environment) error {
  return stmt.VisitStmt(env)
}

func executeBlock(statements []Stmt, newEnv environment.Environment) error {
  for _, statement := range statements {
    err := execute(statement, newEnv)
    if err != nil {return err}
  }
  return nil
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
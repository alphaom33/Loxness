package interpret

import (
	"errors"
	"fmt"
	"lox/environment"
	"lox/loxError"
	"time"
)

var GlobalEnv environment.Environment = environment.MakeEnvironment(nil, "asdf")
var locals map[Expr]int = make(map[Expr]int)

type ProtoLoxCallable struct {
  callMethod func(env environment.Environment, arguments []any) (any, error)
  arityMethod func() int
  stringMethod func() string
}

func (p ProtoLoxCallable) Call(env environment.Environment, arguments []any) (any, error) {
  return p.callMethod(env, arguments)
}

func (p ProtoLoxCallable) Arity() int {
  return p.arityMethod()
}

func (p ProtoLoxCallable) String() string {
  return p.stringMethod()
}

func Interpret(statements []Stmt) {
  environment.Define(&GlobalEnv, "clock", ProtoLoxCallable{
    arityMethod: func() int {
      return 0
    },
    
    callMethod: func(env environment.Environment, arguments []any) (any, error) {
      return time.Now().UnixMilli() / 1000, nil
    },
    stringMethod: func() string {
      return "<native fn>"
    },
  })
  
  for _, statement := range statements {
    err := execute(statement, GlobalEnv)
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

func (e Function) VisitStmt(env environment.Environment) error {
  function := LoxFunction{e, env, false}
  return environment.Define(&env, e.Name.Lexeme, function)
}

func (e Print) VisitStmt(env environment.Environment) error {
  val, err := e.Expression.VisitExpr(env)
  if err != nil {return err}
  fmt.Println(Stringify(val))

  return nil
}

type ReturnError struct {
  Value any
}
func (e ReturnError) Error() string {
  return "Return statement can only be used inside a function"
}
func (e Return) VisitStmt(env environment.Environment) error {
  var value any = nil
  if e.Value != nil {
    var err error
    value, err = evaluate(e.Value, env)
    if err != nil {return err}
  }
  
  return ReturnError{value}
}

func (e Var) VisitStmt(env environment.Environment) error {
  var value any
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
    if err != nil {
      if err.Error() == "break" {
        break
      }
      return err
    }
    
    val, err = evaluate(e.Condition, env)
    if err != nil {return err}
  }

  return nil
}

func (e Break) VisitStmt(env environment.Environment) error {
  return errors.New("break")
}

func (e Block) VisitStmt(env environment.Environment) error {
  return executeBlock(e.Statements, environment.MakeEnvironment(&env, ""))
}

func (e Class) VisitStmt(env environment.Environment) error {
  var superclass any
  if e.Superclass != nil {
    var err error
    superclass, err = e.Superclass.VisitExpr(env)
    if err != nil {return err}

    var ok bool
    superclass, ok = superclass.(LoxClass)
    if !ok {
      loxError.TokenError(e.Superclass.Name, "Superclass must be a class.")
    }
  }
  
  environment.Define(&env, e.Name.Lexeme, nil)

  envy := env
  if e.Superclass != nil {
    env = environment.MakeEnvironment(&env, "a")
    environment.Define(&env, "super", superclass)
  }

  methods := make(map[string]LoxFunction)
  for _, method := range e.Methods {
    methods[method.Name.Lexeme] = LoxFunction{method, env, method.Name.Lexeme == "init"}
  }

  staticMethods := make(map[string]any)
  for _, method := range e.StaticMethods {
    staticMethods[method.Name.Lexeme] = LoxFunction{method, env, false}
  }

  getters := make(map[string]LoxFunction)
  for _, getter := range e.Getters {
    getters[getter.Name.Lexeme] = LoxFunction{getter, env, false}
  }

  super, _ := superclass.(LoxClass)
  class := LoxClass{LoxInstance{nil, staticMethods}, e.Name, &super, methods, getters}

  if superclass != nil {
    env = envy
  }
  environment.Assign(&env, e.Name, class)
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

func InterpretResolve(expr Expr, depth int) {
  locals[expr] = depth
}
package interpret

import (
	"fmt"
	classtype "lox/classType"
	"lox/environment"
	"lox/functionType"
	"lox/loxError"
	"lox/token"
	"lox/varUsage"
)

type stack []map[string]varusage.VarUsage
var currentFunction functiontype.FunctionType
var currentClass classtype.ClassType = classtype.NONE

func (s stack) Ack(name string, value varusage.VarUsage) {
  s[len(s)-1][name] = value
}

func (s stack) Peek() map[string]varusage.VarUsage {
  return s[len(s)-1]
}

func (s stack) Push(v map[string]varusage.VarUsage) stack {
    return append(s, v)
}

func (s stack) Pop() (stack, map[string]varusage.VarUsage) {
    l := len(s)
    return  s[:l-1], s[l-1]
}

var scopes stack

type Scope interface {
  VisitScope(env environment.Environment)
}

func (e Block) VisitScope(env environment.Environment) {
  beginScope()
  Resolve(env, e.Statements)
  endScope()
}

func (e Class) VisitScope(env environment.Environment) {
  enclosingClass := currentClass
  currentClass = classtype.CLASS
  
  declare(e.Name)
  define(e.Name)

  beginScope()
  scopes.Ack("this", varusage.INITIALIZED)

  for _, method := range e.Methods {
    var declaration functiontype.FunctionType
    
    if method.Name.Lexeme == "init" {
      declaration = functiontype.INITIALIZER
    } else {
      declaration = functiontype.METHOD
    }
    
    resolveFunction(env, method, declaration)  
  }

  endScope()

  currentClass = enclosingClass
}

func (e Var) VisitScope(env environment.Environment) {
  declare(e.Name)
  if e.Initializer != nil {
    resolveExpr(env, e.Initializer)
  }
  define(e.Name)
}

func (e Variable) VisitScope(env environment.Environment) {
  val, ok := scopes.Peek()[e.Name.Lexeme]
  if len(scopes) != 0 && ok && val == varusage.DECLARED {
    loxError.ThrowRuntimeError(loxError.RuntimeError{e.Name, "Can't read local variable in its own initializer."})
  }

  resolveLocal(e, e.Name)
}

func (e Assign) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Value)
  resolveLocal(e, e.Name)
}

func (e Function) VisitScope(env environment.Environment) {
  declare(e.Name)
  define(e.Name)

  resolveFunction(env, e, functiontype.FUNCTION)
}

func (e Expression) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Expression)
}

func (e If) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Condition)
  resolveStmt(env, e.ThenBranch)
  if e.ElseBranch != nil {resolveStmt(env, e.ThenBranch)}
}

func (e Print) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Expression)
}

func (e Return) VisitScope(env environment.Environment) {
  if currentFunction == functiontype.NONE {
    loxError.TokenError(e.Keyword, "Can't return from top-level code")
  }
  
  if e.Value != nil {
    if currentFunction == functiontype.INITIALIZER {
      loxError.TokenError(e.Keyword, "Can't return from an initializer.")
    }
    resolveExpr(env, e.Value)
  }
}

func (e While) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Condition)
  resolveStmt(env, e.Body)
}

func (e Binary) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Left)
  resolveExpr(env, e.Right)
}

func (e Call) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Callee)

  for _, argument := range e.Arguments {
    resolveExpr(env, argument)
  }
}

func (e Get) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Object)
}

func (e Grouping) VisitScope(env environment.Environment) {
}

func (e Literal) VisitScope(env environment.Environment) {
}

func (e Logical) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Right)
}

func (e Set) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Value)
  resolveExpr(env, e.Object)
}

func (e This) VisitScope(env environment.Environment) {
  if currentClass == classtype.NONE {
    loxError.TokenError(e.Keyword, "Can't use 'this' outside of a class")
  }
  
  resolveLocal(e, e.Keyword)
}

func (e Unary) VisitScope(env environment.Environment) {
  resolveExpr(env, e.Right)
}

func InitialResolve(env environment.Environment, statements []Stmt) {
  beginScope()
  Resolve(env, statements)
  endScope()
}

func Resolve(env environment.Environment, statements []Stmt) {
  for _, statement := range statements {
    resolveStmt(env, statement)
  }
}

func resolveStmt(env environment.Environment, statement Stmt) {
  s, ok := statement.(Scope)
  if ok {
    s.VisitScope(env) 
  }
}

func resolveExpr(env environment.Environment, expr Expr) {
  s, _ := expr.(Scope)
  s.VisitScope(env) 
}

func resolveFunction(env environment.Environment, function Function, typey functiontype.FunctionType) {
  enclosingFunction := currentFunction
  currentFunction = typey
    
  beginScope()
  for _, param := range function.Params {
    declare(param)
    define(param)
  }
  Resolve(env, function.Body)
  endScope()

  currentFunction = enclosingFunction
}

func beginScope() {
  scopes = scopes.Push(make(map[string]varusage.VarUsage))
}

func endScope() {
  var scope map[string]varusage.VarUsage
  scopes, scope = scopes.Pop()

  for k, v := range scope {
    if v != varusage.USED && k != "this" {
      fmt.Println("Warning: " + k + " is never used")
    }
  }
}

func declare(name token.Token) {
 if len(scopes) == 0 {return} 
  _, ok := scopes.Peek()[name.Lexeme]
  if ok {
    loxError.TokenError(name, "Already a variable with this name in this scope.")
  }
  scopes.Ack(name.Lexeme, varusage.DECLARED)
}

func define(name token.Token) {
  if len(scopes) == 0 {return}
  scopes.Ack(name.Lexeme, varusage.INITIALIZED)
}

func resolveLocal(expr Expr, name token.Token) {
  for i := len(scopes) - 1; i >= 0; i-- {
    _, ok := scopes[i][name.Lexeme]
    if ok {
      scopes[i][name.Lexeme] = varusage.USED
      InterpretResolve(expr, len(scopes) - 1 - i)
      return
    }
  }
}
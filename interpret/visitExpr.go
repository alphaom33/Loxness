package interpret

import (
	"fmt"
	"lox/environment"
	"lox/loxError"
	"lox/token"
)

func (e Literal) VisitExpr(_ environment.Environment) (any, error) {
  return e.Value, nil
}

func (e Grouping) VisitExpr(env environment.Environment) (any, error) {
  return evaluate(e.Expression, env)
}

func (e Unary) VisitExpr(env environment.Environment) (any, error) {
  right, err := evaluate(e.Right, env)
  if err != nil {return nil, err}

  switch e.Operator.TokenType {
  case token.BANG:
    b, _ := right.(bool)
    return !isTruthy(b), nil
  case token.MINUS:
    f, err := checkNumberOperand(e.Operator, right)
    return -f, err
  }

  return nil, nil
}

func (e Ternary) VisitExpr(_ environment.Environment) (any, error) {
  return nil, nil
}

func (e Binary) VisitExpr(env environment.Environment) (any, error) {
  left, err := evaluate(e.Left, env)
  if err != nil {return nil, err}
  right, err := evaluate(e.Right, env)
  if err != nil {return nil, err}

  switch e.Operator.TokenType {
  case token.GREATER:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    return fL > fR, nil
  case token.GREATER_EQUAL:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    return fL >= fR, nil
  case token.LESS:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    return fL < fR, nil
  case token.LESS_EQUAL:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    return fL <= fR, nil

  case token.BANG_EQUAL:
    return !isEqual(left, right), nil
  case token.EQUAL_EQUAL:
    return isEqual(left, right), nil

    
  case token.MINUS:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    return fL - fR, nil
  case token.SLASH:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    if fR == 0 {return nil, loxError.RuntimeError{e.Operator, "Division by zero."}}
    return fL / fR, nil
  case token.STAR:
    fL, fR, err := checkNumberOperands(e.Operator, left, right)
    if err != nil {return nil, err}
    return fL * fR, nil
  case token.PLUS:
    fL, okL := left.(float64)
    fR, okR := right.(float64)
    if okL && okR {return fL + fR, nil}

    sL, okL := left.(string)
    sR, okR := right.(string)
    if okL || okR {
        if okL != okR {
          if okL {
            return fmt.Sprintf("%s%s", sL, Stringify(fR)), nil
          } else {
            return fmt.Sprintf("%s%s", Stringify(fL), sR), nil
          }
        }
      return sL + sR, nil
    }

    return nil, loxError.RuntimeError{e.Operator, "Operands must be two numbers or two strings"}
  }

  return nil, nil
}

func (e Variable) VisitExpr(env environment.Environment) (any, error) {
  return environment.Get(&env, e.Name)
}

func (e Assign) VisitExpr(env environment.Environment) (any, error) {
  value, err := evaluate(e.Value, env)
  if err != nil {return nil, err}
  environment.Assign(&env, e.Name, value)
  return value, nil
}

func evaluate(expression Expr, env environment.Environment) (any, error) {
  return expression.VisitExpr(env)
}

func isTruthy(thing any) bool {
  if thing == nil {return false}
  b, ok := thing.(bool)
  if ok {return b}

  return true
}

func isEqual(a any, b any) bool {
  if a == nil && b == nil {return true}
  if a == nil {return false}

  return a == b
}

func checkNumberOperand(operator token.Token, right any) (float64, error) {
  f, ok := right.(float64)
  if ok {
    return f, nil
  } else {
    return 0, loxError.RuntimeError{operator, "Operand must be a number."}
  }
}

func checkNumberOperands(operator token.Token, left any, right any) (float64, float64, error) {
  fL, okL := left.(float64)
  fR, okR := right.(float64)
  if okL && okR {return fL, fR, nil}
  return 0, 0, loxError.RuntimeError{operator, "Operands must be numbers."}
}
package interpret

import (
	"fmt"
	"strings"
)

func (e Ternary) AstPrint() string {
  return parenthesize("?", e.Condition, e.OnTrue, e.OnFalse)
}

func (e Binary) AstPrint() string {
 return parenthesize(e.Operator.Lexeme, e.Left, e.Right) 
}

func (e Call) AstPrint() string {
  return parenthesize(e.Paren.Lexeme, append([]Expr{e.Callee}, e.Arguments...)...)
}

func (e Get) AstPrint() string {
  return parenthesize(e.Name.Lexeme, e.Object)
}

func (e Logical) AstPrint() string {
  return parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (e Set) AstPrint() string {
  return parenthesize(e.Name.Lexeme, e.Object, e.Value)
}

func (e This) AstPrint() string {
  return parenthesize(e.Keyword.Lexeme)
}

func (e Grouping) AstPrint() string {
  return parenthesize("group", e.Expression)
}

func (e Literal) AstPrint() string {
  // if e == nil {return "nil"}
  return parenthesize(fmt.Sprintf("%+v", e.Value))
}

func (e Unary) AstPrint() string {
  return parenthesize(e.Operator.Lexeme, e.Right)
}

func (e Variable) AstPrint() string {
  return parenthesize(e.Name.Lexeme)
}

func (e Assign) AstPrint() string {
  return parenthesize(fmt.Sprintf("= %s", e.Name.Lexeme), e.Value)
}

func parenthesize(name string, exprs... Expr) string {
  builder := strings.Builder{}
  builder.WriteString("(")
  builder.WriteString(name)

  for _, expr := range exprs {
    builder.WriteString(" ")
    builder.WriteString(expr.AstPrint())
  }
  builder.WriteString(")")

  return builder.String()
}
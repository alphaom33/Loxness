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
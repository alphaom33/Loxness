package token

import (
  "fmt"
)

type Token struct {
  TokenType TokenType
  Lexeme string
  Literal any
  Line int
}

func (token Token) String() string {
   return fmt.Sprintf("type %s lexeme %s literal %s", token.TokenType, token.Lexeme, token.Literal)
}
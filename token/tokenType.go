package token

import (
  "fmt"
)

type TokenType int

const (
  LEFT_PAREN TokenType = iota
  RIGHT_PAREN
  LEFT_BRACE
  RIGHT_BRACE
  COMMA
  QUESTION
  COLON
  DOT
  MINUS
  PLUS
  SEMICOLON
  SLASH
  STAR

  BANG
  BANG_EQUAL
  EQUAL
  EQUAL_EQUAL
  GREATER
  GREATER_EQUAL
  LESS
  LESS_EQUAL

  IDENTIFIER
  STRING
  NUMBER

  AND
  CLASS
  ELSE
  FALSE
  FUN
  FOR
  IF
  NIL
  OR
  PRINT
  RETURN
  SUPER
  THIS
  TRUE
  VAR
  WHILE

  EOF
)

func (e TokenType) String() string {
  switch e {
  case LEFT_PAREN:
    return "LEFT_PAREN"
  case RIGHT_PAREN:
    return "RIGHT_PAREN"
  case LEFT_BRACE:
    return "LEFT_BRACE"
  case RIGHT_BRACE:
    return "RIGHT_BRACE"
  case COMMA:
    return "COMMA"
  case QUESTION:
    return "QUESTION"
  case COLON:
    return "COLON"
  case DOT:
    return "DOT"
  case PLUS:
    return "PLUS"
  case SLASH:
    return "SLASH"
  case STAR:
    return "STAR"
  case BANG:
    return "BANG"
  case BANG_EQUAL:
    return "BANG_EQUAL"
  case EQUAL:
    return "EQUAL"
  case EQUAL_EQUAL:
    return "EQUAL_EQUAL"
  case GREATER:
    return "GREATER"
  case GREATER_EQUAL:
    return "GREATER_EQUAL"
  case LESS:
    return "LESS"
  case LESS_EQUAL:
    return "LESS_EQUAL"
  case IDENTIFIER:
    return "IDENTIFIER"
  case STRING:
    return "STRING"
  case NUMBER:
    return "NUMBER"
  case AND:
    return "AND"
  case ELSE:
    return "ELSE"
  case FUN:
    return "FUN"
  case IF:
    return "IF"
  case OR:
    return "OR"
  case RETURN:
    return "RETURN"
  case THIS:
    return "THIS"
  case  VAR:
    return "VAR"
  case WHILE:
    return "WHILE"
  case EOF:
    return "EOF"
  default:
    return fmt.Sprintf("%d", int(e))
  }
}
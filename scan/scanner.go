package scan

import (
	"fmt"
	"strconv"
  "lox/loxError"
  . "lox/token"
)

type scanner struct {
  source string 
  tokens []Token
  start int
  current int
  line int
}

var keywords = map[string]TokenType{
  "and": AND,
  "class": CLASS,
  "else": ELSE,
  "false": FALSE,
  "for": FOR,
  "fun": FUN,
  "if": IF,
  "nil": NIL,
  "or": OR,
  "print": PRINT,
  "return": RETURN,
  "super": SUPER,
  "this": THIS,
  "true": TRUE,
  "var": VAR,
  "while": WHILE,
}

func NewScanner(source string) *scanner {
  return &scanner{
    source,
    nil,
    0,
    0,
    1,
  }
}

func ScanTokens(scanner *scanner) []Token {
  for !isAtEnd(scanner) {
    scanner.start = scanner.current
    scanToken(scanner)
  }

  scanner.tokens = append(scanner.tokens, Token{
    EOF,
    "",
    nil,
    scanner.line,
  })
  return scanner.tokens;
}

func scanToken(scanner *scanner) {
  c := advance(scanner)

  switch c {
  case '(': addToken(scanner, LEFT_PAREN, nil); break
  case ')': addToken(scanner, RIGHT_PAREN, nil); break
  case '{': addToken(scanner, LEFT_BRACE, nil); break
  case '}': addToken(scanner, RIGHT_BRACE, nil); break
  case ',': addToken(scanner, COMMA, nil); break
  case '.': addToken(scanner, DOT, nil); break
  case '-': addToken(scanner, MINUS, nil); break
  case '+': addToken(scanner, PLUS, nil); break
  case ';': addToken(scanner, SEMICOLON, nil); break
  case '*': addToken(scanner, STAR, nil); break
  case '?': addToken(scanner, QUESTION, nil); break
  case ':': addToken(scanner, COLON, nil); break
  case '!':
    addToken(scanner, ifThenElse(match(scanner, '='), BANG_EQUAL, BANG), nil)
    break
  case '=':
    addToken(scanner, ifThenElse(match(scanner, '='), EQUAL_EQUAL, EQUAL), nil)
    break
  case '>':
    addToken(scanner, ifThenElse(match(scanner, '='), GREATER_EQUAL, GREATER), nil)
    break
  case '<':
    addToken(scanner, ifThenElse(match(scanner, '='), LESS_EQUAL, LESS), nil)
    break
  case '/':
    if match(scanner, '/') {
      for peek(scanner) != '\n' && !isAtEnd(scanner) {advance(scanner)}
    } else if match(scanner, '*') {blockComment(scanner)} else {
      addToken(scanner, SLASH, nil)
    }
    break
  case ' ', '\r', '\t':
    break
  case '\n':
    scanner.line++
    break
  case '"': lexString(scanner); break
  default:
    if isDigit(c) {
      number(scanner)
    } else if isAlpha(c) {
      identifier(scanner)
    } else {
      loxError.Error(scanner.line, "unexpected character")
    }
    break
  }
}

func number(scanner *scanner) {
  for isDigit(peek(scanner)) {advance(scanner)}

  if peek(scanner) == '.' && isDigit(peekNext(scanner)) {
    advance(scanner)

    for (isDigit(peek(scanner))) {advance(scanner)}
  }
  number, _ := strconv.ParseFloat(scanner.source[scanner.start:scanner.current], 16)
  addToken(scanner, NUMBER, number)
}

func identifier(scanner *scanner) {
  for isAlphanumeric(peek(scanner)) {advance(scanner)}

  text := scanner.source[scanner.start:scanner.current]
  tokenType := keywords[text]
  if tokenType == 0 {tokenType = IDENTIFIER}

  addToken(scanner, tokenType, text)
}

func lexString(scanner *scanner) {
  for peek(scanner) != '"' && !isAtEnd(scanner) {
    if (peek(scanner) == '\n') {scanner.line++}
    advance(scanner)
  }

  if (isAtEnd(scanner)) {
    loxError.Error(scanner.line, "Unterminated string.")
    return
  }

  advance(scanner)

  value := scanner.source[scanner.start + 1: scanner.current - 1]
  addToken(scanner, STRING, value)
}

func blockComment(scanner *scanner) {
  numComments := 1
  for numComments > 0 && scanner.current < len(scanner.source) - 1 {
    if peek(scanner) == '/' && peekNext(scanner) == '*' {
      numComments++
      fmt.Println("asdf")
    } else if peek(scanner) == '*' && peekNext(scanner) == '/' {
      numComments--
    }
    advance(scanner)
  }
  advance(scanner)
}

func match(scanner *scanner, expected rune) bool {
  if isAtEnd(scanner) {return false}
  if rune(scanner.source[scanner.current]) != expected {return false}

  scanner.current++
  return true
}

func peek(scanner *scanner) rune {
  if isAtEnd(scanner) {return 0}

  return rune(scanner.source[scanner.current])
}

func peekNext(scanner *scanner) rune {
  if (scanner.current + 1 >= len(scanner.source)) {return 0}
  return rune(scanner.source[scanner.current + 1])
}

func isAlpha(c rune) bool {
  return (c >= 'a' && c <= 'z') ||
    (c >= 'A' && c <= 'Z') ||
    c == '_'
}

func isAlphanumeric(c rune) bool {
  return isAlpha(c) || isDigit(c)
}

func isDigit(c rune) bool {
  return c >= '0' && c <= '9'
}

func ifThenElse(condition bool, then TokenType, elsey TokenType) TokenType {
  out := elsey
  if condition {
    out = then
  }
  return out
}

func advance(scanner *scanner) rune {
  tmp := scanner.source[scanner.current]
  scanner.current++
  return rune(tmp)
}

func addToken(scanner *scanner, tokenType TokenType, literal any) {
  text := scanner.source[scanner.start:scanner.current]
  scanner.tokens = append(scanner.tokens, Token{
    tokenType,
    text,
    literal,
    scanner.line,
  })
}

func isAtEnd(scanner *scanner) bool {
  return scanner.current >= len(scanner.source)
}

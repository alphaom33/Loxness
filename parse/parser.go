package parse

import (
	"errors"
	"lox/loxError"
	"lox/token"
    . "lox/interpret"
)

var tokens []token.Token
var current int

func Parse(p_tokens []token.Token) []Stmt {
    tokens = p_tokens
    current = 0
    var statements []Stmt
    for !isAtEnd() {
        statements = append(statements, declaration())
    }

    return statements
}

func declaration() Stmt {
    var err error
    var out Stmt

    if match(token.VAR) {
        out, err = varDeclaration()
    } else {
        out = statement()
    }
    
    if err != nil {
        synchronize()     
        return nil
    }

    return out
}

func varDeclaration() (Stmt, error) {
    name, err := consume(token.IDENTIFIER, "Expect variable name.")
    if err != nil {return nil, err}

    var initializer Expr
    if match(token.EQUAL) {
        initializer, err = comma()
        if err != nil {return nil, err}
    }

    _, err = consume(token.SEMICOLON, "Expect ';' after variable declaration.")
    if err != nil {return nil, err}

    return Var{name, initializer}, nil
}



func statement() Stmt {
    if (match(token.PRINT)) {return printStatement()}
    return expressionStatement()
}

func printStatement() Stmt {
    value, _ := comma()
    consume(token.SEMICOLON, "Expect ';' after value.")
    return Print{value}
}

func expressionStatement() Stmt {
    expre, _ := comma()
    consume(token.SEMICOLON, "Expect ';' after expression.")
    return Expression{expre}
}

func comma() (Expr, error) {
    expression, err := ternary()
    if err != nil {return expression, err}

    for (match(token.COMMA)) {
        operator := previous()
        right, err := ternary()
        if err != nil {return right, err}
        expression = Binary{expression, operator, right}
    }

    return expression, nil
}


func expression() (Expr, error) {
    return ternary()
}

func ternary() (Expr, error) {
    expression, err := equality()    
    if err != nil {return expression, err}

    for (match(token.QUESTION)) {
        onTrue, err := ternary()
        if err != nil {return onTrue, err}

        match(token.COLON)

        onFalse, err := ternary()
        if err != nil {return onFalse, err}

        expression = Ternary{expression, onTrue, onFalse}
    }

    return expression, nil
}

func equality() (Expr, error) {
    expression, err := comparison()
    if err != nil {return expression, err}

    for match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
        operator := previous()
        right, err := comparison()
        if err != nil {return right, err}
        expression = Binary{expression, operator, right}
    }

    return expression, nil
}

func comparison() (Expr, error) {
    expression, err := term()
    if err != nil {return expression, err}

    for match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
        operator := previous()
        right, err := term()
        if err != nil {return right, err}
        expression = Binary{expression, operator, right}
    }

    return expression, nil
}

func term() (Expr, error) {
    expression, err := factor()
    if err != nil {return expression, err}

    for match(token.MINUS, token.PLUS) {
        operator := previous()
        right, err := factor()
        if err != nil {return right, err}
        expression = Binary{expression, operator, right}
    }

    return expression, nil
}

func factor() (Expr, error) {
    expression, err := unary()    
    if err != nil {return expression, err}

    for match(token.SLASH, token.STAR) {
        operator := previous()
        right, err := unary()
        if err != nil {return right, err}
        expression = Binary{expression, operator, right}
    }

    return expression, nil
}

func unary() (Expr, error) {
    if match(token.BANG, token.MINUS) {
        operator := previous()
        right, err := unary()
        if err != nil {return right, err}
        return Unary{operator, right}, nil
    }

    if match(token.PLUS, token.STAR, token.SLASH, token.EQUAL_EQUAL, token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
        operator := previous()
        return nil, parseError(operator, "Binary operator without left-hand operarand")
    }

    return primary()
}

func primary() (Expr, error) {
    if match(token.FALSE) {return Literal{false}, nil}
    if match(token.TRUE) {return Literal{true}, nil}
    if match(token.NIL) {return Literal{nil}, nil}

    if match(token.NUMBER, token.STRING) {
        return Literal{previous().Literal}, nil
    }

    if match(token.IDENTIFIER) {
        return Variable{previous()}, nil
    }
    
    if match(token.LEFT_PAREN) {
        expression, err := expression()
        if err != nil {return expression, err}
        consume(token.RIGHT_PAREN, "Expect ')' after expression.")
        return Grouping{expression}, nil
    }

    return Grouping{}, parseError(peek(), "Expect expression.")
}

func match(types... token.TokenType) bool {
    for _, t := range types {
        if check(t) {
            advance()
            return true
        }
    } 
    return false
}

func consume(tokenType token.TokenType, message string) (token.Token, error) {
    if check(tokenType) {return advance(), nil}

    return peek(), errors.New(message)
}

func check(t token.TokenType) bool {
    if isAtEnd() {return false}

    return peek().TokenType == t
}

func advance() token.Token {
    if !isAtEnd() {current++}
    return previous()
}

func isAtEnd() bool {
    return peek().TokenType == token.EOF
}

func peek() token.Token {
    return tokens[current]
}

func previous() token.Token {
    return tokens[current - 1]
}

func parseError(token token.Token, message string) error {
    Error(token, message)
    return ParseError{}
}

func Error(tokeny token.Token, message string) {
	if tokeny.TokenType == token.EOF {
		loxError.Report(tokeny.Line, " at end", message)
	} else {
        loxError.Report(tokeny.Line, " at '" + tokeny.Lexeme + "'", message)
    }
}

func synchronize() {
    advance()

    for !isAtEnd() {
        if previous().TokenType == token.SEMICOLON {return}

        switch peek().TokenType {
        case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
            return
        }

        advance()
    }
}
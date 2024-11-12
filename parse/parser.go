package parse

import (
	"errors"
	"fmt"
	. "lox/interpret"
	"lox/loxError"
	"lox/token"
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

    if match(token.CLASS) {
        out, err = classDeclaration()
    } else if match(token.FUN) {
        out, err = function("function")
    } else if match(token.VAR) {
        out, err = varDeclaration()
    } else {
        out, err = statement()
    }
    
    if err != nil {
        synchronize()     
        return nil
    }

    return out
}

func classDeclaration() (Stmt, error) {
    name, err := consume(token.IDENTIFIER, "Expect class name.")
    if err != nil {return nil, err}
    
    _, err = consume(token.LEFT_BRACE, "Expect '{' before class body.")
    if err != nil {return nil, err}

    var methods []Function
    var staticMethods []Function
    var getters []Function
    for !check(token.RIGHT_BRACE) && !isAtEnd() {
        var class = check(token.CLASS)
        if class {consume(token.CLASS, "")}

        funcType := "method"
        var getter = doublePeek().TokenType != token.LEFT_PAREN
        if getter {
            funcType = "getter"
        }
        fun, err := function(funcType)
        if err != nil {return fun, err}
        
        if class {
            staticMethods = append(staticMethods, fun)
        } else if getter  {
            getters = append(getters, fun)
        } else {
            methods = append(methods, fun)
        }
    }

    _, err = consume(token.RIGHT_BRACE, "Expect '}' after class body.")
    if err != nil {return nil, err}

    return Class{name, methods, staticMethods, getters}, nil
}

func function(kind string) (Function, error) {
    name, err := consume(token.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
    if err != nil {return Function{}, err}

    var parameters []token.Token
    if kind != "getter" {
        consume(token.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name", kind))
        if !check(token.RIGHT_PAREN) {
            for commad := true; commad; commad = match(token.COMMA) {
                if len(parameters) >= 255 {
                    loxError.TokenError(peek(), "Can't have more than 255 parameters.")
                }

                toAdd, err := consume(token.IDENTIFIER, "Expect parameter name.")
                if err != nil {return Function{}, err}
                parameters = append(parameters, toAdd)
            }
        }
        _, err = consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
        if err != nil {return Function{}, err}
    }

    _, err = consume(token.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body", kind))
    if err != nil {return Function{}, err}

    body := block()
    return Function{name, parameters, body}, nil
}

func varDeclaration() (Stmt, error) {
    name, err := consume(token.IDENTIFIER, "Expect variable name.")
    if err != nil {return nil, err}

    var initializer Expr
    if match(token.EQUAL) {
        initializer, err = expression()
        if err != nil {return nil, err}
    }

    _, err = consume(token.SEMICOLON, "Expect ';' after variable declaration.")
    if err != nil {return nil, err}

    return Var{name, initializer}, nil
}



func statement() (Stmt, error) {
    if match(token.FOR) {return forStatement()}
    if match(token.IF) {return ifStatement()}
    if match(token.PRINT) {return printStatement()}
    if match(token.RETURN) {return returnStatement()}
    if match(token.WHILE) {return whileStatement()}
    if match(token.LEFT_BRACE) {return Block{block()}, nil}
    if match(token.BREAK) {return breakStatement()}
    return expressionStatement()
}

func breakStatement() (Stmt, error) {
    _, err := consume(token.SEMICOLON, "Expect ';' after break statement.")
    if err != nil {return nil, err}

    return Break{}, nil
}

func forStatement() (Stmt, error) {
    consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

    var initializer Stmt
    var err error
    if (match(token.SEMICOLON)) {
        initializer = nil
    } else if (match(token.VAR)) {
        initializer, err = varDeclaration()
    } else {
        initializer, err = expressionStatement()
    }
    if err != nil {return initializer, err}

    var condition Expr
    if !check(token.SEMICOLON) {
        condition, err = expression()
        if err != nil {return nil, err}
    }
    consume(token.SEMICOLON, "Expect ';' after loop condition")

    var increment Expr
    if !check(token.RIGHT_PAREN) {
        increment, err = expression()
        if err != nil {return nil, err}
    }
    consume(token.RIGHT_PAREN, "Expect ')' after for clauses")
    
    body, err := statement()
    if err != nil {return body, err}

    if increment != nil {
        body = Block{[]Stmt{body, Expression{increment}}}
    }
    
    if condition == nil {condition = Literal{true}}
    body = While{condition, body}

    if initializer != nil {
        body = Block{[]Stmt{initializer, body}}
    }
    
    return body, nil
}

func whileStatement() (Stmt, error) {
    consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
    condition, err := expression()
    if err != nil {return nil, err}
    
    consume(token.RIGHT_PAREN, "Expect ')' after 'condition'.")
    body, err := statement()
    if err != nil {return body, err}

    return While{condition, body}, nil
}

func ifStatement() (Stmt, error) {
    consume(token.LEFT_PAREN, "Expect '(' after 'if'.") 
    condition, err := expression() 
    if err != nil {return nil, err}
    consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

    thenBranch, err := statement()
    if err != nil {return nil, err}
    var elseBranch Stmt = nil
    if (match(token.ELSE)) {
        elseBranch, err = statement()
        if err != nil {return nil, err}
    }

    return If{condition, thenBranch, elseBranch}, nil
}

func block() []Stmt {
    var statements []Stmt

    for !check(token.RIGHT_BRACE) && !isAtEnd() {
        statements = append(statements, declaration())
    }

    consume(token.RIGHT_BRACE, "Expect '}' after block.")
    return statements
}

func printStatement() (Stmt, error) {
    value, err := expression()
    if err != nil {return nil, err}
    _, err = consume(token.SEMICOLON, "Expect ';' after value.")
    if err != nil {return nil, err}
    return Print{value}, nil
}

func returnStatement() (Stmt, error) {
    keyword := previous()
    var value Expr = nil
    if !check(token.SEMICOLON) {
        var err error
        value, err = expression()
        if err != nil {return nil, err}
    }

    consume(token.SEMICOLON, "Expect ';' after return value.")
    return Return{keyword, value}, nil
}

func expressionStatement() (Stmt, error) {
    expre, err := expression()
    if err != nil {return nil, err}
    
    _, err = consume(token.SEMICOLON, "Expect ';' after expression.")
    if err != nil {return nil, err}
    
    return Expression{expre}, nil
}

func assignment() (Expr, error) {
    expr, err := or()
    if err != nil {return nil, err}

    if match(token.EQUAL) {
        equals := previous()
        value, err := assignment()
        if err != nil {return value, err}

        v, okv := expr.(Variable)
        i, oki := expr.(Get)
        if okv {
            name := v.Name
            return Assign{name, value}, nil
        } else if oki {
            return Set{i.Object, i.Name, value}, nil
        }

        parseError(equals, "Invalid assignment target.")
    }

    return expr, nil
}

func or() (Expr, error) {
    expr, err := and()
    if err != nil {return expr, err}

    for match(token.OR) {
        operator := previous()
        right, err := and()
        if err != nil {return right, err}

        expr = Logical{expr, operator, right}
    }

    return expr, nil
}

func and() (Expr, error) {
    expr, err := equality()
    if err != nil {return expr, err}

    for match(token.AND) {
        operator := previous()
        right, err := equality()
        if err != nil {return right, err}
        expr = Logical{expr, operator, right}
    }

    return expr, nil
}

func expression() (Expr, error) {
    return assignment()
}

func ternary() (Expr, error) {
    expression, err := equality()    
    if err != nil {return expression, err}

    for match(token.QUESTION) {
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

    return call()
}

func call() (Expr, error) {
    expr, err := primary()
    if err != nil {return expr, err}

    for true {
        if match(token.LEFT_PAREN) {
            expr, err = finishCall(expr)
            if err != nil {return expr, err}
        } else if match(token.DOT) {
            name, err := consume(token.IDENTIFIER, "Expect property name after '.'.")
            if err != nil {return nil, err}
            expr = Get{expr, name}
        } else {
            break
        }
    }

    return expr, nil
}

func finishCall(callee Expr) (Expr, error) {
    var arguments []Expr
    if !check(token.RIGHT_PAREN) {
        for commad := true; commad; commad = match(token.COMMA) {
            if (len(arguments) >= 255) {
                loxError.TokenError(peek(), "Can't have more than 255 arguments.")
            }
            expr, err := expression()
            if err != nil {return expr, err}
            arguments = append(arguments, expr)
        }
    }

    paren, err := consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
    if err != nil {return nil, err}

    return Call{callee, paren, arguments}, nil
}

func primary() (Expr, error) {
    if match(token.FALSE) {return Literal{false}, nil}
    if match(token.TRUE) {return Literal{true}, nil}
    if match(token.NIL) {return Literal{nil}, nil}

    if match(token.NUMBER, token.STRING) {
        return Literal{previous().Literal}, nil
    }

    if match(token.THIS) {
        return This{previous()}, nil
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

    loxError.Error(tokens[current].Line, message)
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

func doublePeek() token.Token {
    return tokens[current + 1]
}

func previous() token.Token {
    return tokens[current - 1]
}

func parseError(token token.Token, message string) error {
    loxError.TokenError(token, message)
    return ParseError{}
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
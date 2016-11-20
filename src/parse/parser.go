package parse

import (
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	toks       chan Token
	buf        [3]Token
	bufIndex   int
	errors     [10]string
	errorIndex int
}

func (p *parser) next() Token {
	if p.bufIndex == 0 {
		p.buf[2] = p.buf[1]
		p.buf[1] = p.buf[0]
		p.buf[0] = <-p.toks
	} else {
		p.bufIndex--
	}

	return p.buf[p.bufIndex]
}

func (p *parser) backup() {
	p.bufIndex++
	if p.bufIndex >= len(p.buf) {
		panic("Backup exceeded buffer size")
	}
}

func (p *parser) peek() Token {
	tok := p.next()
	p.backup()

	return tok
}

func (p *parser) errorf(tok Token, message string, args ...interface{}) Expr {
	panic(SyntaxError{tok, fmt.Sprintf(message, args...)})

	return nil
}

func (p *parser) unexpected(tok Token) Expr {
	if tok.Type == tokError {
		return p.errorf(tok, tok.Value)
	}

	return p.errorf(tok, "Unexpected token %v", tok)
}

func (p *parser) accept(typ TokenType) string {
	tok := p.next()
	if tok.Type != typ {
		p.unexpected(tok)
		return ""
	}

	return tok.Value
}

func Parse(tokens chan Token) (stmts StmtCollection, er error) {
	defer func() {
		if recoveredEr := recover(); recoveredEr != nil {
			if _, ok := recoveredEr.(SyntaxError); ok {
				er = recoveredEr.(error)
			} else {
				panic(recoveredEr)
			}
		}
	}()

	stmts = parseStmts(&parser{toks: tokens}, tokEOF)
	return stmts, er
}

func parseStmts(p *parser, endToken TokenType) StmtCollection {
	statements := make([]Stmt, 0, 10)
	for p.peek().Type != endToken {
		statements = append(statements, parseStmt(p))
	}

	p.accept(endToken)
	return statements
}

func parseStmt(p *parser) Stmt {
	return Stmt{Expr: parseExprStatement(p)}
}

func parseExprStatement(p *parser) Expr {
	switch tok := p.next(); tok.Type {
	case tokVar:
		return parseVarDef(p)
	case tokPrint:
		return parsePrintStmt(p)
	case tokDef:
		return parseFuncDef(p)
	case tokIf:
		return parseIfStmt(p)
	case tokWhile:
		return parseWhileStmt(p)
	case tokIdent: //Could be arithmatic or a function call
		return parseFuncCallOrAssignment(p)
	default:
		p.unexpected(tok)
		return nil
	}
}

func parseFuncCallOrAssignment(p *parser) Expr {
	nextTok := p.peek()
	p.backup() // Now at the start of the line

	switch nextTok.Type {
	case tokOpenParen:
		return parseFuncCall(p)
	case tokIdent, tokEquals:
		return parseAssignment(p)
	default:
		p.unexpected(nextTok)
		return nil
	}
}

func parseFuncCall(p *parser) Expr {
	ident := parseIdent(p)
	p.accept(tokOpenParen)

	args := parseIdentifierList(p, tokCloseParen)
	p.accept(tokSemicolon)

	return FuncCall{Func: ident, Args: args}
}

func parseAssignment(p *parser) Expr {
	lhs := parseIdentifierList(p, tokEquals)
	rhs := parseAssignmentRhs(p)
	p.accept(tokSemicolon)

	return Assignment{Lhs: lhs, Rhs: rhs}
}

func parseAssignmentRhs(p *parser) Expr {
	switch tok := p.next(); tok.Type {
	case tokNum:
		c, _ := strconv.Atoi(tok.Value) // Validated by the lexer already
		return Lit{Val: c}
	case tokChar:
		return Lit{Val: int(tok.Value[0])}
	case tokIdent:
		ident := asIdent(tok.Value)
		if ident.Op != None && ident.Op != Floor {
			p.unexpected(tok)
		}
		return ident
	default:
		p.unexpected(tok)
		return nil
	}
}

func parseIfStmt(p *parser) Expr {
	subject := parseIdent(p)
	p.accept(tokOpenBrace)
	body := parseStmts(p, tokCloseBrace)

	return IfStmt{Subject: subject, Body: body}
}

func parseWhileStmt(p *parser) Expr {
	subject := parseIdent(p)
	p.accept(tokOpenBrace)
	body := parseStmts(p, tokCloseBrace)

	return WhileStmt{Subject: subject, Body: body}
}

func parseFuncDef(p *parser) Expr {
	funcName := parseIdent(p)
	p.accept(tokOpenParen)

	args := parseIdentifierList(p, tokCloseParen)
	p.accept(tokOpenBrace)

	body := parseStmts(p, tokCloseBrace)

	return FuncDec{Name: funcName, Args: args, Body: body}
}

func parseVarDef(p *parser) Expr {
	return VarDef{Idents: parseIdentifierList(p, tokSemicolon)}
}

func parsePrintStmt(p *parser) Expr {
	return PrintStmt{Idents: parseIdentifierList(p, tokSemicolon)}
}

func parseIdentifierList(p *parser, endToken TokenType) []Ident {
	args := make([]Ident, 0, 10)
	for tok := p.next(); tok.Type != endToken; tok = p.next() {
		if tok.Type != tokIdent {
			p.unexpected(tok)
		}

		args = append(args, asIdent(tok.Value))
	}

	return args
}

func parseIdent(p *parser) Ident {
	return asIdent(p.accept(tokIdent))
}

func asIdent(value string) Ident {
	if len(value) == 0 {
		return Ident{}
	}

	return Ident{Op: getOp(value), Id: trimOp(value)}
}

func getOp(identifier string) IdentifierOp {
	switch identifier[0] {
	case '_':
		return Floor
	case '+':
		return Add
	case '-':
		return Sub
	default:
		return None
	}
}

func trimOp(identifier string) string {
	return strings.TrimLeft(identifier, "+-_")
}

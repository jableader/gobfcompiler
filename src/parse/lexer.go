package parse

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type TokenType int

const (
	tokError TokenType = iota
	tokEOF

	// Brackets / Parenthesis
	tokOpenParen
	tokCloseParen
	tokOpenBrace
	tokCloseBrace

	// Punctuation
	tokSemicolon

	// Literals
	tokNum
	tokChar

	// Assorted
	tokEquals
	tokIdent

	// Keywords
	tokKeyword // Used to distinguish keywords for print method
	tokDef
	tokWhile
	tokIf
	tokElse
	tokPrint
	tokVar
)

type Token struct {
	Type   TokenType
	Value  string
	endPos int
	lexer  *lexer
}

func (t Token) LineNumber() int {
	return strings.Count(t.lexer.input[0:t.endPos], "\n") + 1
}

func (t Token) String() string {
	switch {
	case t.Type == tokEOF:
		return "EOF"
	case t.Type == tokError:
		return fmt.Sprintf("Err: %s", t.Value)
	case t.Type > tokKeyword:
		return fmt.Sprintf("<%s>", t.Value)
	case t.Type == tokNum:
		return fmt.Sprintf("D(%s)", t.Value)
	case t.Type == tokIdent:
		return fmt.Sprintf("I(%s)", t.Value)
	case t.Type == tokChar:
		return fmt.Sprintf("C(%s)", t.Value)
	default:
		return t.Value
	}
}

// lexer holds the state of the scanner.
type lexer struct {
	name  string // used only for error reports.
	input string // the string being scanned.

	tokens chan Token // channel of scanned items.

	start int // start position of this item.
	pos   int // current position in the input.
	width int // width of last rune read from input.
}

const eof = -1

func (l *lexer) emit(tokenType TokenType) {
	l.tokens <- Token{
		Value:  l.input[l.start:l.pos],
		Type:   tokenType,
		lexer:  l,
		endPos: l.pos,
	}

	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	next := l.next()
	l.backup()
	return next
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) int {
	i := 0
	for ; l.accept(valid); i++ {
	}
	return i
}

func (l *lexer) skipWhitespace() {
	const whitespace = " \t\r\n"

	l.acceptRun(whitespace)
	l.ignore()
}

func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) errorf(message string, args ...interface{}) stateFn {
	parserMessage := fmt.Sprintf(message, args...)

	l.tokens <- Token{
		Type:   tokError,
		Value:  fmt.Sprintf("Message: %s\nToken: %v", parserMessage, l.current()),
		lexer:  l,
		endPos: l.pos,
	}

	return nil
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

const letterChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberChars = "0123456789"
const newLine = '\n'

func isLetter(c rune) bool {
	return strings.IndexRune(letterChars, c) >= 0
}

// run lexes the input by executing state functions
// until the state is nil.
func run(l *lexer) {
	for state := lexStatement; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func grabIdentifier(l *lexer, prefixes string) bool {
	l.skipWhitespace()
	l.accept(prefixes)
	if !l.accept("$") || l.acceptRun(letterChars) == 0 {
		return false
	}

	l.emit(tokIdent)
	return true
}

func lexStatement(l *lexer) stateFn {
	l.skipWhitespace()

	next := l.next()
	switch next {
	case eof:
		l.emit(tokEOF)
		return nil
	case '}':
		l.emit(tokCloseBrace)
		return lexStatement
	case '#':
		return lexComment
	case '$', '+', '-':
		l.backup()
		return lexIdentifier
	default:
		if isLetter(next) {
			l.backup()
			return lexKeyword
		}

		return l.errorf("Unexpected character at statement start: %c", next)
	}
}

func lexComment(l *lexer) stateFn {
	for l.next() != newLine {
	}
	l.ignore()

	return lexStatement
}

func lexEndStatement(l *lexer) stateFn {
	l.skipWhitespace()
	if !l.accept(";") {
		return l.errorf("Expected end statement. Are you missing a semicolon?")
	}

	l.emit(tokSemicolon)
	return lexStatement
}

func lexKeyword(l *lexer) stateFn {
	l.acceptRun(letterChars)
	switch l.current() {
	case "if":
		l.emit(tokIf)
		return lexControlStatement
	case "while":
		l.emit(tokWhile)
		return lexControlStatement
	case "def":
		l.emit(tokDef)
		return lexFunctionDefinition
	case "print":
		l.emit(tokPrint)
		return lexVar
	case "var":
		l.emit(tokVar)
		return lexVar
	default:
		return l.errorf("Unknown keyword (%v)", l.current())
	}
}

func lexVar(l *lexer) stateFn {
	varsGrabbed := grabCommaSeperatedArgs(l, "")
	if varsGrabbed == 0 {
		return l.errorf("Expected comma seperated arguments list")
	}

	return lexEndStatement
}

func lexFunctionDefinition(l *lexer) stateFn {
	l.skipWhitespace()
	if !grabIdentifier(l, "") {
		return l.errorf("Expected an identifier")
	}

	l.skipWhitespace()
	if l.next() != '(' {
		return l.errorf("Expected open bracket")
	}

	l.emit(tokOpenParen)
	grabCommaSeperatedArgs(l, "")

	l.skipWhitespace()
	if l.next() != ')' {
		return l.errorf("Expected close bracket")
	}

	l.emit(tokCloseParen)
	l.skipWhitespace()
	if l.next() != '{' {
		l.errorf("Expected opening brace")
	}
	l.emit(tokOpenBrace)

	return lexStatement
}

func lexControlStatement(l *lexer) stateFn {
	if !grabIdentifier(l, "_") {
		return l.errorf("Expected an identifier")
	}

	l.skipWhitespace()
	if l.next() == '{' {
		l.emit(tokOpenBrace)
		return lexStatement
	} else {
		return l.errorf("Expected open brace")
	}
}

func lexRhs(l *lexer) stateFn {
	l.skipWhitespace()
	if l.acceptRun(numberChars) > 0 {
		l.emit(tokNum)
		return lexEndStatement
	}

	if l.next() == '\'' {
		l.ignore()
		l.next()
		l.emit(tokChar)
		if l.next() != '\'' {
			return l.errorf("Expected closing quote")
		}
		l.ignore()
		return lexEndStatement
	}

	l.backup()
	if grabIdentifier(l, "_") {
		return lexEndStatement
	}

	return l.errorf("Expected identifier or literal")
}

func lexIdentifier(l *lexer) stateFn {
	firstWasNotOp := l.peek() == '$'
	argsGrabbed := grabCommaSeperatedArgs(l, "+-")

	switch l.next() {
	case '=':
		l.emit(tokEquals)
		return lexRhs

	case '(':
		if firstWasNotOp && argsGrabbed == 1 {
			l.emit(tokOpenParen)

			grabCommaSeperatedArgs(l, "")
			l.skipWhitespace()
			if l.next() == ')' {
				l.emit(tokCloseParen)
				return lexEndStatement
			} else {
				return l.errorf("Expected close paren")
			}
		}
	}

	return l.errorf("Unexpected")
}

func grabCommaSeperatedArgs(l *lexer, prefixes string) int {
	for count := 0; ; {
		if !grabIdentifier(l, prefixes) {
			return count
		}
		count++

		l.skipWhitespace()
		if l.next() != ',' {
			l.backup()
			return count
		}
	}
}

func Lex(input string) chan Token {
	l := lexer{
		name:   "barry",
		input:  input,
		tokens: make(chan Token),
	}
	go run(&l)

	return l.tokens
}

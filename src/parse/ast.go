package parse

import (
	"fmt"
	"strings"
)

type Expr interface {
	String() string
}
type Stmt struct {
	Expr Expr
}

func (s Stmt) String() string {
	return s.Expr.String() + ";"
}

type StmtCollection []Stmt

func (s StmtCollection) String() string {
	strs := make([]string, len(s))
	for i, v := range s {
		strs[i] = v.String()
	}

	return strings.Join(strs, "")
}

type IdentifierOp int

const (
	None IdentifierOp = iota
	Add
	Sub
	Floor
)

type Ident struct {
	Op IdentifierOp
	Id string
}

func (i Ident) String() string {
	switch i.Op {
	case Add:
		return "+" + i.Id
	case Sub:
		return "-" + i.Id
	case Floor:
		return "_" + i.Id
	default:
		return i.Id
	}
}

type Lit struct {
	Val int
}

func (l Lit) String() string {
	return fmt.Sprintf("%v", l.Val)
}

type Assignment struct {
	Lhs []Ident
	Rhs Expr
}

func (a Assignment) String() string {
	return fmt.Sprintf("%v = %v", a.Lhs, a.Rhs)
}

type VarDef struct {
	Idents []Ident
}

func (v VarDef) String() string {
	return "var " + fmt.Sprintf("%v", v.Idents)
}

type PrintStmt struct {
	Idents []Ident
}

func (v PrintStmt) String() string {
	return fmt.Sprintf("print %v", v.Idents)
}

type IfStmt struct {
	Subject Ident
	Body    StmtCollection
}

func (i IfStmt) String() string {
	return fmt.Sprintf("if %v { %v }", i.Subject, i.Body)
}

type WhileStmt struct {
	Subject Ident
	Body    StmtCollection
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("while %v { %v }", w.Subject, w.Body)
}

type FuncDec struct {
	Name Ident
	Args []Ident
	Body StmtCollection
}

func (w FuncDec) String() string {
	return fmt.Sprintf("def %v(%v) { %v }", w.Name, w.Args, w.Body)
}

type FuncCall struct {
	Func Ident
	Args []Ident
}

func (f FuncCall) String() string {
	return fmt.Sprintf("%v(%v)", f.Func, f.Args)
}

type SyntaxError struct {
	Token   Token
	Message string
}

func (s SyntaxError) String() string {
	return fmt.Sprintf("Syntax Error at %d. %s", s.Token.LineNumber(), s.Message)
}

func (s SyntaxError) Error() string {
	return s.String()
}

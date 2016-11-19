package parse

import (
	"reflect"
	"fmt"
)

type Pointer struct {
	Loc int
	Name *string
}

type Program interface {
	EnterScope()
	ExitScope()

	OpenLoop()
	CloseLoop()

	CreatePt(id *string, near int) (Pointer, error)
	DestroyPt(pt Pointer)
	GetPt(id string) (Pointer, bool)

	Print(pt *Pointer)
	MoveTo(pt Pointer)
	Add(pt Pointer, n int)
	Err(msg string, expr Expr)
}

type compilable interface {
	compile(prog Program)
}

func (expr VarDef) compile(p Program) {
	for _, ident := range expr.Idents {
		if _, exists := p.CreatePt(&ident.Id, -1); exists != nil {
			p.Err("Cannot redefine variable within the same scope", ident)
		}
	}
}

func (expr Assignment) compile(p Program) {
	var rhs Pointer
	switch val := expr.Rhs.(type) {
	case Lit:
		rhs, _ := p.CreatePt(nil, -1)
		p.Add(rhs, val.Val)
	case Ident:
		pt, ok := p.GetPt(val.Id)
		if !ok {
			p.Err("RHS of expression is undefined", expr)
			return
		}

		rhs = pt
	}

	p.MoveTo(rhs)
	p.OpenLoop()

	for _, v := range expr.Lhs {
		pt, ok := p.GetPt(v.Id)
		if !ok {
			p.Err("Identifier is undefined", v)
		} else {
			switch v.Op {
			case Add: p.Add(pt, 1)
			case Sub: p.Add(pt, -1)
			default: p.Err("Invalid operator.", v)
			}
		}
	}
	p.MoveTo(rhs)
	p.CloseLoop()
}

func (expr PrintStmt) compile(p Program) {
	for _, v := range expr.Idents {
		if v.Op != None {
			p.Err("Unexpected operator in print statement", v)
		}

		pt, ok := p.GetPt(v.Id)
		if !ok {
			p.Err("Variable is undefined", v)
		} else {
			p.Print(&pt)
		}
	}
}

func (expr *WhileStmt) compile(p Program) {
	pt, subjectExists := p.GetPt(expr.Subject.Id)
	if subjectExists {
		p.MoveTo(pt)
	}	else {
		p.Err("Subject of loop is undefined", expr.Subject)
	}


	p.OpenLoop()
	defer p.CloseLoop()

	p.EnterScope()
	defer p.ExitScope()

	expr.Body.compile(p)

	if subjectExists {
		switch expr.Subject.Op {
		case Add: p.Add(pt, 1)
		case Sub, Floor: p.Add(pt, -1)
		}

		p.MoveTo(pt)
	}
}

func (expr FuncDec) compile(p Program) {
	p.Err("Not implemented", expr)
}

func (expr SyntaxError) compile(p Program) {
	p.Err(expr.String(), expr)
}

func (expr Stmt) compile(p Program) {
	val, ok := expr.Expr.(compilable)
	if !ok {
		p.Err(fmt.Sprintf("%v is not compilable", reflect.TypeOf(expr.Expr)), expr)
		return
	}

	val.compile(p)
}

func (stmts StmtCollection) compile(p Program) {
	for _, stmt := range stmts {
		stmt.compile(p)
	}
}

func Compile(p Program, stmts StmtCollection) {
	stmts.compile(p)
}

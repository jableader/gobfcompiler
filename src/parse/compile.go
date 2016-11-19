package parse

import (
	"reflect"
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
	Err(expr Expr, s string, args... interface{})
}

type compilable interface {
	compile(prog Program)
}

func (expr VarDef) compile(p Program) {
	for _, ident := range expr.Idents {
		if _, exists := p.CreatePt(&ident.Id, -1); exists != nil {
			p.Err(ident, "Cannot redefine variable within the same scope")
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
			p.Err(val, "%v is undefined", val.Id)
			return
		}

		rhs = pt
	}

	for _, v := range expr.Lhs {
		if v.Op == None {
			pt, _ := p.GetPt(v.Id) // Ignoring OK: It's handled below
			p.MoveTo(pt)
			p.OpenLoop()
			p.Add(pt, -1)
			p.CloseLoop()
		}
	}

	p.MoveTo(rhs)
	p.OpenLoop()

	for _, v := range expr.Lhs {
		pt, ok := p.GetPt(v.Id)
		if !ok {
			p.Err(v, "Identifier is undefined")
		} else {
			switch v.Op {
			case Add, None: p.Add(pt, 1)
			case Sub: p.Add(pt, -1)
			default: p.Err(v, "Invalid operator")
			}
		}
	}
	p.MoveTo(rhs)
	p.Add(rhs, -1)
	p.CloseLoop()
}

func (expr PrintStmt) compile(p Program) {
	for _, v := range expr.Idents {
		if v.Op != None {
			p.Err(v, "Unexpected operator in print statement")
		}

		pt, ok := p.GetPt(v.Id)
		if !ok {
			p.Err(v, "Variable is undefined")
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
		p.Err(expr.Subject, "Subject of loop is undefined")
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
	p.Err(expr, "Functions are not implemented... yet")
}

func (expr SyntaxError) compile(p Program) {
	p.Err(expr, expr.String())
}

func (expr Stmt) compile(p Program) {
	val, ok := expr.Expr.(compilable)
	if !ok {
		p.Err(expr, "%v is not compilable", reflect.TypeOf(expr.Expr))
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

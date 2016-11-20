package compiler

import (
	"asm"
	"memory"
	"parse"
	"reflect"
	"scope"
	"sort"
)

type program struct {
	sc  *scope.Scope
	mem *memory.Memory
	asm asm.Assembler
}

func (p *program) GetPt(id parse.Ident) (asm.Pointer, bool) {
	variable, ok := p.sc.Get(id.Id)
	if !ok {
		p.asm.Err(id, "%v is not defined", id.Id)
		return asm.NullPointer, false
	}

	pt, ok := variable.Value.(int)
	if !ok {
		p.asm.Err(id, "Expected pointer, got %v", reflect.TypeOf(variable.Value))
		return asm.NullPointer, false
	}

	return asm.Pointer(pt), true
}

func (p *program) DefPt(id *string, near int) (asm.Pointer, error) {
	pt := p.mem.Malloc(near)
	_, err := p.sc.Define(id, pt)

	return asm.Pointer(pt), err
}

func (p *program) EnterScope() {
	p.sc = p.sc.Enter()
}

func (p *program) ExitScope() {
	p.sc = p.sc.Exit()
}

func compileVarDef(p *program, expr parse.VarDef) {
	for _, ident := range expr.Idents {
		if _, exists := p.DefPt(&ident.Id, -1); exists != nil {
			p.asm.Err(ident, "Cannot redefine variable within the same scope")
		}
	}
}

type ptIdentWrapper struct {
	id parse.Ident
	pt asm.Pointer
}
type byPt []ptIdentWrapper

func (s byPt) Len() int {
	return len(s)
}
func (s byPt) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPt) Less(i, j int) bool {
	return s[i].pt < s[j].pt
}

func getAndSort(p *program, idents []parse.Ident) []ptIdentWrapper {
	res := make([]ptIdentWrapper, len(idents))
	for i, id := range idents {
		pt, _ := p.GetPt(id)
		res[i] = ptIdentWrapper{id, pt}
	}

	sort.Sort(byPt(res))
	return res
}

func compileAssignment(p *program, expr parse.Assignment) {
	p.asm.Comment(expr.String())

	var rhs asm.Pointer
	switch val := expr.Rhs.(type) {
	case parse.Lit:
		rhs, _ = p.DefPt(nil, -1)
		p.asm.Add(rhs, val.Val)
	case parse.Ident:
		rhs, _ = p.GetPt(val)
	}

	lhs := getAndSort(p, expr.Lhs)
	for _, v := range lhs {
		if v.id.Op == parse.None {
			p.asm.OpenLoop(v.pt)
			p.asm.Add(v.pt, -1)
			p.asm.CloseLoop()
		}
	}

	p.asm.OpenLoop(rhs)
	p.asm.Add(rhs, -1)

	for _, v := range lhs {
		switch v.id.Op {
		case parse.Add, parse.None:
			p.asm.Add(v.pt, 1)
		case parse.Sub:
			p.asm.Add(v.pt, -1)
		default:
			p.asm.Err(v.id, "Invalid operator")
		}
	}

	p.asm.CloseLoop()
}

func compilePrintStmt(p *program, expr parse.PrintStmt) {
	for _, v := range expr.Idents {
		if v.Op != parse.None {
			p.asm.Err(v, "Unexpected operator in print statement")
		}

		pt, ok := p.GetPt(v)
		if ok {
			p.asm.Print(pt)
		}
	}
}

func compileWhileStmt(p *program, expr parse.WhileStmt) {
	pt, subjectExists := p.GetPt(expr.Subject)
	p.asm.OpenLoop(pt)
	defer p.asm.CloseLoop()

	p.EnterScope()
	defer p.ExitScope()

	compileStmtCollection(p, expr.Body)

	if subjectExists {
		switch expr.Subject.Op {
		case parse.Add:
			p.asm.Add(pt, 1)
		case parse.Sub, parse.Floor:
			p.asm.Add(pt, -1)
		}
	}
}

func compileFuncDec(p *program, expr parse.FuncDec) {
	p.asm.Err(expr, "Functions are not implemented... yet")
}

func compileSyntaxError(p *program, expr parse.SyntaxError) {
	p.asm.Err(expr, expr.String())
}

func compileStmt(p *program, expr parse.Stmt) {
	switch val := expr.Expr.(type) {
	case parse.VarDef:
		compileVarDef(p, val)
	case parse.Assignment:
		compileAssignment(p, val)
	case parse.PrintStmt:
		compilePrintStmt(p, val)
	case parse.WhileStmt:
		compileWhileStmt(p, val)
	case parse.FuncDec:
		compileFuncDec(p, val)
	case parse.SyntaxError:
		compileSyntaxError(p, val)
	case parse.Stmt:
		compileStmt(p, val)
	case parse.StmtCollection:
		compileStmtCollection(p, val)
	default:
		p.asm.Err(expr, "%v is not compilable", reflect.TypeOf(expr.Expr))
	}
}

func compileStmtCollection(p *program, stmts parse.StmtCollection) {
	for _, stmt := range stmts {
		compileStmt(p, stmt)
	}
}

func Compile(a asm.Assembler, stmts parse.StmtCollection) {
	compileStmtCollection(&program{scope.New(), memory.New(), a}, stmts)
}

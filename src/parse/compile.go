package parse

import (
	"errors"
)

var ErrorAlreadyExists = errors.New("Pointer already exists")

type program struct {
	sc *scope
	mem *memory
	pc int

	outp chan BfNode
}
func (prog *program) enter(ex Expr) *scope {
	prog.sc = &scope{prog.sc, prog, make(map[string]pointer)}
	return prog.sc
}

type pointer struct {
	loc int
	name string
}

type memory []bool
func (m memory) malloc(near int) int {
	for i := 0; i < len(m); i++ {
		if !m[i] {
			m[i] = true
			return i
		}
	}

	panic("Memory is full!")
}

func (m memory) free(p int) {
	m[p] = false
}

type scope struct {
	parent *scope
	prog *program
	pts map[string]pointer
}

func (s *scope) get(name string) (pointer, bool) {
	for sc := s; sc != nil; sc = s.parent {
		pt, wasFound := sc.pts[name]
		if wasFound {
			return pt, true
		}
	}

	return pointer{}, false
}

func (s *scope) create(name string, near int) (pointer, error) {
	if pt, exists := s.pts[name]; exists {
		return pt, ErrorAlreadyExists
	}

	return pointer{s.prog.mem.malloc(near), name}, nil
}

func (s *scope) destroy(p pointer) {
	s.prog.mem.free(p.loc)
	delete(s.pts, p.name)
}

func compile(prog *program, stmts StmtCollection) {
	for _, stmt := range stmts {
		switch s := stmt.Expr.(type) {
			case Stmt:
			case Ident:
			case Lit:
			case Assignment:
			case VarDef:
			case PrintStmt:
			case IfStmt:
			case WhileStmt:
			case FuncDec:
			case FuncCall:
			case SyntaxError:

		}
	}
}

func Compile(stmts StmtCollection) chan BfNode {
	prog := program {
		outp: make(chan BfNode),
	}

	go compile(&prog, stmts)
	return prog.outp
}

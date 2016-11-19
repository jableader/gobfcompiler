package program

import (
	"fmt"
	"parse"
	"errors"
)

type Assembler interface {
	Move(from, to int)
	Add(n int)

	StartLoop()
	EndLoop()

	Print()
	Err(msg string, expr parse.Expr)
}

var ErrorAlreadyExists = errors.New("parse.Pointer already exists")

func New(asm Assembler) parse.Program {
	prog := program {
		sc: &scope{},
		mem: &memory{},
		asm: asm,
	}
	prog.EnterScope()

	return &prog
}

type program struct {
	sc *scope
	mem *memory
	pc int

	asm Assembler
}

func (prog *program) EnterScope() {
	prog.sc = &scope{prog.sc, prog, make(map[string]parse.Pointer)}
}

func (prog *program) ExitScope() {
	prog.sc = prog.sc.parent
}

func (prog *program) OpenLoop() {
	prog.asm.StartLoop()
}

func (prog *program) CloseLoop() {
	prog.asm.EndLoop()
}

func (prog *program) CreatePt(name *string, near int) (parse.Pointer, error) {
	return prog.sc.create(name, near)
}

func (prog *program) DestroyPt(pt parse.Pointer) {
	prog.sc.destroy(pt)
}

func (prog *program) GetPt(name string) (parse.Pointer, bool) {
	return prog.sc.get(name)
}

func (prog *program) Print(pt *parse.Pointer) {
	if pt != nil {
		prog.MoveTo(*pt)
	}
	prog.asm.Print()
}

func (prog *program) MoveTo(pt parse.Pointer) {
	if prog.pc != pt.Loc {
		prog.asm.Move(prog.pc, pt.Loc)
		prog.pc = pt.Loc
	}
}

func (prog *program) Add(pt parse.Pointer, n int) {
	if n != 0 {
		prog.MoveTo(pt)
		prog.asm.Add(n)
	}
}

func (prog *program) Err(expr parse.Expr, msg string, args... interface{}) {
	prog.asm.Err(fmt.Sprintf(msg, args...), expr)
}

type memory [100]bool
func (m *memory) malloc(near int) int {
	for i := 0; i < len(m); i++ {
		if !m[i] {
			m[i] = true
			return i
		}
	}

	panic("Memory is full!")
}

func (m *memory) free(p int) {
	m[p] = false
}

type scope struct {
	parent *scope
	prog *program
	pts map[string]parse.Pointer
}

func (s *scope) get(name string) (parse.Pointer, bool) {
	for sc := s; sc != nil; sc = sc.parent {
		pt, wasFound := sc.pts[name]
		if wasFound {
			return pt, true
		}
	}

	return parse.Pointer{}, false
}

func (s *scope) create(name *string, near int) (parse.Pointer, error) {
	if name != nil {
		if pt, exists := s.pts[*name]; exists {
			return pt, ErrorAlreadyExists
		}
	}

	pt := parse.Pointer{s.prog.mem.malloc(near), name}
	if name != nil {
		s.pts[*name] = pt
	}

	return pt, nil
}

func (s *scope) destroy(pt parse.Pointer) {
	s.prog.mem.free(pt.Loc)
	if pt.Name != nil {
		delete(s.pts, *pt.Name)
	}
}

package asm

import (
	"fmt"
	"strings"
	"parse"
)

type Pointer int
const (
	NullPointer Pointer = -1
	ZeroPointer Pointer = 0
)

type Assembler interface {
	Read(pt Pointer)
	Print(pt Pointer)

	OpenLoop(pt Pointer)
	CloseLoop()

	Add(pt Pointer, n int)

	Comment(s string)
	Err(expr parse.Expr, msg string, args... interface{})
}

type assembler struct {
	loops []Pointer
	pc Pointer
	output chan BfNode
}

func New() (Assembler, chan BfNode) {
	channel := make(chan BfNode)

	return &assembler {
		loops: make([]Pointer, 0, 10),
		pc: ZeroPointer,
		output: channel,
	}, channel
}

func (a *assembler) move(to Pointer) {
	if a.pc != to {
		a.output <- bfMov { int(a.pc), int(to) }
		a.pc = to
	}
}

func (a *assembler) Add(pt Pointer, n int) {
	a.move(pt)
	a.output <- bfAdd { n }
}

func (a *assembler) OpenLoop(pt Pointer) {
	a.loops = append(a.loops, pt)
	a.move(pt)
	a.output <- bfStartLoop{}
}

func (a *assembler) CloseLoop() {
	pt := a.loops[len(a.loops) -1]
	a.move(pt)
	a.output <- bfEndLoop{}

	a.loops = a.loops[:len(a.loops) - 1]
}

func (a *assembler) Print(pt Pointer) {
	a.move(pt)
	a.output <- bfPrint{}
}

func (a *assembler) Read(pt Pointer) {
	a.move(pt)
	a.output <- bfRead{}
}

func (a *assembler) Comment(s string) {
	for _, ch := range strings.Split("+-<>[].,", "") {
		s = strings.Replace(s, ch, "_", -1)
	}
	a.output <- bfComment{s}
}

func (a *assembler) Err(expr parse.Expr, msg string, args... interface{}) {
	a.output <- bfErr {expr, fmt.Sprintf(msg, args...)}
}

type BfNode interface {
	ToBF() string
	String() string
}

type bfMov struct {
	from, to int
}

func (b bfMov) ToBF() string {
	if b.from > b.to {
		return strings.Repeat("<", b.from - b.to)
	} else {
		return strings.Repeat(">", b.to - b.from)
	}
}
func (b bfMov) String() string {
	return fmt.Sprintf("MOV %d %d", b.from, b.to)
}

type bfAdd struct {
	num int
}
func (b bfAdd) ToBF() string {
	if b.num > 0 {
		return strings.Repeat("+", b.num)
	} else {
		return strings.Repeat("-", -b.num)
	}
}
func (b bfAdd) String() string {
	return fmt.Sprintf("ADD %d", b.num)
}

type bfStartLoop struct { } // Nothing to put inside it...
func (b bfStartLoop) ToBF() string {
	return "["
}
func (b bfStartLoop) String() string {
	return "SLOOP"
}

type bfEndLoop struct { }
func (b bfEndLoop) ToBF() string {
	return "]"
}
func (b bfEndLoop) String() string {
	return "ELOOP"
}

type bfPrint struct { }
func (b bfPrint) ToBF() string {
	return "."
}
func (b bfPrint) String() string {
	return "PRINT"
}

type bfRead struct {}
func (b bfRead) ToBF() string {
	return ","
}
func (b bfRead) String() string {
	return "READ"
}

type bfErr struct {
	branch parse.Expr
	reason string
}
func (b bfErr) ToBF() string {
	return b.String()
}
func (b bfErr) String() string {
	return fmt.Sprintf("\nCompiler Error: %v, %v\n", b.reason, b.branch)
}

type bfComment struct {
	s string
}
func (b bfComment) ToBF() string {
	return fmt.Sprintf("  %v\n", b.s)
}
func (b bfComment) String() string {
	return b.s
}

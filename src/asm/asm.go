package asm

import (
	"fmt"
	"strings"
	"parse"
	"program"
)

func New() (program.Assembler, chan BfNode) {
	channel := make(chan BfNode)
	return assembler { channel }, channel
}

type assembler struct {
	output chan BfNode
}

func (a assembler) Move(from, to int) {
	a.output <- bfMov { from, to }
}

func (a assembler) Add(n int) {
	a.output <- bfAdd { n }
}

func (a assembler) StartLoop() {
	a.output <- bfStartLoop {}
}

func (a assembler) EndLoop() {
	a.output <- bfEndLoop{}
}

func (a assembler) Print() {
	a.output <- bfPrint{}
}

func (a assembler) Err(msg string, expr parse.Expr) {
	a.output <- bfErr {expr, msg}
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
		return strings.Repeat("-", b.num)
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

type bfErr struct {
	branch parse.Expr
	reason string
}
func (b bfErr) ToBF() string {
	return b.String()
}
func (b bfErr) String() string {
	return fmt.Sprintf("Compiler Error: %v, %v", b.reason, b.branch)
}

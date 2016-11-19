package parse

import (
	"fmt"
	"strings"
)

type BfNode interface {
	ToBF() string
	Children() []BfNode
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
func (b *bfStartLoop) ToBF() string {
	return "["
}
func (b *bfStartLoop) String() string {
	return "SLOOP"
}

type bfEndLoop struct { }
func (b *bfEndLoop) ToBF() string {
	return "]"
}
func (b *bfEndLoop) String() string {
	return "ELOOP"
}

type bfPrint struct { }
func (b *bfPrint) ToBF() string {
	return "."
}
func (b *bfPrint) String() string {
	return "PRINT"
}

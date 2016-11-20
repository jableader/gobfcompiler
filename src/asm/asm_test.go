package asm_test

import (
	"asm"
	"testing"
)

type assemblerMethod func(assembler asm.Assembler)

func expectBf(t *testing.T, expected string, m assemblerMethod) {
	assembler, ch := asm.New()

	go func() {
		m(assembler)
		close(ch)
	}()

	result := ""
	for node := range ch {
		result += node.ToBF()
	}

	if result != expected {
		t.Errorf("\nExpect:\t%v\nActual:\t%v", expected, result)
	}
}

func TestAddAndSubtract(t *testing.T) {
	expectBf(t, ">>>>>++<<-->,<.", func(assembler asm.Assembler) {
		assembler.Add(5, 2)
		assembler.Add(3, -2)
		assembler.Read(4)
		assembler.Print(3)
	})
}

func TestSimpleLoop(t *testing.T) {
	expectBf(t, ">+++++[<++>-]", func(assembler asm.Assembler) {
		assembler.Add(1, 5)
		assembler.OpenLoop(1)
		assembler.Add(0, 2)
		assembler.Add(1, -1)
		assembler.CloseLoop()
	})
}

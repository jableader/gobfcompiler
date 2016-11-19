package program

import (
	"asm"
	"program"
	"testing"
)

func TestCreateGetDestroy(t *testing.T) {
	assembler, ch := asm.New()
	defer close(ch)
	
	p := program.New(assembler)

	id := "Johnny"
	pt, err := p.CreatePt(&id, -1)
	if err != nil {
		t.Error(err)
	}

	rpt, ok := p.GetPt(id)
	if !ok {
		t.Error("The pointer was not retrieved")
	}

	if rpt.Loc != pt.Loc {
		t.Error("Pointer position is different")
	}

	p.DestroyPt(pt)
	_, okNow := p.GetPt(id)
	if okNow {
		t.Error("Pointer should have been destroyed")
	}
}

package program

import (
	"parse"
	"program"
	"testing"
)

type DummyAsm struct{}

func (d DummyAsm) Move(from, to int)               {}
func (d DummyAsm) Add(n int)                       {}
func (d DummyAsm) StartLoop()                      {}
func (d DummyAsm) EndLoop()                        {}
func (d DummyAsm) Print()                          {}
func (d DummyAsm) Err(msg string, expr parse.Expr) {}

func TestCreateGetDestroy(t *testing.T) {
	p := program.New(DummyAsm{})

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

func TestCreatePtDoesntDuplicate(t *testing.T) {
	p := program.New(DummyAsm{})

	name1, name2 := "NAME", "OTHER"
	pt1, ok1 := p.CreatePt(&name1, -1)
	pt2, ok2 := p.CreatePt(&name2, -1)

	if ok1 != nil || ok2 != nil {
		t.Fatalf("A create was unsuccessful. \n%v\n%v", ok1, ok2)
	}

	if pt1.Loc == pt2.Loc {
		t.Errorf("Both pointers have the same location of %v", pt1.Loc)
	}
}

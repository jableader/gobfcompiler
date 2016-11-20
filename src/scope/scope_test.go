package scope_test

import (
	"scope"
	"testing"
)

func TestScopeDefineGetUndefineGet(t *testing.T) {
	scope := scope.New()
	name := "VarName"

	_, exists := scope.Get(name)
	if exists {
		t.Error("When we haven't defined a variable Get should return false")
	}

	val := 69
	vDefine, err := scope.Define(&name, val)
	if err != nil {
		t.Errorf("Unexpected error from Define: %v", err)
	}

	vGet, exists := scope.Get(name)
	if !exists {
		t.Error("Get returned false after we defined a var")
	}
	if vGet != vDefine {
		t.Error("Get returned a different variable")
	}

	err = scope.Undefine(vGet)
	if err != nil {
		t.Errorf("Unexpected error from Undefine: %v", err)
	}
}

func TestScopeDifferentLevels(t *testing.T) {
	sc := scope.New()

	globalName, globalValue := "GLOB", 69
	localName, localValue := "LOCAL", 96

	if _, err := sc.Define(&globalName, globalValue); err != nil {
		t.Errorf("Unexpected error from Define: %v", err)
	}

	sc = sc.Enter()
	if _, err := sc.Define(&localName, localValue); err != nil {
		t.Errorf("Unexpected error from Define: %v", err)
	}

	if _, exists := sc.Get(globalName); !exists {
		t.Error("Couldnt find parent value from local scope")
	}

	if _, exists := sc.Get(localName); !exists {
		t.Error("Couldnt find local value from local scope")
	}

	sc = sc.Exit()
	if _, exists := sc.Get(localName); exists {
		t.Error("Found local value in global scope")
	}

	if _, exists := sc.Get(globalName); !exists {
		t.Error("Couldnt find global after exiting scope")
	}
}

func TestScopeNamePointers(t *testing.T) {
	ids := [2]string{"First", "Second"}
	sc := scope.New()
	for _, id := range ids {
		_, err := sc.Define(&id, 69)
		if err != nil {
			t.Errorf("Shouldnt be an issue. %v", err)
		}
	}
}

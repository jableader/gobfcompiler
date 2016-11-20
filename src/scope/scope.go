package scope

import "errors"

type Variable struct {
	Name *string
	Value interface{}
}

type Scope struct {
	parent *Scope
	vars []Variable
}

var (
	ErrAlreadyDefined error = errors.New("Variable is already defined")
	ErrDoesNotExist error = errors.New("Variable does not exist at this scope")

	VarUndefined Variable = Variable{}
)

func New() *Scope {
	return &Scope{nil, make([]Variable, 0, 5)}
}

func (s *Scope) Define(id *string, value interface{}) (Variable, error) {
	if id != nil {
		if _, alreadyExists := s.getWithoutParents(*id); alreadyExists {
			return VarUndefined, ErrAlreadyDefined
		}
	}

	v := Variable{id, value}
	s.vars = append(s.vars, v)
	return v, nil
}

func (s *Scope) Get(id string) (Variable, bool) {
	for sc := s; sc != nil; sc = sc.parent {
		if val, found := sc.getWithoutParents(id); found {
			return val, true
		}
	}

	return VarUndefined, false
}

func (s *Scope) Undefine(v Variable) error {
	for i, value := range s.vars {
		if value == v {
			s.vars[i] = s.vars[len(s.vars) - 1]
			s.vars = s.vars[:len(s.vars) - 1]
			return nil
		}
	}

	return ErrDoesNotExist
}

func (s *Scope) Enter() *Scope {
	return &Scope{s, make([]Variable, 0, 5)}
}

func (s *Scope) Exit() *Scope {
	return s.parent
}

func (s *Scope) getWithoutParents(id string) (Variable, bool) {
	for _, v := range s.vars {
		if v.Name != nil && *v.Name == id {
			return v, true
		}
	}

	return VarUndefined, false
}

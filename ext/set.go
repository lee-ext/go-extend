package ext

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
)

type Set[E comparable] map[E]Unit

func Set_[E comparable](cap int) Set[E] {
	return make(map[E]Unit, cap)
}

func SetOf[E comparable](es ...E) Set[E] {
	s := make(map[E]Unit, len(es))
	for _, e := range es {
		s[e] = Unit{}
	}
	return s
}

func (s Set[E]) Foreach(fn func(E)) {
	for e := range s {
		fn(e)
	}
}

func (s Set[E]) Len() int {
	return len(s)
}

func (s Set[E]) Empty() bool {
	return len(s) == 0
}

func (s Set[E]) Insert(element E) {
	s[element] = Unit{}
}

func (s Set[E]) Remove(element E) {
	delete(s, element)
}

func (s Set[E]) Contains(element E) bool {
	_, b := s[element]
	return b
}

// Or Union |
func (s Set[E]) Or(other Set[E]) Set[E] {
	s0, s1 := s, other
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := maps.Clone(s0)
	for e := range s1 {
		s_[e] = Unit{}
	}
	return s_
}

// And Intersection &
func (s Set[E]) And(other Set[E]) Set[E] {
	s0, s1 := s, other
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := Set_[E](s1.Len() / 2)
	for e, u := range s1 {
		if s0.Contains(e) {
			s_[e] = u
		}
	}
	return s_
}

// Sub Difference -
func (s Set[E]) Sub(other Set[E]) Set[E] {
	if s.Len() < other.Len()*2 {
		s_ := Set_[E](s.Len() / 2)
		for e := range s {
			if !other.Contains(e) {
				s_[e] = Unit{}
			}
		}
		return s_
	} else {
		s_ := maps.Clone(s)
		for e := range other {
			s_.Remove(e)
		}
		return s_
	}
}

// Xor SymmetricDifference ^
func (s Set[E]) Xor(other Set[E]) Set[E] {
	s0, s1 := s, other
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := maps.Clone(s0)
	for e := range s1 {
		if s0.Contains(e) {
			s_.Remove(e)
		} else {
			s_[e] = Unit{}
		}
	}
	return s_
}

func (s Set[E]) ToVec() Vec[E] {
	v := Vec_[E](len(s))
	for e := range s {
		v.Append(e)
	}
	return v
}

func (s Set[E]) Clear() {
	clear(s)
}

func (s Set[E]) AppendSelf(element E) Set[E] {
	s[element] = Unit{}
	return s
}

func (s Set[E]) String() string {
	return fmt.Sprintf("set%v", s.ToVec())
}

// MarshalJSON returns m as the JSON encoding of m.
func (s Set[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToVec())
}

// UnmarshalJSON sets *m to a copy of data.
func (s *Set[E]) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}
	vec := new(Vec[E])
	err := json.Unmarshal(data, vec)
	if err == nil {
		*s = SetOf[E](*vec...)
	}
	return err
}

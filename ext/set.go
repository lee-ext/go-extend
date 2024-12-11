package ext

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Set Define a generic Set[E], based on the map
type Set[E comparable] map[E]Unit

// Set_ Create a Set[E] with a specified capacity
func Set_[E comparable](cap int) Set[E] {
	return make(map[E]Unit, cap)
}

// SetOf Create a Set[E] that contains the specified elements
func SetOf[E comparable](es ...E) Set[E] {
	s := make(map[E]Unit, len(es))
	for _, e := range es {
		s[e] = Unit{}
	}
	return s
}

// ForEach Traverse the Set[E]
func (s Set[E]) ForEach(fn func(E)) {
	for e := range s {
		fn(e)
	}
}

// Len Gets the number of elements in the Set[E]
func (s Set[E]) Len() int {
	return len(s)
}

// Empty Determine if the Set[E] is empty
func (s Set[E]) Empty() bool {
	return len(s) == 0
}

// Insert Inserts an element into the Set[E]
func (s Set[E]) Insert(element E) {
	s[element] = Unit{}
}

// Remove Removes an element from the Set[E]
func (s Set[E]) Remove(element E) {
	delete(s, element)
}

// Contains Determines whether the Set[E] contains an element
func (s Set[E]) Contains(element E) bool {
	_, ok := s[element]
	return ok
}

// Or Find the union of two Set[E]
func (s Set[E]) Or(other Set[E]) Set[E] {
	s0, s1 := s, other
	// Select a larger Set[E] as a baseline
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := maps.Clone(s0)
	for e := range s1 {
		s_[e] = Unit{}
	}
	return s_
}

// And Find the intersection of two Set[E]
func (s Set[E]) And(other Set[E]) Set[E] {
	s0, s1 := s, other
	// Choose a smaller Set[E] as the baseline for iteration
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := Set_[E]((s1.Len() + 1) / 2)
	for e := range s1 {
		if s0.Contains(e) {
			s_[e] = Unit{}
		}
	}
	return s_
}

// Sub Find the difference between two Set[E]
func (s Set[E]) Sub(other Set[E]) Set[E] {
	if s.Len() < other.Len()*2 {
		// When the Set[E] is small, check them one by one
		s_ := Set_[E]((s.Len() + 1) / 2)
		for e := range s {
			if !other.Contains(e) {
				s_[e] = Unit{}
			}
		}
		return s_
	} else {
		// Clone the current Set[E] and remove elements from another Set[E]
		s_ := maps.Clone(s)
		for e := range other {
			s_.Remove(e)
		}
		return s_
	}
}

// Xor Find the symmetrical difference between two Set[E]
func (s Set[E]) Xor(other Set[E]) Set[E] {
	s0, s1 := s, other
	// Select a larger Set[E] as a baseline
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := maps.Clone(s0)
	for e := range s1 {
		if s_.Contains(e) {
			s_.Remove(e)
		} else {
			s_[e] = Unit{}
		}
	}
	return s_
}

// ToVec Convert a Set[E] to a Vec[E]
func (s Set[E]) ToVec() Vec[E] {
	v := Vec_[E](len(s))
	for e := range s {
		v.Append(e)
	}
	return v
}

// Clear Empty the Set[E]
func (s Set[E]) Clear() {
	clear(s)
}

// AppendSelf Inserts an element into the Set[E] and returns self
func (s Set[E]) AppendSelf(element E) Set[E] {
	s[element] = Unit{}
	return s
}

// String Convert a Set[E] to a string
func (s Set[E]) String() string {
	return fmt.Sprintf("set%v", s.ToVec())
}

// MarshalJSON Implement JSON serialization
func (s Set[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToVec())
}

// UnmarshalJSON Implement JSON deserialization
func (s *Set[E]) UnmarshalJSON(data []byte) error {
	vec := new(Vec[E])
	err := json.Unmarshal(data, vec)
	if err == nil {
		*s = SetOf[E](*vec...)
	}
	return err
}

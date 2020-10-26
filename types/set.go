package types

var exists = struct{}{}

// Set denotes the set data structure
type Set struct {
	m map[string]struct{}
}

// NewSet initializes a new set
func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s
}

// Add inserts an element into the set
func (s *Set) Add(value string) {
	s.m[value] = exists
}

// Remove removes an element from the set
func (s *Set) Remove(value string) {
	delete(s.m, value)
}

// Contains checks if an element is present within the set or not
func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

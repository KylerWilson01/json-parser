package internal

// Stack is a simple implementation of a stack datastructure
type Stack[T any] struct {
	state []T
}

// NewStack creates a new stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{state: []T{}}
}

// Peek looks at the last inserted item in the stack
func (s *Stack[T]) Peek() T {
	if s.IsEmpty() {
		panic("empty stack")
	}
	return s.state[len(s.state)-1]
}

// Push inserts an item to the top of the stack
func (s *Stack[T]) Push(ts T) {
	s.state = append(s.state, ts)
}

// Pop removes and returns the last insereted item
func (s *Stack[T]) Pop() T {
	r := s.Peek()

	s.state = s.state[:len(s.state)-1]

	return r
}

// IsEmpty returns whether the stack is empty or not
func (s *Stack[T]) IsEmpty() bool {
	return len(s.state) == 0
}

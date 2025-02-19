package internal

type Stack[T any] struct {
	state []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{state: []T{}}
}

func (s *Stack[T]) Peek() *T {
	if len(s.state) == 0 {
		return nil
	}
	return &s.state[len(s.state)-1]
}

func (s *Stack[T]) Push(ts T) {
	s.state = append(s.state, ts)
}

func (s *Stack[T]) Pop() *T {
	r := s.Peek()

	s.state = s.state[:len(s.state)-1]

	return r
}

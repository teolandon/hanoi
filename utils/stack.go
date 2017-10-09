package utils

type Stack struct {
	top *node
}

type node struct {
	next  *node
	value interface{}
}

func (s Stack) Pop() interface{} {
	if s.top == nil {
		return nil
	}

	ret := s.top.value
	s.top = s.top.next
	return ret
}

func (s Stack) Push(val interface{}) {
	newNode := node{s.top, val}
	s.top = &newNode
}

func (s Stack) Current() interface{} {
	if s.top == nil {
		return nil
	}

	return s.top.value
}

func (s Stack) Clear() {
	s.top = nil
}

func NewStack() Stack {
	ret := Stack{nil}
	return ret
}

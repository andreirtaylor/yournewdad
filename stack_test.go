package kaa

import (
	"testing"
)

func TestEmptyStack(t *testing.T) {
	stack := new(Stack)

	if stack.Len() > 0 {
		t.Errorf("Should not start with a length > 0 found, ", stack.Len())
	}

	if stack.Pop() != nil {
		t.Errorf("should not be able to pop anything off the empty stack")
	}
}

func TestPushPopStack(t *testing.T) {
	stack := new(Stack)
	stack.Push(&Point{X: 5, Y: 6})
	stack.Push(&Point{X: 5, Y: 6})
	stack.Push(&Point{X: 5, Y: 6})

	for stack.Len() > 0 {
		// We have to do a type assertion because we get back a variable of type
		// interface{} while the underlying type is a string.
		val := stack.Pop()
		if val == nil {
			t.Errorf("should not be able to get nils while the stack still has items")
		}
	}
}

func TestNilPushStack(t *testing.T) {
	stack := new(Stack)
	stack.Push(nil)

	if stack.Len() > 0 {
		t.Errorf("Should not be able to push empty things on the stack ")
	}
}

package queue

import (
	"container/list"
)

type Comparable interface {
	Less(other Comparable) bool
}

type PriorityQueue[T Comparable] interface {
	Push(value T) error
	Get(index int32) *T
	Add(value T) 
	Pop() T
}

type orderedList[T Comparable] struct {
	*list.List
}

func NewPriorityQueue[T Comparable]() PriorityQueue[T] {
	list := list.New()
	return &orderedList[T]{list}
}

func (l *orderedList[T]) Push(value T) error {
	return nil
}

func (l orderedList[T]) Get(index int32) *T {
	return new(T)
}

func (l *orderedList[T]) Add(value T) {
}

func (l orderedList[T]) Pop() T {
	return *new(T)
}
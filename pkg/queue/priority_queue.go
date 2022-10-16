package queue

import (
	"container/list"
	"errors"
)

type Comparable interface {
	Less(other Comparable) bool
}

type PriorityQueue[T Comparable] interface {
	Push(value T) error
	Front() *QueueElement[T]
	Back() *QueueElement[T]
	Add(value T)
	Pop() T
}

type QueueElement[T Comparable] struct {
	ListElement *list.Element
}

func (e QueueElement[T]) Next() *QueueElement[T] {
	if e.ListElement.Next() == nil {
		return nil
	}
	return &QueueElement[T]{e.ListElement.Next()}
}

func (e QueueElement[T]) Prev() *QueueElement[T] {
	if e.ListElement.Prev() == nil {
		return nil
	}
	return &QueueElement[T]{e.ListElement.Prev()}
}

func (e QueueElement[T]) Get() T {
	return e.ListElement.Value.(T)
}

type orderedList[T Comparable] struct {
	list *list.List
}

func NewPriorityQueue[T Comparable]() PriorityQueue[T] {
	list := list.New()
	return &orderedList[T]{list}
}

func (l *orderedList[T]) Push(value T) error {
	back := l.list.Back()
	if (back == nil) || back.Value.(T).Less(value) {
		l.list.PushBack(value)
		return nil
	}

	return errors.New("Value is smaller than the last element of queue")
}

func (l orderedList[T]) Front() *QueueElement[T] {
	if l.list.Front() == nil {
		return nil
	}
	return &QueueElement[T]{l.list.Front()}
}

func (l orderedList[T]) Back() *QueueElement[T] {
	if l.list.Back() == nil {
		return nil
	}
	return &QueueElement[T]{l.list.Back()}
}

func (l *orderedList[T]) Add(value T) {
	for el := l.Front(); el != nil; el = el.Next() {
		if value.Less(el.Get()) {
			l.list.InsertBefore(value, el.ListElement)
			return
		}
	}
	l.list.PushBack(value)
}

func (l orderedList[T]) Pop() T {
	front := l.list.Front()
	return l.list.Remove(front).(T)
}

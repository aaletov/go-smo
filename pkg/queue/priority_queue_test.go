package queue

import (
	"testing"
)

type TestInt int32

func (t TestInt) Less(other Comparable) bool {
	return t < other.(TestInt)
}

func getElements() []TestInt {
	return []TestInt{1, 3, 5, 6, 10}
}

func TestEmptyAdd(t *testing.T) {
	queue := NewPriorityQueue[TestInt]()
	queue.Add(1)
	last := queue.Back().Get()
	if last != 1 {
		t.Fatalf("Expected %v, got %v", 1, last)
	}
}

func TestEmptyPush(t *testing.T) {
	queue := NewPriorityQueue[TestInt]()
	queue.Push(1)
	last := queue.Back().Get()
	if last != 1 {
		t.Fatalf("Expected %v, got %v", 1, last)
	}
}

func TestEmptyBack(t *testing.T) {
	queue := NewPriorityQueue[TestInt]()
	if queue.Back() != nil {
		t.Fatal("Back is not nil")
	}
}

func TestEmptyFront(t *testing.T) {
	queue := NewPriorityQueue[TestInt]()
	if queue.Front() != nil {
		t.Fatal("Front is not nil")
	}
}

func getQueue() PriorityQueue[TestInt] {
	queue := NewPriorityQueue[TestInt]()
	for _, value := range getElements() {
		queue.Add(value)
	}
	return queue
}

func TestBack(t *testing.T) {
	queue := getQueue()
	elements := getElements()
	queueBackValue := queue.Back().Get()
	assertBackValue := elements[len(elements)-1]
	if queueBackValue != assertBackValue {
		t.Fatalf("queue.Back().Get() returned %v instead of %v", queueBackValue, assertBackValue)
	}
}

func TestFront(t *testing.T) {
	queue := getQueue()
	elements := getElements()
	queueBackValue := queue.Front().Get()
	assertBackValue := elements[0]
	if queueBackValue != assertBackValue {
		t.Fatalf("queue.Back().Get() returned %v instead of %v", queueBackValue, assertBackValue)
	}
}

func TestPush(t *testing.T) {
	queue := getQueue()
	backValue := queue.Back().Get()
	newValue := backValue + 1
	queue.Push(newValue)
	if queue.Back().Get() != newValue {
		t.Fatalf("queue.Back().Get() returned %v instead of %v", queue.Back().Get(), newValue)
	}
}

func TestAdd(t *testing.T) {
	var testValue TestInt = 4
	queue := getQueue()
	queue.Add(testValue)
	var el *QueueElement[TestInt]
	for el = queue.Front(); el != nil; el = el.Next() {
		if el.Get() == testValue {
			break
		}
	}
	if !(el.Prev().Get().Less(testValue) && !el.Next().Get().Less(testValue)) {
		t.Fatal("Value added incorrectly")
	}
}

func TestPop(t *testing.T) {
	queue := getQueue()
	if queue.Pop() != 1 {
		t.Fatal("Pop() returned not the minimal value")
	}
}

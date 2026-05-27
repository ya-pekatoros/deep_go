package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type QueueElement interface {
	int | int8 | int16 | int32 | int64
}

type CircularQueue[T QueueElement] struct {
	len    int
	head   int
	tail   int
	values []T
}

func NewCircularQueue[T QueueElement](size int, _ T) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.len == len(q.values) {
		// очередь полностью заполнена
		return false
	}
	q.values[q.tail] = value
	// в tail храним положение следующей вставки
	q.tail = (q.tail + 1) % len(q.values)
	q.len++

	return true
}
func (q *CircularQueue[T]) Pop() bool {
	if q.len == 0 {
		// очередь пуста
		return false
	}

	q.len--
	q.head = (q.head + 1) % len(q.values)

	return true
}

func (q *CircularQueue[T]) Front() (T, bool) {
	if q.len == 0 {
		return T(0), false
	}
	return q.values[q.head], true
}

func (q *CircularQueue[T]) Back() (T, bool) {
	if q.len == 0 {
		return T(0), false
	}
	backIndex := (q.tail - 1 + len(q.values)) % len(q.values)
	return q.values[backIndex], true
}

func (q *CircularQueue[T]) Empty() bool {
	return q.len == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.len == len(q.values)
}

func runTestQueue[T QueueElement](t *testing.T, _ T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize, T(0))

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	front, ok := queue.Front()
	assert.Equal(t, T(0), front)
	assert.False(t, ok)

	back, ok := queue.Back()
	assert.Equal(t, T(0), back)
	assert.False(t, ok)

	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(T(1)))
	assert.True(t, queue.Push(T(2)))
	assert.True(t, queue.Push(T(3)))
	assert.False(t, queue.Push(T(4)))

	assert.True(t, reflect.DeepEqual([]T{T(1), T(2), T(3)}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	front, ok = queue.Front()
	assert.Equal(t, T(1), front)
	assert.True(t, ok)

	back, ok = queue.Back()
	assert.Equal(t, T(3), back)
	assert.True(t, ok)

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(T(4)))

	assert.True(t, reflect.DeepEqual([]T{T(4), T(2), T(3)}, queue.values))

	front, ok = queue.Front()
	assert.Equal(t, T(2), front)
	assert.True(t, ok)

	back, ok = queue.Back()
	assert.Equal(t, T(4), back)
	assert.True(t, ok)

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

func TestCircularQueue(t *testing.T) {
	runTestQueue(t, int8(0))
	runTestQueue(t, int16(0))
	runTestQueue(t, int32(0))
	runTestQueue(t, int64(0))
	runTestQueue(t, int(0))
}

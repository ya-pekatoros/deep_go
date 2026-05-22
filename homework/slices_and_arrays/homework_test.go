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
	// тут мы для удобства вызова передаем тип элемента очереди
	// еще как альтернатива - передавать сразу элементы для инициализации
	// или же можно вынести отдельные инициализаторы под каждый тип элемента
	// например:
	// 	type Int64Queue = CircularQueue[int64]

	// func NewInt64Queue(size int) Int64Queue {
	// 	return CircularQueue[int64]{
	// 		values: make([]int64, size),
	// 	}
	// }
	return CircularQueue[T]{
		values: make([]T, size),
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.len == len(q.values) {
		// очередь полностью заполнена
		return false
	}

	// в tail храним положение последнего элемента
	// но если буфер еще не заполнялся
	// то вставлять нужно в начало
	// если буфер стал пустым - можно сделать также
	if q.len == 0 {
		q.tail = 0
	} else {
		q.tail = (q.tail + 1) % len(q.values)
	}
	q.values[q.tail] = value
	q.len++

	return true
}
func (q *CircularQueue[T]) Pop() bool {
	if q.len == 0 {
		// очередь пуста
		return false
	}

	q.len--
	q.values[q.head] = T(0)
	if q.len == 0 {
		// если буфер опустел переводим head в начало
		q.head = 0
	} else {
		q.head = (q.head + 1) % len(q.values)
	}

	return true
}

func (q *CircularQueue[T]) Front() (T, bool) {
	if q.len == 0 {
		// я тут отступлю от задания, мне не очень понятно как мы -1 будет отличать от реального значения элемента
		// очереди, так что я решил, что будем отдавать значение + bool
		// тем более что вернуть int когда мы работает с дженериками в одном кейсе и T - в другом
		// все равно надо нарушать контракт из задания
		return T(0), false
	}
	return q.values[q.head], true
}

func (q *CircularQueue[T]) Back() (T, bool) {
	if q.len == 0 {
		// я тут отступлю от задания, мне не очень понятно как мы -1 будет отличать от реального значения элемента
		// очереди, так что я решил, что будем отдавать значение + bool
		// тем более что вернуть int когда мы работает с дженериками в одном кейсе и T - в другом
		// все равно надо нарушать контракт из задания
		return T(0), false
	}
	// тут все удобно, тк в tail храним положение последнего элемента
	// если бы хранили место вставки нового, то пришлось бы делать
	// q.values[q.tail - 1] и обрабатывать кейс когда tail == 0
	return q.values[q.tail], true
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

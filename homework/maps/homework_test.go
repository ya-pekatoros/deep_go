package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap[K cmp.Ordered, V any] struct {
	key   K
	value V

	size int

	left  *OrderedMap[K, V]
	right *OrderedMap[K, V]
}

func NewOrderedMap[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.size == 0 {
		m.key = key
		m.value = value
		m.size++
		return
	}
	current := m
	for current != nil {
		if current.key == key {
			current.value = value
			return
		}
		if current.key > key {
			if current.left == nil {
				current.left = &OrderedMap[K, V]{key: key, value: value, size: 1}
				m.size++
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &OrderedMap[K, V]{key: key, value: value, size: 1}
				m.size++
				return
			}
			current = current.right
		}
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	if m.size == 0 {
		return
	}

	var parent *OrderedMap[K, V]
	current := m

	for current != nil {
		if current.key == key {
			break
		}

		parent = current
		if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}

	if current == nil {
		return
	}

	if current.left != nil && current.right != nil {
		leftMostParent := current
		leftMost := current.right

		// ищем минимальное значение в правом поддереве, им заменим удаляемый элемент
		for leftMost.left != nil {
			leftMostParent = leftMost
			leftMost = leftMost.left
		}

		current.key = leftMost.key
		current.value = leftMost.value

		// теперь будем удалять тот элемент, который ранее исполльзовался для замены удаляемого
		parent = leftMostParent
		current = leftMost
	}

	var child *OrderedMap[K, V]
	if current.left != nil {
		child = current.left
	} else {
		child = current.right
	}

	if parent == nil {
		if child == nil {
			// кейс мапы из одного элемента
			var zeroK K
			var zeroV V

			m.key = zeroK
			m.value = zeroV
			m.left = nil
			m.right = nil
		} else {
			// кейс удаления корня мапы с одним поддеревом
			m.key = child.key
			m.value = child.value
			m.left = child.left
			m.right = child.right
		}
	} else if parent.left == current {
		parent.left = child
	} else {
		parent.right = child
	}

	m.size--
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	if m.size == 0 {
		return false
	}

	current := m
	for current != nil {
		if current.key == key {
			return true
		}
		if current.key > key {
			current = current.left
		} else {
			current = current.right
		}
	}

	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	if m.size == 0 {
		return
	}

	if m.left != nil {
		m.left.ForEach(action)
	}
	action(m.key, m.value)
	if m.right != nil {
		m.right.ForEach(action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

func TestCircularQueueWithSlices(t *testing.T) {
	data := NewOrderedMap[string, []map[int]int]() // наввернем что-нибудь покруче
	assert.Zero(t, data.Size())

	data.Insert("D", []map[int]int{
		{10: 10},
		{5: 5},
	})
	data.Insert("C", []map[int]int{
		{10: 10},
		{5: 5},
	})
	data.Insert("G", []map[int]int{
		{10: 10},
		{5: 5},
	})
	data.Insert("A", []map[int]int{
		{10: 10},
		{5: 5},
	})
	data.Insert("B", []map[int]int{
		{10: 10},
		{5: 5},
	})
	data.Insert("E", []map[int]int{
		{10: 10},
		{5: 5},
	})
	data.Insert("F", []map[int]int{
		{10: 10},
		{5: 5},
	})

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains("A"))
	assert.True(t, data.Contains("F"))
	assert.False(t, data.Contains("Z"))
	assert.False(t, data.Contains("Y"))

	var keys []string
	expectedKeys := []string{"A", "B", "C", "D", "E", "F", "G"}
	data.ForEach(func(key string, _ []map[int]int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase("C")
	data.Erase("F")
	data.Erase("A")

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains("B"))
	assert.True(t, data.Contains("E"))
	assert.False(t, data.Contains("C"))
	assert.False(t, data.Contains("F"))

	keys = nil
	expectedKeys = []string{"B", "D", "E", "G"}
	data.ForEach(func(key string, _ []map[int]int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type orderedMapNode[K cmp.Ordered] struct {
	key   K
	value any

	left  *orderedMapNode[K]
	right *orderedMapNode[K]
}

type OrderedMap[K cmp.Ordered] struct {
	root *orderedMapNode[K]
	size int
}

func NewOrderedMap[K cmp.Ordered]() *OrderedMap[K] {
	return &OrderedMap[K]{}
}

func (m *OrderedMap[K]) Insert(key K, value any) {
	var sizeIncreased bool
	m.root, sizeIncreased = insertNode(m.root, key, value)
	if sizeIncreased {
		m.size++
	}
}

func insertNode[K cmp.Ordered](node *orderedMapNode[K], key K, value any) (*orderedMapNode[K], bool) {
	if node == nil {
		return &orderedMapNode[K]{key: key, value: value}, true
	}

	sizeIncreased := false

	if node.key == key {
		node.value = value
		return node, sizeIncreased
	}

	if key < node.key {
		node.left, sizeIncreased = insertNode(node.left, key, value)
		return node, sizeIncreased
	}

	node.right, sizeIncreased = insertNode(node.right, key, value)
	return node, sizeIncreased
}

func (m *OrderedMap[K]) Erase(key K) {
	var sizeDecreased bool
	m.root, sizeDecreased = eraseNode(m.root, key)
	if sizeDecreased {
		m.size--
	}
}

func eraseNode[K cmp.Ordered](node *orderedMapNode[K], key K) (*orderedMapNode[K], bool) {
	if node == nil {
		return nil, false
	}

	sizeDecreased := false

	if key < node.key {
		node.left, sizeDecreased = eraseNode(node.left, key)
		return node, sizeDecreased
	}

	if key > node.key {
		node.right, sizeDecreased = eraseNode(node.right, key)
		return node, sizeDecreased
	}

	if node.left == nil {
		return node.right, true
	}

	if node.right == nil {
		return node.left, true
	}

	replacementRight, replacement := detachMinNode(node.right)
	replacement.left = node.left
	replacement.right = replacementRight

	return replacement, true
}

func detachMinNode[K cmp.Ordered](node *orderedMapNode[K]) (*orderedMapNode[K], *orderedMapNode[K]) {
	if node.left == nil {
		right := node.right
		node.right = nil
		return right, node
	}

	var min *orderedMapNode[K]
	node.left, min = detachMinNode(node.left)
	return node, min
}

func (m *OrderedMap[K]) Contains(key K) bool {
	current := m.root
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

func (m *OrderedMap[K]) Size() int {
	return m.size
}

func (m *OrderedMap[K]) ForEach(action func(K, any)) {
	forEachNode(m.root, action)
}

func forEachNode[K cmp.Ordered](node *orderedMapNode[K], action func(K, any)) {
	if node == nil {
		return
	}

	forEachNode(node.left, action)
	action(node.key, node.value)
	forEachNode(node.right, action)
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap[int]()
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
	data.ForEach(func(key int, _ any) {
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
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Insert(6, 6)
	data.Erase(5)

	assert.False(t, data.Contains(5))
	assert.True(t, data.Contains(6))

	keys = nil
	expectedKeys = []int{4, 6, 10, 12}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

func TestOrderedMapEraseRootWithDeepSuccessor(t *testing.T) {
	data := NewOrderedMap[int]()
	for _, key := range []int{10, 5, 20, 15, 30, 12, 17} {
		data.Insert(key, key)
	}

	data.Erase(10)

	assert.Equal(t, 6, data.Size())
	assert.False(t, data.Contains(10))

	var keys []int
	expectedKeys := []int{5, 12, 15, 17, 20, 30}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

func TestOrderedMapWithSlices(t *testing.T) {
	data := NewOrderedMap[string]() // наввернем что-нибудь покруче
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
	data.ForEach(func(key string, _ any) {
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
	data.ForEach(func(key string, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

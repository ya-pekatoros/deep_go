package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Map[T, R any](data []T, action func(T) R) []R {
	if data == nil {
		return nil
	}
	if len(data) == 0 {
		return []R{}
	}

	result := make([]R, len(data))
	for i, item := range data {
		result[i] = action(item)
	}

	return result
}

func Filter[T any](data []T, action func(T) bool) []T {
	if len(data) == 0 {
		return data
	}

	result := make([]T, 0, len(data))

	for _, item := range data {
		if action(item) {
			result = append(result, item)
		}
	}

	return result
}

func Reduce[T any](data []T, initial T, action func(T, T) T) T {
	result := initial

	for _, item := range data {
		result = action(result, item)
	}

	return result
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) int
		result []int
	}{
		"nil numbers": {
			action: func(number int) int {
				return -number
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) int {
				return -number
			},
			result: []int{},
		},
		"inc numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number + 1
			},
			result: []int{2, 3, 4, 5, 6},
		},
		"double numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number * number
			},
			result: []int{1, 4, 9, 16, 25},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Map(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestMapStringsToLengths(t *testing.T) {
	data := []string{"go", "rust", "java"}

	result := Map(data, func(word string) int {
		return len(word)
	})

	assert.Equal(t, []int{2, 4, 4}, result)
	assert.Equal(t, []string{"go", "rust", "java"}, data)
}

func TestMapEmptyStringsToLengths(t *testing.T) {
	data := []string{}

	result := Map(data, func(word string) int {
		return len(word)
	})

	assert.Equal(t, []int{}, result)
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) bool
		result []int
	}{
		"nil numbers": {
			action: func(number int) bool {
				return number == 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) bool {
				return number == 1
			},
			result: []int{},
		},
		"even numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) bool {
				return number%2 == 0
			},
			result: []int{2, 4},
		},
		"positive numbers": {
			data: []int{-1, -2, 1, 2},
			action: func(number int) bool {
				return number > 0
			},
			result: []int{1, 2},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Filter(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestFilterStrings(t *testing.T) {
	data := []string{"go", "python", "c", "rust"}

	result := Filter(data, func(word string) bool {
		return len(word) > 3
	})

	assert.Equal(t, []string{"python", "rust"}, result)
	assert.Equal(t, []string{"go", "python", "c", "rust"}, data)
}

func TestReduce(t *testing.T) {
	tests := map[string]struct {
		initial int
		data    []int
		action  func(int, int) int
		result  int
	}{
		"nil numbers": {
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"sum of numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 15,
		},
		"sum of numbers with initial value": {
			initial: 10,
			data:    []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 25,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reduce(test.data, test.initial, test.action)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestReduceStrings(t *testing.T) {
	result := Reduce([]string{"go", "lang"}, "", func(result, word string) string {
		return result + word
	})

	assert.Equal(t, "golang", result)
}

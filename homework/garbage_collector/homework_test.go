package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type color int

const (
	white color = iota
	gray
	black
)

func Trace(stacks [][]uintptr) []uintptr {
	roots := make(map[uintptr]struct{})
	for _, stack := range stacks {
		for _, pointer := range stack {
			if pointer != 0 {
				roots[pointer] = struct{}{}
			}
		}
	}

	colors := make(map[uintptr]color, len(roots))
	for pointer := range roots {
		colors[pointer] = gray
	}

	var pointers []uintptr
	added := make(map[uintptr]struct{}, len(roots))

	var scan func(uintptr)
	scan = func(pointer uintptr) {
		if colors[pointer] == black {
			return
		}

		nextPointer := *(*uintptr)(unsafe.Pointer(pointer))
		if nextPointer != 0 && colors[nextPointer] == white {
			colors[nextPointer] = gray
			pointers = append(pointers, nextPointer)
			scan(nextPointer)
		}

		colors[pointer] = black
	}

	for _, stack := range stacks {
		for _, pointer := range stack {
			if pointer == 0 {
				continue
			}

			if _, ok := added[pointer]; ok {
				continue
			}

			added[pointer] = struct{}{}
			pointers = append(pointers, pointer)
			scan(pointer)
		}
	}

	return pointers
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}

func TestTraceSkipsDuplicatesAndCycles(t *testing.T) {
	var root uintptr
	var child uintptr

	root = uintptr(unsafe.Pointer(&child))
	child = uintptr(unsafe.Pointer(&root))

	stacks := [][]uintptr{
		{
			uintptr(unsafe.Pointer(&root)),
			uintptr(unsafe.Pointer(&root)),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&root)),
		uintptr(unsafe.Pointer(&child)),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}

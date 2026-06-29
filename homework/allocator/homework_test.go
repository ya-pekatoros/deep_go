package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	DefragmentWithSize(memory, pointers, 1)
}

func DefragmentWithSize(memory []byte, pointers []unsafe.Pointer, valueSize int) {
	if len(memory) == 0 {
		return
	}

	if valueSize <= 0 {
		return
	}

	if len(pointers) == 0 {
		clear(memory)
		return
	}

	writeIndex := 0

	for readIndex := 0; readIndex+valueSize <= len(memory); readIndex++ {
		currentPointer := unsafe.Pointer(&memory[readIndex])
		occupied := false

		for _, pointer := range pointers {
			if pointer == currentPointer {
				occupied = true
				break
			}
		}

		if !occupied {
			continue
		}

		copy(memory[writeIndex:writeIndex+valueSize], memory[readIndex:readIndex+valueSize])
		newPointer := unsafe.Pointer(&memory[writeIndex])

		for pointerIndex, pointer := range pointers {
			if pointer == currentPointer {
				pointers[pointerIndex] = newPointer
			}
		}

		writeIndex += valueSize
	}

	clear(memory[writeIndex:])
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationUsesPointersNotValues(t *testing.T) {
	var fragmentedMemory = []byte{0x10, 0x20, 0x30, 0x40}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
	}

	var defragmentedMemory = []byte{0x20, 0x40, 0x00, 0x00}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationKeepsPointerOrder(t *testing.T) {
	var fragmentedMemory = []byte{0x10, 0x20, 0x30, 0x40}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[3]),
		unsafe.Pointer(&fragmentedMemory[1]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[0]),
	}

	var defragmentedMemory = []byte{0x20, 0x40, 0x00, 0x00}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationUpdatesDuplicatePointers(t *testing.T) {
	var fragmentedMemory = []byte{0x10, 0x20, 0x30, 0x40}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[2]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[0]),
	}

	var defragmentedMemory = []byte{0x30, 0x00, 0x00, 0x00}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationWithSizedValues(t *testing.T) {
	tests := []struct {
		name                 string
		valueSize            int
		memory               []byte
		pointers             []int
		defragmentedMemory   []byte
		defragmentedPointers []int
	}{
		{
			name:                 "two byte values",
			valueSize:            2,
			memory:               []byte{0x00, 0x00, 0x11, 0x12, 0x00, 0x00, 0x00, 0x00, 0x21, 0x22, 0x00, 0x00},
			pointers:             []int{2, 8},
			defragmentedMemory:   []byte{0x11, 0x12, 0x21, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			defragmentedPointers: []int{0, 2},
		},
		{
			name:                 "four byte values",
			valueSize:            4,
			memory:               []byte{0x00, 0x00, 0x00, 0x00, 0x11, 0x12, 0x13, 0x14, 0x00, 0x00, 0x00, 0x00, 0x21, 0x22, 0x23, 0x24},
			pointers:             []int{4, 12},
			defragmentedMemory:   []byte{0x11, 0x12, 0x13, 0x14, 0x21, 0x22, 0x23, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			defragmentedPointers: []int{0, 4},
		},
		{
			name:                 "eight byte values",
			valueSize:            8,
			memory:               []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28},
			pointers:             []int{24, 8},
			defragmentedMemory:   []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			defragmentedPointers: []int{8, 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memory := append([]byte(nil), test.memory...)
			pointers := make([]unsafe.Pointer, len(test.pointers))

			for index, pointerAt := range test.pointers {
				pointers[index] = unsafe.Pointer(&memory[pointerAt])
			}

			DefragmentWithSize(memory, pointers, test.valueSize)

			assert.True(t, reflect.DeepEqual(test.defragmentedMemory, memory))

			for index, pointerAt := range test.defragmentedPointers {
				assert.Equal(t, unsafe.Pointer(&memory[pointerAt]), pointers[index])
			}
		})
	}
}

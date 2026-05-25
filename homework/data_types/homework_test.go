package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndian[T uint16 | uint32 | uint64](number T) T {
	bytesCount := int(unsafe.Sizeof(number))

	var result T

	for i := 0; i < bytesCount; i++ {
		b := (number >> (i * 8)) & T(0xFF)
		result |= b << ((bytesCount - 1 - i) * 8)
	}

	return result
}

func runTests[T uint16 | uint32 | uint64](t *testing.T, tests []struct {
	name   string
	number T
	result T
}) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.result, ToLittleEndian(tc.number))
		})
	}
}

func TestСonversion(t *testing.T) {
	tests32 := []struct {
		name   string
		number uint32
		result uint32
	}{
		{"test case #1", 0x00000000, 0x00000000},
		{"test case #2", 0xFFFFFFFF, 0xFFFFFFFF},
		{"test case #3", 0x00FF00FF, 0xFF00FF00},
		{"test case #4", 0x0000FFFF, 0xFFFF0000},
		{"test case #5", 0x01020304, 0x04030201},
	}

	tests16 := []struct {
		name   string
		number uint16
		result uint16
	}{
		{"test case #6", 0x0000, 0x0000},
		{"test case #7", 0xFFFF, 0xFFFF},
		{"test case #8", 0x00FF, 0xFF00},
	}

	tests64 := []struct {
		name   string
		number uint64
		result uint64
	}{
		{"test case #9", 0x0000000000000000, 0x0000000000000000},
		{"test case #10", 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
		{"test case #11", 0x00000000FFFFFFFF, 0xFFFFFFFF00000000},
		{"test case #12", 0x0000FFFF0000FFFF, 0xFFFF0000FFFF0000},
		{"test case #13", 0x0102030405060708, 0x0807060504030201},
	}

	runTests(t, tests32)
	runTests(t, tests16)
	runTests(t, tests64)
}

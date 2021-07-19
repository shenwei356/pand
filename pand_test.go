package pand

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/grailbio/base/simd"
)

var data [][2][]byte

var data2 [][2][]byte

func init() {
	// for test
	sizes := []int{0, 1, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33, 47, 48, 49, 71, 72, 73, 127, 128, 129,
		1 << 7, 1 << 10, 1 << 16, (1 << 16) + 1, (1 << 16) - 1}
	// sizes := []int{33}

	data = make([][2][]byte, len(sizes))
	for i, s := range sizes {
		data[i] = [2][]byte{randByteSlice(s), randByteSlice(s)}
	}

	// for benchmark
	sizes2 := []int{8, 16, 32, 128, 512, 1 << 16}

	data2 = make([][2][]byte, len(sizes2))
	for i, s := range sizes2 {
		data2[i] = [2][]byte{randByteSlice(s), randByteSlice(s)}
	}
}

func randByteSlice(n int) []byte {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = byte(rand.Intn(255))
	}
	return s
}

func TestGoAsm(t *testing.T) {
	x := []byte{0b01, 0b11, 0b101} // 1, 3, 5
	y := []byte{0b10, 0b10, 0b111} // 2, 2, 7

	r := make([]byte, len(x))
	And(r, x, y)
	fmt.Println(r) // [0 2 5]

	AndInplace(x, y)
	fmt.Println(x) // [0 2 5]

	a := []byte{
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
		3, 3, 3,
	}
	b := []byte{
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2,
	}

	// fmt.Println(len(a), cap(a), a)
	AndInplace(a, b)
	// fmt.Println(AND(a, b))
	// fmt.Println(len(a), cap(a), a)
}

func TestAll(t *testing.T) {
	for i := range data {
		x := data[i][0]
		size := len(x)
		y := data[i][1]

		fmt.Printf("testing size: %d ... ", size)

		and1 := make([]byte, size)
		copy(and1, x)
		andInplaceGeneric0(and1, y)

		and2 := make([]byte, size)
		copy(and2, x)
		andInplaceGeneric(and2, y)

		if !bytes.Equal(and1, and2) {
			t.Errorf("oh no")
		}

		and3 := make([]byte, size)
		copy(and3, x)
		simd.AndUnsafeInplace(and3, y)

		if !bytes.Equal(and1, and3) {
			t.Errorf("oh no, grailbio error")
		}

		and4 := make([]byte, size)
		And(and4, x, y)
		if !bytes.Equal(and1, and4) {
			t.Errorf("oh no, goasmAnd error")
			return
		}

		and5 := make([]byte, size)
		copy(and5, x)
		AndInplace(and5, y)

		if !bytes.Equal(and1, and5) {
			t.Errorf("oh no, goasmAndInplace error")
			return
		}

		fmt.Println("ok")
	}
}

package pand

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/grailbio/base/simd"
	"github.com/shenwei356/util/bytesize"
)

var data [][2][]byte

var data2 [][2][]byte

func init() {
	// for test
	sizes := []int{0, 1, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33, 47, 48, 49, 71, 72, 73, 127, 128, 129, 1 << 7, 1 << 10, 1 << 16}
	// sizes := []int{33}

	data = make([][2][]byte, len(sizes))
	for i, s := range sizes {
		data[i] = [2][]byte{randByteSlice(s), randByteSlice(s)}
	}

	// for benchmark
	sizes2 := []int{8, 16, 32, 128, 1 << 8, 1 << 10, 1 << 14}

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

func TestAndInplace(t *testing.T) {
	for i := range data {
		x := data[i][0]
		size := len(x)
		y := data[i][1]

		fmt.Printf("tesing size: %d ... ", size)

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
		copy(and4, x)
		AndInplace(and4, y)

		if !bytes.Equal(and1, and4) {
			t.Errorf("oh no, goasm error")
			// fmt.Println(len(and1), and1)
			// fmt.Println(len(and4), and4)
			return
		}

		fmt.Println("ok")
	}
}

func BenchmarkLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				andInplaceGeneric0(and, y)
			}
		})
	}
}

func BenchmarkUnrollLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				andInplaceGeneric(and, y)
			}
		})
	}
}

func BenchmarkGrailbio(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				simd.AndUnsafeInplace(and, y)
			}
		})
	}
}

func BenchmarkGoAsm(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				AndInplace(and, y)
			}
		})
	}
}

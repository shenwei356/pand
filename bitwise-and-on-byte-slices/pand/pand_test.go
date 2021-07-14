package pand

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/grailbio/base/simd"
	"github.com/shenwei356/util/bytesize"
)

var data [][2][]byte

var data2 [][2][]byte

func init() {
	// for test
	sizes := []int{0, 1, 7, 8, 9, 31, 32, 127, 128, 129, 1 << 7, 1 << 10, 1 << 16}

	data = make([][2][]byte, len(sizes))
	for i, s := range sizes {
		data[i] = [2][]byte{randByteSlice(s), randByteSlice(s)}
	}

	// for benchmark
	sizes2 := []int{7, 32, 128, 1 << 8, 1 << 10, 1 << 16}

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

func TestAll(t *testing.T) {
	for i := range data {
		x := data[i][0]
		size := len(x)
		y := data[i][1]

		and1 := make([]byte, size)
		copy(and1, x)
		loop(and1, y)

		and2 := make([]byte, size)
		copy(and2, x)
		unroll(and2, y)

		if !bytes.Equal(and1, and2) {
			t.Errorf("oh no")
		}

		and3 := make([]byte, size)
		copy(and3, x)
		grailbio(and3, y)

		if !bytes.Equal(and1, and3) {
			t.Errorf("oh no, grailbio error")
		}

		and4 := make([]byte, size)
		copy(and4, x)
		goasm(and4, y)

		if !bytes.Equal(and1, and4) {
			t.Errorf("oh no,  error")
			// fmt.Println(len(and1), and1)
			// fmt.Println(len(and4), and4)
			return
		}
	}
}

func loop(x, y []byte) {
	for k, b := range y {
		x[k] &= b
	}
}

func unroll(x, y []byte) {
	k := 0
	for len(y) >= 8 { // unroll loop
		x[k] &= y[0]
		k++
		x[k] &= y[1]
		k++
		x[k] &= y[2]
		k++
		x[k] &= y[3]
		k++
		x[k] &= y[4]
		k++
		x[k] &= y[5]
		k++
		x[k] &= y[6]
		k++
		x[k] &= y[7]
		k++
		y = y[8:]
	}
	for _, b := range y {
		x[k] &= b
		k++
	}
}

func grailbio(x, y []byte) {
	simd.AndUnsafeInplace(x, y)
}

func goasm(x, y []byte) {
	PAND(x, y)
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

				loop(and, y)
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

				unroll(and, y)
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

				grailbio(and, y)
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

				goasm(and, y)
			}
		})
	}
}

package bitwiseandonbyteslices

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/grailbio/base/simd"
	"github.com/shenwei356/bench/bitwise-and-on-byte-slices/pand"
	"github.com/shenwei356/util/bytesize"
)

var data2 [][2][]byte

func init() {
	sizes := []int{1 << 7, 1 << 10, 1 << 16}

	data2 = make([][2][]byte, len(sizes))
	for i, s := range sizes {
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

func TestAll2(t *testing.T) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]

		and1 := make([]byte, size)
		copy(and1, x)
		loop2(and1, y)

		and2 := make([]byte, size)
		copy(and2, x)
		unroll2(and2, y)

		if !bytes.Equal(and1, and2) {
			t.Errorf("oh no")
		}

		and3 := make([]byte, size)
		copy(and3, x)
		grailbio2(and3, y)

		if !bytes.Equal(and1, and3) {
			t.Errorf("oh no, grailbio error")
		}

		and4 := make([]byte, size)
		copy(and4, x)
		goasmavx2(and4, y)

		if !bytes.Equal(and1, and4) {
			t.Errorf("oh no, avx2 error")
		}
	}
}

func loop2(x, y []byte) {
	for k, b := range y {
		x[k] &= b
	}
}

func unroll2(x, y []byte) {
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

func grailbio2(x, y []byte) {
	simd.AndUnsafeInplace(x, y)
}

func goasmavx2(x, y []byte) {
	pand.PAND(x, y)
}

func Benchmark2Loop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				loop2(and, y)
			}
		})
	}
}

func Benchmark2UnrollLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				unroll2(and, y)
			}
		})
	}
}

func Benchmark2Grailbio(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				grailbio2(and, y)
			}
		})
	}
}

func Benchmark2GoAsmAvx2(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				copy(and, x)

				goasmavx2(and, y)
			}
		})
	}
}

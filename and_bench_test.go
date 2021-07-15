package pand

import (
	"testing"

	"github.com/grailbio/base/simd"
	"github.com/shenwei356/util/bytesize"
)

func BenchmarkAndLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				andGeneric0(and, x, y)
			}
		})
	}
}

func BenchmarkAndUnrollLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				andGeneric(and, x, y)
			}
		})
	}
}

func BenchmarkAndGrailbio(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				simd.AndUnsafe(and, x, y)
			}
		})
	}
}

func BenchmarkAndGoAsm(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		and := make([]byte, size)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				AndUnsafe(and, x, y)
			}
		})
	}
}

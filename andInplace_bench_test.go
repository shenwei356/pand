package pand

import (
	"testing"

	"github.com/grailbio/base/simd"
	"github.com/shenwei356/util/bytesize"
)

func BenchmarkAndInplaceLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]

		and := make([]byte, size)
		copy(and, x)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				andInplaceGeneric0(and, y)
			}
		})
	}
}

func BenchmarkAndInplaceUnrollLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]

		and := make([]byte, size)
		copy(and, x)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				andInplaceGeneric(and, y)
			}
		})
	}
}

func BenchmarkAndInplaceGrailbio(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]

		and := make([]byte, size)
		copy(and, x)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				simd.AndUnsafeInplace(and, y)
			}
		})
	}
}

func BenchmarkAndInplaceGoAsm(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]

		and := make([]byte, size)
		copy(and, x)
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				AndUnsafeInplace(and, y)
			}
		})
	}
}

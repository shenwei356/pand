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
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				andInplaceGeneric0(x, y)
			}
		})
	}
}

func BenchmarkAndInplaceUnrollLoop(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				andInplaceGeneric(x, y)
			}
		})
	}
}

func BenchmarkAndInplaceGrailbio(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				simd.AndUnsafeInplace(x, y)
			}
		})
	}
}

func BenchmarkAndInplaceGoAsm(b *testing.B) {
	for i := range data2 {
		size := len(data2[i][0])
		x := data2[i][0]
		y := data2[i][1]
		b.Run(bytesize.ByteSize(size).String(), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				AndUnsafeInplace(x, y)
			}
		})
	}
}

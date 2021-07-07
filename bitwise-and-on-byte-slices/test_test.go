package bitwiseandonbyteslices

import (
	"bytes"
	"math/rand"
	"testing"
)

var matrix [][]byte
var positions [][]int

func init() {
	ncols := 256
	nrows := 1024

	rand.Seed(1)

	// generate matrix
	matrix = make([][]byte, nrows)
	for i := 0; i < nrows; i++ {
		row := make([]byte, ncols)
		for j := 0; j < ncols; j++ {
			row[j] = byte(rand.Intn(256))
		}
		matrix[i] = row
	}

	npos := 128
	// positions for extract subset (2-4 rows) of matrix
	positions = make([][]int, npos)
	for i := 0; i < npos; i++ {
		pos := make([]int, 0, 8)
		for j := 0; j < rand.Intn(3)+2; j++ { // [2,4]
			pos = append(pos, rand.Intn(nrows))
		}
		positions[i] = pos
	}
}

var AND []byte

func Test(t *testing.T) {
	data := make([][]byte, 0, 8)

	var p int
	for _, pos := range positions {
		data = data[:0] // subset of matrix
		for _, p = range pos {
			data = append(data, matrix[p])
		}

		if !bytes.Equal(loop(data), unroll(data)) {
			t.Errorf("error")
		}
	}
}

func BenchmarkLoop(b *testing.B) {
	var and []byte
	data := make([][]byte, 0, 8)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p int
		for _, pos := range positions {
			data = data[:0] // subset of matrix
			for _, p = range pos {
				data = append(data, matrix[p])
			}

			and = loop(data)
		}
	}
	AND = and
}

func BenchmarkUnrollLoop(b *testing.B) {
	var and []byte
	data := make([][]byte, 0, 8)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p int
		for _, pos := range positions {
			data = data[:0] // subset of matrix
			for _, p = range pos {
				data = append(data, matrix[p])
			}

			and = unroll(data)
		}
	}
	AND = and
}

// the input data is a byte matrix,
// this function proceeds bitwise AND operation on every column.
func loop(data [][]byte) []byte {
	// if len(data) < 2 {
	// 	panic("input matrix should have >=2 rows")
	// }
	// if len(data[0]) == 0 {
	// 	panic("input matrix should have >0 columns")
	// }

	and := make([]byte, len(data[0]))
	copy(and, data[0]) // copy the first row

	var i int
	var b byte
	for _, row := range data[1:] { // left rows
		// if len(row) != len(data[0]) {
		// 	panic("column lengths should be consistant")
		// }

		for i, b = range row { // every byte at a row
			and[i] &= b // bitwise AND, update the value
		}
	}
	return and
}

func unroll(data [][]byte) []byte {
	// if len(data) < 2 {
	// 	panic("input matrix should have >=2 rows")
	// }
	// if len(data[0]) == 0 {
	// 	panic("input matrix should have >0 columns")
	// }

	and := make([]byte, len(data[0]))
	copy(and, data[0])

	var i int
	var b byte
	for _, row := range data[1:] {
		// if len(row) != len(data[0]) {
		// 	panic("column lengths should be consistant")
		// }

		i = 0
		for len(row) >= 8 { // unroll loop
			and[i] &= row[0]
			i++
			and[i] &= row[1]
			i++
			and[i] &= row[2]
			i++
			and[i] &= row[3]
			i++
			and[i] &= row[4]
			i++
			and[i] &= row[5]
			i++
			and[i] &= row[6]
			i++
			and[i] &= row[7]
			i++
			row = row[8:]
		}
		for _, b = range row {
			and[i] &= b
			i++
		}
	}
	return and
}

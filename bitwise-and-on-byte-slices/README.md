## Bitwise AND operation on byte slices

The core code:

```
// the input data is a byte matrix,
// this function proceeds bitwise AND operation on every column.
func loop(data [][]byte) []byte {
	and := make([]byte, len(data[0]))
	copy(and, data[0]) // copy the first row

	var i int
	var b byte
	for _, row := range data[1:] { // left rows
		for i, b = range row { // every byte at a row
			and[i] &= b // bitwise AND, update the value
		}
	}
	return and
}
```

Benchmark result

```
goos: linux
goarch: amd64
pkg: github.com/shenwei356/bench/bitwise-and-on-byte-slices
cpu: AMD Ryzen 7 2700X Eight-Core Processor         
BenchmarkLoop-16                   14467             84265 ns/op
BenchmarkUnrollLoop-16             17668             67550 ns/op

```

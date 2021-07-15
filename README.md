# Bitwise AND operation on two []byte

[![GoDoc](https://godoc.org/github.com/shenwei356/pand?status.svg)](https://pkg.go.dev/github.com/shenwei356/pand)
[![Go Report Card](https://goreportcard.com/badge/github.com/shenwei356/pand)](https://goreportcard.com/report/github.com/shenwei356/pand)

## Introduction

This package provides a vectorised function which performs
bitwise AND operation on every pair of elements in two byte-slices.

The generic Go code is below, and unrolling the `for` loop could
increase the speed. 


```
func AND(x, y []byte) {
	for k, b := range y {
		x[k] &= b
	}
}
```

[grailbio/base](https://github.com/grailbio/base/blob/master/simd/and_amd64.go)
provides a faster pure Go implementation `AndUnsafeInplace` ultlizing `unsafe` package.

The solution here (`pand.AND`) is faster than `AndUnsafeInplace` for `[]byte`  with 32 or more elements.

see [benchmark](#benchmark).

## Getting started

```
go get -u github.com/shenwei356/pand

x := []byte{0b01, 0b11, 0b101} // 1, 3, 5
y := []byte{0b10, 0b10, 0b111} // 2, 2, 7

r := make([]byte, len(x))
pand.And(r, x, y)
fmt.Println(r) // [0 2 5]

pand.AndInplace(x, y)
fmt.Println(x) // [0 2 5]

```

Generate Go assembly code

```
go run asm-AndInplaceAvx.go -out andInplaceAvx_amd64.s -stubs andInplaceAvx.go 

go run asm-AndAvx.go -out andAvx_amd64.s -stubs andAvx.go 

go test .

```

## Benchmark

```
go test . -bench=Benchmark* | tee t.txt

cat t.txt \
    | grep Bench \
    | sed -r 's/\s\s+/\t/g' \
    | csvtk cut -Ht -f 1,3 \
    | csvtk add-header -t -n test,time \
    | csvtk mutate -t -n data-size -p "/(.+)-" \
    | csvtk replace -t -p "(.+)/.+" -r "\$1" \
    | csvtk cut -t -f test,data-size,time \
    | csvtk sort -t -k data-size:N -k time:N \
    | csvtk pretty -t -s "        "
rm t.txt

test                       data-size        time
-------------------        ---------        -----------
BenchmarkGrailbio          8.00_B           4.654 ns/op
BenchmarkGoAsm             8.00_B           4.824 ns/op
BenchmarkUnrollLoop        8.00_B           6.851 ns/op
BenchmarkLoop              8.00_B           8.683 ns/op

BenchmarkGrailbio          16.00_B          5.363 ns/op
BenchmarkGoAsm             16.00_B          6.369 ns/op
BenchmarkUnrollLoop        16.00_B          10.47 ns/op
BenchmarkLoop              16.00_B          13.48 ns/op

BenchmarkGoAsm             32.00_B          6.079 ns/op
BenchmarkGrailbio          32.00_B          6.497 ns/op
BenchmarkUnrollLoop        32.00_B          17.46 ns/op
BenchmarkLoop              32.00_B          21.09 ns/op

BenchmarkGoAsm             128.00_B         10.52 ns/op
BenchmarkGrailbio          128.00_B         14.40 ns/op
BenchmarkUnrollLoop        128.00_B         56.97 ns/op
BenchmarkLoop              128.00_B         80.12 ns/op

BenchmarkGoAsm             256.00_B         15.48 ns/op
BenchmarkGrailbio          256.00_B         23.76 ns/op
BenchmarkUnrollLoop        256.00_B         110.8 ns/op
BenchmarkLoop              256.00_B         147.5 ns/op

BenchmarkGoAsm             1.00_KB          47.16 ns/op
BenchmarkGrailbio          1.00_KB          87.75 ns/op
BenchmarkUnrollLoop        1.00_KB          443.1 ns/op
BenchmarkLoop              1.00_KB          540.5 ns/op

BenchmarkGoAsm             16.00_KB         751.6 ns/op
BenchmarkGrailbio          16.00_KB         1342 ns/op
BenchmarkUnrollLoop        16.00_KB         7007 ns/op
BenchmarkLoop              16.00_KB         8623 ns/op

```

## Credits

- Go assembly code was generated with [avo](https://github.com/mmcloughlin/avo).
- [Peter Cordes](https://stackoverflow.com/users/224132/peter-cordes)
  provided valuable suggestions on the Assembly language
  in this [post](https://stackoverflow.com/questions/68280854/).
- We copied and edited dispatching code from [pospop](https://github.com/clausecker/pospop).

## License

[MIT License](https://github.com/shenwei356/pand/blob/master/LICENSE)
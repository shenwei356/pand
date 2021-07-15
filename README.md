# Bitwise AND operation on two []byte

[![GoDoc](https://godoc.org/github.com/shenwei356/pand?status.svg)](https://pkg.go.dev/github.com/shenwei356/pand)
[![Go Report Card](https://goreportcard.com/badge/github.com/shenwei356/pand)](https://goreportcard.com/report/github.com/shenwei356/pand)

## Introduction

This package provides a vectorised function which performs
bitwise AND operation on every pair of elements in two byte-slices.

The generic Go code is below, and unrolling the `for` loop could
increase the speed. 


```
func AndInplace(x, y []byte) {
	for k, b := range y {
		x[k] &= b
	}
}

func And(r, x, y []byte) {
	for i, b := range x {
		r[i] = b & y[i]
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

## Benchmark

```
go test . -bench=BenchmarkAnd* | tee t.txt

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

test                                 data-size        time
-----------------------------        ---------        -----------
BenchmarkAndInplaceGrailbio          8.00_B           3.013 ns/op
BenchmarkAndInplaceGoAsm             8.00_B           4.210 ns/op
BenchmarkAndInplaceLoop              8.00_B           5.096 ns/op
BenchmarkAndInplaceUnrollLoop        8.00_B           5.237 ns/op

BenchmarkAndInplaceGrailbio          16.00_B          3.825 ns/op
BenchmarkAndInplaceGoAsm             16.00_B          4.648 ns/op
BenchmarkAndInplaceLoop              16.00_B          8.092 ns/op
BenchmarkAndInplaceUnrollLoop        16.00_B          8.363 ns/op

BenchmarkAndInplaceGrailbio          32.00_B          4.772 ns/op
BenchmarkAndInplaceGoAsm             32.00_B          5.161 ns/op
BenchmarkAndInplaceLoop              32.00_B          14.21 ns/op
BenchmarkAndInplaceUnrollLoop        32.00_B          14.55 ns/op

BenchmarkAndInplaceGoAsm             128.00_B         9.315 ns/op
BenchmarkAndInplaceGrailbio          128.00_B         11.18 ns/op
BenchmarkAndInplaceUnrollLoop        128.00_B         51.11 ns/op
BenchmarkAndInplaceLoop              128.00_B         55.52 ns/op

BenchmarkAndInplaceGoAsm             512.00_B         16.50 ns/op
BenchmarkAndInplaceGrailbio          512.00_B         39.31 ns/op
BenchmarkAndInplaceLoop              512.00_B         206.4 ns/op
BenchmarkAndInplaceUnrollLoop        512.00_B         211.3 ns/op

BenchmarkAndInplaceGoAsm             64.00_KB         1589 ns/op
BenchmarkAndInplaceGrailbio          64.00_KB         3914 ns/op
BenchmarkAndInplaceUnrollLoop        64.00_KB         24819 ns/op
BenchmarkAndInplaceLoop              64.00_KB         25852 ns/op

```

### For developers

Generate Go assembly code with [avo](https://github.com/mmcloughlin/avo)

```
go run asm-AndInplaceAvx2.go -out andInplaceAvx2_amd64.s -stubs andInplaceAvx2.go
go run asm-AndAvx2.go -out andAvx2_amd64.s -stubs andAvx2.go

go test .
```

***Attention: since avo does not support Avx512 yet, we need to manuall edit
`andAvx512_amd64.s` and `andInplaceAvx512_amd64.s`***

```
go run asm-AndInplaceAvx512.go -out andInplaceAvx512_amd64.s -stubs andInplaceAvx512.go
go run asm-AndAvx512.go -out andAvx512_amd64.s -stubs andAvx512.go

```

For `andInplaceAvx512_amd64.s`, Change

```
loop64:
	// compute x & y, and save value to x
	VMOVDQU (AX), Y0
	VPAND   (DX), T0, T0
	VMOVDQU T0, (AX)
```

to

```
loop64:
	// compute x & y, and save value to x
	VMOVDQU64 (AX), Z0
	VPANDQ   (DX), Z0, Z0
	VMOVDQU64 Z0, (AX)
```

For `andAvx512_amd64.s`, Change

```
loop64:
	// compute x & y, and save value to x
	VMOVDQU (CX), Y0
	VPAND   (BX), Y0, Y0
	VMOVDQU Y0, (CX)
```

to

```
loop64:
	// compute x & y, and save value to x
	VMOVDQU64 (CX), Z0
	VPANDQ   (BX), Z0, Z0
	VMOVDQU64 Z0, (AX)
```


## Credits

- Go assembly code was generated with [avo](https://github.com/mmcloughlin/avo).
- [Peter Cordes](https://stackoverflow.com/users/224132/peter-cordes)
  provided valuable suggestions on the Assembly language
  in this [post](https://stackoverflow.com/questions/68280854/).
- We copied and edited dispatching code from [pospop](https://github.com/clausecker/pospop).

## License

[MIT License](https://github.com/shenwei356/pand/blob/master/LICENSE)
# Bitwise AND operation on two []byte

[![GoDoc](https://godoc.org/github.com/shenwei356/pand?status.svg)](https://pkg.go.dev/github.com/shenwei356/pand)
[![Go Report Card](https://goreportcard.com/badge/github.com/shenwei356/pand)](https://goreportcard.com/report/github.com/shenwei356/pand)

## Introduction

This package provides a vectorised function which performs
bitwise AND operation on all pairs of elements in two byte-slices.
It detects CPU instruction set and chooses the available best one (AVX512, AVX2, SSE2).

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

The solution here (`pand.AndUnsafeInplace`) is faster than `AndUnsafeInplace` for `[]byte`  with 32 or more elements.

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

CPU: AMD Ryzen 7 2700X Eight-Core Processor
Instruction set: AVX2

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
    | csvtk pretty -t -s "        " \
    | perl -pe 's/\n/\n\n/ if /AndLoop/'

rm t.txt

test                                 data-size        time
-----------------------------        ---------        -----------
BenchmarkAndInplaceGrailbio          8.00_B           3.073 ns/op
BenchmarkAndGrailbio                 8.00_B           3.952 ns/op
BenchmarkAndInplaceGoAsm             8.00_B           4.060 ns/op
BenchmarkAndGoAsm                    8.00_B           5.017 ns/op
BenchmarkAndInplaceLoop              8.00_B           5.271 ns/op
BenchmarkAndInplaceUnrollLoop        8.00_B           5.287 ns/op
BenchmarkAndUnrollLoop               8.00_B           6.001 ns/op
BenchmarkAndLoop                     8.00_B           6.012 ns/op

BenchmarkAndInplaceGrailbio          16.00_B          3.836 ns/op
BenchmarkAndInplaceGoAsm             16.00_B          4.308 ns/op
BenchmarkAndGrailbio                 16.00_B          4.718 ns/op
BenchmarkAndGoAsm                    16.00_B          5.370 ns/op
BenchmarkAndInplaceLoop              16.00_B          8.158 ns/op
BenchmarkAndInplaceUnrollLoop        16.00_B          8.357 ns/op
BenchmarkAndUnrollLoop               16.00_B          9.291 ns/op
BenchmarkAndLoop                     16.00_B          10.69 ns/op

BenchmarkAndInplaceGrailbio          32.00_B          4.706 ns/op
BenchmarkAndInplaceGoAsm             32.00_B          4.798 ns/op
BenchmarkAndGoAsm                    32.00_B          5.562 ns/op
BenchmarkAndGrailbio                 32.00_B          5.843 ns/op
BenchmarkAndInplaceLoop              32.00_B          14.06 ns/op
BenchmarkAndInplaceUnrollLoop        32.00_B          14.18 ns/op
BenchmarkAndUnrollLoop               32.00_B          15.76 ns/op
BenchmarkAndLoop                     32.00_B          18.65 ns/op

BenchmarkAndInplaceGoAsm             64.00_KB         1543 ns/op
BenchmarkAndGoAsm                    64.00_KB         1640 ns/op
BenchmarkAndInplaceGrailbio          64.00_KB         3999 ns/op
BenchmarkAndGrailbio                 64.00_KB         4391 ns/op
BenchmarkAndInplaceLoop              64.00_KB         25215 ns/op
BenchmarkAndInplaceUnrollLoop        64.00_KB         25570 ns/op
BenchmarkAndUnrollLoop               64.00_KB         26882 ns/op
BenchmarkAndLoop                     64.00_KB         32291 ns/op

BenchmarkAndInplaceGoAsm             128.00_B         7.589 ns/op
BenchmarkAndGoAsm                    128.00_B         8.157 ns/op
BenchmarkAndInplaceGrailbio          128.00_B         10.95 ns/op
BenchmarkAndGrailbio                 128.00_B         11.82 ns/op
BenchmarkAndInplaceUnrollLoop        128.00_B         51.60 ns/op
BenchmarkAndUnrollLoop               128.00_B         53.78 ns/op
BenchmarkAndInplaceLoop              128.00_B         58.02 ns/op
BenchmarkAndLoop                     128.00_B         79.00 ns/op

BenchmarkAndInplaceGoAsm             512.00_B         16.90 ns/op
BenchmarkAndGoAsm                    512.00_B         17.08 ns/op
BenchmarkAndInplaceGrailbio          512.00_B         38.57 ns/op
BenchmarkAndGrailbio                 512.00_B         41.69 ns/op
BenchmarkAndInplaceUnrollLoop        512.00_B         205.7 ns/op
BenchmarkAndInplaceLoop              512.00_B         207.3 ns/op
BenchmarkAndUnrollLoop               512.00_B         214.3 ns/op
BenchmarkAndLoop                     512.00_B         271.2 ns/op
```

### For developers

Generate Go assembly code with [avo](https://github.com/mmcloughlin/avo)

```
go run asm-AndInplaceSSE2.go -out andInplaceSSE2_amd64.s -stubs andInplaceSSE2.go
go run asm-AndSSE2.go -out andSSE2_amd64.s -stubs andSSE2.go

go run asm-AndInplaceAVX2.go -out andInplaceAVX2_amd64.s -stubs andInplaceAVX2.go
go run asm-AndAVX2.go -out andAVX2_amd64.s -stubs andAVX2.go

go test . -count=1

```

***Attention: since avo does not support AVX512 yet, we need to manually edit
`andAVX512_amd64.s` and `andInplaceAVX512_amd64.s`***

```
go run asm-AndInplaceAVX512.go -out andInplaceAVX512_amd64.s -stubs andInplaceAVX512.go
go run asm-AndAVX512.go -out andAVX512_amd64.s -stubs andAVX512.go

```

For `andInplaceAVX512_amd64.s`, Change

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

For `andAVX512_amd64.s`, Change

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
- [Robert Clausecker](https://github.com/clausecker/) gave some great
  [advices](https://github.com/shenwei356/pand/issues/1) on the Assembly language.
- [Peter Cordes](https://stackoverflow.com/users/224132/peter-cordes)
  provided valuable suggestions on the Assembly language
  in this [post](https://stackoverflow.com/questions/68280854/).
- We copied and edited dispatching code from [pospop](https://github.com/clausecker/pospop).

## License

[MIT License](https://github.com/shenwei356/pand/blob/master/LICENSE)

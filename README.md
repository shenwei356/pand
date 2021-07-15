# Bitwise AND operation on two []byte

The generic code:

```
func AND(x, y []byte) {
	for k, b := range y {
		x[k] &= b
	}
}

```

Generate Go assembly code:

```
go run asm.go -out pand.s -stubs stub.go && go test .

```

Benchmark

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

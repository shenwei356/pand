# Bitwise AND operation on byte slices

The generic code:

```
func loop(x, y []byte) {
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
go test . -bench=Benchmark* -benchmem | tee t.txt

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
BenchmarkGrailbio          7.00_B           5.983 ns/op
BenchmarkGoAsm             7.00_B           6.112 ns/op
BenchmarkLoop              7.00_B           6.977 ns/op
BenchmarkUnrollLoop        7.00_B           8.819 ns/op

BenchmarkGrailbio          32.00_B          6.438 ns/op
BenchmarkGoAsm             32.00_B          6.755 ns/op
BenchmarkUnrollLoop        32.00_B          16.84 ns/op
BenchmarkLoop              32.00_B          18.82 ns/op

BenchmarkGoAsm             128.00_B         11.94 ns/op
BenchmarkGrailbio          128.00_B         14.67 ns/op
BenchmarkUnrollLoop        128.00_B         57.21 ns/op
BenchmarkLoop              128.00_B         70.42 ns/op

BenchmarkGoAsm             256.00_B         17.09 ns/op
BenchmarkGrailbio          256.00_B         23.29 ns/op
BenchmarkUnrollLoop        256.00_B         108.0 ns/op
BenchmarkLoop              256.00_B         135.7 ns/op

BenchmarkGoAsm             1.00_KB          47.78 ns/op
BenchmarkGrailbio          1.00_KB          89.01 ns/op
BenchmarkUnrollLoop        1.00_KB          438.5 ns/op
BenchmarkLoop              1.00_KB          531.5 ns/op

BenchmarkGoAsm             64.00_KB         2655 ns/op
BenchmarkGrailbio          64.00_KB         5049 ns/op
BenchmarkUnrollLoop        64.00_KB         27747 ns/op
BenchmarkLoop              64.00_KB         33881 ns/op
```

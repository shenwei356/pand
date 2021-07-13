## Bitwise AND operation on byte slices

The generic code:

```
func loop(x, y []byte) {
	for k, b := range y {
		x[k] &= b
	}
}

```

Benchmark result

```
go test . -bench=Benchmark2* -benchmem | tee t.txt

cat t.txt \
    | grep Bench \
    | sed -r 's/\s\s+/\t/g' \
    | csvtk cut -Ht -f 1,3 \
    | csvtk add-header -t -n test,time \
    | csvtk mutate -t -n data-size -p "/(.+)-" \
    | csvtk replace -t -p "(.+)/.+" -r "\$1" \
    | csvtk cut -t -f test,data-size,time \
    | csvtk sort -t -k data-size:N -k time:N \
    | csvtk pretty -t -s "      "

test                      data-size      time
--------------------      ---------      -----------
Benchmark2GoAsmAvx2       1.00_KB        51.02 ns/op
Benchmark2Grailbio        1.00_KB        90.20 ns/op
Benchmark2UnrollLoop      1.00_KB        427.2 ns/op
Benchmark2Loop            1.00_KB        550.7 ns/op

Benchmark2GoAsmAvx2       64.00_KB       2616 ns/op
Benchmark2Grailbio        64.00_KB       4979 ns/op
Benchmark2UnrollLoop      64.00_KB       27835 ns/op
Benchmark2Loop            64.00_KB       33436 ns/op

Benchmark2GoAsmAvx2       128.00_B       13.07 ns/op
Benchmark2Grailbio        128.00_B       15.90 ns/op
Benchmark2UnrollLoop      128.00_B       59.58 ns/op
Benchmark2Loop            128.00_B       81.20 ns/op


```

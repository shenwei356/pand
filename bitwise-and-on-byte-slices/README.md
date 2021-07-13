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
Benchmark2GoAsmAvx2       1.00_KB        49.22 ns/op
Benchmark2Grailbio        1.00_KB        89.67 ns/op
Benchmark2UnrollLoop      1.00_KB        436.6 ns/op
Benchmark2Loop            1.00_KB        536.7 ns/op

Benchmark2GoAsmAvx2       64.00_KB       2623 ns/op
Benchmark2Grailbio        64.00_KB       5011 ns/op
Benchmark2UnrollLoop      64.00_KB       26378 ns/op
Benchmark2Loop            64.00_KB       34118 ns/op

Benchmark2GoAsmAvx2       128.00_B       13.28 ns/op
Benchmark2Grailbio        128.00_B       15.83 ns/op
Benchmark2UnrollLoop      128.00_B       56.17 ns/op
Benchmark2Loop            128.00_B       81.70 ns/op

```

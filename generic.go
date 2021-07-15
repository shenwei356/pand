package pand

func andGeneric0(x []byte, y []byte) {
	if len(x) != len(y) {
		panic("x and y should have equal length")
	}

	for i, b := range x {
		x[i] = b & y[i]
	}
}

func andGeneric(x []byte, y []byte) {
	if len(x) != len(y) {
		panic("x and y should have equal length")
	}

	k := 0
	for len(y) >= 8 { // unroll loop
		x[k] &= y[0]
		k++
		x[k] &= y[1]
		k++
		x[k] &= y[2]
		k++
		x[k] &= y[3]
		k++
		x[k] &= y[4]
		k++
		x[k] &= y[5]
		k++
		x[k] &= y[6]
		k++
		x[k] &= y[7]
		k++
		y = y[8:]
	}
	for _, b := range y {
		x[k] &= b
		k++
	}
}

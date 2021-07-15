package pand

func andInplaceGeneric0(x []byte, y []byte) {
	if len(x) != len(y) {
		panic("pand: byte slices should have equal length")
	}

	for i, b := range x {
		x[i] = b & y[i]
	}
}

func andGeneric0(r []byte, x []byte, y []byte) {
	if !(len(x) == len(y) && len(r) == len(x)) {
		panic("pand: byte slices should have equal length")
	}

	for i, b := range x {
		r[i] = b & y[i]
	}
}

func andInplaceGeneric(x []byte, y []byte) {
	if len(x) != len(y) {
		panic("pand: byte slices should have equal length")
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

func andGeneric(r []byte, x []byte, y []byte) {
	if !(len(x) == len(y) && len(r) == len(x)) {
		panic("pand: byte slices should have equal length")
	}

	k := 0
	for len(y) >= 8 { // unroll loop
		r[k] = x[k] & y[0]
		k++
		r[k] = x[k] & y[1]
		k++
		r[k] = x[k] & y[2]
		k++
		r[k] = x[k] & y[3]
		k++
		r[k] = x[k] & y[4]
		k++
		r[k] = x[k] & y[5]
		k++
		r[k] = x[k] & y[6]
		k++
		r[k] = x[k] & y[7]
		k++
		y = y[8:]
	}
	for _, b := range y {
		r[k] = x[k] & b
		k++
	}
}

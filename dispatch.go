// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

package pand

type andInplaceImpl struct {
	function  func(x []byte, y []byte)
	name      string
	available bool
}

type andImpl struct {
	function  func(r []byte, x []byte, y []byte)
	name      string
	available bool
}

var andInplaceFunc = func() func(x []byte, y []byte) {
	for _, f := range andInplaceFuncs {
		if f.available {
			return f.function
		}
	}

	panic("no implementation available")
}()

var andFunc = func() func(r []byte, x []byte, y []byte) {
	for _, f := range andFuncs {
		if f.available {
			return f.function
		}
	}

	panic("no implementation available")
}()

func AndInplace(x []byte, y []byte) {
	if len(x) != len(y) {
		panic("pand: byte slices should have equal length")
	}

	andInplaceFunc(x, y)
}

func AndUnsafeInplace(x []byte, y []byte) {
	andInplaceFunc(x, y)
}

func And(r []byte, x []byte, y []byte) {
	if !(len(x) == len(y) && len(r) == len(x)) {
		panic("pand: byte slices should have equal length")
	}

	andFunc(r, x, y)
}

func AndUnsafe(r []byte, x []byte, y []byte) {
	andFunc(r, x, y)
}

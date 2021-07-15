// we adopt the similar method to choose platform relevant function from:
// https://github.com/clausecker/pospop

package pand

type andImpl struct {
	function  func(x []byte, y []byte)
	name      string
	available bool
}

var andFunc = func() func(x []byte, y []byte) {
	for _, f := range andFuncs {
		if f.available {
			return f.function
		}
	}

	panic("no implementation available")
}()

func AND(x []byte, y []byte) {
	andFunc(x, y)
}

//go:build ignore
// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	TEXT("PAND", NOSPLIT|NOPTR, "func(x []byte, y []byte)")

	Comment("pointer of x")
	x := Mem{Base: Load(Param("x").Base(), GP64())}
	Comment("length of x")
	n := Load(Param("x").Len(), GP64())

	Comment("pointer of y")
	y := Mem{Base: Load(Param("y").Base(), GP64())}

	Comment("--------------------------------------------")

	Comment("end address of x")
	end0 := GP64()
	MOVQ(x.Base, end0)
	ADDQ(n, end0)

	Comment("end address for loop")
	end := GP64()
	MOVQ(end0, end)

	Comment("n < 32, jump to loop8")
	CMPQ(n, U32(32))
	JB(LabelRef("loop8_start"))

	Comment("--------------------------------------------")

	Comment("end address for loop32")
	SUBQ(U32(31), end)

	s := YMM()

	Label("loop32")

	Comment("compute x & y, and save value to x")
	VMOVDQA(x.Offset(0), s)
	VANDPS(y.Offset(0), s, s)
	VMOVDQA(s, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(32), x.Base)
	ADDQ(U32(32), y.Base)

	CMPQ(x.Base, end)
	JB(LabelRef("loop32"))

	Comment("--------------------------------------------")

	Label("loop8_start")

	Comment("end address for loop8")
	SUBQ(U32(8), end0)

	t := GP64()

	Label("loop8")

	Comment("compute x & y, and save value to x")
	MOVQ(x.Offset(0), t)
	ANDQ(y.Offset(0), t)
	MOVQ(t, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(8), x.Base)
	ADDQ(U32(8), y.Base)

	CMPQ(x.Base, end0)
	JB(LabelRef("loop8"))

	Comment("--------------------------------------------")

	Comment("left elements (<8)")
	MOVQ(x.Offset(0), t)
	ANDQ(y.Offset(0), t)
	MOVQ(t, x.Offset(0))

	RET()

	Generate()
}

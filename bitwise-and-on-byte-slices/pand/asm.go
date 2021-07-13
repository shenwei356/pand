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

	// ------------------

	Comment("--------------------------------------------")

	s := YMM()
	VXORPD(s, s, s)

	Label("loop32")

	Comment("check number of left elements")
	CMPQ(n, U32(32)) // 256/8
	JL(LabelRef("loop8"))

	Comment("compute bitwise AND and save the value back to *x")
	VMOVAPD(x.Offset(0), s)
	VPAND(y.Offset(0), s, s)
	// VANDPD(y.Offset(0), s, s)
	VMOVAPD(s, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(32), x.Base)
	ADDQ(U32(32), y.Base)

	Comment("number of left elements")
	SUBQ(U32(32), n) // 256/8
	JMP(LabelRef("loop32"))

	// ------------------

	Comment("--------------------------------------------")

	w := GP64()
	XORQ(w, w)

	Label("loop8")

	Comment("check number of left elements")
	CMPQ(n, U32(8)) // 256/8/8
	JL(LabelRef("end"))

	Comment("compute bitwise AND and save the value back to *x")
	MOVQ(x.Offset(0), w)
	ANDQ(y.Offset(0), w)
	MOVQ(w, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(8), x.Base)
	ADDQ(U32(8), y.Base)

	Comment("number of left elements")
	SUBQ(U32(8), n) // 256/8//8
	JMP(LabelRef("loop8"))

	// ------------------

	Comment("--------------------------------------------")

	Label("end")

	MOVQ(x.Offset(0), w)
	ANDQ(y.Offset(0), w)
	MOVQ(w, x.Offset(0))

	RET()

	Generate()
}
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

	Comment("check number of left elements")
	CMPQ(n, U32(8))
	JL(LabelRef("tail"))

	CMPQ(n, U32(16))
	JL(LabelRef("loop8"))

	CMPQ(n, U32(32))
	JL(LabelRef("loop16"))

	Comment("--------------------------------------------")

	Label("loop32")

	s := YMM() // 256 bits

	Comment("check number of left elements")
	CMPQ(n, U32(32))
	JL(LabelRef("loop16"))

	Comment("compute x & y, and save value to x")
	VMOVDQA(x.Offset(0), s)
	VANDPS(y.Offset(0), s, s)
	VMOVDQA(s, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(32), x.Base)
	ADDQ(U32(32), y.Base)

	SUBQ(U32(32), n)
	JMP(LabelRef("loop32"))

	Comment("--------------------------------------------")

	Label("loop16")

	h := XMM() // 128 bits

	Comment("check number of left elements")
	CMPQ(n, U32(16))
	JL(LabelRef("loop8"))

	Comment("compute x & y, and save value to x")
	VMOVDQA(x.Offset(0), h)
	VANDPS(y.Offset(0), h, h)
	VMOVDQA(h, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(16), x.Base)
	ADDQ(U32(16), y.Base)

	SUBQ(U32(16), n)
	JMP(LabelRef("loop16"))

	Comment("--------------------------------------------")

	Label("loop8")

	t := GP64() // 64 bits

	Comment("check number of left elements")
	CMPQ(n, U32(8))
	JL(LabelRef("tail"))

	Comment("compute x & y, and save value to x")
	MOVQ(x.Offset(0), t)
	ANDQ(y.Offset(0), t)
	MOVQ(t, x.Offset(0))

	Comment("move pointer")
	ADDQ(U32(8), x.Base)
	ADDQ(U32(8), y.Base)

	SUBQ(U32(8), n)
	JMP(LabelRef("loop8"))

	Comment("--------------------------------------------")

	Label("tail")

	Comment("left elements (<=8)")
	MOVQ(x.Offset(0), t)
	ANDQ(y.Offset(0), t)
	MOVQ(t, x.Offset(0))

	RET()

	Generate()
}

//go:build ignore
// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	TEXT("andAVX2", NOSPLIT|NOPTR, "func(r []byte, x []byte, y []byte)")

	Comment("pointer of r")
	r := Mem{Base: Load(Param("r").Base(), GP64())}
	Comment("pointer of x")
	x := Mem{Base: Load(Param("x").Base(), GP64())}
	Comment("length of x")
	n := Load(Param("x").Len(), GP64())

	Comment("pointer of y")
	y := Mem{Base: Load(Param("y").Base(), GP64())}

	Comment("--------------------------------------------")

	Comment("end address of x, will not change: p + n")
	end0 := GP64()
	MOVQ(x.Base, end0)
	ADDQ(n, end0)

	Comment("end address for loop")
	end := GP64()

	Comment("n < 8, jump to tail")
	CMPQ(n, U32(8))
	JL(LabelRef("tail"))

	Comment("n < 16, jump to loop8")
	CMPQ(n, U32(16))
	JL(LabelRef("loop8_start"))

	Comment("n < 32, jump to loop16")
	CMPQ(n, U32(32))
	JL(LabelRef("loop16_start"))

	left := GP64()

	Comment("--------------------------------------------")

	Comment("end address for loop32")
	MOVQ(end0, end)
	SUBQ(U32(31), end)

	Label("loop32")

	s := YMM() // 256 bits

	Comment("compute x & y, and save value to x")
	VMOVDQU(x.Offset(0), s)
	// VANDPS(y.Offset(0), s, s)
	VPAND(y.Offset(0), s, s)
	VMOVDQU(s, r.Offset(0))

	Comment("move pointer")
	ADDQ(U32(32), x.Base)
	ADDQ(U32(32), y.Base)
	ADDQ(U32(32), r.Base)

	CMPQ(x.Base, end)
	JL(LabelRef("loop32"))

	Comment("n < 8, jump to tail")
	MOVQ(end0, left)
	SUBQ(x.Base, left)
	CMPQ(left, U32(8))
	JL(LabelRef("tail"))

	Comment("n < 16, jump to loop8")
	CMPQ(left, U32(16))
	JL(LabelRef("loop8_start"))

	Comment("--------------------------------------------")

	Label("loop16_start")

	Comment("end address for loop16")
	MOVQ(end0, end)
	SUBQ(U32(15), end)

	Label("loop16")

	h := XMM() // 128 bits

	Comment("compute x & y, and save value to x")
	VMOVDQU(x.Offset(0), h)
	// VANDPS(y.Offset(0), h, h)
	VPAND(y.Offset(0), h, h)
	VMOVDQU(h, r.Offset(0))

	Comment("move pointer")
	ADDQ(U32(16), x.Base)
	ADDQ(U32(16), y.Base)
	ADDQ(U32(16), r.Base)

	CMPQ(x.Base, end)
	JL(LabelRef("loop16"))

	Comment("n < 8, jump to tail")
	MOVQ(end0, left)
	SUBQ(x.Base, left)
	CMPQ(left, U32(8))
	JL(LabelRef("tail"))

	Comment("--------------------------------------------")

	Label("loop8_start")

	Comment("end address for loop8")
	MOVQ(end0, end)
	SUBQ(U32(7), end)

	Label("loop8")

	t := GP64() // 64 bits

	Comment("compute x & y, and save value to x")
	MOVQ(x.Offset(0), t)
	ANDQ(y.Offset(0), t)
	MOVQ(t, r.Offset(0))

	Comment("move pointer")
	ADDQ(U32(8), x.Base)
	ADDQ(U32(8), y.Base)
	ADDQ(U32(8), r.Base)

	CMPQ(x.Base, end)
	JL(LabelRef("loop8"))

	Comment("--------------------------------------------")

	Label("tail")
	Comment("left elements (<8)")

	o := GP8()

	CMPQ(x.Base, end0)
	JE(LabelRef("end"))

	MOVB(x.Offset(0), o)
	ANDB(y.Offset(0), o)
	MOVB(o, r.Offset(0))

	ADDQ(U32(1), x.Base)
	ADDQ(U32(1), y.Base)
	ADDQ(U32(1), r.Base)
	JMP(LabelRef("tail"))

	Label("end")

	RET()

	Generate()
}

// Code generated by command: go run asm.go -out pand.s -stubs stub.go. DO NOT EDIT.

#include "textflag.h"

// func PAND(x []byte, y []byte)
// Requires: AVX
TEXT ·PAND(SB), NOSPLIT|NOPTR, $0-48
	// pointer of x
	MOVQ x_base+0(FP), AX

	// length of x
	MOVQ x_len+8(FP), CX

	// pointer of y
	MOVQ y_base+24(FP), DX

	// --------------------------------------------
	// end address of x
	MOVQ AX, BX
	ADDQ CX, BX

	// end address for loop
	MOVQ BX, SI

	// n < 32, jump to loop8
	CMPQ CX, $0x00000020
	JB   loop8_start

	// --------------------------------------------
	// end address for loop32
	SUBQ $0x0000001f, SI

loop32:
	// compute x & y, and save value to x
	VMOVDQA (AX), Y0
	VANDPS  (DX), Y0, Y0
	VMOVDQA Y0, (AX)

	// move pointer
	ADDQ $0x00000020, AX
	ADDQ $0x00000020, DX
	CMPQ AX, SI
	JB   loop32

	// --------------------------------------------
loop8_start:
	// end address for loop8
	SUBQ $0x00000008, BX

loop8:
	// compute x & y, and save value to x
	MOVQ (AX), CX
	ANDQ (DX), CX
	MOVQ CX, (AX)

	// move pointer
	ADDQ $0x00000008, AX
	ADDQ $0x00000008, DX
	CMPQ AX, BX
	JB   loop8

	// --------------------------------------------
	// left elements (<8)
	MOVQ (AX), CX
	ANDQ (DX), CX
	MOVQ CX, (AX)
	RET

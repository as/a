// Copyright 2018 (as). Added avx and avx2 support for capable CPUs
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

DATA ·AVX2_swizzletab<>+0x00(SB)/8, $0x0704050603000102
DATA ·AVX2_swizzletab<>+0x08(SB)/8, $0x0f0c0d0e0b08090a
DATA ·AVX2_swizzletab<>+0x10(SB)/8, $0x1714151613101112
DATA ·AVX2_swizzletab<>+0x18(SB)/8, $0x1f1c1d1e1b18191a
GLOBL ·AVX2_swizzletab<>(SB), (NOPTR+RODATA), $32

// func haveSSSE3() bool
TEXT ·haveSSSE3(SB),NOSPLIT,$0
	MOVQ	$1, AX
	CPUID
	SHRQ	$9, CX
	ANDQ	$1, CX
	MOVB	CX, ret+0(FP)
	RET

// func haveAVX() bool
TEXT ·haveAVX(SB),NOSPLIT,$0
	MOVQ	$1, AX
	CPUID
	SHRQ	$28, CX
	ANDQ	$1, CX
	MOVB	CX, ret+0(FP)
	RET
	
// func haveAVX2() bool
TEXT ·haveAVX2(SB),NOSPLIT,$0
	MOVQ	$7, AX
	MOVQ	$0, CX
	CPUID
	SHRQ	$5, BX
	ANDQ	$1, BX
	MOVB	BX, ret+0(FP)
	RET

// func bgra256sd(p, q []byte)
TEXT ·bgra256sd(SB),NOSPLIT,$0
	MOVQ	p+0(FP), SI
	MOVQ	len+8(FP), CX
	MOVQ	q+24(FP), DI
	
	VMOVDQU ·AVX2_swizzletab<>(SB), Y0
	ADDQ SI, CX
	ADDQ $256, SI
	CMPQ CX, SI
	JL prep32
	SUBQ	$256, SI
	
loop256:
	VMOVDQU 	(0*32)(SI),Y1 
	VMOVDQU 	(1*32)(SI),Y2 
	VMOVDQU 	(2*32)(SI),Y3 
	VMOVDQU 	(3*32)(SI),Y4 
	VMOVDQU 	(4*32)(SI),Y5 
	VMOVDQU 	(5*32)(SI),Y6 
	VMOVDQU 	(6*32)(SI),Y7 
	VMOVDQU 	(7*32)(SI),Y8 
	VPSHUFB Y0, Y1,  Y1
	VPSHUFB Y0, Y2,  Y2
	VPSHUFB Y0, Y3,  Y3
	VPSHUFB Y0, Y4,  Y4
	VPSHUFB Y0, Y5,  Y5
	VPSHUFB Y0, Y6,  Y6
	VPSHUFB Y0, Y7,  Y7
	VPSHUFB Y0, Y8,  Y8
	VMOVDQU	Y1, (0*32)(DI)
	VMOVDQU	Y2, (1*32)(DI)
	VMOVDQU	Y3, (2*32)(DI)
	VMOVDQU	Y4, (3*32)(DI)
	VMOVDQU	Y5, (4*32)(DI)
	VMOVDQU	Y6, (5*32)(DI)
	VMOVDQU	Y7, (6*32)(DI)
	VMOVDQU	Y8, (7*32)(DI)
	ADDQ	$256, SI
	ADDQ	$256, DI
	CMPQ CX, SI
	JGT loop256
	JEQ done
	
	SUBQ	$256, DI
prep32:
	SUBQ	$256, SI
	ADDQ $32, SI
	CMPQ CX, SI
	JL prep4
	SUBQ	$32, SI

loop32:
	VMOVDQU 	(0*32)(SI),Y1 
	VPSHUFB Y0, Y1,  Y1
	VMOVDQU	Y1, (0*32)(DI)
	ADDQ	$32, SI
	ADDQ	$32, DI
	CMPQ CX, SI
	JGT	loop32
	JEQ done
	
	SUBQ	$32, DI
prep4:
	SUBQ	$32, SI
	
loop4:
	MOVD	0(SI), AX	// r g b a
	BSWAPL AX   // a b g r 
	RORL	$8, AX 	// b g r a 
	MOVD	AX, (DI)

	ADDQ	$4, SI
	ADDQ	$4, DI
	CMPQ CX, SI
	JGT	loop4

done:
	RET

// func bgra128sd(p, q []byte)
TEXT ·bgra128sd(SB),NOSPLIT,$0
	MOVQ	p+0(FP), SI
	MOVQ	len+8(FP), CX
	MOVQ	q+24(FP), DI
	
	VMOVDQU ·AVX2_swizzletab<>(SB), X0
	ADDQ SI, CX
	ADDQ $128, SI
	CMPQ CX, SI
	JL prep16
	SUBQ	$128, SI
	
loop128:
	VMOVDQU 	(0*16)(SI),X1 
	VMOVDQU 	(1*16)(SI),X2 
	VMOVDQU 	(2*16)(SI),X3 
	VMOVDQU 	(3*16)(SI),X4 
	VMOVDQU 	(4*16)(SI),X5 
	VMOVDQU 	(5*16)(SI),X6 
	VMOVDQU 	(6*16)(SI),X7 
	VMOVDQU 	(7*16)(SI),X8 
	VPSHUFB X0, X1,  X1
	VPSHUFB X0, X2,  X2
	VPSHUFB X0, X3,  X3
	VPSHUFB X0, X4,  X4
	VPSHUFB X0, X5,  X5
	VPSHUFB X0, X6,  X6
	VPSHUFB X0, X7,  X7
	VPSHUFB X0, X8,  X8
	VMOVDQU	X1, (0*16)(DI)
	VMOVDQU	X2, (1*16)(DI)
	VMOVDQU	X3, (2*16)(DI)
	VMOVDQU	X4, (3*16)(DI)
	VMOVDQU	X5, (4*16)(DI)
	VMOVDQU	X6, (5*16)(DI)
	VMOVDQU	X7, (6*16)(DI)
	VMOVDQU	X8, (7*16)(DI)
	ADDQ	$128, SI
	ADDQ	$128, DI
	CMPQ CX, SI
	JGT loop128
	JEQ done
	
	SUBQ	$128, DI
prep16:
	SUBQ	$128, SI
	ADDQ $16, SI
	CMPQ CX, SI
	JL prep4
	SUBQ	$16, SI

loop16:
	VMOVDQU 	(0*16)(SI),X1 
	VPSHUFB X0, X1,  X1
	VMOVDQU	X1, (0*16)(DI)
	ADDQ	$16, SI
	ADDQ	$16, DI
	CMPQ CX, SI
	JGT	loop16
	JEQ done
	
	SUBQ	$16, DI
prep4:
	SUBQ	$16, SI
	
loop4:
	MOVD	0(SI), AX	// r g b a
	BSWAPL AX   // a b g r 
	RORL	$8, AX 	// b g r a 
	MOVD	AX, (DI)

	ADDQ	$4, SI
	ADDQ	$4, DI
	CMPQ CX, SI
	JGT	loop4

done:
	RET
	
// func bgra16sd(p, q []byte)
TEXT ·bgra16sd(SB),NOSPLIT,$0
	MOVQ	p+0(FP), SI
	MOVQ	len+8(FP), CX
	MOVQ	q+24(FP), DI

	// Sanity check that len is a multiple of 16.
	//	MOVQ	CX, AX
	//	ANDQ	$15, AX
	//	JNZ	done
	ADDQ SI, CX

	// Make the shuffle control mask (16-byte register X0) look like this,
	// where the low order byte comes first:
	//
	// 02 01 00 03  06 05 04 07  0a 09 08 0b  0e 0d 0c 0f
	//
	// Load the bottom 8 bytes into X0, the top into X1, then interleave them
	// into X0.
	MOVQ	$0x0704050603000102, AX
	MOVQ	AX, X0
	MOVQ	$0x0f0c0d0e0b08090a, AX
	MOVQ	AX, X1
	PUNPCKLQDQ	X1, X0

loop16:
	MOVOU	(SI), X1
	PSHUFB	X0, X1
	MOVOU	X1, (DI)

	ADDQ	$16, SI
	ADDQ	$16, DI
	CMPQ CX, SI
	JGT	loop16
	JEQ done

prep4:
	SUBQ	$16, DI
	SUBQ	$16, SI
loop4:
	MOVD	0(SI), AX	// r g b a
	BSWAPL AX   // a b g r 
	RORL	$8, AX 	// b g r a 
	MOVD	AX, (DI)

	ADDQ	$4, SI
	ADDQ	$4, DI
	CMPQ CX, SI
	JGT	loop4
done:
	RET

// func bgra4sd(p, q []byte)
TEXT ·bgra4sd(SB),NOSPLIT,$0
	MOVQ	p+0(FP), SI
	MOVQ	len+8(FP), CX
	MOVQ	q+24(FP), DI

	ADDQ SI, CX
loop:
	CMPQ	SI, CX
	JEQ	done

	MOVD	0(SI), AX	// r g b a
	BSWAPL AX   // a b g r 
	RORL	$8, AX 	// b g r a 
	MOVD	AX, (DI)

	ADDQ	$4, SI
	ADDQ	$4, DI
	JMP	loop
done:
	RET

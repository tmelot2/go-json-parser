// Original assembly by David Terei: https://github.com/dterei/gotsc
// I unified the begin/end calls into 1 single call.

#include "textflag.h"

// func ReadCPUTimer() uint64
TEXT ·ReadCPUTimer(SB),NOSPLIT,$0-8
	CPUID
	RDTSC
	SHLQ	$32, DX
	ADDQ	DX, AX
	MOVQ	AX, ret+0(FP)
	RET

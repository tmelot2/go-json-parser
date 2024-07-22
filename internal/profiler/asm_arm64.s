#include "textflag.h"

// Empty package before · is current package.
// SB means that the function is an offset from virtual register SB.
// NOSPLIT means "do not insert stack-split preamble", we do not need stack
// $0 is the stack size needed
// 8 is the size of arguments+return values. Can be addressed using FP
//
// See NOTE-1 in timer_darwin.go for info on why this timestamp is not useful.
TEXT ·ReadCPUTimer(SB),NOSPLIT,$0-8
	ISB $1
	MRS CNTVCT_EL0, R0

	MOVD R0, ret+0(FP)
	RET

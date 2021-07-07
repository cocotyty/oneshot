#include "textflag.h"

// func getG() uintptr
TEXT ·getG(SB), NOSPLIT, $0-8
    BL	runtime·load_g(SB)
	MOVD g, R0
	MOVD R0, ret+0(FP)
	RET

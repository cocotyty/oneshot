#include "textflag.h"
#include "go_tls.h"

// func getG() uintptr
TEXT ·getG(SB), NOSPLIT, $0-8
    get_tls(CX)
	MOVD g(CX), AX
	MOVD AX, ret+0(FP)
	RET

#include "textflag.h"

TEXT ·fastCaller(SB), NOSPLIT, $0-16
	MOVQ	$0, DX
	MOVQ	BP, BX          // BX = BP
	MOVQ	s+0(FP), AX     // AX = s
	INCQ	AX
loop:
	NOP
	CMPQ	BX, DX
	JEQ	nil
	DECQ	AX
	JEQ	done
	MOVQ	(BX), BX
	JMP	loop
done:
	NOP
	MOVQ	8(BX), BX
	MOVQ	BX, c+8(FP) // ret = BX
	RET
nil:
	NOP
	MOVQ	BX, c+8(FP) // ret = BX
	RET

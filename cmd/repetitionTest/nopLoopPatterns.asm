global NOP3x1AllBytes
global NOP1x3AllBytes
global NOP1x9AllBytes

section .text

; Targets Win64 ABI.
; 1st param is RCX (count), 2nd param is RDX (data pointer).

NOP3x1AllBytes:
	xor rax, rax
.loop:
	db 0x0f, 0x1f, 0x00 ; Byte sequence for 3-byte NOP
	inc rax
	cmp rax, rcx
	jb .loop
	ret

NOP1x3AllBytes:
	xor rax, rax
.loop:
	nop
	nop
	nop
	inc rax
	cmp rax, rcx
	jb .loop
	ret

NOP1x9AllBytes:
	xor rax, rax
.loop:
	nop
	nop
	nop
	nop
	nop
	nop
	nop
	nop
	nop
	inc rax
	cmp rax, rcx
	jb .loop
	ret

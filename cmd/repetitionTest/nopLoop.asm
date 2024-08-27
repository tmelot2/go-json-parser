global MOVAllBytesASM
global NOPAllBytesASM
global CMPAllBytesASM
global DECAllBytesASM

section .text

; Targets Win64 ABI.
; 1st param is RCX (count), 2nd param is RDX (data pointer).

MOVAllBytesASM:
	xor rax, rax
.loop:
	mov [rdx + rax], al
	inc rax
	cmp rax, rcx
	jb .loop
	ret

NOPAllBytesASM:
	xor rax, rax
.loop:
	db 0x0f, 0x1f, 0x00 ; Byte sequence for 3-byte NOP
	inc rax
	cmp rax, rcx
	jb .loop
	ret

CMPAllBytesASM:
	xor rax, rax
.loop:
	inc rax
	cmp rax, rcx
	jb .loop
	ret

DECAllBytesASM:
.loop:
	dec rcx
	jnz .loop
	ret

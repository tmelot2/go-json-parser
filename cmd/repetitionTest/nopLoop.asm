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
	nop
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

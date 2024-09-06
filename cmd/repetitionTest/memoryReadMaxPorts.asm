global Read_x1
global Read_x2
global Read_x3
global Read_x4

section .text

; Targets Win64 ABI.
; 1st param is RCX (count), 2nd param is RDX (data pointer).

Read_x1:
.loop:
	mov rbx, [rdx]
	dec rcx
	jnz .loop
	ret

Read_x2:
.loop:
	mov rbx, [rdx]
	mov rbx, [rdx]
	dec rcx
	jnz .loop
	ret

Read_x3:
.loop:
	mov rbx, [rdx]
	mov rbx, [rdx]
	mov rbx, [rdx]
	dec rcx
	jnz .loop
	ret

Read_x4:
.loop:
	mov rbx, [rdx]
	mov rbx, [rdx]
	mov rbx, [rdx]
	mov rbx, [rdx]
	dec rcx
	jnz .loop
	ret

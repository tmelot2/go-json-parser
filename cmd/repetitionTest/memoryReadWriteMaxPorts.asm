global Read_x1
global Read_x2
global Read_x3
global Read_x4
global Write_x1
global Write_x2
global Write_x3
global Write_x4

section .text

; Targets Win64 ABI.
; 1st param is RCX (count), 2nd param is RDX (data pointer).

Read_x1:
	align 64
.loop:
	mov rbx, [rdx]
	sub rcx, 1
	jnle .loop
	ret

Read_x2:
	align 64
.loop:
	mov rbx, [rdx]
	mov rbx, [rdx]
	sub rcx, 2
	jnle .loop
	ret

Read_x3:
	align 64
.loop:
	mov rbx, [rdx]
	mov rbx, [rdx]
	mov rbx, [rdx]
	sub rcx, 3
	jnle .loop
	ret

Read_x4:
	align 64
.loop:
	mov rbx, [rdx]
	mov rbx, [rdx]
	mov rbx, [rdx]
	mov rbx, [rdx]
	sub rcx, 4
	jnle .loop
	ret

Write_x1:
	xor rbx, rbx
	align 64
.loop:
	mov [rdx], rbx
	sub rcx, 1
	jnle .loop
	ret

Write_x2:
	xor rbx, rbx
	align 64
.loop:
	mov [rdx], rbx
	mov [rdx], rbx
	sub rcx, 2
	jnle .loop
	ret

Write_x3:
	xor rbx, rbx
	align 64
.loop:
	mov [rdx], rbx
	mov [rdx], rbx
	mov [rdx], rbx
	sub rcx, 3
	jnle .loop
	ret

Write_x4:
	xor rbx, rbx
	align 64
.loop:
	mov [rdx], rbx
	mov [rdx], rbx
	mov [rdx], rbx
	mov [rdx], rbx
	sub rcx, 4
	jnle .loop
	ret

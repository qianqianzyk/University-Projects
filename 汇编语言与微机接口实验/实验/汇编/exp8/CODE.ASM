data segment
    promptInput db "Enter a string: $"          ; 输入提示
    resultUpper db 0Dh,0Ah,"Uppercase letters: $" ; 大写结果前缀
    resultLower db 0Dh,0Ah,"Lowercase letters: $" ; 小写结果前缀
    newline     db 0Dh,0Ah,"$"                  ; 换行

    buffer      db 100,0,100 dup(?)             ; 输入缓冲区，最大 100 字符
data ends

code segment
main proc far
    assume cs:code, ds:data
start:
    ; --- 初始化 ---
    push ds
    sub ax,ax
    push ax
    mov ax,data
    mov ds,ax

InputLoop:
    ; --- 提示输入 ---
    lea dx,promptInput
    mov ah,9
    int 21h

    ; --- 读取字符串 ---
    lea dx,buffer
    mov ah,0Ah
    int 21h

    ; --- 打印换行 ---
    lea dx,newline
    mov ah,9
    int 21h

CountString:
    ; --- 准备统计 ---
    mov cl, buffer+1        ; 获取实际输入长度
    lea si, buffer+2        ; 指向字符串内容
    
    xor bx, bx              ; BX 清零，用于统计大写字母
    xor di, di              ; DI 清零，用于统计小写字母

CountLoop:
    cmp cl, 0               ; 检查剩余长度
    jz ShowResult           ; 如果长度为 0，统计结束，去显示结果

    mov al, [si]            ; 读取当前字符
    
    ; --- 检查大写字母 ---
    cmp al, 'A'
    jb NotUpper             ; 小于 'A'，肯定不是大写
    cmp al, 'Z'
    ja NotUpper             ; 大于 'Z'，肯定不是大写
    inc bx                  ; 是大写字母，BX 加 1
    jmp NextChar            ; 处理下一个

NotUpper:
    ; --- 检查小写字母 ---
    cmp al, 'a'
    jb NextChar             ; 小于 'a'，不是小写
    cmp al, 'z'
    ja NextChar             ; 大于 'z'，不是小写
    inc di                  ; 是小写字母，DI 加 1

NextChar:
    inc si                  ; 指针后移
    dec cl                  ; 剩余长度减 1
    jmp CountLoop           ; 继续循环

ShowResult:
    ; --- 显示大写结果 ---
    lea dx, resultUpper
    mov ah, 9
    int 21h

    mov ax, bx              ; 将大写计数放入 AX
    call PrintAX            ; 打印

    ; --- 显示小写结果 ---
    lea dx, resultLower
    mov ah, 9
    int 21h

    mov ax, di              ; 将小写计数放入 AX
    call PrintAX            ; 打印

    ; --- 结尾换行 ---
    lea dx, newline
    mov ah, 9
    int 21h

    jmp InputLoop           ; 跳回开始，允许重复测试

ExitProgram:
    ; 此标签未被跳转使用，作为程序结构的一部分保留
    mov ah, 4Ch
    int 21h

main endp

; --- 子程序：打印 AX ---
PrintAX proc near
    push ax
    push bx
    push cx
    push dx

    mov bx, 10
    xor cx, cx

convertLoop:
    xor dx, dx
    div bx 
    push dx
    inc cx
    test ax, ax
    jnz convertLoop

printLoop:
    pop dx
    add dl, '0'
    mov ah, 2
    int 21h
    loop printLoop

    pop dx
    pop cx
    pop bx
    pop ax
    ret
PrintAX endp

code ends
end start
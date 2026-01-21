data segment
    ; --- 提示信息 ---
    msg         db 'Enter: $'                   ; 输入提示
    msg1        db 0Dh,0Ah,'Decimal = $'        ; 十进制输出前缀
    msg2        db 0Dh,0Ah,'Hex = $'            ; 十六进制输出前缀
    msg3        db 0Dh,0Ah,'Hex(chars) = $'     ; 字符ASCII输出前缀
    newline     db 0Dh,0Ah,'$'

    ; --- 输入缓冲区 ---
    buffer      db 10,0,10 dup(?)
data ends

code segment
main proc far
    assume cs:code, ds:data
start:
    ; --- 初始化 ---
    push    ds
    sub     ax, ax
    push    ax
    mov     ax, data
    mov     ds, ax

    ; --- 读取输入 ---
    lea     dx, msg
    mov     ah, 9
    int     21h

    lea     dx, buffer
    mov     ah, 0Ah
    int     21h

    ; 第一步：检查是否为二进制字符串
    lea     si, buffer+2        ; 指向字符串内容
    mov     cx, 4               ; 只检查前 4 个字符
    mov     bh, 1               ; 设置标志位 BH = 1 (默认为是二进制)

Check1:
    mov     al, [si]            ; 取字符
    cmp     al, '0'
    jz      OK1                 ; 是 '0'，检查通过
    cmp     al, '1'
    jz      OK1                 ; 是 '1'，检查通过
    
    mov     bh, 0               ; 既不是 '0' 也不是 '1'，标志位 BH 置 0
OK1:
    inc     si                  ; 下一个字符
    loop    Check1              ; 循环 4 次

    ; --- 根据标志位跳转 ---
    cmp     bh, 1
    jz      IsBin               ; BH=1，跳转处理二进制逻辑
    jmp     IsChars             ; BH=0，跳转处理字符逻辑

    ; 分支 A：二进制处理
IsBin:
    ; 打印 "Decimal = "
    lea     dx, msg1
    mov     ah, 9
    int     21h

    ; --- 将 "0101" 字符串转换为数值 ---
    lea     si, buffer+2        ; 指向字符串头
    xor     ax, ax              ; AX 清零，用于累加结果
    mov     cx, 4               ; 循环 4 次
Bin1:
    shl     ax, 1               ; AX 左移一位 (乘以 2)
    mov     dl, [si]            ; 取当前字符 ('0' 或 '1')
    sub     dl, '0'             ; 转为数字 (0 或 1)
    or      al, dl              ; 将当前位加到 AL 最低位
    inc     si
    loop    Bin1                ; 循环结束，AX 中即为数值

    ; --- 打印十进制 ---
    call    PrintAX             ; 调用子程序打印十进制

    ; --- 打印十六进制 ---
    push    ax                  ; 保存数值 AX
    lea     dx, msg2            ; 打印 "Hex = "
    mov     ah, 9
    int     21h
    pop     ax                  ; 恢复数值 AX

    call    PrintHex            ; 调用子程序打印十六进制

    jmp     Exit                ; 结束

    ; 分支 B：字符 ASCII 码处理
IsChars:
    ; 打印 "Hex(chars) = "
    lea     dx, msg3
    mov     ah, 9
    int     21h

    lea     si, buffer+2        ; 指向字符串头
    mov     cx, 4               ; 处理前 4 个字符
Loop1:
    xor     ax, ax
    mov     al, [si]            ; 取出字符 ASCII 码到 AL
    call    PrintHex            ; 打印该 ASCII 码的十六进制值

    ; 打印空格分隔
    mov     dl, ' '
    mov     ah, 2
    int     21h

    inc     si
    loop    Loop1

    jmp     Exit

; 子程序定义
; --- 打印 AX 的十进制值 ---
PrintAX proc near
    push    ax
    push    bx
    push    cx
    push    dx

    mov     bx, 10
    xor     cx, cx

Convert1:
    xor     dx, dx
    div     bx
    push    dx
    inc     cx
    test    ax, ax
    jnz     Convert1

Print1:
    pop     dx
    add     dl, '0'
    mov     ah, 2
    int     21h
    loop    Print1

    pop     dx
    pop     cx
    pop     bx
    pop     ax
    ret
PrintAX     endp

; --- 打印 AX 的十六进制值 (打印4位Hex) ---
PrintHex proc near
    push    ax
    push    bx
    push    cx
    push    dx

    mov     bx, ax              ; 将数值移入 BX
    mov     ch, 4               ; 循环 4 次 (打印 4 个 Hex 位)
    mov     cl, 4               ; 每次移位 4 bit
HexLoop:
    rol     bx, cl              ; 循环左移 4 位，将最高 4 位移到最低位
    mov     dl, bl              ; 取出低 8 位
    and     dl, 0Fh             ; 屏蔽高 4 位，只保留最低 4 位
    
    ; --- Hex 转 ASCII ---
    cmp     dl, 9
    jbe     HexNum              ; 如果是 0-9
    add     dl, 7               ; 如果是 A-F (需额外加 7，因为 'A'-'9' != 1)
HexNum:
    add     dl, '0'             ; 转 ASCII
    mov     ah, 2               ; 打印
    int     21h

    dec     ch
    jnz     HexLoop             ; 循环

    pop     dx
    pop     cx
    pop     bx
    pop     ax
    ret
PrintHex    endp

Exit:
    mov     ah, 4Ch
    int     21h

main    endp
code    ends
        end     start
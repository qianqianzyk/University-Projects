data segment
    prompt      db 'Enter n (1-20): $'          ; 输入 N 的提示
    err         db 0Dh,0Ah,'Try again.$'        ; 错误提示
    inputValue  db 0Dh,0Ah,'Enter value (10-99): $' ; 输入数值的提示
    errValue    db 0Dh,0Ah,'Error 3 times. Exit.$'  ; 错误次数过多退出提示
    result      db 0Dh,0Ah,'Sum = $'            ; 结果前缀

    ; --- 变量定义 ---
    n           db 0                            ; 存储 N
    errCount    db 0                            ; 错误计数器
    
    ; --- 输入缓冲区 ---
    ; 格式：最大长度(5), 实际长度, 内容...
    buffer      db 5,0,5 dup(?)
data ends

code segment
main proc far
    assume ds:data, cs:code 
start:
    ; --- 初始化数据段 ---
    push    ds
    sub     ax, ax
    push    ax
    mov     ax, data
    mov     ds, ax

    mov errCount, 0             ; 初始化错误计数器

; 第一阶段：读取 N (范围 1-20)
ReadLoop:
    ; 打印提示 "Enter n (1-20): "
    lea dx, prompt
    mov ah, 9
    int 21h

    ; 读取输入字符串
    lea dx, buffer
    mov ah, 0Ah
    int 21h

    ; --- 校验输入长度 ---
    mov cl, buffer+1            ; 获取实际输入长度
    cmp cl, 1
    jb BadNumber                ; 长度 < 1，非法
    cmp cl, 2
    ja BadNumber                ; 长度 > 2，非法

    ; --- 字符串转数值 ---
    lea si, buffer+2            ; 指向字符串起始
    xor ax, ax
    mov al, [si]                ; 取第一个字符
    sub al, '0'                 ; 转数字

    cmp cl, 2                   ; 如果长度不是 2 (即长度为 1)
    jne Check                   ; 直接去检查范围

    ; 如果长度为 2 (例如 "15")
    mov dl, [si+1]              ; 取第二个字符
    sub dl, '0'                 ; 转数字
    mov ah, al                  ; 将高位暂存到 AH
    mov al, 10                  ; AL = 10
    mul ah                      ; AL = 十位 * 10
    add al, dl                  ; AL = AL + 个位

Check:
    ; --- 校验数值范围 (1-20) ---
    cmp al, 1
    jb BadNumber                ; < 1，非法
    cmp al, 20
    ja BadNumber                ; > 20，非法

    mov n, al                   ; 输入合法，保存 N
    jmp InputValues             ; 进入第二阶段

BadNumber:
    ; --- 错误处理 ---
    inc errCount                ; 错误计数 +1
    cmp errCount, 3             ; 检查是否达到 3 次
    je Exit1                    ; 达到 3 次，跳转退出

    ; 打印 "Try again."
    lea dx, err
    mov ah, 9
    int 21h

    jmp ReadLoop                ; 重试

Exit1:
    ; 打印错误过多提示并退出
    lea dx, errValue
    mov ah, 9
    int 21h
    jmp Exit2

; 第二阶段：循环读取 N 个数值并求和
InputValues:
    xor bx, bx                  ; BX 用于累加总和 (清零)
    mov cl, n                   ; 循环次数设为 N

ValueLoop:
    mov errCount, 0             ; 每次读新数时，重置错误计数

ReadValue:
    ; 打印 "Enter value (10-99): "
    lea dx, inputValue
    mov ah, 9
    int 21h

    ; 读取输入
    lea dx, buffer
    mov ah, 0Ah
    int 21h

    ; --- 校验长度 (必须是两位数) ---
    mov ch, buffer+1
    cmp ch, 2
    jne BadValue                ; 长度不为 2，非法

    ; --- 字符串转数值 ---
    lea si, buffer+2
    mov al, [si]                ; 十位字符
    sub al, '0'
    mov dl, [si+1]              ; 个位字符
    sub dl, '0'

    mov ah, al                  ; 计算数值
    mov al, 10
    mul ah                      ; 十位 * 10
    add al, dl                  ; + 个位

    ; --- 校验数值范围 (10-99) ---
    cmp al, 10
    jb BadValue
    cmp al, 99
    ja BadValue

    ; --- 累加 ---
    add bx, ax                  ; 将当前数值加到 BX
    dec cl                      ; 剩余待输入个数 -1
    jnz ValueLoop               ; 如果还有剩余，继续循环
    jmp ShowResult              ; 全部读完，去显示结果

BadValue:
    ; --- 数值错误处理 ---
    inc errCount
    cmp errCount, 3
    je ExitBadValue

    lea dx, errValue            
    mov ah, 9
    int 21h                     ; 打印错误提示

    jmp ReadValue               ; 重试

ExitBadValue:
    ; 最终失败退出
    lea dx, err
    mov ah, 9
    int 21h
    jmp Exit2

; 第三阶段：显示结果
ShowResult:
    lea dx, result              ; 打印 "Sum = "
    mov ah, 9
    int 21h

    mov ax, bx                  ; 将总和移入 AX
    call PrintAX                ; 打印十进制数值

Exit2:
    mov ah, 4Ch                 ; 退出程序
    int 21h

; --- 子程序：打印 AX (十进制) ---
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
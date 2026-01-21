data    segment
        space   db      " $"            ; 定义空格字符串
        newline db      0Dh,0Ah,"$"     ; 定义换行字符串
data    ends

code    segment
main    proc    far
        assume  cs:code,ds:data
start:
        ; --- 初始化 ---
        push    ds
        sub     ax,ax
        push    ax
        mov     ax,data
        mov     ds,ax

        ; --- 主循环准备 ---
        xor     ax,ax           ; AX = 0，作为当前要打印的数字
        mov     cx,101          ; 循环 101 次 (打印 0 到 100)
        mov     bx,0            ; BX 作为行计数器，记录当前行已打印个数

PrintLoop:
        ; --- 打印数字 ---
        push    ax              ; 保护 AX
        call    PrintAX         ; 调用子程序打印 AX 的十进制值
        pop     ax              ; 恢复 AX

        ; --- 打印空格 ---
        lea     dx,space
        mov     ah,9
        int     21h

        ; --- 处理换行逻辑 ---
        inc     bx              ; 当前行打印个数 +1
        cmp     bx,10           ; 检查是否满 10 个
        jnz     NotNewline      ; 如果不为 0 (未满 10 个)，跳过换行

        ; 满 10 个，打印换行符
        lea     dx,newline
        mov     ah,9
        int     21h
        mov     bx,0            ; 重置行计数器
NotNewline:
        inc     ax              ; 准备打印下一个数字
        loop    PrintLoop       ; CX--，循环

        ; --- 结束换行 ---
        lea     dx,newline
        mov     ah,9
        int     21h

        ; --- 退出 ---
        mov     ah,4Ch
        int     21h
main    endp

; --- 子程序：打印 AX 中的十进制数 ---
PrintAX proc    near
        push    ax
        push    bx
        push    cx
        push    dx

        mov     bx,10           ; 除数
        xor     cx,cx           ; 计数器清零
again:
        mov     dx,0            ; DX:AX 为被除数，DX 置 0
        xor     ah,ah           
        div     bx              ; 除法
        push    dx              ; 余数入栈
        inc     cx              ; 位数 +1
        test    ax,ax           ; 商是否为 0
        jnz     again           ; 不为 0 继续

print:
        pop     dx              ; 弹出数字
        add     dl,"0"          ; 转 ASCII
        mov     ah,2            ; 打印
        int     21h
        loop    print           ; 循环

        pop     dx
        pop     cx
        pop     bx
        pop     ax
        ret
PrintAX endp

code    ends
        end     start
data    segment
        A dw 15 dup(?)          ; 数组 A，容量 15
        B dw 20 dup(?)          ; 数组 B，容量 20
        C dw 15 dup(?)          ; 数组 C，存放结果
data    ends

code    segment
main proc far
        assume cs:code,ds:data
start:
        ; --- 初始化 ---
        push    ds
        sub     ax,ax
        push    ax
        mov     ax,data
        mov     ds,ax

        ; --- 初始化数组 A ---
        ; 逻辑：生成奇数序列 1, 3, 5, 7 ...
        lea     si,A
        mov     cx,15
        mov     bx,1            ; 起始值为 1
fillA:
        mov     [si],bx         ; 存入值
        add     bx,2            ; 值 +2
        add     si,2            ; 指针 +2
        loop    fillA

        ; --- 初始化数组 B ---
        ; 逻辑：生成连续序列 5, 6, 7, 8 ...
        lea     si,B
        mov     cx,20
        mov     bx,5            ; 起始值为 5
fillB:
        mov     [si],bx
        add     bx,1            ; 值 +1
        add     si,2
        loop    fillB
        
        ; -------------------------------
        ; --- 查找交集 (Intersection) ---
        ; -------------------------------
        lea     si,A            ; SI 指向 A
        mov     cx,15           ; 外层循环次数 (A 的长度)
        lea     di,C            ; DI 指向 C (存放位置)
        mov     bp,di           ; BP 备份 C 的起始地址，用于最后打印

outer:
        push    cx              ; 保存外层循环计数器
        mov     ax,[si]         ; 取出 A 的当前元素 -> AX
        
        lea     bx,B            ; BX 指向 B 的头部
        mov     cx,20           ; 内层循环次数 (B 的长度)

inner:
        cmp     ax,[bx]         ; 比较 A[i] (AX) 和 B[j] ([BX])
        jnz     not_equal       ; 如果不相等，跳过存储步骤

        ; -- 找到相同元素 --
        mov     [di],ax         ; 将该元素存入 C
        add     di,2            ; C 的指针后移
        jmp     found           ; 找到一个即可跳出内层循环

not_equal:
        add     bx,2            ; B 的指针后移
        loop    inner           ; 继续比较 B 的下一个元素

found:
        add     si,2            ; A 的指针后移
        pop     cx              ; 恢复外层循环计数器
        loop    outer           ; 继续处理 A 的下一个元素

        ; --------------------------
        ; --- 打印数组 C ---
        ; --------------------------
        mov     si,bp           ; SI 回到 C 的起始地址
printC:
        cmp     si,di           ; 比较当前打印指针 SI 和 结束指针 DI
        jz      exit            ; 如果相遇，说明打印完毕，退出

        mov     ax,[si]         ; 取出要打印的数
        call    printAX         ; 调用子程序打印十进制数
        
        mov     dl,' '          ; 打印一个空格分隔
        mov     ah,2
        int     21h

        add     si,2            ; 指针后移
        jmp     printC          ; 继续打印下一个

exit:
        mov     ah,4Ch
        int     21h
main    endp

; --- 子程序：以十进制打印 AX ---
printAX proc near
        push    ax
        push    bx
        push    cx
        push    dx

        mov     bx,10
        mov     cx,0
again:
        mov     dx,0
        div     bx
        push    dx
        inc     cx
        test    ax,ax
        jnz     again
print:
        pop     dx
        add     dl,'0'
        mov     ah,2
        int     21h
        loop    print

        pop     dx
        pop     cx
        pop     bx
        pop     ax
        ret
printAX endp
code    ends
        end     start
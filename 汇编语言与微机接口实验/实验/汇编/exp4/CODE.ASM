datarea    segment
        DATA dw 100 dup(?)      ; 定义长度为 100 的字数组
datarea    ends

code    segment
main proc far
        assume cs:code,ds:datarea
start:
        ; --- 初始化 ---
        push    ds
        sub     ax,ax
        push    ax
        mov     ax,datarea
        mov     ds,ax

        ; --- 步骤 1: 填充测试数据 (100, 99, ... 1) ---
        mov     cx,100          ; 循环 100 次
        lea     si,DATA         ; SI 指向数组首地址
fill:
        mov     ax,cx           ; 将 CX 的值 (100, 99...) 放入 AX
        mov     [si],ax         ; 存入数组
        add     si,2            ; 移动指针到下一个字
        loop    fill

        ; --- 步骤 2: 查找满足条件的数 ---
        mov     cx,100          ; 重置循环计数
        lea     si,DATA         ; 重置指针
        mov     ax,0FFFFh       ; 将 AX 初始化为最大可能值 (作为当前最小值保存者)

next:
        mov     bx,[si]         ; 读取数组当前元素到 BX
        
        ; 逻辑判断 1: 奇偶校验
        test    bx,1            ; 测试最低位
        jnz     skip            ; 如果最低位是 1 (结果不为 0)，说明是奇数，跳转跳过

        ; 逻辑判断 2: 大小比较
        cmp     bx,ax           ; 比较当前偶数 (BX) 和当前记录的最小值 (AX)
        ja      skip            ; 如果 BX > AX (ja 为无符号大于)，跳转跳过
                                ; 这意味着我们只保留比 AX 更小的数

        mov     ax,bx           ; 更新 AX 为当前找到的更小的偶数

skip:
        add     si,2            ; 指针后移
        loop    next            ; 继续下一个

        ; --- 步骤 3: 打印结果 ---
        call    printAX         ; 打印 AX 中的值

        ; --- 退出 ---
        pop     ds
        mov     ah,4Ch
        int     21h
main    endp

; --- 子程序：以十进制打印 AX 中的数值 ---
printAX proc near
        push    ax
        push    bx
        push    cx
        push    dx

        mov     bx,10           ; 除数 10
        mov     cx,0            ; 计数器清零
again:
        mov     dx,0            ; DX:AX 做被除数，DX 置 0
        div     bx              ; AX / 10, 商在 AX, 余数在 DX
        push    dx              ; 保存余数
        inc     cx              ; 计数器 +1
        test    ax,ax           ; 检查商是否为 0
        jnz     again           ; 不为 0 继续除

print:
        pop     dx              ; 弹出余数 (从高位到低位)
        add     dl,'0'          ; 转换为 ASCII 字符
        mov     ah,2            ; 打印
        int     21h
        loop    print           ; 循环 CX 次

        pop     dx
        pop     cx
        pop     bx
        pop     ax
        ret
printAX endp

code    ends
        end     start
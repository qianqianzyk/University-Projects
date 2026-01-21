data segment
data ends

code segment
main proc far
        assume cs:code,ds:data
start:
        push    ds              ; 保存旧的数据段寄存器
        sub     ax,ax           ; AX 清零
        push    ax              ; 压栈 0，为了后续 ret 返回 DOS
        mov     ax,data         ; 加载数据段地址
        mov     ds,ax           ; 初始化 DS

        ; --- 主逻辑开始 ---
        mov     ax,0ABCDh       ; 将要处理的 16 进制数装入 AX (AH=AB, AL=CD)

        ; 1. 提取第 4 部分 (低位 AL 的低 4 位 -> D)
        mov     dl,al           ; 复制 AL (CD) 到 DL
        and     dl,0Fh          ; 屏蔽高 4 位，只保留低 4 位 (得到 0D)

        ; 2. 提取第 3 部分 (低位 AL 的高 4 位 -> C)
        mov     cl,al           ; 复制 AL (CD) 到 CL
        and     cl,0F0h         ; 屏蔽低 4 位 (得到 C0)
        shr     cl,1            ; 右移 4 次，将 C0 变成 0C
        shr     cl,1
        shr     cl,1
        shr     cl,1

        ; 3. 提取第 2 部分 (高位 AH 的低 4 位 -> B)
        mov     bl,ah           ; 复制 AH (AB) 到 BL
        and     bl,0Fh          ; 屏蔽高 4 位，只保留低 4 位 (得到 0B)

        ; 4. 提取第 1 部分 (高位 AH 的高 4 位 -> A)
        mov     al,ah           ; 复制 AH (AB) 到 AL
        and     al,0F0h         ; 屏蔽低 4 位 (得到 A0)
        shr     al,1            ; 右移 4 次，将 A0 变成 0A
        shr     al,1
        shr     al,1
        shr     al,1            ; 现在 AL 中存放的是第一个数字 A

        ; --- 依次打印 4 个部分 ---
        
        ; 打印第 1 部分 (A)
        call    printBin        ; 调用子程序打印 AL 中的二进制
        call    printSpace      ; 打印空格

        ; 打印第 2 部分 (B)
        mov     al,bl           ; 将保存的 B (在 BL 中) 移入 AL
        call    printBin        ; 打印
        call    printSpace      ; 打印空格

        ; 打印第 3 部分 (C)
        mov     al,cl           ; 将保存的 C (在 CL 中) 移入 AL
        call    printBin        ; 打印
        call    printSpace      ; 打印空格

        ; 打印第 4 部分 (D)
        mov     al,dl           ; 将保存的 D (在 DL 中) 移入 AL
        call    printBin        ; 打印
        call    printSpace      ; 打印空格

        ; --- 打印回车换行 ---
        mov     dl,0Dh          ; 回车符 CR
        mov     ah,2
        int     21h
        mov     dl,0Ah          ; 换行符 LF
        mov     ah,2
        int     21h

        ; --- 程序退出 ---
        pop     ds              ; 恢复 DS
        mov     ah,4Ch          ; DOS 功能调用：带返回码结束
        int     21h
main endp

; --- 子程序：打印 4 位二进制 ---
printBin proc near
        push    bx              ; 保护寄存器
        push    cx
        push    dx
        
        mov     cx,4            ; 循环 4 次 (打印 4 位)
        mov     bl,al           ; 将数值存入 BL
        
        ; 将低 4 位移到高 4 位，以便使用 SHL 指令将最高位移入 CF 标志位
        shl     bl,1
        shl     bl,1
        shl     bl,1
        shl     bl,1
        
bin_loop:
        shl     bl,1            ; 左移一位，最高位进入 CF
        jc      bin_one         ; 如果 CF=1，跳转去打印 '1'
        
        ; 如果 CF=0，打印 '0'
        mov     dl,30h          ; '0' 的 ASCII 码
        jmp     bin_print
        
bin_one:
        mov     dl,31h          ; '1' 的 ASCII 码
        
bin_print:
        mov     ah,2            ; DOS 功能：打印字符
        int     21h
        loop    bin_loop        ; CX--, 如果 CX!=0 继续循环
        
        pop     dx              ; 恢复寄存器
        pop     cx
        pop     bx
        ret
printBin endp
        
; --- 子程序：打印空格 ---
printSpace proc near
        push    ax
        push    dx
        mov     dl,20h          ; 空格的 ASCII 码
        mov     ah,2
        int     21h
        pop     dx
        pop     ax
        ret
printSpace endp

code ends
     end start
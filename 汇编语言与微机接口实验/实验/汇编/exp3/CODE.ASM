data segment
        N       dw      ?       ; 定义变量 N，用于存储输入的次数
data ends

code segment
main proc far
        assume cs:code,ds:data
start:
        ; --- 初始化 DS ---
        push    ds
        sub     ax,ax
        push    ax
        mov     ax,data
        mov     ds,ax

        ; --- 读取输入 ---
        mov     ah,1            ; DOS 功能 1：读取一个字符到 AL
        int     21h
        
        ; --- 校验输入 ---
        sub     al,30h          ; 将 ASCII 码转换为数字 (如 '3' -> 3)
        jb      stop            ; 如果结果 < 0 (即原字符 < '0')，跳转退出
        cmp     al,9
        ja      stop            ; 如果结果 > 9 (即原字符 > '9')，跳转退出
        
        ; --- 准备循环 ---
        cbw                     ; 将 AL 扩展为 AX (Clear Byte into Word)
        mov     N,ax            ; 将数字存入 N
        mov     cx,N            ; 将循环计数器 CX 设置为 N

bell:
        mov     dl,07h          ; 07h 是响铃字符 (BEL) 的 ASCII 码
        mov     ah,2            ; DOS 功能 2：输出字符
        int     21h             ; 执行输出（响铃）
        loop    bell            ; CX--, 若 CX!=0 则跳转继续响铃

stop:
        ; --- 退出程序 ---
        pop     ds
        mov     ah,4Ch
        int     21h
main    endp
code    ends
        end     start
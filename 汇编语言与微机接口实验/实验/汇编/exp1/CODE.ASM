data segment
data ends

code segment
main proc far
        assume cs:code,ds:data
start:
        push    ds              ; 保存旧的DS
        sub     ax,ax
        push    ax              ; 压入0，以便ret返回
        mov     ax,data
        mov     ds,ax           ; 初始化DS

read_char:
        mov     ah,1            ; DOS中断：读取一个字符到AL
        int     21h
        cmp     al,1Bh          ; 比较是否为 ESC
        je      exit            ; 如果是ESC，跳转到退出
        cmp     al,61h          ; 比较是否小于 'a'
        jl      exit            ; 非小写字母，退出
        cmp     al,7Ah          ; 比较是否大于 'z'
        jg      exit            ; 非小写字母，退出
        sub     al,20h          ; 小写转大写的核心：ASCII码减去20h
exec:
        mov     dl,al           ; 将处理后的字符放入DL准备输出
        mov     ah,2            ; DOS中断：输出DL中的字符
        int     21h

        ; 输出回车换行
        mov     dl,0Dh
        mov     ah,2
        int     21h
        mov     dl,0Ah
        mov     al,2
        int     21h

        jmp     read_char       ; 循环继续读取下一个字符
exit:
        pop     ds
        mov     ah,4Ch          ; DOS中断：程序终止
        int     21h

main    endp
code    ends
        end     start
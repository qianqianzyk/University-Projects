data segment
        ; --- 定义菜单字符串 ---
        menu    db      "====== Menu ======",0Dh,0Ah
                db      "1. Convert lowercase to uppercase",0Dh,0Ah
                db      "2. Convert uppercase to lowercase",0Dh,0Ah
                db      "3. Exit",0Dh,0Ah,'$'
        
        ; --- 定义提示信息 ---
        option  db      0Dh,0Ah,"Enter option:$"
        prompt  db      0Dh,0Ah,"Enter a string:$"
        output  db      0Dh,0Ah,"Converted result:$"
        newline db      0Dh,0Ah,"$"

        ; --- 定义输入缓冲区 ---
        ; 第一个字节 50 表示缓冲区最大长度
        ; 第二个字节 0 用于 DOS 填入实际读取的字符数
        ; 后面的 50 dup(?) 用于存放字符串内容
        buffer  db      50,0,50 dup(?)
data    ends

code    segment
main    proc    far
        assume  cs:code,ds:data
start:
        ; --- 初始化数据段 ---
        push    ds
        sub     ax,ax
        push    ax
        mov     ax,data
        mov     ds,ax

MenuLoop:
        ; --- 显示主菜单 ---
        lea     dx,menu         ; 加载菜单字符串地址
        mov     ah,9            ; DOS 功能 9：打印字符串
        int     21h

        ; --- 读取用户选项 ---
        mov     ah,1            ; DOS 功能 1：读取单个字符
        int     21h
        mov     bl,al           ; 将输入的字符保存到 BL

        ; --- 打印换行 ---
        lea     dx,newline
        mov     ah,9
        int     21h

        ; --- 判断选项 ---
        cmp     bl,"1"
        jz      GetString       ; 选 1，跳转去读取字符串
        cmp     bl,"2"
        jz      GetString       ; 选 2，也跳转去读取字符串（处理逻辑在后面区分）
        cmp     bl,"3"
        jz      ExitProgram     ; 选 3，退出程序
        jmp     MenuLoop        ; 输入无效，重新显示菜单

GetString:
        ; --- 提示输入字符串 ---
        lea     dx,prompt
        mov     ah,9
        int     21h

        ; --- 读取字符串 ---
        lea     dx,buffer       ; DX 指向缓冲区结构首地址
        mov     ah,0Ah          ; DOS 功能 0Ah：缓冲输入
        int     21h

        ; --- 打印输出前缀 ---
        lea     dx,output
        mov     ah,9
        int     21h

        ; --- 再次根据保存的选项 BL 跳转到具体处理逻辑 ---
        cmp     bl,"1"
        jz      LowerToUpper    ; 小写转大写
        cmp     bl,"2"
        jz      UpperToLower    ; 大写转小写

LowerToUpper:
        ; --- 小写转大写逻辑 ---
        mov     cl,buffer+1     ; 获取实际输入的字符长度
        lea     si,buffer+2     ; SI 指向字符串内容的起始位置
again1:
        mov     al,[si]         ; 读取当前字符
        cmp     al,"a"
        jb      skip1           ; 小于 'a'，不是小写字母，跳过
        cmp     al,"z"
        ja      skip1           ; 大于 'z'，不是小写字母，跳过
        
        sub     al,20h          ; 核心转换：减去 20H (例如 'a'-20h = 'A')
        mov     [si],al         ; 写回缓冲区
skip1:
        inc     si              ; 指针移动
        loop    again1          ; 循环处理下一个字符
        jmp     PrintResult     ; 处理完毕，跳转打印

UpperToLower:
        ; --- 大写转小写逻辑 ---
        mov     cl,buffer+1     ; 获取长度
        lea     si,buffer+2     ; 获取内容首地址
again2:
        mov     al,[si]         ; 读取字符
        cmp     al,"A"
        jb      skip2           ; 小于 'A'，跳过
        cmp     al,"Z"
        ja      skip2           ; 大于 'Z'，跳过
        
        add     al,20h          ; 核心转换：加上 20H (例如 'A'+20h = 'a')
        mov     [si],al         ; 写回
skip2:
        inc     si
        loop    again2
        jmp     PrintResult

PrintResult:
        ; --- 打印处理后的字符串 ---
        ; 为了使用 DOS 功能 9 打印，需要在字符串末尾添加 '$'
        lea     si,buffer+2     ; 字符串开始
        mov     cl,buffer+1     ; 长度
        xor     ch,ch           ; 清空高位，确保 CX 为纯长度
        add     si,cx           ; SI 指向字符串末尾的下一个字节
        mov     byte ptr [si],"$" ; 写入结束符 '$'

        lea     dx,buffer+2     ; 指向字符串内容
        mov     ah,9
        int     21h

        ; --- 打印换行 ---
        lea     dx,newline
        mov     ah,9
        int     21h

        jmp     MenuLoop        ; 处理完成，回到主菜单

ExitProgram:
        ; --- 退出程序 ---
        mov     ah,4Ch
        int     21h
main    endp
code    ends
        end     start
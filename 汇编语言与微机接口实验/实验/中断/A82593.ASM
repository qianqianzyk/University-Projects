;=========================================================
; 文件名: A82593.ASM
; 功能描述: 8259级联中断实验，中断源为主片8259的IR7,
; 从片8259的IR1。从片8259通过主片8259的IR2
; 进行级联。
; 主片每产生一次中断输出显示一个字符M7,从片
; 每产生一次中断输出显示一个字符S1。
;=========================================================
SSTACK SEGMENT STACK
    DW 32 DUP(?)
SSTACK ENDS

CODE SEGMENT
    ASSUME CS:CODE
START: 
    PUSH DS
    MOV AX, 0000H
    MOV DS, AX
    
    ;设置主片 IRQ7 中断向量
    MOV AX, OFFSET MIR7
    MOV SI, 003CH       ;取中断入口地址
    MOV [SI], AX        ;填IRQ7的偏移矢量
    MOV AX, CS          ;段地址
    MOV SI, 003EH
    MOV [SI], AX        ;填IRQ7的段地址矢量
    
    ;设置从片 IRQ1 中断向量 (基址 C0H + 1*4 = C4H)
    MOV AX, OFFSET SIR1
    MOV SI, 00C4H
    MOV [SI], AX
    MOV AX, CS
    MOV SI, 00C6H
    MOV [SI], AX
    
    CLI
    POP DS              
    
    ;初始化主片8259
    MOV AL, 11H
    OUT 20H, AL         ;ICW1
    MOV AL, 08H
    OUT 21H, AL         ;ICW2
    MOV AL, 04H
    OUT 21H, AL         ;ICW3 (主片IR2接从片)
    MOV AL, 01H
    OUT 21H, AL         ;ICW4
    
    ;初始化从片8259
    MOV AL, 11H
    OUT 0A0H, AL        ;ICW1
    MOV AL, 30H         ;ICW2 (从片中断基址 30H? 
                        ;注意:这里代码用30H对应C0H矢量地址? 
                        ;原文注释写的是C0H~C3H, 此处配置需根据实验系统实际情况)
                        ;从代码看向量填的是C4H，说明中断号应是31H? 
                        ;或者中断向量表基址设置为了 0C0H? 
                        ;通常30H对应中断号30H，向量地址30H*4=C0H
    OUT 0A1H, AL        
    MOV AL, 02H
    OUT 0A1H, AL        ;ICW3 (从片接主片IR2 -> ID=02H)
    MOV AL, 01H
    OUT 0A1H, AL        ;ICW4
    
    MOV AL, 0FDH        ;从片 OCW1 (1111 1101, 开放IR1)
    OUT 0A1H, AL
    
    MOV AL, 6BH         ;主片 OCW1 (0110 1011, 开放IR2, IR4, IR7)
    OUT 21H, AL
    
    STI

AA1: 
    NOP
    JMP AA1

MIR7: 
    STI
    CALL DELAY
    MOV AX, 014DH
    INT 10H             ;显示字符 M
    MOV AX, 0137H
    INT 10H             ;显示字符 7
    MOV AX, 0120H
    INT 10H
    MOV AL, 20H
    OUT 20H, AL         ;中断结束命令
    IRET

SIR1: 
    STI
    CALL DELAY
    MOV AX, 0153H
    INT 10H             ;显示字符 S
    MOV AX, 0131H
    INT 10H             ;显示字符 1
    MOV AX, 0120H
    INT 10H
    MOV AL, 20H
    OUT 0A0H, AL        ;发送EOI给从片
    OUT 20H, AL         ;发送EOI给主片
    IRET

DELAY: 
    PUSH CX
    MOV CX, 0F00H
AA0: 
    PUSH AX
    POP AX
    LOOP AA0
    POP CX
    RET

CODE ENDS
END START
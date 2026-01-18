#include <iostream>
using namespace std;

// x * y 采用移位加法实现
int MUL(int x, int y)
{
    int res = 0;

    while (y != 0)
    {
        if (y & 1) // 判断 y 是否是奇数
        {
            res += x;
        }
        x = x << 1; // x 左移一位
        y = y >> 1; // y 右移一位
    }

    return res;
}

// x / y 用减法实现
int DIV(int x, int y)
{
    int res = 1;

    while (x - y != 0)
    {
        x = x - y;
        res = res + 1;
    }

    return res;
}

// C(n, m) = n! / (m! * (n - m)!)
int C(int n, int m)
{
    int res = 1;

    for (int i = 1; i <= m; ++i)
    {
        res = DIV(MUL(res, (n - m + i)), i); // res = res * (n - m + i) / i
    }

    return res;
}

// 寄存器
int R0, R1, R2, R3;
// 主存储器的三个地址单元
int P60H, P61H, P62H;

int main()
{
    cout << DIV(4, 1);
    cin >> R0; // 模拟IN R0, 00H（输入 n）
    cin >> R1; // 模拟IN R1, 00H（输入 m）
    P60H = R1; // 模拟STA R1, 60H（将 m 保存到主存 60H）

    R0 = R0 - R1; // 模拟SUB R0, R1（R0 = n - m）

    R1 = 0x01; // 模拟LDI R1, 01H（R1 初始化为1，存储结果 res）
    R3 = 0x01; // 模拟LDI R3, 01H（R3 初始化为1，存储计数 i）

LOOP:
    P61H = R3; // 模拟STA R3, 61H（把当前 i 值存入主存地址 61H）
    R2 = P60H; // 模拟LAD R2, 60H（取回 m 值，放入 R2）
    R2++;      // 模拟INC R2（R2 = m + 1）
    // 模拟CMP R2, R3（i > m）
    if (R2 - R3 == 0)
    {
        goto RESULT; // 模拟BZC RESULT
    }

    R0++;      // 模拟INC R0（R0 = n - m + i）
    P62H = R0; // 模拟STA R0, 62H（将 R0 当前值备份至主存 62H）

    // 开始乘法计算
    R2 = 0x00; // 模拟LDI R2, 00H（R2 = 0，用于存储乘法结果）

MULLOOP:
    R3 = 0x00; // 模拟LDI R3, 00H（清零 R3，用于比较是否乘法结束）
    // 模拟CMP R0, R3
    if ((R0 - R3) == 0)
    {
        goto MULEND; // 模拟BZC MULEND
    }

    R3 = 0x01; // 模拟LDI R3, 01H（R3 = 1，用于判断乘数奇偶性）
    // 模拟TEST R0, R3
    if ((R0 & R3) == 0) // 若 R0 为偶数
    {
        goto EVEN; // 模拟BZC EVEN（跳到偶数处理部分）
    }

    R2 = R2 + R1; // 模拟ADD R2, R1（累加乘法）

EVEN:
    R1 = R1 << 1; // 模拟SHL R1（被乘数左移1位）
    R0 = R0 >> 1; // 模拟SHR R0（乘数右移1位）
    goto MULLOOP; // 模拟JMP MULLOOP（重复乘法过程）

MULEND:
    R1 = R2; // 模拟MOV R1 R2（把乘法结果放入 R1）

    R0 = P62H; // 模拟LAD R0, 62H（还原 R0 原值）
    R3 = P61H; // 模拟LAD R3, 61H（还原 R3 原值）
    // 乘法计算结束

    // 开始除法计算
    R2 = 0x01; // 模拟LDI R2, 01H（R2 初始为 1，用作除法结果）

DIVLOOP:
    R1 = R1 - R3; // 模拟SUB R1, R3（被除数减去除数）
    if (R1 == 0)
    {
        goto DIVEND; // 模拟BZC DIVEND（除尽）
    }

    R2++;         // 模拟INC R2（除法结果加 1）
    goto DIVLOOP; // 模拟JMP DIVLOOP（重复除法过程）

DIVEND:
    R1 = R2; // 模拟MOV R1 R2（除法结果保存在 R1）
    // 除法计算结束

    R3++; // 模拟INC R3（i++）

    goto LOOP; // 模拟JMP LOOP（进入下一轮循环）

RESULT:
    cout << R1; // 模拟OUT R1, 40H（输出结果）

    system("pause");
    return 0;
}

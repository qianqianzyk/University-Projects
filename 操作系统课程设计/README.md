# 操作系统课程设计

## 项目简介

本项目包含两个核心部分的实现，旨在从用户态工具（Shell）到内核态功能（Pintos User Programs）全方位深入理解操作系统的核心机制：

1. **简易 Shell 实现 (Simple Shell)**: 基于 CS 162 HW1，实现了一个支持基本命令解析、内建命令、重定向和信号处理的命令行解释器
2. **Pintos 用户程序 (Pintos User Programs)**: 基于 CS 162 Project 1，完善了 Pintos 教学操作系统的内核，使其能够加载并运行用户态程序，并提供标准的系统调用接口

## 开发环境

- **操作系统**: Ubuntu 18.04 LTS
- **模拟器**: QEMU (用于运行 Pintos 内核)
- **调试器**: GDB (配合 QEMU 进行内核调试)
- **编译器**: GCC
- **构建工具**: Make

## 项目一：简易 Shell 实现 (Simple Shell)

### 核心功能

- **命令解析**: 使用 `tokenizer` 对用户输入进行分词，支持带参数的程序执行
- **内建命令 (Built-ins)**:
  - `cd`: 切换当前工作目录 (使用 `chdir`)
  - `pwd`: 显示当前工作目录 (使用 `getcwd`)
  - `exit`: 退出 Shell
- **外部程序执行**: 通过 `fork()` 创建子进程，并使用 `execv()` 加载外部可执行文件
- **路径解析**: 自动遍历 `PATH` 环境变量，查找可执行文件路径
- **I/O 重定向**: 支持 `>` (输出重定向) 和 `<` (输入重定向)，利用 `dup2()` 接管标准输入输出流
- **信号处理**: 能够忽略或处理 `SIGINT` (Ctrl-C) 等信号，防止 Shell 被意外终止

------

## 项目二：Pintos 用户程序 (Pintos User Programs)

### 核心功能

在此项目中，我们修改了 Pintos 内核，使其具备了运行 x86 用户程序的能力

1. **参数传递**:
   - 重写了 `process_execute` 和 `start_process`
   - 实现了符合 x86 栈帧规范（80x86 Calling Convention）的**参数压栈**机制
   - 使用 `strtok_r` 分割命令行参数，并正确对齐内存地址，确保参数能被用户程序正确读取
2. **系统调用**:
   - 实现了 `syscall_handler`，支持以下核心调用：
     - **进程控制**: `exec` (执行新程序), `wait` (等待子进程), `exit` (进程退出)
     - **文件操作**: `create`, `remove`, `open`, `filesize`, `read`, `write`, `seek`, `tell`, `close`
   - **同步机制**:
     - 引入全局文件锁 `filesys_lock` 保证文件系统操作的线程安全
     - 使用二元信号量 `load_done` 实现父进程阻塞等待子进程加载结果
3. **内存安全**:
   - 实现了防御性编程，验证用户指针的合法性
   - 设计了 `check_ptr_and_size` 和 `check_string` 函数，防止用户程序访问内核空间或未映射的内存区域，避免 OS 崩溃
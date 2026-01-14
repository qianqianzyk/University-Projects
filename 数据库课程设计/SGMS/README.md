# SGMS - Student Grade Management System

SGMS（学生成绩管理系统）是一个基于Go语言开发的后端服务，旨在为学校或教育机构提供高效、可靠的学生成绩、课程、教师、班级等信息的管理能力。项目采用Gin作为Web框架，GORM作为ORM，支持PostgreSQL数据库。

## 主要功能

- 学生、教师、管理员多角色支持
- 学生信息管理
- 教师信息管理
- 课程与班级管理
- 成绩录入与查询
- 地区（省、市）管理
- 用户登录与权限控制
- 跨域请求支持（CORS）

## 技术栈

- **语言**：Go 1.23.x
- **Web框架**：Gin v1.10.1
- **ORM**：GORM v1.25.10
- **数据库驱动**：gorm.io/driver/postgres v1.6.0
- **配置管理**：Viper
- **依赖管理**：Go Modules

## 目录结构

```
SGMS/
  ├── app/
  │   ├── apiException/         # 自定义异常处理
  │   ├── controllers/          # 控制器
  │   ├── midwares/             # 中间件（CORS、错误处理等）
  │   ├── model/                # 数据模型
  │   └── utils/                # 工具类
  ├── config/
  │   ├── config/               # 配置加载
  │   ├── database/             # 数据库连接
  │   └── router/               # 路由注册
  ├── main.go                   # 程序入口
  ├── go.mod                    # 依赖管理
  ├── config.example.yaml       # 配置文件示例
  └── README.md                 # 项目说明
```

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/qianqianzyk/database_project.git
cd SGMS
```

### 2. 配置环境

复制配置文件模板并根据实际情况修改：

```bash
cp config.example.yaml config.yaml
```

编辑`config.yaml`，填写数据库等相关配置。

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 运行项目

```bash
go run main.go
```

或编译后运行：

```bash
go build -o sgms main.go
./sgms
```

### 5. 访问接口

默认服务端口可在main文件中设置，启动后可通过Apifox等工具访问API接口。

## 注意事项

如需进一步了解接口文档、部署方式或有其他问题，请联系项目维护者。 
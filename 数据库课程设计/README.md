# database_project

2024/2025数据库课程设计

## 项目简介

本项目为“学生成绩管理系统（SGMS）”，包含后端（SGMS，基于Go）和前端（SGMS-FRONT，基于Vue 3 + TypeScript）。系统支持学生、教师、管理员多角色，提供学生成绩、课程、班级、教师、地区等信息的高效管理。

- **后端**：Go + Gin + GORM + PostgreSQL，提供RESTful API。
- **前端**：Vue 3 + TypeScript + Element Plus，提供现代化Web界面。

---

## 主要功能
- 学生、教师、管理员多角色支持
- 学生/教师/课程/班级/成绩/地区（省市）信息管理
- 成绩录入与查询
- 用户登录与权限控制
- 跨域请求支持（CORS）
- 数据统计与可视化（前端）

---

## 技术栈

### 后端
- Go 1.23.x
- Gin v1.10.1（Web框架）
- GORM v1.25.10（ORM）
- PostgreSQL（数据库）
- Viper（配置管理）

### 前端
- Vue 3 ^3.5.13
- TypeScript ~5.7.2
- Vite ^6.2.0
- Element Plus ^2.10.2
- Pinia ^3.0.3
- Vue Router ^4.5.1
- Axios ^1.10.0
- ECharts ^5.6.0

---

## 目录结构

```
database_project/
├── SGMS/           # 后端Go服务
│   ├── app/
│   ├── config/
│   ├── main.go
│   ├── go.mod
│   └── ...
├── SGMS-FRONT/     # 前端Vue项目
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── ...
└── README.md       # 项目总说明
```

---

## 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/qianqianzyk/database_project.git
cd database_project
```

### 2. 后端（SGMS）启动
```bash
cd SGMS
cp config.example.yaml config.yaml   # 复制配置模板并根据实际情况修改
# 编辑config.yaml，填写数据库等信息

go mod tidy                         # 安装依赖
go run main.go                      # 启动后端服务
# 或编译后运行：
go build -o sgms main.go
./sgms
```

### 3. 前端（SGMS-FRONT）启动
```bash
cd SGMS-FRONT
npm install                         # 安装依赖
npm run dev                         # 启动开发服务器
# 构建生产包：
npm run build
```

### 4. 访问系统
- 前端开发环境默认：http://localhost:5173
- 后端API默认：http://localhost:8888

前端通过Vite代理将`/api`请求转发到后端，跨域配置见`SGMS-FRONT/vite.config.ts`。

---

## 其他说明
- 后端API接口文档、详细部署方式等请见`SGMS/README.md`
- 前端开发细节、页面结构等请见`SGMS-FRONT/README.md`
- 如有问题请联系项目维护者。

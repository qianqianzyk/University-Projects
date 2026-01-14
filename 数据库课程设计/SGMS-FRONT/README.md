# SGMS-FRONT 前端项目

## 项目简介
本项目是一个基于 Vue 3 + TypeScript + Element Plus 的学生成绩管理系统前端项目，支持学生、教师、管理员多角色的成绩、课程、班级等信息的管理与统计分析。

## 主要技术栈及版本
- [Vue 3](https://vuejs.org/) ^3.5.13
- [TypeScript](https://www.typescriptlang.org/) ~5.7.2
- [Vite](https://vitejs.dev/) ^6.2.0
- [Element Plus](https://element-plus.org/) ^2.10.2
- [Pinia](https://pinia.vuejs.org/) ^3.0.3
- [Vue Router](https://router.vuejs.org/) ^4.5.1
- [Axios](https://axios-http.com/) ^1.10.0
- [ECharts](https://echarts.apache.org/) ^5.6.0
- [Vue-ECharts](https://github.com/ecomfe/vue-echarts) ^7.0.3

## 启动和构建方法

1. 安装依赖：
   ```bash
   npm install
   ```
2. 启动开发服务器：
   ```bash
   npm run dev
   ```
3. 构建生产包：
   ```bash
   npm run build
   ```

## 跨域处理说明
本项目通过 Vite 的 dev server 代理功能，将前端的 `/api` 请求代理到后端服务器（如 `http://localhost:8888`），实现开发环境下的跨域请求，相关配置见 `vite.config.ts`。

## 目录结构简述
```
SGMS-FRONT/
├── index.html                # 入口 HTML 文件
├── package.json              # 项目依赖与脚本
├── src/                      # 源码目录
│   ├── api/                  # API 请求封装
│   ├── assets/               # 静态资源
│   ├── components/           # 通用组件
│   ├── router/               # 路由配置
│   ├── stores/               # 状态管理（Pinia）
│   ├── types/                # TypeScript 类型定义
│   ├── views/                # 页面视图（按角色分子目录）
│   ├── App.vue               # 根组件
│   └── main.ts               # 入口文件
├── vite.config.ts            # Vite 配置
└── ...
```

## 其他
如需后端配合或详细接口文档，请联系项目开发者。 
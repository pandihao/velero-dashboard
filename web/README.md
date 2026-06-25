# Velero API Console

基于 Vue 3 + Element Plus 的 Velero 管理控制台。

## 功能特性

### 1. 备份管理
- 创建备份
  - 支持自定义备份名称
  - 支持选择命名空间（多选，支持 all 选项）
  - 支持选择资源类型（多选，支持 all 选项）
  - 支持设置 TTL 保留时间
  - 支持卷快照选项
- 查看备份列表（实时状态、错误/警告统计）
- 查看备份详情
- 删除备份

### 2. 恢复管理
- 创建恢复任务
  - 从已有备份中选择
  - 支持恢复到原命名空间或指定命名空间
  - 支持选择恢复的资源类型（默认全部）
  - 支持持久卷恢复选项
- 查看恢复任务列表和状态
- 查看恢复详情

### 3. 定时备份任务
- 创建定时备份
  - 支持 Cron 表达式（提供常用模板）
  - 支持选择命名空间和资源类型
  - 支持设置 TTL
- 查看定时任务列表
- 查看最后备份时间
- 删除定时任务

## 快速开始

### 安装依赖

```bash
cd web
npm install
```

### 开发模式

```bash
npm run dev
```

访问: http://localhost:3000

开发模式会自动代理 API 请求到 `http://localhost:8080`

### 生产构建

```bash
npm run build
```

构建产物在 `dist/` 目录。

## 配置说明

### API 代理配置

编辑 `vite.config.js`:

```javascript
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // 修改为你的后端地址
        changeOrigin: true
      }
    }
  }
})
```

## 项目结构

```
web/
├── src/
│   ├── api/              # API 请求封装
│   │   └── velero.js
│   ├── views/            # 页面组件
│   │   ├── BackupView.vue
│   │   ├── RestoreView.vue
│   │   └── ScheduleView.vue
│   ├── App.vue           # 根组件
│   ├── router.js         # 路由配置
│   └── main.js           # 入口文件
├── index.html
├── vite.config.js        # Vite 配置
└── package.json
```

## 使用说明

### 1. 启动后端服务

```bash
# 在项目根目录
./velero-api-server --kubeconfig ~/.kube/config --port 8080
```

### 2. 启动前端服务

```bash
cd web
npm run dev
```

### 3. 访问控制台

浏览器打开: http://localhost:3000

## 功能截图

### 备份管理
- 支持多命名空间选择（all 选项表示全部）
- 支持多资源类型选择（all 选项表示全部）
- 实时显示备份状态（新建/进行中/已完成/失败）

### 恢复管理
- 从备份列表中选择要恢复的备份
- 自动生成恢复任务名称
- 留空命名空间则恢复到原位置
- 留空资源类型则恢复全部资源

### 定时任务
- 提供常用 Cron 表达式模板
- 显示最后备份时间
- 支持启用/禁用状态查看

## 技术栈

- Vue 3 (Composition API)
- Element Plus (UI 组件库)
- Vue Router (路由)
- Axios (HTTP 客户端)
- Vite (构建工具)

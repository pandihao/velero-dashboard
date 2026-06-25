# Velero API Server

完整的 Velero 备份管理解决方案，包含 REST API 服务和 Web 控制台。支持多集群管理，可以从单一界面管理多个 Kubernetes 集群的备份和恢复操作。

## 项目组成

### 1. REST API 服务（Go）
- 基于 Gin 框架
- 封装 Velero CRD 操作
- 多集群支持（通过查询参数或 HTTP Header 选择目标集群）
- 提供完整的 Swagger 文档

### 2. Web 控制台（Vue 3）
- 基于 Vue 3 + Element Plus
- 可视化管理界面
- 支持备份、恢复、定时任务管理
- 顶部集群选择器，轻松切换目标集群
<img width="1893" height="781" alt="image" src="https://github.com/user-attachments/assets/61987f7d-b65a-4780-a696-64b2dafddad6" />
<img width="1896" height="841" alt="image" src="https://github.com/user-attachments/assets/d59e2d48-ddcb-4516-bb73-5b1f2e9e1987" />
<img width="1887" height="892" alt="image" src="https://github.com/user-attachments/assets/ad9fb9a8-fb5e-4882-818f-33a76d4bd020" />





## 快速开始

### 前置要求

- Go 1.22+
- Node.js 16+
- Kubernetes 集群（已安装 Velero）
- kubeconfig 文件（单集群或多集群）

### 1. 启动 API 服务

#### 单集群模式

```bash
# 编译
go build -o velero-api-server cmd/server/main.go

# 启动服务（指定单个 kubeconfig）
./velero-api-server --kubeconfig ~/.kube/config --port 8080

# 或使用环境变量
export KUBECONFIG=~/.kube/config
export VELERO_NAMESPACE=velero
./velero-api-server
```

#### 多集群模式

```bash
# 准备多个 kubeconfig 文件
mkdir -p .kube
cp /path/to/cluster1-config .kube/config-fat01
cp /path/to/cluster2-config .kube/config-fat02

# 启动服务（自动加载目录中的所有 kubeconfig）
./velero-api-server --kubeconfig-dir .kube --port 8080

# 或使用环境变量
export KUBECONFIG_DIR=.kube
export VELERO_NAMESPACE=velero
./velero-api-server
```

**说明**：
- 单集群模式：使用 `--kubeconfig` 参数指定单个配置文件
- 多集群模式：使用 `--kubeconfig-dir` 参数指定包含多个配置文件的目录
- 目录中的文件名（去掉扩展名）将作为集群名称
- 第一个加载的集群将成为默认集群

查看 API 文档: http://localhost:8080/docs

### 2. 启动 Web 控制台

#### 开发环境

```bash
cd web
npm install
npm run dev
```

访问控制台: http://localhost:3000

#### 生产构建

```bash
cd web
npm install
npm run build
```

构建产物位于 `web/dist` 目录。

## 功能特性

### 多集群管理
- ✅ 支持同时管理多个 Kubernetes 集群
- ✅ 通过查询参数 `?cluster=cluster-name` 指定目标集群
- ✅ 通过 HTTP Header `X-Cluster: cluster-name` 指定目标集群
- ✅ Web 界面顶部集群选择器
- ✅ 自动加载目录中的所有 kubeconfig 文件

### 备份管理
- ✅ 创建备份（支持多命名空间、多资源类型）
- ✅ 查看备份列表和状态
- ✅ 查看备份详情
- ✅ 删除备份
- ✅ 支持 TTL 设置
- ✅ 支持卷快照

### 恢复管理
- ✅ 从备份创建恢复任务
- ✅ 恢复到原命名空间或指定命名空间
- ✅ 选择性恢复资源类型
- ✅ 查看恢复状态和详情
- ✅ 支持 PV 恢复

### 定时备份
- ✅ 创建定时备份任务（Cron 表达式）
- ✅ 查看任务列表和执行状态
- ✅ 查看最后备份时间
- ✅ 删除定时任务

### 存储位置
- ✅ 查看备份存储位置列表
- ✅ 查看存储位置详情

## API 端点

| 方法 | 路径 | 说明 | 多集群支持 |
|------|------|------|-----------|
| GET | `/healthz` | 健康检查 | ❌ |
| GET | `/docs` | Swagger 文档 | ❌ |
| GET | `/api/v1/clusters` | 列出可用集群 | ❌ |
| POST | `/api/v1/backups` | 创建备份 | ✅ |
| GET | `/api/v1/backups` | 列出备份 | ✅ |
| GET | `/api/v1/backups/:name` | 获取备份详情 | ✅ |
| DELETE | `/api/v1/backups/:name` | 删除备份 | ✅ |
| POST | `/api/v1/restores` | 创建恢复 | ✅ |
| GET | `/api/v1/restores` | 列出恢复 | ✅ |
| GET | `/api/v1/restores/:name` | 获取恢复详情 | ✅ |
| POST | `/api/v1/schedules` | 创建定时任务 | ✅ |
| GET | `/api/v1/schedules` | 列出定时任务 | ✅ |
| GET | `/api/v1/schedules/:name` | 获取定时任务详情 | ✅ |
| DELETE | `/api/v1/schedules/:name` | 删除定时任务 | ✅ |
| GET | `/api/v1/storage-locations` | 列出存储位置 | ✅ |
| GET | `/api/v1/storage-locations/:name` | 获取存储位置详情 | ✅ |

**多集群支持说明**：
- 标记 ✅ 的端点支持通过 `?cluster=cluster-name` 查询参数或 `X-Cluster: cluster-name` HTTP Header 指定目标集群
- 如果不指定集群，将使用默认集群（第一个加载的集群）

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--port` | 8080 | HTTP 服务端口 |
| `--kubeconfig` | "" | 单个 kubeconfig 文件路径（优先级高于 kubeconfig-dir） |
| `--kubeconfig-dir` | .kube | 包含多个 kubeconfig 文件的目录（多集群模式） |
| `--namespace` | velero | Velero 安装的命名空间 |
| `--insecure-skip-tls` | false | 跳过 TLS 证书验证 |

## 环境变量

| 变量 | 说明 |
|------|------|
| `VELERO_NAMESPACE` | Velero 命名空间，覆盖 `--namespace` 参数 |
| `KUBECONFIG` | 单个 kubeconfig 文件路径，当 `--kubeconfig` 未指定时使用 |
| `KUBECONFIG_DIR` | kubeconfig 目录路径，当 `--kubeconfig-dir` 为默认值时使用 |

## 构建

### 本地构建

```bash
go build -o velero-api-server cmd/server/main.go
```

### 交叉编译（Linux）

```bash
# 标准编译
GOOS=linux GOARCH=amd64 go build -o velero-api-server-linux cmd/server/main.go

# 优化编译（减小体积，去除符号表和调试信息）
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o velero-api-server-linux cmd/server/main.go
```

### 前端构建

```bash
cd web
npm install
npm run build
```

构建产物位于 `web/dist` 目录。

## 部署

### 方案 1：单端口部署（推荐）

后端服务提供静态文件服务，前端和后端使用同一端口。

```bash
# 1. 构建前端
cd web
npm run build
cd ..

# 2. 将前端构建产物复制到后端静态目录
mkdir -p static
cp -r web/dist/* static/

# 3. 启动后端服务（需要在代码中添加静态文件服务）
./velero-api-server --kubeconfig-dir .kube --port 8080
```

访问: http://your-server:8080

### 方案 2：Nginx 反向代理

前端和后端分别部署，通过 Nginx 反向代理统一入口。

#### 1. 部署后端

```bash
# 启动后端服务
./velero-api-server --kubeconfig-dir .kube --port 8080
```

#### 2. 部署前端

```bash
# 构建前端
cd web
npm run build

# 将 dist 目录内容部署到 Nginx 静态文件目录
cp -r dist/* /usr/share/nginx/html/
```

#### 3. 配置 Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 健康检查
    location /healthz {
        proxy_pass http://localhost:8080;
    }

    # Swagger 文档
    location ~ ^/(docs|swagger) {
        proxy_pass http://localhost:8080;
    }
}
```

#### 4. 重启 Nginx

```bash
nginx -t
nginx -s reload
```

### 生产环境建议

1. **使用 systemd 管理后端服务**

创建 `/etc/systemd/system/velero-api-server.service`:

```ini
[Unit]
Description=Velero API Server
After=network.target

[Service]
Type=simple
User=velero
WorkingDirectory=/opt/velero-api-server
Environment="KUBECONFIG_DIR=/opt/velero-api-server/.kube"
Environment="VELERO_NAMESPACE=velero"
ExecStart=/opt/velero-api-server/velero-api-server --port 8080
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
systemctl daemon-reload
systemctl enable velero-api-server
systemctl start velero-api-server
systemctl status velero-api-server
```

2. **配置 HTTPS**（推荐使用 Let's Encrypt）

```bash
# 安装 certbot
apt-get install certbot python3-certbot-nginx

# 获取证书
certbot --nginx -d your-domain.com

# 证书自动续期
certbot renew --dry-run
```

3. **防火墙配置**

```bash
# 开放 HTTP/HTTPS 端口
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

## 项目结构

```
velero-api-server/
├── cmd/server/              # 应用入口
│   └── main.go
├── internal/
│   ├── handler/             # HTTP 处理器
│   │   ├── backup.go        # 备份相关接口
│   │   ├── restore.go       # 恢复相关接口
│   │   ├── schedule.go      # 定时任务相关接口
│   │   ├── bsl.go           # 存储位置相关接口
│   │   └── routes.go        # 路由注册和中间件
│   ├── model/               # 数据模型
│   │   └── types.go         # 请求/响应结构体
│   └── service/             # 业务逻辑
│       └── velero.go        # Velero CRD 操作封装
├── pkg/
│   ├── cluster/             # 多集群管理
│   │   └── manager.go       # 集群配置管理器
│   └── k8s/                 # Kubernetes 客户端
│       └── client.go        # 客户端工具函数
├── docs/                    # 文档
│   ├── swagger.yaml         # OpenAPI 规范
│   ├── swagger.html         # Swagger UI
│   └── README.md            # 文档说明
├── web/                     # 前端控制台
│   ├── src/
│   │   ├── api/             # API 客户端
│   │   │   └── velero.js    # Velero API 封装
│   │   ├── views/           # 页面组件
│   │   │   ├── BackupView.vue
│   │   │   ├── RestoreView.vue
│   │   │   ├── ScheduleView.vue
│   │   │   └── StorageLocationView.vue
│   │   ├── App.vue          # 根组件
│   │   ├── router.js        # 路由配置
│   │   └── main.js          # 应用入口
│   ├── package.json
│   ├── vite.config.js
│   └── index.html
├── go.mod
├── go.sum
└── README.md
```

## 使用示例

### 列出可用集群

```bash
curl http://localhost:8080/api/v1/clusters
```

响应:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "clusters": ["config-fat01", "config-fat02"],
    "default": "config-fat01"
  }
}
```

### 创建备份

#### 使用查询参数指定集群

```bash
curl -X POST "http://localhost:8080/api/v1/backups?cluster=config-fat02" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-default-all",
    "includedNamespaces": ["default"],
    "includedResources": ["deployments.apps", "services"],
    "ttl": "720h"
  }'
```

#### 使用 HTTP Header 指定集群

```bash
curl -X POST "http://localhost:8080/api/v1/backups" \
  -H "Content-Type: application/json" \
  -H "X-Cluster: config-fat02" \
  -d '{
    "name": "backup-default-all",
    "includedNamespaces": ["default"],
    "ttl": "720h"
  }'
```

#### 使用默认集群

```bash
curl -X POST http://localhost:8080/api/v1/backups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-default-all",
    "includedNamespaces": ["default"],
    "ttl": "720h"
  }'
```

### 恢复备份

```bash
curl -X POST "http://localhost:8080/api/v1/restores?cluster=config-fat02" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "restore-default-20260623",
    "backupName": "backup-default-all",
    "includedNamespaces": ["default"]
  }'
```

### 创建定时备份

```bash
curl -X POST "http://localhost:8080/api/v1/schedules?cluster=config-fat02" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "daily-backup",
    "schedule": "0 2 * * *",
    "includedNamespaces": ["default"],
    "ttl": "720h"
  }'
```

### 查看备份列表

```bash
curl "http://localhost:8080/api/v1/backups?cluster=config-fat02"
```

### 删除备份

```bash
curl -X DELETE "http://localhost:8080/api/v1/backups/backup-default-all?cluster=config-fat02"
```

## 开发

### 运行测试

```bash
go test ./...
```

### 代码检查

```bash
go vet ./...
```

### 前端开发

```bash
cd web
npm run dev
```

前端开发服务器会自动代理 API 请求到后端（配置在 `web/vite.config.js`）。

## 常见问题

### 1. 命名空间拼写错误

**问题**：在备份中指定了不存在的命名空间（如 "defalut" 而非 "default"）。

**现象**：Velero 不会报错，备份会正常完成，但不会包含任何资源。

**解决**：
- 仔细检查命名空间名称的拼写
- 使用 `kubectl get namespaces` 确认命名空间存在
- 查看备份详情，检查 `includedNamespaces` 字段

### 2. 多集群配置

**问题**：如何配置多个集群？

**解决**：
```bash
# 创建配置目录
mkdir -p .kube

# 复制多个 kubeconfig 文件，文件名即为集群名称
cp /path/to/cluster1.yaml .kube/config-fat01
cp /path/to/cluster2.yaml .kube/config-fat02

# 启动服务
./velero-api-server --kubeconfig-dir .kube
```

### 3. 跨集群恢复

**问题**：能否从集群 A 的备份恢复到集群 B？

**解决**：可以，前提是两个集群使用相同的备份存储位置（BackupStorageLocation）。步骤：
1. 在集群 A 创建备份
2. 确保集群 B 配置了相同的存储位置
3. 在集群 B 创建恢复，指定集群 A 的备份名称

### 4. 前端无法连接后端

**问题**：前端页面加载正常，但 API 调用失败。

**解决**：
- 检查后端服务是否正常运行：`curl http://localhost:8080/healthz`
- 检查前端 API 配置（`web/src/api/velero.js` 中的 `baseURL`）
- 如果使用 Nginx，检查代理配置是否正确

### 5. TLS 证书验证失败

**问题**：连接 Kubernetes 集群时出现证书验证错误。

**解决**：
```bash
# 启动时跳过 TLS 验证（仅开发环境）
./velero-api-server --insecure-skip-tls
```

**注意**：生产环境不建议跳过 TLS 验证。

## 最佳实践

### 1. 备份策略

- **全量备份**：每周进行一次全量备份
- **增量备份**：每天备份关键命名空间
- **TTL 设置**：根据合规要求设置合理的保留期（如 30 天、90 天）

示例定时任务：
```bash
# 每天凌晨 2 点备份 default 命名空间，保留 30 天
curl -X POST http://localhost:8080/api/v1/schedules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "daily-default-backup",
    "schedule": "0 2 * * *",
    "includedNamespaces": ["default"],
    "ttl": "720h"
  }'
```

### 2. 命名规范

- **备份名称**：`backup-<namespace>-<date>` 或 `backup-<purpose>-<date>`
- **恢复名称**：`restore-<backup-name>-<date>`
- **定时任务名称**：`<frequency>-<namespace>-backup`（如 `daily-default-backup`）

### 3. 监控和告警

- 定期检查备份状态（`phase: Completed`）
- 监控备份存储位置的可用性（`phase: Available`）
- 设置告警：备份失败、存储空间不足

### 4. 安全建议

- 使用 HTTPS 部署生产环境
- 限制 API 访问权限（通过防火墙或认证中间件）
- 定期轮换 kubeconfig 凭证
- 备份数据加密（使用 Velero 的加密功能）

## 技术栈

### 后端
- **Go 1.22+**：编程语言
- **Gin**：Web 框架
- **controller-runtime**：Kubernetes 客户端
- **Velero CRDs**：v1.velero.io API

### 前端
- **Vue 3**：前端框架
- **Element Plus**：UI 组件库
- **Axios**：HTTP 客户端
- **Vue Router**：路由管理
- **Vite**：构建工具

## 相关资源

- [Velero 官方文档](https://velero.io/docs/)
- [Kubernetes API 参考](https://kubernetes.io/docs/reference/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [Vue 3 文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)

## 贡献

欢迎提交 Issue 和 Pull Request！

## License

MIT
# velero-dashboard

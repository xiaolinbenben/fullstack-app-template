# Fullstack App Template

个人使用的全栈开发模板。

## 技术栈

- 公共前端：React + Vite + TypeScript + shadcn/ui
- 管理端：React + Vite + Ant Design + `@ant-design/pro-components`
- 后端：Go + Gin
- 部署：单个 Go 应用镜像，由 1Panel 负责域名、HTTPS 证书和反向代理
- CI/CD：GitHub Actions 构建并推送 GHCR 镜像，通过 SSH 更新 Compose

模板只包含健康检查、静态资源托管和工程规范，不预设 SQLite、PostgreSQL、Redis、MinIO 或其他持久化方案。

## 项目结构

```text
web/                    # 公共前端，访问 /
admin/                  # 管理端，访问 /admin/
server/cmd/server/      # Go 程序入口
server/internal/config/ # 环境变量配置
server/internal/httpapi/ # Gin 路由和静态资源
server/internal/middleware/
server/internal/response/
server/web/dist/        # 两套前端的构建产物
deploy/                 # Compose 和环境变量模板
```

## 本地开发

安装依赖（需要 pnpm 11）：

```bash
corepack enable
make install
```

启动 Go 服务：

```bash
make server-dev
```

在另外两个终端启动前端：

```bash
make web-dev
make admin-dev
```

开发地址：

- 公共前端：`http://localhost:3000/`
- 管理端：`http://localhost:5173/admin/`
- Go 服务：`http://localhost:8000/`
- 健康检查：`http://localhost:8000/healthz`

两个 Vite 开发服务器都会将 `/api` 和 `/healthz` 通过同源代理转发到 `http://127.0.0.1:8000`，因此本地开发不需要给 Gin 添加 CORS 白名单，也不需要额外的前端 `.env` 文件。启动 Go 服务后，再分别运行 `make web-dev` 和 `make admin-dev` 即可。

## API 规范

当前只提供：

```text
GET /healthz
```

响应统一为：

```json
{
  "success": true,
  "message": "",
  "data": {
    "status": "ok"
  }
}
```

后续业务接口统一放在 `/api/*` 或 `/api/admin/*` 下，只允许使用 `GET` 和 `POST`。不得未经说明加入 `PATCH`、`PUT` 或 `DELETE`。业务接口必须通过后端的统一响应封装返回数据。

## 环境变量

唯一的环境变量模板是：

```text
deploy/.env.example
```

复制为本地或生产配置：

```bash
cp deploy/.env.example deploy/.env
```

`deploy/.env` 被 Git 忽略。当前只保留数据库连接串和加密密钥：

```dotenv
DATABASE_URL=replace-with-database-connection-string
ENCRYPTION_KEY=replace-with-a-long-random-key
```

当前后端只读取这两项配置并保留给后续基础设施接入使用，暂不建立数据库连接或实现加密逻辑。端口和镜像地址固定在代码与 Compose 中，不通过环境变量覆盖。

## Docker Compose

本模板的 Compose 只编排 Go 应用，不启动任何数据库或基础服务：

```bash
cp deploy/.env.example deploy/.env
docker compose \
  --env-file deploy/.env \
  -f deploy/docker-compose.yml \
  up -d
```

本地构建镜像：

```bash
docker build -t ghcr.io/xiaolinbenben/fullstack-app-template:latest .
docker compose \
  --env-file deploy/.env \
  -f deploy/docker-compose.yml \
  up -d
```

打包镜像固定监听容器端口 `3000`，Compose 绑定到宿主机回环地址 `127.0.0.1:3000`，适合由 1Panel 反向代理。

## 1Panel 部署

1. 在服务器安装 Docker、Compose 插件和 1Panel。
2. 将 `deploy/docker-compose.yml` 和 `deploy/.env.example` 上传到 `/opt/<app-name>/deploy/`。
3. 在服务器复制 `deploy/.env.example` 为 `deploy/.env`，填写数据库连接串和加密密钥。
4. 执行 `docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d`。
5. 在 1Panel 创建网站、配置 HTTPS 证书，并将反向代理目标设置为 `127.0.0.1:3000`。

HTTPS 证书只由 1Panel 管理，应用镜像不包含站点证书，也不负责 TLS 终止。数据库、缓存和对象存储由具体项目自行选择和接入，本模板不绑定相关服务。

## 构建检查

```bash
make test
make build
docker build -t ghcr.io/xiaolinbenben/fullstack-app-template:latest .
```

## CI/CD

GitHub Actions 会在 Pull Request 和推送时执行 Go 测试、两套前端 typecheck/build 和 Docker 构建。前端依赖统一使用 pnpm 11.9.0 和 `pnpm-lock.yaml`。推送 `main` 时，工作流会将镜像推送到 GHCR，并通过 SSH 更新服务器上的 Compose 文件。

生产部署需要配置 GitHub Environment `production`：

- `SERVER_HOST`
- `SERVER_USER`
- `SERVER_SSH_KEY`

服务器上的 `deploy/.env` 由服务器管理员维护，CI 不覆盖真实配置。

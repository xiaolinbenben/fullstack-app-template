# 项目开发规范

## 架构边界

- `web/` 是公共前端，构建后访问 `/`。
- `admin/` 是管理端，构建后访问 `/admin/`。
- `server/cmd/server/` 是 Go 入口。
- Gin 路由放在 `server/internal/httpapi/`。
- 中间件放在 `server/internal/middleware/`。
- 统一响应放在 `server/internal/response/`。
- 当前模板不预设任何数据库、缓存、对象存储或持久化实现；`DATABASE_URL` 只作为后续接入预留。
- 前端包管理统一使用 pnpm 11.9.0，不使用 npm、yarn 或其他包管理器。

## API 规范

- 业务 API 使用 `/api/*` 或 `/api/admin/*`。
- 只允许 `GET` 和 `POST`。
- 不得未经明确设计加入 `PATCH`、`PUT`、`DELETE`。
- 所有 API 返回 `success`、`message`、`data` 三个字段。
- Handler 不直接拼接 JSON，统一使用 response 包。

## 配置与部署

- 唯一环境变量模板是 `deploy/.env.example`。
- `deploy/.env` 只用于本地或生产环境，禁止提交。
- `.env.example` 只保留 `DATABASE_URL` 和 `ENCRYPTION_KEY`。
- 本地 Go 服务固定为 `8000`，打包镜像固定为 `3000`，不得通过环境变量覆盖。
- Vite 本地开发必须通过 `server.proxy` 将 `/api` 和 `/healthz` 转发到 `http://127.0.0.1:8000`；不要用 CORS 配置替代开发代理。
- HTTPS 证书由 1Panel 管理，不放入仓库或应用镜像。
- Compose 只编排当前应用，不预设基础服务。
- 新增数据库、缓存或对象存储前，必须先确定技术选型，再同步更新 README、Compose 和 CI/CD。
- 1Panel 只负责域名、HTTPS 和反向代理。

## 质量检查

提交前至少运行：

```bash
make test
make build
docker build -t fullstack-app-template:local .
```

前端依赖变更必须提交对应的 `pnpm-lock.yaml`。不要提交 `node_modules`、`package-lock.json`、构建产物、本地 `.env` 或运行数据。

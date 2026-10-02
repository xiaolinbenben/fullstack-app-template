# 项目开发规范

## 架构边界

- `web/` 是公共前端，构建后访问 `/`。
- `admin/` 是管理端，构建后访问 `/admin/`。
- `server/cmd/server/` 是 Go 入口。
- Gin 路由放在 `server/internal/httpapi/`。
- 中间件放在 `server/internal/middleware/`。
- 统一响应放在 `server/internal/response/`。
- GORM 连接和迁移放在 `server/internal/database/`。
- GORM 模型放在 `server/internal/model/`。
- 查询放在 `server/internal/store/`。
- 前端包管理统一使用 pnpm 11.9.0，不使用 npm、yarn 或其他包管理器。

## API 规范

- 业务 API 使用 `/api/*` 或 `/api/admin/*`。
- 只允许 `GET` 和 `POST`。
- 不得未经明确设计加入 `PATCH`、`PUT`、`DELETE`。
- 所有 API 返回 `success`、`message`、`data` 三个字段。
- Handler 不直接拼接 JSON，统一使用 response 包。
- `server/internal/httpapi/woodfish.go` 是 GORM 业务接口的实现。新接口按模型、store、路由这三层添加，不在 Handler 里写查询。
- 参观管理端的账号密码写死为 `admin` / `admin123`。没有修改密码的接口。

## 数据访问

- 持久化只使用 GORM，数据库固定为 PostgreSQL。驱动固定为 `gorm.io/driver/postgres`。不引入 SQLite 或其他 ORM。
- 应用从 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` 拼出 `postgres://<user>:<password>@postgres:5432/<db>?sslmode=disable`。这三个值只允许字母、数字和连字符。主机名、端口和 `sslmode=disable` 写在代码里。连接池为最多 10 个连接。
- 构建继续使用 `CGO_ENABLED=0`。
- 模型使用 GORM 默认表名。主键用自增 `uint`。包含 `CreatedAt` 和 `UpdatedAt`。只有业务需要软删除时才加 `gorm.DeletedAt`。
- 新增模型后，只在 `database.Migrate` 里注册并执行 `AutoMigrate`。不在请求路径里迁移。
- `AutoMigrate` 只建表和加列。删列、改列类型、改约束要单独写迁移。
- 查询使用 `WithContext`。列表在 store 里写明排序和上限，上限不交给客户端随意放大。
- 请求体使用独立结构体，再显式赋给模型。不要把 JSON 直接绑定到 GORM 模型。
- 字符串长度按字符数在 Handler 校验，并与模型 `size` 保持一致。
- 不要把 `0` 值主键传给 `First`。GORM 会忽略零值条件。
- 未找到记录用 `errors.Is(err, gorm.ErrRecordNotFound)` 转成统一响应。数据库错误只写日志，不返回给客户端。
- 日志里不打印数据库连接串和密码。
- 同一次请求里的多条写入放在 `Transaction` 里。更新用 `Updates` 并指定字段，不用 `Save` 回写整行。
- 需要把用户输入放进 SQL 时使用占位参数，不拼接 SQL。

## 配置与部署

- 开发时前后端分开运行。部署时前端构建产物由 `go:embed` 打进同一个 Go 二进制，再封装成一个应用镜像。
- 唯一环境变量模板是 `deploy/.env.example`。
- `deploy/.env` 只用于本地或生产环境，禁止提交。
- `.env.example` 保留 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` 和 `ENCRYPTION_KEY`。用户名和数据库名固定为 `app`。密码只使用字母、数字和连字符。应用用这三个值连接主机名 `postgres`。
- PostgreSQL 和应用容器都通过 `env_file` 读取 `deploy/.env`。不要在 Compose 里用 `environment` 插值这些变量。
- 本地 `make server-dev` 和 `make test` 读取 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB`。这三个变量已经在环境里时直接使用，否则读 `deploy/.env`。传给 Go 进程时主机名用 `127.0.0.1`。`make server-dev` 另外只传入 `ENCRYPTION_KEY`。
- 本地 Go 服务固定为 `8000`，打包镜像固定为 `3000`，不得通过环境变量覆盖。
- Vite 本地开发必须通过 `server.proxy` 将 `/api` 和 `/healthz` 转发到 `http://127.0.0.1:8000`；不要用 CORS 配置替代开发代理。
- HTTPS 证书由 1Panel 管理，不放入仓库或应用镜像。
- 只使用 `deploy/docker-compose.yml`。它同时编排应用和 PostgreSQL。PostgreSQL 绑定 `127.0.0.1:5432`。应用容器用 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` 连接主机名 `postgres`。
- 应用镜像里不包含数据库。缓存和对象存储不预设。
- 1Panel 只负责域名、HTTPS 和反向代理。

## 质量检查

提交前至少运行：

```bash
make test
make build
docker build -t fullstack-app-template:local .
```

前端依赖变更必须提交对应的 `pnpm-lock.yaml`。不要提交 `node_modules`、`package-lock.json`、构建产物、本地 `.env` 或运行数据。

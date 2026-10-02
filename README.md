# Fullstack App Template

个人使用的全栈开发模板。

在线演示：<https://fullstack-app-template.linzhiqing.dev/>

管理端：<https://fullstack-app-template.linzhiqing.dev/admin/>（账号 `admin`，密码 `admin123`）

## 技术栈

- 公共前端：React + Vite + TypeScript + shadcn/ui
- 管理端：React + Vite + Ant Design + `@ant-design/pro-components`
- 后端：Go + Gin + GORM
- 部署：单个 Go 应用镜像，由 1Panel 负责域名、HTTPS 证书和反向代理
- CI/CD：GitHub Actions 构建并推送 GHCR 镜像，再以默认用户 root 登录服务器，更新 `/opt/<仓库名>/deploy`

开发时前后端分开运行，部署时用 `go:embed` 合成一个二进制和一个应用镜像。数据库使用 PostgreSQL，和 Go 应用写在同一个 `deploy/docker-compose.yml` 里。缓存和对象存储不预设。

## 项目结构

```text
web/                    # 公共前端，访问 /
admin/                  # 管理端，访问 /admin/
server/cmd/server/      # Go 程序入口
server/internal/config/ # 环境变量配置
server/internal/database/ # GORM 连接和迁移
server/internal/model/  # GORM 模型
server/internal/store/  # 查询
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

先启动 PostgreSQL，再启动 Go 服务：

```bash
cp deploy/.env.example deploy/.env
docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d postgres
make server-dev
```

`make server-dev` 把 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` 和 `ENCRYPTION_KEY` 传给 Go 进程。连接主机名在本机改成 `127.0.0.1`。

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

当前提供：

```text
GET  /healthz
GET  /api/woodfish
POST /api/woodfish
POST /api/admin/login
GET  /api/admin/session
POST /api/admin/logout
```

`GET /api/woodfish` 返回木鱼被敲过的总次数，`POST /api/woodfish` 把次数加一。次数存在 PostgreSQL 里。参观管理端的账号和密码写死为 `admin` / `admin123`，首页上能看到，不能在页面里修改。

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

`deploy/.env` 被 Git 忽略。PostgreSQL 容器和应用容器都读取这个文件：

```dotenv
POSTGRES_USER=app
POSTGRES_PASSWORD=replace-with-a-long-random-password
POSTGRES_DB=app
ENCRYPTION_KEY=replace-with-a-long-random-key
# APP_HOST_PORT=3000
# POSTGRES_HOST_PORT=5432
```

用户名和数据库名固定为 `app`。这三个值只使用字母、数字和连字符。服务启动时用它们拼出 PostgreSQL 连接串，主机名是 `postgres`，容器端口是 `5432`，`sslmode=disable`，再用 GORM 连接并执行 `AutoMigrate`。数据库时区由 Compose 固定为 `Asia/Shanghai`。本地 Go 监听 `8000`，镜像监听容器端口 `3000`，这两个监听端口不通过环境变量改变。`APP_HOST_PORT` 和 `POSTGRES_HOST_PORT` 只决定映射到宿主机哪个回环端口。`ENCRYPTION_KEY` 仍只读入配置，加密逻辑由具体业务实现。

## Docker Compose

`deploy/docker-compose.yml` 同时启动 PostgreSQL 和应用。两个服务都用 `env_file` 读取 `deploy/.env`。端口只绑定宿主机回环地址：应用默认 `127.0.0.1:3000`，PostgreSQL 默认 `127.0.0.1:5432`。应用容器用 `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` 连接主机名 `postgres` 的容器端口 `5432`，不走宿主机映射。

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

打包镜像固定监听容器端口 `3000`。Compose 默认把这个端口映射到宿主机 `127.0.0.1:3000`，适合由 1Panel 反向代理。服务器上 `3000` 或 `5432` 已被占用时，在 `.env` 里去掉对应行的注释并改成空闲端口。应用容器仍然连接 `postgres:5432`。

## 1Panel 部署

服务器目录为 `/opt/<仓库名>/deploy/`。本仓库对应：

```text
/opt/fullstack-app-template/deploy/
```

运维在这个目录里工作，并且只提前放置 `.env`。内容按仓库里的 `deploy/.env.example` 填写数据库账号和 `ENCRYPTION_KEY`。`docker-compose.yml` 由 CI 在每次部署时写入该目录并覆盖同名文件。`.env` 保持运维填写的内容。

1. 在服务器安装 Docker 和 1Panel。现在安装 Docker 会自带 Compose。
2. 创建 `/opt/fullstack-app-template/deploy/`，把填好的 `.env` 放进去。
3. 在 GitHub Environment `production` 配好下方三个 Secret 后，推送 `main`。工作流会写入 Compose、拉取镜像并在该目录执行 `docker compose up -d`。
4. 在 1Panel 创建网站、配置 HTTPS 证书，并将反向代理目标设置为 `127.0.0.1:3000`。`.env` 里改过 `APP_HOST_PORT` 时，用改后的端口。

HTTPS 证书只由 1Panel 管理，应用镜像不包含站点证书，也不负责 TLS 终止。PostgreSQL 由同一个 Compose 文件启动。缓存和对象存储不预设。之后查看容器或修改 `.env`，都在 `/opt/fullstack-app-template/deploy/` 里进行。

## 构建检查

```bash
make test
make build
docker build -t ghcr.io/xiaolinbenben/fullstack-app-template:latest .
```

## CI/CD

GitHub Actions 会在 Pull Request 和推送时构建两套前端、Go 服务和 Docker 镜像，不运行测试。前端依赖统一使用 pnpm 11.9.0 和 `pnpm-lock.yaml`。推送 `main` 时，工作流把镜像推送到 GHCR，再以默认用户 `root` 登录服务器，把仓库里的 `deploy/docker-compose.yml` 写到 `/opt/<仓库名>/deploy/docker-compose.yml`，用下面的 PAT 登录 `ghcr.io` 后执行 `docker compose pull` 和 `docker compose up -d`。

登录用户默认为 `root`，不放入部署 Secret。要换成其他用户时，在 GitHub Environment `production` 设置 Variable `SERVER_USER`。

生产部署使用 GitHub Environment `production`。部署 Secret 固定为这三个：

| Secret | 内容 |
| --- | --- |
| `SERVER_HOST` | 服务器 IP |
| `SERVER_PASSWORD` | root 密码 |
| `GHCR_PAT` | 可拉取 GHCR 镜像的 GitHub Personal Access Token |

`GHCR_PAT` 需要 `read:packages`。镜像为私有时，经典 PAT 还需要 `repo`。工作流用仓库所有者的用户名登录 `ghcr.io`，这个 PAT 属于仓库所有者。服务器上的 `.env` 由运维在 `/opt/<仓库名>/deploy/` 维护，CI 只写入 `docker-compose.yml`。

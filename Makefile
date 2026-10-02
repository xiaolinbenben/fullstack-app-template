.PHONY: install web-install admin-install web-dev admin-dev server-dev web-build admin-build server-build test build docker-build docker-up

# 只给本机的 server-dev 和 test 使用。三个变量已在环境里时直接使用，否则读 deploy/.env。
# 本次进程的主机名用 127.0.0.1，容器里仍是 postgres。
define postgres_env
if [ -n "$$POSTGRES_USER" ] && [ -n "$$POSTGRES_PASSWORD" ] && [ -n "$$POSTGRES_DB" ]; then user="$$POSTGRES_USER"; pass="$$POSTGRES_PASSWORD"; name="$$POSTGRES_DB"; else [ -f deploy/.env ] || { echo "缺少 deploy/.env，请先复制 deploy/.env.example"; exit 1; }; user=$$(sed -n 's/^POSTGRES_USER=//p' deploy/.env | head -n 1); pass=$$(sed -n 's/^POSTGRES_PASSWORD=//p' deploy/.env | head -n 1); name=$$(sed -n 's/^POSTGRES_DB=//p' deploy/.env | head -n 1); [ -n "$$user" ] && [ -n "$$pass" ] && [ -n "$$name" ] || { echo "deploy/.env 缺少 POSTGRES_USER、POSTGRES_PASSWORD 或 POSTGRES_DB"; exit 1; }; fi
endef

install: web-install admin-install

web-install:
	cd web && pnpm install --frozen-lockfile

admin-install:
	cd admin && pnpm install --frozen-lockfile

web-dev:
	cd web && pnpm run dev

admin-dev:
	cd admin && pnpm run dev

server-dev:
	@$(postgres_env); if [ -n "$$ENCRYPTION_KEY" ]; then key="$$ENCRYPTION_KEY"; elif [ -f deploy/.env ]; then key=$$(sed -n 's/^ENCRYPTION_KEY=//p' deploy/.env | head -n 1); else key=""; fi; cd server && POSTGRES_USER="$$user" POSTGRES_PASSWORD="$$pass" POSTGRES_DB="$$name" POSTGRES_HOST=127.0.0.1 ENCRYPTION_KEY="$$key" go run ./cmd/server

web-build:
	cd web && pnpm run typecheck && pnpm run build

admin-build:
	cd admin && pnpm run typecheck && pnpm run build

server-build:
	cd server && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

test:
	@$(postgres_env); cd server && POSTGRES_USER="$$user" POSTGRES_PASSWORD="$$pass" POSTGRES_DB="$$name" POSTGRES_HOST=127.0.0.1 go test ./...

build: web-build admin-build server-build

docker-build:
	docker build -t ghcr.io/xiaolinbenben/fullstack-app-template:latest .

docker-up:
	cp -n deploy/.env.example deploy/.env 2>/dev/null || true
	docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d

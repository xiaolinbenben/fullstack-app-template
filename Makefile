.PHONY: install web-install admin-install web-dev admin-dev server-dev web-build admin-build server-build test build docker-build docker-up

# 只给本机的 server-dev 和 test 使用。容器内主机名 postgres 在本机改成 127.0.0.1。
define database_url
if [ -n "$$DATABASE_URL" ]; then url="$$DATABASE_URL"; else [ -f deploy/.env ] || { echo "缺少 deploy/.env，请先复制 deploy/.env.example"; exit 1; }; url=$$(sed -n 's/^DATABASE_URL=//p' deploy/.env | head -n 1); [ -n "$$url" ] || { echo "deploy/.env 缺少 DATABASE_URL"; exit 1; }; fi; url=$$(printf '%s' "$$url" | sed 's/@postgres:/@127.0.0.1:/')
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
	@$(database_url); if [ -n "$$ENCRYPTION_KEY" ]; then key="$$ENCRYPTION_KEY"; elif [ -f deploy/.env ]; then key=$$(sed -n 's/^ENCRYPTION_KEY=//p' deploy/.env | head -n 1); else key=""; fi; cd server && DATABASE_URL="$$url" ENCRYPTION_KEY="$$key" go run ./cmd/server

web-build:
	cd web && pnpm run typecheck && pnpm run build

admin-build:
	cd admin && pnpm run typecheck && pnpm run build

server-build:
	cd server && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

test:
	@$(database_url); cd server && DATABASE_URL="$$url" go test ./...

build: web-build admin-build server-build

docker-build:
	docker build -t ghcr.io/xiaolinbenben/fullstack-app-template:latest .

docker-up:
	cp -n deploy/.env.example deploy/.env 2>/dev/null || true
	docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d

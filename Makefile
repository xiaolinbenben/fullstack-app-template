.PHONY: install web-install admin-install web-dev admin-dev server-dev web-build admin-build server-build test build docker-build docker-up

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
	cd server && go run ./cmd/server

web-build:
	cd web && pnpm run typecheck && pnpm run build

admin-build:
	cd admin && pnpm run typecheck && pnpm run build

server-build:
	cd server && go test ./... && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

test:
	cd server && go test ./...

build: web-build admin-build server-build

docker-build:
	docker build -t ghcr.io/xiaolinbenben/fullstack-app-template:latest .

docker-up:
	cp -n deploy/.env.example deploy/.env 2>/dev/null || true
	docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d

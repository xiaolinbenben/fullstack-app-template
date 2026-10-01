FROM node:22-alpine AS web-builder

WORKDIR /src/web
RUN corepack enable && corepack prepare pnpm@11.9.0 --activate
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm run typecheck && pnpm run build

FROM node:22-alpine AS admin-builder

WORKDIR /src/admin
RUN corepack enable && corepack prepare pnpm@11.9.0 --activate
COPY admin/package.json admin/pnpm-lock.yaml admin/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY admin/ ./
RUN pnpm run typecheck && pnpm run build

FROM golang:1.26-alpine AS server-builder

WORKDIR /src/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
COPY --from=web-builder /src/server/web/dist/public ./web/dist/public
COPY --from=admin-builder /src/server/web/dist/admin ./web/dist/admin
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.listenAddr=:3000" -o /out/server ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache tzdata wget \
    && addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app

COPY --from=server-builder /out/server /usr/local/bin/server

USER app
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/server"]

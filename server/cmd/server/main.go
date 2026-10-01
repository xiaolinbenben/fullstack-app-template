package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fullstack-app-template/server/internal/config"
	"fullstack-app-template/server/internal/database"
	"fullstack-app-template/server/internal/httpapi"
	webassets "fullstack-app-template/server/web"
)

// listenAddr is fixed at build time. Local builds use 8000; the container
// build overrides it to 3000 with a linker flag.
var listenAddr = ":8000"

func main() {
	cfg := config.Load(listenAddr)

	webFS, err := fs.Sub(webassets.Assets, "dist")
	if err != nil {
		log.Fatalf("读取前端资源失败: %v", err)
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
	}()
	if err := database.Migrate(db); err != nil {
		log.Fatalf("迁移数据库失败: %v", err)
	}
	log.Printf("数据库已连接")

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.New(cfg, db, webFS),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server 已启动，监听 %s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务异常退出: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("关闭服务失败: %v", err)
	}
}

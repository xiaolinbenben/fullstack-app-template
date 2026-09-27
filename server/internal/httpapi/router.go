package httpapi

import (
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	"fullstack-app-template/server/internal/config"
	"fullstack-app-template/server/internal/middleware"
	"fullstack-app-template/server/internal/response"
	"github.com/gin-gonic/gin"
)

func New(_ config.Config, webFS fs.FS) *gin.Engine {
	engine := gin.New()
	engine.Use(
		gin.CustomRecovery(func(c *gin.Context, _ any) {
			response.Error(c, http.StatusInternalServerError, "服务器内部错误")
			c.Abort()
		}),
		middleware.RequestLog(),
	)
	engine.HandleMethodNotAllowed = true

	engine.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})
	engine.NoMethod(func(c *gin.Context) {
		response.Error(c, http.StatusMethodNotAllowed, "请求方法不支持")
	})
	engine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			response.Error(c, http.StatusNotFound, "接口不存在")
			return
		}
		serveFrontend(c, webFS)
	})

	return engine
}

func serveFrontend(c *gin.Context, webFS fs.FS) {
	fsPrefix, urlPrefix := frontendTarget(c.Request.URL.Path)
	relative := strings.TrimPrefix(c.Request.URL.Path, urlPrefix)
	name := strings.TrimPrefix(relative, "/")
	if name == "" || strings.HasPrefix(name, ".") {
		name = "index.html"
	}

	assetName := path.Join(fsPrefix, name)
	info, err := fs.Stat(webFS, assetName)
	if err != nil || info.IsDir() {
		name = "index.html"
		assetName = path.Join(fsPrefix, name)
	}

	content, err := fs.ReadFile(webFS, assetName)
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, "前端构建产物不存在，请先执行前端构建")
		return
	}

	if name == "index.html" {
		c.Header("Cache-Control", "no-cache")
	} else {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	}
	if contentType := mime.TypeByExtension(path.Ext(name)); contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Data(http.StatusOK, contentType(name), content)
}

func frontendTarget(requestPath string) (string, string) {
	if requestPath == "/admin" || strings.HasPrefix(requestPath, "/admin/") {
		return "admin", "/admin"
	}
	return "public", ""
}

func contentType(name string) string {
	if value := mime.TypeByExtension(path.Ext(name)); value != "" {
		return value
	}
	return "application/octet-stream"
}

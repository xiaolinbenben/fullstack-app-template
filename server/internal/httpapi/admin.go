package httpapi

import (
	"net/http"
	"strings"

	"fullstack-app-template/server/internal/response"
	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(engine *gin.Engine, encryptionKey string) {
	engine.POST("/api/admin/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "请求格式错误")
			return
		}
		if strings.TrimSpace(req.Username) != adminUsername || req.Password != adminPassword {
			response.Error(c, http.StatusUnauthorized, "账号或密码错误")
			return
		}
		writeSession(c, encryptionKey)
		response.OK(c, gin.H{"username": adminUsername})
	})
	engine.GET("/api/admin/session", func(c *gin.Context) {
		name, ok := currentAdmin(c, encryptionKey)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "未登录")
			return
		}
		response.OK(c, gin.H{"username": name})
	})
	engine.POST("/api/admin/logout", func(c *gin.Context) {
		clearSession(c)
		response.OK(c, gin.H{})
	})
}

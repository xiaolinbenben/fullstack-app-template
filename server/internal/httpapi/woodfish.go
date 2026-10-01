package httpapi

import (
	"errors"
	"log"
	"net/http"

	"fullstack-app-template/server/internal/model"
	"fullstack-app-template/server/internal/response"
	"fullstack-app-template/server/internal/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerWoodfishRoutes(engine *gin.Engine, db *gorm.DB) {
	counters := store.NewCounterStore(db)
	engine.GET("/api/woodfish", func(c *gin.Context) {
		count, err := counters.Count(c.Request.Context(), model.WoodfishName)
		if err != nil {
			log.Printf("查询木鱼次数失败: %v", err)
			response.Error(c, http.StatusInternalServerError, "服务器内部错误")
			return
		}
		response.OK(c, gin.H{"count": count})
	})
	engine.POST("/api/woodfish", func(c *gin.Context) {
		count, err := counters.Increment(c.Request.Context(), model.WoodfishName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("记录木鱼次数失败: 计数行不存在")
			response.Error(c, http.StatusInternalServerError, "服务器内部错误")
			return
		}
		if err != nil {
			log.Printf("记录木鱼次数失败: %v", err)
			response.Error(c, http.StatusInternalServerError, "服务器内部错误")
			return
		}
		response.OK(c, gin.H{"count": count})
	})
}

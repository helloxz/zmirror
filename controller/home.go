package controller

import (
	"net/http"
	"zmirror/internal/handler"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	// 获取版本日期
	version_date := handler.DATE_VERSION
	// 渲染主页，使用前端应用
	c.HTML(http.StatusOK, "index.html", gin.H{
		"date": version_date,
	})
}

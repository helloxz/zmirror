package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	// 渲染主页，使用前端应用
	c.HTML(http.StatusOK, "index.html", nil)
}

package routers

import (
	"EG/handler"

	"github.com/gin-gonic/gin"
)

//Router .
func Router(r *gin.Engine) {
	r.POST("/api/v1", handler.Login)
	r.POST("/api/v1/crawler", handler.Crawler)
}

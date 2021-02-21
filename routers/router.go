package routers

import (
	"EG/handler"

	"github.com/gin-gonic/gin"
)

//Router .
func Router(r *gin.Engine) {

	//user:
	g1 := r.Group("/api/v1/user")
	{
		//登陆
		g1.POST("/", handler.Login)

		//修改用户信息
		g1.PUT("/", handler.ChangeUserInfo)
	}
	//作业
	r.POST("/api/v1/homework", handler.Crawler)
}

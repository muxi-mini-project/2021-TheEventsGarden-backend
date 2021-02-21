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

		//获取用户信息
		g1.GET("/", handler.GetUserInfo)
	}
	//记事本
	g2 := r.Group("/api/v1/notepad")
	{
		//新建待办
		g2.POST("/create", handler.CreateBackpad)

		//取消待办
		g2.PUT("/", handler.ChangeBackpad)

		//查询待办
		g2.GET("/", handler.Getbackpads)

		//消除未完成待办
		g2.PUT("/clear", handler.ClearBackpad)

		//完成待办
		g2.POST("/", handler.CompleteBackpad)
	}
	//花园
	g3 := r.Group("/api/v1/garden")
	{
		//获取用户花园皮肤
		g3.GET("/", handler.GetSkins)
		//g3.GET("/skin", handler.GetSkins)

		//新增皮肤
		g3.POST("/", handler.BuySkin)
		//g3.POST("/buy", handler.BuySkin)

		//买花
		g3.PUT("/", handler.BuyFlower)
		//g3.PUT("/buy", handler.BuyFlower)
	}
	//作业
	r.POST("/api/v1/homework", handler.Crawler)
}

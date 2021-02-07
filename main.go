package main

import (
	"EG/routers"

	"github.com/gin-gonic/gin"
)

// @title TheEventsGarden API
// @version 1.0.0
// @description 记事园API
// @termsOfService http://swagger.io/terrms/
// @contact.name TAODEI
// @contact.email 864978550@qq.com
// @host localhost
// @BasePath: /api/v1
func main() {
	r := gin.Default()
	routers.Router(r)
	r.Run(":3333")
}

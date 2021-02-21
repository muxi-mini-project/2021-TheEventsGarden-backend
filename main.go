package main

import (
	"EG/model"
	"EG/routers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var err error

// @title TheEventsGarden API
// @version 1.0.0
// @description 记事园API
// @termsOfService http://swagger.io/terrms/
// @contact.name TAODEI
// @contact.email 864978550@qq.com
// @host 124.71.184.107
// @BasePath: /api/v1
// @Schemes http
func main() {
	model.DB, err = gorm.Open("mysql", "tao:12345678@/EG?parseTime=True")
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	routers.Router(r)
	r.Run(":1333")
	defer model.DB.Close()
}

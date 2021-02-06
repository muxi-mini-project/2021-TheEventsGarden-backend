package main

import (
	"EG/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.Router(r)
	r.Run(":3333")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//result, err := model.MakeAccountPreflightRequest()
	//fmt.Println(result, err)
	//fmt.Println()
	//fmt.Println(ho)
}

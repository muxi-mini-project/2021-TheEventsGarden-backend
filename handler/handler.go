package handler

import (
	"EG/model"

	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var p model.User
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"message": "输入有误，格式错误"})
		return
	}
	_, err := model.GetUserInfoFormOne(p.StudentID, p.Password)
	if err != nil {
		c.JSON(400, "用户名或密码错误")
		return
	}
	claims := &model.Jwt{StudentID: p.StudentID}

	claims.ExpiresAt = time.Now().Add(200 * time.Hour).Unix()
	claims.IssuedAt = time.Now().Unix()

	var Secret = "vinegar" //加醋

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		log.Println(err)
	}
	c.JSON(200, gin.H{"token": signedToken})
}

func Crawler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(404, gin.H{"message": "找不到该用户信息，请先登录"})
		return
	}

	var p model.User
	c.BindJSON(&p)
	results := model.GetHomework(id, p.Password)
	log.Println(results)
	c.JSON(200, results)
}

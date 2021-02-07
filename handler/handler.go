package handler

import (
	"EG/model"

	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Token struct {
	Token string `json:"token"`
}

// @Summary  登录
// @Description 学号密码登录
// @Accept application/json
// @Produce application/json
// @Param object body model.User true "登录的用户信息"
// @Success 200 {object} Token "将student_id作为token保留"
// @Failure 401 "用户名或密码错误"
// @Failure 400 "输入有误，格式错误"
// @Router / [post]
func Login(c *gin.Context) {
	var p model.User
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"message": "输入有误，格式错误"})
		return
	}
	_, err := model.GetUserInfoFormOne(p.StudentID, p.Password)
	if err != nil {
		c.JSON(401, "用户名或密码错误")
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
	var Token Token
	Token.Token = signedToken
	c.JSON(200, Token)
	//c.JSON(200, gin.H{"token": signedToken})
}

// @Summary  爬取作业
// @Description 爬取用户云课堂作业
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param password body model.User true "password"
// @Success 200 {object} []model.homework "获取成功"
// @Failure 400 "输入有误，格式错误"
// @Failure 401 "找不到该用户信息，请先登录"
// @Router /crawler [post]
func Crawler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "找不到该用户信息，请先登录"})
		return
	}

	var p model.User
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"message": "输入有误，格式错误"})
		return
	}
	results := model.GetHomework(id, p.Password)
	//log.Println(results)
	c.JSON(200, results)
}

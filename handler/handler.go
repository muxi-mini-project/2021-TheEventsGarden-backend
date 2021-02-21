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
// @Tags user
// @Description 学号密码登录
// @Accept application/json
// @Produce application/json
// @Param object body model.User true "登录的用户信息"
// @Success 200 {object} Token "将student_id作为token保留"
// @Failure 401 "用户名或密码错误"
// @Failure 400 "Lack Necessary_Param."
// @Failure 500 "Fail."
// @Router /user [post]
func Login(c *gin.Context) {
	var p model.User
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	_, err := model.GetUserInfoFormOne(p.StudentID, p.Password)
	if err != nil {
		c.JSON(401, "用户名或密码错误")
		return
	}
	if resu := model.DB.Where("student_id = ?", p.StudentID).First(&p); resu.Error != nil {
		p.Summary = "这个家伙很懒，还没有写简介哦"
		p.Name = "小樨"
		//性别未设置 默认为0
		p.Gold = 0
		p.Sex = 0
		p.UserPicture = "www.baidu.com"
		model.DB.Create(&p)
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
}

// @Summary  爬取作业
// @Description 爬取用户云课堂作业
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param password body model.User true "password"
// @Success 200 {object} []model.homework "获取成功"
// @Failure 400 "Lack Necessary_Param."
// @Failure 401 "Token Invalid."
// @Failure 500 "Fail."
// @Router /homework [post]
func Crawler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var p model.User
	if err := c.BindJSON(&p); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	results := model.GetHomework(id, p.Password)
	c.JSON(200, results)
}

// @Summary  修改用户信息
// @Tags user
// @Description 接收新的User结构体来修改用户信息
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param User body model.User true "需要修改的用户信息"
// @Success 200 "修改成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "Lack Necessary_Param." or "Sex参数错误(0 = 未设置， 1 = 男， 2 = 女)"
// @Failure 500 "Fail."
// @Router /user [put]
func ChangeUserInfo(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	if user.Sex != 0 || user.Sex != 1 || user.Sex != 2 {
		c.JSON(400, gin.H{"message": "Sex参数错误(0 = 未设置， 1 = 男， 2 = 女)"})
		return
	}
	user.StudentID = id
	if err := model.UpdateUserInfo(user); err != nil {
		c.JSON(400, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(200, gin.H{"message": "修改成功"})
}

func CreateBackpad(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var b model.Backpad
	if err := c.BindJSON(&b); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	if err := model.CreateBackpad(id, b); err != nil {
		c.JSON(400, gin.H{"message": "新增待办失败"})
		return
	}
	c.JSON(200, gin.H{"message": "新增待办成功"})
}

func ChangeBackpad(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var b model.Backpad
	if err := c.BindJSON(&b); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	if err := model.ChangeBackpad(id, b); err != nil {
		c.JSON(400, gin.H{"message": "修改待办失败"})
		return
	}
	c.JSON(200, gin.H{"message": "修改待办成功"})
}

func Getbackpads(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	backpads, err := model.GetBackpads(id)
	if err != nil {
		c.JSON(400, gin.H{"message": "查询失败"})
		return
	}
	c.JSON(200, backpads)
}

func GetUserInfo(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	user, err := model.GetUserInfo(id)
	if err != nil {
		c.JSON(400, gin.H{"message": "查询失败"})
		return
	}
	c.JSON(200, user)
}

func Clear(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var b model.Backpad
	if err := c.BindJSON(&b); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	if err, str := model.ClearBackpad(id, b); str != "" {
		c.JSON(400, gin.H{"message": str})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"message": "清除待办失败"})
		return
	}
	c.JSON(200, gin.H{"message": "清除待办成功"})
}

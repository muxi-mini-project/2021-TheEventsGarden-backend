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
		p.Flower = 0
		p.UserPicture = "www.baidu.com"
		model.DB.Create(&p)
	}

	claims := &model.Jwt{StudentID: p.StudentID}

	claims.ExpiresAt = time.Now().Add(200 * time.Hour).Unix()
	claims.IssuedAt = time.Now().Unix()

	var Secret = "sugar" //加糖

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
// @Success 200 {object} []model.Homework "获取成功"
// @Failure 400 "Fail."
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

	u, err := model.GetUserInfo(id)
	if err != nil {
		c.JSON(400, gin.H{"message": "Fail."})
		return
	}
	results := model.GetHomework(id, u.Password)
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
	if user.Sex != 0 && user.Sex != 1 && user.Sex != 2 {
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

// @Summary 获取用户信息
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Success 200 {object} model.User
// @Failure 401 "Token Invalid."
// @Failure 400 "查询失败"
// @Failure 500 "Fail."
// @Router /user [get]
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

// @Summary 新建待办
// @Tags notepad
// @Description 接收新的Backpad结构体来新建待办
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param Backpad body model.Backpad true "name 必需 hour和minute可选"
// @Success 200 "新增待办成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "新增待办失败" or "Lack Necessary_Param."
// @Failure 203 "失败，该用户今日已使用该待办名"
// @Failure 500 "Fail."
// @Router /notepad/create [post]
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
	if err, str := model.CreateBackpad(id, b); str != "" {
		c.JSON(203, gin.H{"message": str})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"message": "新增待办失败"})
		return
	}
	c.JSON(200, gin.H{"message": "新增待办成功"})
}

// @Summary 取消待办
// @Tags notepad
// @Description
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param Backpad body model.Backpad true "只需要该待办名 name"
// @Success 200 "修改待办成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "新增待办失败" or "Lack Necessary_Param."
// @Failure 500 "Fail."
// @Router /notepad [put]
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

// @Summary 查询待办
// @Tags notepad
// @Description 获取该用户所有待办
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Success 200 {object} []model.Backpad
// @Failure 401 "Token Invalid."
// @Failure 400 "查询失败" or "Lack Necessary_Param."
// @Failure 500 "Fail."
// @Router /notepad [get]
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

// @Summary 消除未完成待办
// @Tags notepad
// @Description
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param Backpad body model.Backpad true "只需要该待办名 name"
// @Success 200 "清除待办成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "新增待办失败" or "Lack Necessary_Param."
// @Failure 203 "该待办已完成或已取消" or "金币不足"
// @Failure 500 "Fail."
// @Router /notepad/clear [put]
func ClearBackpad(c *gin.Context) {
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
		c.JSON(203, gin.H{"message": str})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"message": "清除待办失败"})
		return
	}
	c.JSON(200, gin.H{"message": "清除待办成功"})
}

type T struct {
	FinishTime int `json:"finish_time"`
}

// @Summary 完成待办
// @Tags notepad
// @Description 需要该待办名 name 和
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param Backpad body model.Backpad true "只需要该待办名 name"
// @Param time body T true "需要完成时间 finish_time"
// @Success 200 "完成待办成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "完成待办失败" or "Lack Necessary_Param."
// @Failure 203 "失败，该用户今日已使用该待办名"
// @Failure 500 "Fail."
// @Router /notepad [post]
func CompleteBackpad(c *gin.Context) {
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
	var t T
	c.BindJSON(&t)
	if err := model.CompleteBackpad(id, b, t.FinishTime); err != nil {
		c.JSON(400, gin.H{"message": "完成待办失败"})
		return
	}
	c.JSON(200, gin.H{"message": "完成待办成功"})

}

// @Summary 获取用户花园皮肤
// @Tags garden
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Success 200 {object} []model.Skin
// @Failure 401 "Token Invalid."
// @Failure 400 "获取皮肤失败"
// @Failure 500 "Fail."
// @Router /garden [get]
func GetSkins(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	skins, err := model.GetSkins(id)
	if err != nil {
		c.JSON(400, gin.H{"message": "获取皮肤失败"})
		return
	}
	c.JSON(200, skins)
}

// @Summary 新增皮肤
// @Tags garden
// @Description
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param Skin body model.Skin true "只需要 skin_id"
// @Success 200 "购买成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "购买皮肤失败" or "Lack Necessary_Param."
// @Failure 203 "未找到该皮肤" or "已拥有" or "未购买x号皮肤" or "金币不足"
// @Failure 500 "Fail."
// @Router /garden [post]
func BuySkin(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var skin model.Skin
	if err := c.BindJSON(&skin); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	if err, str := model.BuySkin(id, skin); str != "" {
		c.JSON(203, gin.H{"message": str})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"message": "购买皮肤失败"})
		return
	}
	c.JSON(200, gin.H{"message": "购买成功"})
}

type N struct {
	Number int `json:"number"`
}

// @Summary 买花
// @Tags garden
// @Description
// @Accept application/json
// @Produce application/json
// @Param token header string true "token"
// @Param Skin body N true "需要 number"
// @Success 200 "摘花成功"
// @Failure 401 "Token Invalid."
// @Failure 400 "摘花失败" or "Lack Necessary_Param."
// @Failure 203 "金币不足"
// @Failure 500 "Fail."
// @Router /garden [put]
func BuyFlower(c *gin.Context) {
	token := c.Request.Header.Get("token")
	id, err := model.VerifyToken(token)
	if err != nil {
		c.JSON(401, gin.H{"message": "Token Invalid."})
		return
	}

	var n N
	if err := c.BindJSON(&n); err != nil {
		c.JSON(400, gin.H{"message": "Lack Necessary_Param."})
		return
	}
	if err, str := model.BuyFlower(id, n.Number); str != "" {
		c.JSON(203, gin.H{"message": str})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"message": "摘花失败"})
		return
	}
	c.JSON(200, gin.H{"message": "摘花成功"})
}

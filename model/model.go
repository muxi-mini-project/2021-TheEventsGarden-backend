package model

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

const (
	ErrorReasonServerBusy = "服务器繁忙"
	ErrorReasonReLogin    = "请重新登陆"
)

//DB 全局变量
var DB *gorm.DB

func VerifyToken(strToken string) (string, error) {
	//解析
	token, err := jwt.ParseWithClaims(strToken, &Jwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("vinegar"), nil
	})

	if err != nil {
		return "", errors.New(ErrorReasonServerBusy + ",或token解析失败")
	}
	claims, ok := token.Claims.(*Jwt)
	if !ok {
		return "", errors.New(ErrorReasonReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return "", errors.New(ErrorReasonReLogin)
	}
	return claims.StudentID, nil
}

func UpdateUserInfo(user User) error {
	var u User
	result := DB.Model(&u).Update(user)
	return result.Error
}

func CreateBackpad(id string, backpad Backpad) error {
	backpad.Time = time.Now()
	backpad.StudentID = id
	// 0 = 未完成
	backpad.State = 0
	if result := DB.Create(&backpad); result.Error != nil {
		return result.Error
	}
	return nil
}

func ChangeBackpad(id string, backpad Backpad) error {
	var b Backpad
	DB.Where("student_id = ? AND name = ? ", id, backpad.Name).First(&b)
	result := DB.Model(&b).Update(backpad)
	return result.Error
}

func GetBackpads(id string) ([]Backpad, error) {
	var backdrops []Backpad
	result := DB.Where("student_id = ?", id).Find(&backdrops)
	if result.Error != nil {
		return backdrops, result.Error
	}
	return backdrops, nil
}

func GetUserInfo(id string) (User, error) {
	var u User
	result := DB.Where("student_id = ?", id).First(&u)
	if result.Error != nil {
		return u, result.Error
	}
	return u, nil
}

func ClearBackpad(id string, backpad Backpad) (error, string) {
	var b Backpad
	DB.Where("student_id = ? AND name = ? ", id, backpad.Name).First(&b)
	backpad.State = 1
	if result := DB.Model(&b).Update(backpad); result.Error != nil {
		return result.Error, ""
	}
	u, _ := GetUserInfo(id)
	if u.Gold < 500 {
		return nil, "金币不足"
	}
	u.Gold -= 500
	err := UpdateUserInfo(u)
	return err, ""
}

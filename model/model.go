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
		return []byte("sugar"), nil
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

func CreateBackpad(id string, backpad Backpad) (error, string) {
	var Backpad Backpad
	today := time.Now().YearDay()
	result := DB.Where("student_id = ? AND name = ? AND day = ?", id, backpad.Name, today).First(&Backpad)
	if result.Error == nil {
		return nil, "失败，该用户今日已使用该待办名"
	}
	backpad.Day = today
	backpad.Time = time.Now()
	backpad.StudentID = id
	// 0 = 未完成
	backpad.State = 0
	if result := DB.Create(&backpad); result.Error != nil {
		return result.Error, ""
	}
	return nil, ""
}

func ChangeBackpad(id string, backpad Backpad) error {
	var b Backpad
	DB.Where("student_id = ? AND name = ? ", id, backpad.Name).First(&b)
	backpad.State = 2
	result := DB.Model(&b).Update(backpad)
	return result.Error
}

func GetBackpads(id string) ([]Backpad, error) {
	var backpads []Backpad
	result := DB.Where("student_id = ?", id).Find(&backpads)
	if result.Error != nil {
		return backpads, result.Error
	}
	return backpads, nil
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
	if b.State == 1 || b.State == 2 {
		return nil, "该待办已完成或已取消"
	}
	u, _ := GetUserInfo(id)
	if u.Gold < 500 {
		return nil, "金币不足"
	}
	backpad.State = 1
	if result := DB.Model(&b).Update(backpad); result.Error != nil {
		return result.Error, ""
	}
	u.Gold -= 500
	err := UpdateUserInfo(u)
	return err, ""
}

func CompleteBackpad(id string, backpad Backpad, time int) error {
	backpad.State = 1
	var b Backpad
	DB.Where("student_id = ? AND name = ? ", id, backpad.Name).First(&b)
	if result := DB.Model(&b).Update(backpad); result.Error != nil {
		return result.Error
	}
	u, _ := GetUserInfo(id)
	u.Gold += time
	err := UpdateUserInfo(u)
	return err
}

func GetSkins(id string) ([]Skin, error) {
	var skins []Skin
	result := DB.Where("student_id = ? ", id).Find(&skins)
	return skins, result.Error
}

func BuySkin(id string, skin Skin) (error, string) {
	if skin.SkinID > 3 || skin.SkinID < 1 {
		return nil, "未找到该皮肤"
	}
	//判断是否已购买前面的皮肤
	var skins []Skin
	DB.Where("student_id = ?", id).Find(&skins)
	if skin.SkinID == 1 {
		skin.Price = 500
		if len(skins) == 1 {
			return nil, "已拥有"
		}
	}
	if skin.SkinID == 2 {
		skin.Price = 1500
		if len(skins) == 0 {
			return nil, "未购买1号皮肤"
		}
		if len(skins) == 2 || len(skins) == 3 {
			return nil, "已拥有"
		}
	}
	if skin.SkinID == 3 {
		skin.Price = 3000
		if len(skins) < 2 {
			return nil, "未购买2号皮肤"
		}
		if len(skins) == 3 {
			return nil, "已拥有"
		}
	}
	u, _ := GetUserInfo(id)
	if u.Gold < skin.Price {
		return nil, "金币不足"
	}
	skin.StudentID = id
	if result := DB.Create(skin); result.Error != nil {
		return result.Error, ""
	}
	u.Gold -= skin.Price
	err := UpdateUserInfo(u)
	return err, ""
}

func BuyFlower(id string, n int) (error, string) {
	u, _ := GetUserInfo(id)
	if u.Gold < 240*n {
		return nil, "金币不足"
	}
	u.Gold -= 240 * n
	u.Flower += n
	err := UpdateUserInfo(u)
	return err, ""
}

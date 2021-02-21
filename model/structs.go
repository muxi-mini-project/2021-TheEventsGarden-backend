package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	StudentID string `json:"student_id"`
	jwt.StandardClaims
}

type User struct {
	StudentID   string `json:"student_id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	UserPicture string `json:"user_picture"`
	Summary     string `json:"summary"`
	Sex         int    `json:"sex"`
	Gold        int    `json:"gold"`
	Flower      int    `json:"flower"`
}

type Homework struct {
	Teacher string `json:"teacher"`
	Time    string `json:"time"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Class   string `json:"class"`
	Content string `json:"content"`
}

type Backpad struct {
	StudentID string    `json:"student_id"`
	Name      string    `json:"name"`
	Time      time.Time `json:"time"`
	Hours     int       `json:"hours"`
	Minutes   int       `json:"minutes"`
	State     int       `json:"state"`
	Day       int       `json:"day"`
}

type Skin struct {
	StudentID string `json:"student_id"`
	SkinID    int    `json:"skin_id"`
	Price     int    `json:"price"`
}

//

//

//

type Response struct {
	Code uint8  `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UserInfoVO struct {
			ID       string `json:"id"`
			Archived bool   `json:"archived"`
			Username string `json:"username"`
			Phone    string `json:"phone"`
			Email    string `json:"email"`
			UserInfo struct {
				ID             string `json:"id"`
				RealName       string `json:"realname"`
				Sex            string `json:"sex"`
				Age            uint   `json:"age"`
				Phone          string `json:"phone"`
				Email          string `json:"email"`
				QQ             string `json:"qq"`
				WeChat         string `json:"wechat"`
				DepartmentCode string `json:"departmentCode"`
				University     string `json:"university"`
				LoginName      string `json:"loginName"`
				HeadImageURL   string `json:"headImageUrl"`
				Sign           string `json:"sign"`
				AddTime        int    `json:"addtime"`
				UpdateTime     int    `json:"updatetime"`
			} `json:"userInfo"`
		} `json:"userInfoVO"`
		RoleDepartment struct {
			ID             string `json:"id"`
			LoginName      string `json:"loginName"`
			DomainCode     string `json:"domainCode"`
			RoleCode       string `json:"roleCode"`
			DepartmentCode string `json:"departmentCode"`
			AddTime        int    `json:"addtime"`
			UpdateTime     int    `json:"updatetime"`
			DomainName     string `json:"domainName"`
			DepartmentName string `json:"departmentName"`
			RoleName       string `json:"roleName"`
			RealName       string `json:"realname"`
		} `json:"roleDepartment"`
	} `json:"data"`
}

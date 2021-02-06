package model

import (
	"github.com/dgrijalva/jwt-go"
)

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

type Jwt struct {
	StudentID string `json:"student_id"`
	jwt.StandardClaims
}

type User struct {
	StudentID string `json:"id"`
	Password  string `json:"password"`
}

type homework struct {
	Teacher string `json:"teacher"`
	Time    string `json:"time"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Class   string `json:"class"`
	Content string `json:"content"`
}

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tidwall/gjson"
	"golang.org/x/net/publicsuffix"
)

const (
	ErrorReasonServerBusy = "服务器繁忙"
	ErrorReasonReLogin    = "请重新登陆"
)

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

func AddHeaders(request *http.Request) *http.Request {
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36")
	request.Header.Set("Host", "spoc.ccnu.edu.cn")
	request.Header.Set("Origin", "http://spoc.ccnu.edu.cn")
	request.Header.Set("Referer", "http://spoc.ccnu.edu.cn/studentHomepage")
	return request
}

func GetHomework(id string, pwd string) []homework {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	fmt.Println(pwd)
	_, Homeworkss, err := LoginSPOC(id, pwd, client)
	if err != nil {
		panic(err)
	}
	//log.Println(response)
	//log.Println(Homeworkss)
	return Homeworkss
}

func NewClient() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Jar:     jar,
	}
	return &client, nil
}

func LoginSPOC(sno, password string, client *http.Client) (*Response, []homework, error) {
	v := url.Values{}
	v.Set("loginName", sno)
	v.Set("password", password)

	/*request, err := http.NewRequest("GET", "http://spoc.ccnu.edu.cn/userLoginController/getVerifCode", nil)
	if err != nil {
		return nil, err
	}
	_, err = client.Do(request)
	if err != nil {
		return nil, err
	}*/
	request, err := http.NewRequest("POST", "http://spoc.ccnu.edu.cn/userLoginController/getUserProfile", strings.NewReader(v.Encode()))
	if err != nil {
		fmt.Println("11111")
		return nil, nil, err
	}

	request = AddHeaders(request)
	_, err = client.Do(request)
	if err != nil {
		fmt.Println("22222")
		return nil, nil, err
	}

	request, err = http.NewRequest("POST", "http://spoc.ccnu.edu.cn/userInfo/getUserInfo", nil)
	if err != nil {
		fmt.Println("33333")
		return nil, nil, err
	}
	request = AddHeaders(request)
	resp, respErr := client.Do(request)
	fmt.Printf("%s\n", resp)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	response2 := Response{}
	err = json.Unmarshal(body, &response2)
	//fmt.Printf("%s\n", response2)
	//fmt.Printf("%s\n", body)
	if err != nil {
		fmt.Println("44444")
		return nil, nil, err
	}

	if respErr != nil {
		fmt.Println("55555")
		return nil, nil, respErr
	}

	// 爬信息
	//payload := strings.NewReader(`{"userId":"` + response2.Data.UserInfoVO.ID + `","termCode":"202101","pageNum":1,"pageSize":4}`)

	//fmt.Println(payload)
	payload := strings.NewReader("")
	//request, err = http.NewRequest("POST", "http://spoc.ccnu.edu.cn/studentHomepage/getMySite", payload)
	request, err = http.NewRequest("GET", "http://spoc.ccnu.edu.cn/studentHomepage/getHistorySite?termCode=202002&domainCode=1", payload)
	request.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(request)
	body, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var sites, names, teachers []string
	for i := 0; ; i++ {
		num := strconv.Itoa(i)
		//siteID := gjson.Get(string(body), "data.list."+num+".siteId")
		siteID := gjson.Get(string(body), "data."+num+".siteId")
		if siteID.String() == "" {
			break
		}
		// siteName := gjson.Get(string(body), "data.list."+num+".siteName")
		siteName := gjson.Get(string(body), "data."+num+".siteName")
		teacher := gjson.Get(string(body), "data."+num+".teacherName")
		sites = append(sites, siteID.String())
		names = append(names, siteName.String())
		teachers = append(teachers, teacher.String())

	}
	var homeworks []homework
	var num string
	for i := 0; i < len(sites); i++ {
		num = strconv.Itoa(i)
		//var contents []string
		payload = strings.NewReader(`{"siteId":"` + sites[i] + `","pageNum":1,"pageSize":5}`)
		request, err = http.NewRequest("POST", "http://spoc.ccnu.edu.cn/assignment/getStudentAssignmentList", payload)
		request.Header.Add("Content-Type", "application/json")
		resp, err = client.Do(request)
		body, _ = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		for j := 0; ; j++ {

			num = strconv.Itoa(j)
			status := gjson.Get(string(body), "data.list."+num+".status")

			if status.Int() == 0 || status.Int() == 2 {
				var Homework homework
				//teacher := gjson.Get(string(body), "data.list."+num+".teacher")
				Homework.Teacher = teachers[i]
				Time := gjson.Get(string(body), "data.list."+num+".endtime")
				x := time.Unix(Time.Int()/1000, 0)
				Homework.Time = x.Format("2006-01-04 03:04:05")
				//Homework.Time = time.Time()
				content := gjson.Get(string(body), "data.list."+num+".content")
				Homework.Content = Content(content.String())
				title := gjson.Get(string(body), "data.list."+num+".title")
				if title.String() == "" {
					break
				}
				Homework.Title = title.String()
				if status.Int() == 0 {
					Homework.Status = "未提交"
				} else {
					Homework.Status = "已驳回"
				}
				Homework.Class = names[i]
				homeworks = append(homeworks, Homework)
			}
		}
	}
	return nil, homeworks, nil
}

//去掉content里面的htmi 如 <> 与 &nbsp
func Content(s string) string {
	ss := []byte(s)
	var s2 []byte
	for i := 0; i < len(ss); i++ {
		if ss[i] == '<' {
		LABEL1:
			for ss[i] != '>' {
				i++
			}
			if i == len(ss)-1 {
				break
			}
			i++
			if ss[i] == '<' {
				goto LABEL1
			}
		}
	LABEL2:
		if ss[i] == '&' {
			i += 6
			if ss[i] == '&' {
				goto LABEL2
			}
		}
		s2 = append(s2, ss[i])
	}
	return string(s2)
}

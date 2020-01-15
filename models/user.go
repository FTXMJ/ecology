package models

import (
	"ecology/consul"
	"encoding/json"

	"github.com/astaxie/beego"

	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//用户表
type User struct {
	Id       int    `orm:"column(id);pk;auto"`
	UserId   string `orm:"column(user_id)"`   //对应 monggodb 的user_id
	FatherId string `orm:"column(father_id)"` //父亲id
	UserName string `orm:"column(user_name)"` //父亲id
}

type Response struct {
	Code int                    `json:"code""`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// 调用远端接口
func PingUserAdmin(token, user_id string) (interface{}, interface{}, error) {
	client := &http.Client{}
	//生成要访问的url
	url := consul.GetAuthApi + beego.AppConfig.String("api::apiurl_auth_get_user")

	//提交请求
	reqest, errnr := http.NewRequest("GET", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", token)
	q := reqest.URL.Query()
	q.Add("user_id", user_id)

	reqest.URL.RawQuery = q.Encode()

	if errnr != nil {
		return "", "", errnr
	}
	//处理返回结果
	response, errdo := client.Do(reqest)

	if errdo != nil {
		return "", "", errdo
	}
	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return "", "", err_read
	}
	fmt.Println(string(bys))
	values := Response{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return "", "", errors.New("解析错误")
	}
	if values.Data == nil {
		return "", "", errors.New(values.Msg)
	}
	response.Body.Close()
	return values.Data["father_id"], values.Data["nick_name"], nil
}

// 调用远端接口
func PingUser(token, user_id string) (interface{}, interface{}, error) {
	client := &http.Client{}
	//生成要访问的url
	url := consul.GetUserApi + beego.AppConfig.String("api::apiurl_user_get_user")

	//提交请求
	reqest, errnr := http.NewRequest("GET", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", token)

	if errnr != nil {
		return "", "", errnr
	}
	//处理返回结果
	response, errdo := client.Do(reqest)

	if errdo != nil {
		return "", "", errdo
	}
	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return "", "", err_read
	}
	fmt.Println(string(bys))
	values := Response{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return "", "", errors.New("解析错误")
	}
	if values.Data == nil {
		return "", "", errors.New(values.Msg)
	}
	response.Body.Close()
	return values.Data["father_id"], values.Data["nick_name"], nil
}

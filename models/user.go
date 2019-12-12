package models

import (
	"ecology1/consul"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

//用户表
type User struct {
	Id       int    `orm:"column(id);pk;auto"`
	Name     string `orm:column(name)`      // 对应 monggodb 的fallname
	UserId   string `orm:column(user_id)`   //对应 monggodb 的user_id
	FatherId string `orm:column(father_id)` //父亲id
}

func (this *User) TableName() string {
	return "user"
}

func (this *User) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *User) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

// 调用远端接口
func PingUser(token string) (interface{}, error) {
	client := &http.Client{}
	//生成要访问的url
	url := consul.GetUserApi + beego.AppConfig.String("consul::apiurl_get_user")
	//提交请求
	reqest, errnr := http.NewRequest("GET", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", token)

	if errnr != nil {
		return "", errnr
	}
	//处理返回结果
	response, errdo := client.Do(reqest)

	if errdo != nil {
		return "", errdo
	}
	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return "", err_read
	}
	values := make(map[string]interface{})
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return "", err
	}
	response.Body.Close()
	return values["father_id"], nil
}

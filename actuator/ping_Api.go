package actuator

import (
	"ecology/consul"
	"ecology/filter"
	"ecology/models"
	"encoding/json"

	"github.com/astaxie/beego"

	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type data_users struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}

// 调用远端接口
func PingWalletAdd(token string, coin_number float64) error {
	//生成要访问的url
	apiurl := consul.GetWalletApi
	resoure := beego.AppConfig.String("api::apiurl_get_all_wallet")
	data := url.Values{}
	data.Set("money", strconv.FormatFloat(coin_number, 'f', -1, 64))
	data.Set("symbol", "USDD")

	u, _ := url.ParseRequestURI(apiurl)
	u.Path = resoure
	urlStr := u.String()

	client := &http.Client{}
	req, err1 := http.NewRequest(`POST`, urlStr, strings.NewReader(data.Encode()))
	if err1 != nil {
		return err1
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	//处理返回结果
	response, errdo := client.Do(req)
	if errdo != nil {
		return errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return err_read
	}
	values := models.Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return errors.New(values.Msg)
	}
	response.Body.Close()
	return nil
}

// TFOR 数量查询
func PingSelectTforNumber(user_id string) (string, float64, error) {
	user := models.User{
		UserId: user_id,
	}
	b, token := filter.GenerateToken(user)
	if b != true {
		return "", 0.0, errors.New("err")
	}
	//生成要访问的url
	apiurl := consul.GetWalletApi
	resoure := beego.AppConfig.String("api::apiurl_tfor_info")
	data := url.Values{}

	u, _ := url.ParseRequestURI(apiurl)
	u.Path = resoure
	urlStr := u.String()

	client := &http.Client{}
	req, err1 := http.NewRequest(`GET`, urlStr, strings.NewReader(data.Encode()))
	if err1 != nil {
		return "", 0.0, err1
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	//处理返回结果
	response, errdo := client.Do(req)
	if errdo != nil {
		return "", 0.0, errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return "", 0.0, err_read
	}
	values := models.Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return "", 0.0, errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return "", 0.0, errors.New(values.Msg)
	}
	response.Body.Close()
	aa := values.Data["balance"].(string)
	time := values.Data["updated_at"].(string)
	bb, err := strconv.ParseFloat(aa, 64)
	return time, bb, nil
}

// 远端连接  -  给定分红收益  释放通用
func PingAddWalletCoin(user_id string, abonus float64) error {
	if abonus == 0 {
		return nil
	}
	user := models.User{
		UserId: user_id,
	}
	b, token := filter.GenerateToken(user)
	if b != true {
		return errors.New("err")
	}
	//生成要访问的url
	apiurl := consul.GetWalletApi
	resoure := beego.AppConfig.String("api::apiurl_share_bonus")
	data := url.Values{}
	data.Set("money", strconv.FormatFloat(abonus, 'f', -1, 64))
	data.Set("symbol", "USDD")

	u, _ := url.ParseRequestURI(apiurl)
	u.Path = resoure
	urlStr := u.String()

	client := &http.Client{}
	req, err1 := http.NewRequest(`POST`, urlStr, strings.NewReader(data.Encode()))
	if err1 != nil {
		return err1
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	//处理返回结果
	response, errdo := client.Do(req)
	if errdo != nil {
		return errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return err_read
	}
	values := models.Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return errors.New(values.Msg)
	}
	response.Body.Close()
	return nil
}

// 从晓东那里获取团队 成员  直推
func GetTeams(user models.User) ([]string, error) {
	client := &http.Client{}
	//生成要访问的url
	url := consul.GetUserApi + beego.AppConfig.String("api::apiurl_get_team")
	//提交请求
	reqest, errnr := http.NewRequest("GET", url, nil)

	b, token := filter.GenerateToken(user)
	if b != true {
		return nil, errors.New("err")
	}

	//增加header选项
	reqest.Header.Add("Authorization", token)
	if errnr != nil {
		return nil, errnr
	}

	//处理返回结果
	response, errdo := client.Do(reqest)
	defer response.Body.Close()
	if errdo != nil {
		return nil, errdo
	}
	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return nil, err_read
	}
	values := data_users{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return values.Data, errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return values.Data, errors.New(values.Msg)
	}
	users := []string{}
	for _, v := range values.Data {
		users = append(users, v)
	}
	return users, nil
}

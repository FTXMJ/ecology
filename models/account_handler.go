package models

import (
	"ecology/common"
	"ecology/consul"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func SendHttpPost(urls string, api string, data map[string]string, token string) (error, *common.ResponseData) {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	u, err := url.ParseRequestURI("http://" + urls)
	if err != nil {
		//log.Log.Error(err)
	}
	u.Path = api
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		r.Header.Set("Authorization", token)
	}
	resp, err := client.Do(r)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var rep common.ResponseData
	json.Unmarshal(body, &rep)
	if rep.Code != 200 {
		return errors.New(rep.Msg), nil
	}
	return nil, &rep
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
	values := Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return errors.New(values.Msg)
	}
	response.Body.Close()
	return nil
}

// 判断是否达到超级节点的要求 --- 页面显示
func SuperLevelSet(user_id string, ec_obj *Ecology_index_obj) {
	s_f_t := []SuperForceTable{}
	NewOrm().QueryTable("super_force_table").All(&s_f_t)
	s_p_t := SuperPeerTable{}
	NewOrm().QueryTable("super_peer_table").Filter("user_id", user_id).One(&s_p_t)

	for i := 0; i < len(s_f_t); i++ {
		for j := 1; j < len(s_f_t)-1; j++ {
			if s_f_t[i].CoinNumberRule > s_f_t[j].CoinNumberRule {
				s_f_t[i].CoinNumberRule, s_f_t[j].CoinNumberRule = s_f_t[j].CoinNumberRule, s_f_t[i].CoinNumberRule
			}
		}
	}
	index := []int{}
	for i, v := range s_f_t {
		if s_p_t.CoinNumber > float64(v.CoinNumberRule) {
			index = append(index, i)
		}
	}
	if len(index) > 0 {
		ec_obj.Super_peer_bool = true
		ec_obj.Super_peer.Level = s_f_t[index[len(index)-1]].Level
		ec_obj.Super_peer.Usdd = s_p_t.CoinNumber
	}
}

//生成凑数信息
func CouShu(a *Ecology_index_obj) {
	b := Formulaindex{
		Id:                  0,
		Level:               "",
		BockedBalance:       0,
		Balance:             0,
		LowHold:             0,
		HighHold:            0,
		ReturnMultiple:      0,
		ToDayRate:           0,
		HoldReturnRate:      0,
		RecommendReturnRate: 0,
		TeamReturnRate:      0,
	}
	a.Ecological_poject = append(a.Ecological_poject, b)
}

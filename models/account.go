package models

import (
	"ecology/common"
	"ecology/consul"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//生态钱包
type Account struct {
	Id            int     `orm:"column(id);pk;auto"`
	UserId        string  `orm:"column(user_id)"`      //用户 Id
	Balance       float64 `orm:"column(balance)"`      //充值交易结余
	Currency      string  `orm:column(currency)`       //货币  USDD
	BockedBalance float64 `orm:column(bocked_balance)` //铸币交易结余
	Level         string  `orm:column(level)`          //等级
}

func (this *Account) TableName() string {
	return "account"
}

func (this *Account) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *Account) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

type Ecology_index_obj struct {
	Usdd                   float64        `json:"usdd"`
	Ecological_poject      []Formulaindex `json:"ecological_poject"` //生态项目
	Ecological_poject_bool bool           `json:"ecological_poject_bool"`
	Super_peer             SuperPeer      `json:"super_peer"` //超级节点
	Super_peer_bool        bool           `json:"super_peer_bool"`
}

//页面显示的　超级节点结构
type SuperPeer struct {
	Usdd        float64 `json:"usdd"`          //总币数ForceTable_test
	Level       string  `json:"level"`         //超级节点的独立属性
	TodayABouns float64 `json:"today_a_bouns"` // 今日分红
}

// 页面显示的　生态仓库结构
type Formulaindex struct {
	Id                  int     `json:"id"`
	Level               string  `json:"level"`
	BockedBalance       float64 `json:"bocked_balance"`        //持币数量
	Balance             float64 `json:"balance"`               //投资总额
	LowHold             int     `json:"low_hold"`              //低位
	HighHold            int     `json:"high_hold"`             //高位
	ReturnMultiple      float64 `json:"return_multiple"`       //杠杆
	ToDayRate           float64 `json:"to_day_rate"`           //今日算力
	HoldReturnRate      float64 `json:"hold_return_rate"`      //本金自由算力
	RecommendReturnRate float64 `json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `json:"team_return_rate"`      //动态算力
}

//从　body.form 读数据的数据格式
type NetValueType struct {
	CoinNumber float64 `form:"coin_number" json:"coin_number"`
	EcologyId  int     `form:"ecology_id" json:"ecology_id"`
	LevelStr   string  `form:"levelstr" json:"level_str"`
}

type Data_wallet struct {
	Code int                    `json:"code""`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// 老罗的钱包数据结构
type WalletInfo struct {
	Balance    float64 `json:"balance"`
	CurrencyId int     `json:"currency_id"`
	Decimals   int     `json:"decimals"`
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
}

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
	fmt.Println(string(bys))
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

type Ecology_index_ob_test struct {
	usdd___usdd数量                             float64
	ecological_poject___生态项目                  []Formulaindex_test //生态项目
	ecological_poject_bool___是否有生态仓库没有就是false bool
	super_peer___超级节点信息                       SuperPeer_test //超级节点
	super_peer_bool__是否显示超级节点                 bool
}

//页面显示的　超级节点结构
type SuperPeer_test struct {
	usdd___总币数ForceTable_test float64 //总币数ForceTable_test
	level___超级节点的独立属性         string  //超级节点的独立属性
	today_a_bouns___今日分红      float64 // 今日分红
}

// 页面显示的　生态仓库结构
type Formulaindex_test struct {
	id__生态仓库id                   int
	level___等级                   string
	bocked_balance___持币数量        float64 //持币数量
	balance___投资总额               float64 //投资总额
	low_hold___低位                int     //低位
	high_hold___高位               int     //高位
	return_multiple___杠杆         float64 //杠杆
	to_day_rate___今日算力           float64 //今日算力
	hold_return_rate___自由算力      float64 //本金自由算力
	recommend_return_rate___直推算力 float64 //直推算力
	team_return_rate____动态算力     float64 //动态算力
}

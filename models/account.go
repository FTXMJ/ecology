package models

import (
	"ecology/consul"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strconv"
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
	Usdd                   float64
	Ecological_poject      []Formulaindex //生态项目
	Ecological_poject_bool bool
	Super_peer             SuperPeer //超级节点
	Super_peer_bool        bool
}

//页面显示的　超级节点结构
type SuperPeer struct {
	Usdd        float64 //总币数ForceTable_test
	Level       string  //超级节点的独立属性
	TodayABouns float64 // 今日分红
}

// 页面显示的　生态仓库结构
type Formulaindex struct {
	Id                  int
	Level               string
	BockedBalance       float64 //持币数量
	LowHold             int     //低位
	HighHold            int     //高位
	ReturnMultiple      float64 //杠杆
	ToDayRate           float64 //今日算力
	HoldReturnRate      float64 //本金自由算力
	RecommendReturnRate float64 //直推算力
	TeamReturnRate      float64 //动态算力
}

//从　body.form 读数据的数据格式
type NetValueType struct {
	CoinNumber float64 `form:"coin_number" json:"coin_number"`
	EcologyId  int     `form:"ecology_id" json:"ecology_id"`
	LevelStr   string  `form:"levelstr" json:"level_str"`
}

type Data_wallet struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// 老罗的钱包数据结构
type WalletInfo struct {
	Balance    float64 `json:"balance"`
	CurrencyId int     `json:"currency_id"`
	Decimals   int     `json:"decimals"`
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
}

// 调用远端接口
func PingWallet(token string, coin_number float64) error {
	client := &http.Client{}
	//生成要访问的url
	url := consul.GetWalletApi + beego.AppConfig.String("api::apiurl_get_all_wallet")
	//提交请求
	reqest, errnr := http.NewRequest("GET", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", token)
	coin := strconv.FormatFloat(coin_number, 'E', -1, 64)
	reqest.Form.Add("money", coin)
	reqest.Form.Add("cion_type", "USDD")
	if errnr != nil {
		return errnr
	}

	//处理返回结果
	response, errdo := client.Do(reqest)
	if errdo != nil {
		return errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return err_read
	}

	values := Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil || values.Code != 200 {
		return err
	}
	response.Body.Close()
	return err
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
	Usdd___usdd数量                             float64
	Ecological_poject___生态项目                  []Formulaindex_test //生态项目
	Ecological_poject_bool___是否有生态仓库没有就是false bool
	Super_peer___超级节点信息                       SuperPeer_test //超级节点
	Super_peer_bool__是否显示超级节点                 bool
}

//页面显示的　超级节点结构
type SuperPeer_test struct {
	Usdd___总币数ForceTable_test float64 //总币数ForceTable_test
	Level___超级节点的独立属性         string  //超级节点的独立属性
	TodayABouns___今日分红        float64 // 今日分红
}

// 页面显示的　生态仓库结构
type Formulaindex_test struct {
	Id__生态仓库id                 int
	Level___等级                 string
	BockedBalance___持币数量       float64 //持币数量
	LowHold___低位               int     //低位
	HighHold___高位              int     //高位
	ReturnMultiple___杠杆        float64 //杠杆
	ToDayRate___今日算力           float64 //今日算力
	HoldReturnRate___自由算力      float64 //本金自由算力
	RecommendReturnRate___直推算力 float64 //直推算力
	TeamReturnRate____动态算力     float64 //动态算力
}

package models

import (
	"encoding/json"
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
	Usdd        float64 //总币数
	Level       string  //超级节点的独立属性
	TodayABouns float64 // 今日分红
}

// 页面显示的　生态仓库结构
type Formulaindex struct {
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
	EcologyId int `form:"ecology_id" json:"ecology_id"`
	LevelStr string `form:"levelstr" json:"level_str"`
}

// 调用远端接口
func PingWallet(token string, coin_id string) (float64,error) {
	client := &http.Client{}
	//生成要访问的url
	url := ""
	//提交请求
	reqest, errnr := http.NewRequest("POST", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", token)

	if errnr != nil {
		return 0,errnr
	}
	//处理返回结果
	response, errdo := client.Do(reqest)
	defer response.Body.Close()
	if errdo!=nil{
		return  0,errdo
	}
	b,_ := ioutil.ReadAll(response.Body)
	m := make(map[string]string)
	json.Unmarshal(b,&m)

	tfor_number,_ := strconv.ParseFloat(m["money"], 64)
	return tfor_number,nil
}

func SuperLevelSet(usdd float64,ec_obj *Ecology_index_obj) {
	if usdd >= 30000 && usdd < 100000 && usdd < 200000 {
		ec_obj.Super_peer_bool = true
		ec_obj.Super_peer.Level = "分红节点"
		ec_obj.Super_peer.Usdd = usdd
		// TODO  分红
		ec_obj.Super_peer.TodayABouns = 998
	}else if usdd >= 100000 && usdd < 200000 {
		ec_obj.Super_peer_bool = true
		ec_obj.Super_peer.Level = "超级节点"
		ec_obj.Super_peer.Usdd = usdd
		// TODO  分红
		ec_obj.Super_peer.TodayABouns = 998
	}else if usdd >= 200000 {
		ec_obj.Super_peer_bool = true
		ec_obj.Super_peer.Level = "创世节点"
		ec_obj.Super_peer.Usdd = usdd
		// TODO  分红
		ec_obj.Super_peer.TodayABouns = 998
	}else {
		ec_obj.Super_peer_bool = false
	}
}


/*
接口地址　 :=    ???
token    :=    在头里面
访问类型　:=    post
返回给我　:=    用户的　ＴＦＯＲ　数量    [键值对返回　　－　放在　body　里]
*/
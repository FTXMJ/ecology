package models

//生态钱包
type Account struct {
	Id             int     `orm:"column(id);pk;auto"`
	UserId         string  `orm:"column(user_id)"`         //用户 Id
	Balance        float64 `orm:"column(balance)"`         //充值交易结余
	Currency       string  `orm:"column(currency)"`        //货币  USDD
	BockedBalance  float64 `orm:"column(bocked_balance)"`  //铸币交易结余
	Level          string  `orm:"column(level)"`           //等级
	CreateDate     string  `orm:"column(create_date)"`     // 创建时间
	DynamicRevenue bool    `orm:"column(dynamic_revenue)"` //动态收益开关
	StaticReturn   bool    `orm:"column(static_return)"`   //静态收益开关
	PeerState      bool    `orm:"column(peer_state)"`
	UpdateDate     string  `orm:"column(update_date)"` //静态收益开关
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

// 生态首页展示
type Ecology_index_obj struct {
	Usdd                   float64        `json:"usdd"`
	Ecological_poject      []Formulaindex `json:"ecological_poject"` //生态项目
	Ecological_poject_bool bool           `json:"ecological_poject_bool"`
	Super_peer             SuperPeer      `json:"super_peer"` //超级节点
	Super_peer_bool        bool           `json:"super_peer_bool"`
}

//　生态首页展示－－超级节点结构
type SuperPeer struct {
	Usdd        float64 `json:"usdd"`          //总币数ForceTable_test
	Level       string  `json:"level"`         //超级节点的独立属性
	TodayABouns float64 `json:"today_a_bouns"` // 今日分红
}

// 生态首页展示－－生态仓库结构
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

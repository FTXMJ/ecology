package models

// COINS table
type BlockedDetail struct {
	Id             int     `gorm:"column:id;primary_key"`
	UserId         string  `gorm:"column:user_id" json:"user_id"`
	CurrentRevenue float64 `gorm:"column:current_revenue" json:"current_revenue"` //本期收入
	CurrentOutlay  float64 `gorm:"column:current_outlay" json:"current_outlay"`   //本期支出
	OpeningBalance float64 `gorm:"column:opening_balance" json:"opening_balance"` //上期余额
	CurrentBalance float64 `gorm:"column:current_balance" json:"current_balance"` //本期余额
	CreateDate     string  `gorm:"column:create_date" json:"create_date"`         //创建时间
	Comment        string  `gorm:"column:comment" json:"comment"`                 //评论
	TxId           string  `gorm:"column:tx_id" json:"tx_id"`                     //任务id
	Account        int     `gorm:"column:account" json:"account"`                 //生态仓库id
	CoinType       string  `gorm:"column:coin_type" json:"coin_type"`             // 币种信息
}

// user ecology information
type UEOBJList struct {
	Items []U_E_OBJ `json:"items"` //数据列表
	Page  Page      `json:"page"`  //分页信息
}

// user coin flow information
type FlowList struct {
	Items []Flow `json:"items"` //数据列表
	Page  Page   `json:"page"`  //分页信息
}

// history information
type HostryFindInfo struct {
	Items []BlockedDetailIndex `json:"items"` //数据列表
	Page  Page                 `json:"page"`  //分页信息
}

// user`s account OFF information
type UserAccountOFF struct {
	Items []AccountOFF `json:"items"` //数据列表
	Page  Page         `json:"page"`  //分页信息
}

// user`s account false table
type UserFalse struct {
	Items []FalseUser `json:"items"` //数据列表
	Page  Page        `json:"page"`  //分页信息
}

// user`s account false table
type MrsfTable struct {
	Items []MrsfStateTable `json:"items"` //数据列表
	Page  Page             `json:"page"`  //分页信息
}

// user`s peer history table
type PeerHistoryList struct {
	Items []PeerHistory `json:"items"` //数据列表
	Page  Page          `json:"page"`  //分页信息
}

// user`s peer table
type PeerUserFalse struct {
	Items []PeerUser `json:"items"` //数据列表
	Page  Page       `json:"page"`  //分页信息
}

type PeerUser struct {
	AccountId  int     `json:"account_id"`
	UserName   string  `json:"user_name"`
	UserId     string  `json:"user_id"`
	Level      string  `json:"level"`
	State      bool    `json:"state"`
	Number     float64 `json:"number"`
	UpdateTime string  `json:"update_time"`
}

type FalseUser struct {
	AccountId  int    `json:"account_id"`
	UserName   string `json:"user_name"`
	UserId     string `json:"user_id"`
	Jintai     bool   `json:"jintai"`
	Dongtai    bool   `json:"dongtai"`
	UpdateTime string `json:"update_time"`
}

// The query object to `user * information`
type FindObj struct {
	UserId    string
	UserName  string
	TxId      string
	StartTime string
	EndTime   string
}

// user ecology information
type Flow struct {
	UserId              string  `json:"user_id"`
	UserName            string  `json:"user_name"`
	HoldReturnRate      float64 `json:"hold_return_rate"`      //本金自由算力
	RecommendReturnRate float64 `json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `json:"team_return_rate"`      //动态算力
	Released            float64 `json:"released"`              //已释放
	UpdateTime          string  `json:"update_time"`           // 最后更新时间
}

// user ecology information object
type U_E_OBJ struct {
	UserId              string  `json:"user_id"`
	UserName            string  `json:"user_name"`
	Level               string  `json:"level"`
	ReturnMultiple      float64 `json:"return_multiple"`       //杠杆
	CoinAll             float64 `json:"coin_all"`              //存币总和
	ToBeReleased        float64 `json:"to_be_released"`        //待释放
	Released            float64 `json:"released"`              //已释放
	HoldReturnRate      float64 `json:"hold_return_rate"`      //本金自由算力
	RecommendReturnRate float64 `json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `json:"team_return_rate"`      //动态算力
}

// user`s account OFF table
type AccountOFF struct {
	UserId         string `json:"user_id"`
	UserName       string `json:"user_name"`
	Account        int    `json:"account"`
	DynamicRevenue bool   `json:"dynamic_revenue"` //动态收益开关
	StaticReturn   bool   `json:"static_return"`   //静态收益开关
	PeerState      bool   `json:"peer_state"`
	CreateDate     string `json:"create_date"` //创建时间
}

type BlockedDetailIndex struct {
	Id                int     `json:"id"`
	UserId            string  `json:"user_id"`
	UserName          string  `json:"user_name"`
	AccCurrentRevenue float64 `json:"acc_current_revenue"` //转入数量
	BloCurrentRevenue float64 `json:"blo_current_revenue"` //铸币数量
	ReturnMultiple    int     `json:"return_multiple"`     // 铸币倍数
	CreateDate        string  `json:"create_date"`         //创建时间
	Comment           string  `json:"comment"`             //评论
	TxId              string  `json:"tx_id"`               //任务id
	Account           int     `json:"account"`             //生态仓库id
	CoinType          string  `json:"coin_type"`           // 币种信息
}

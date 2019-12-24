package models

// COINS table
type BlockedDetail struct {
	Id             int     `orm:"column(id);pk;auto" json:"id"`
	UserId         string  `orm:"column(user_id)" json:"user_id"`
	CurrentRevenue float64 `orm:"column(current_revenue)" json:"current_revenue"` //本期收入
	CurrentOutlay  float64 `orm:"column(current_outlay)" json:"current_outlay"`   //本期支出
	OpeningBalance float64 `orm:"column(opening_balance)" json:"opening_balance"` //上期余额
	CurrentBalance float64 `orm:"column(current_balance)" json:"current_balance"` //本期余额
	CreateDate     string  `orm:"column(create_date)" json:"create_date"`         //创建时间
	Comment        string  `orm:"column(comment)" json:"comment"`                 //评论
	TxId           string  `orm:"column(tx_id)" json:"tx_id"`                     //任务id
	Account        int     `orm:"column(account)" json:"account"`                 //生态仓库id
	CoinType       string  `orm:"column(coin_type)" json:"coin_type"`             // 币种信息
}

func (this *BlockedDetail) TableName() string {
	return "blocked_detail"
}

func (this *BlockedDetail) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *BlockedDetail) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

// user ecology information
type UEOBJList struct {
	Items []U_E_OBJ `json:"items"` //数据列表
	Page  Page      `json:"page"`  //分页信息
}

// user coin flow information
type FlowList struct {
	Items []U_E_OBJ `json:"items"` //数据列表
	Page  Page      `json:"page"`  //分页信息
}

// history information
type HostryFindInfo struct {
	Items []BlockedDetail `json:"items"` //数据列表
	Page  Page            `json:"page"`  //分页信息
}

// The query object to `user * information`
type FindObj struct {
	UserId    string
	TxId      string
	StartTime string
	EndTime   string
}

// user ecology information
type Flow struct {
	UserId              string  `json:"user_id"`
	HoldReturnRate      float64 `json:"hold_return_rate"`      //本金自由算力
	RecommendReturnRate float64 `json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `json:"team_return_rate"`      //动态算力
	Released            float64 `json:"released"`              //已释放
	UpdateTime          string  `json:"update_time"`           // 最后更新时间
}

// user ecology information object
type U_E_OBJ struct {
	UserId              string  `json:"user_id"`
	Level               string  `json:"level"`
	ReturnMultiple      float64 `json:"return_multiple"`       //杠杆
	CoinAll             float64 `json:"coin_all"`              //存币总和
	ToBeReleased        float64 `json:"to_be_released"`        //待释放
	Released            float64 `json:"released"`              //已释放
	HoldReturnRate      float64 `json:"hold_return_rate"`      //本金自由算力
	RecommendReturnRate float64 `json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `json:"team_return_rate"`      //动态算力
}

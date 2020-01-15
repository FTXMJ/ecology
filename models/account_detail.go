package models

// 交易记录表
type AccountDetail struct {
	Id             int     `orm:"column(id);pk;auto"`
	UserId         string  `orm:"column(user_id)"`
	CurrentRevenue float64 `orm:"column(current_revenue)"` //本期收入
	CurrentOutlay  float64 `orm:"column(current_outlay)"`  //本期支出
	OpeningBalance float64 `orm:"column(opening_balance)"` //上期余额
	CurrentBalance float64 `orm:"column(current_balance)"` //本期余额
	CreateDate     string  `orm:"column(create_date)"`     //创建时间
	Comment        string  `orm:"column(comment)"`         //备注
	TxId           string  `orm:"column(tx_id)"`           //任务id
	Account        int     `orm:"column(account)"`         //生态仓库id
	CoinType       string  `orm:"column(coin_type)"`       // 币种信息
}

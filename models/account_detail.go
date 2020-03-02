package models

// 交易记录表
type AccountDetail struct {
	Id             int     `gorm:"column:id;primary_key"`
	UserId         string  `gorm:"column:user_id"`
	CurrentRevenue float64 `gorm:"column:current_revenue"` //本期收入
	CurrentOutlay  float64 `gorm:"column:current_outlay"`  //本期支出
	OpeningBalance float64 `gorm:"column:opening_balance"` //上期余额
	CurrentBalance float64 `gorm:"column:current_balance"` //本期余额
	CreateDate     string  `gorm:"column:create_date"`     //创建时间
	Comment        string  `gorm:"column:comment"`         //备注
	TxId           string  `gorm:"column:tx_id"`           //任务id
	Account        int     `gorm:"column:account"`         //生态仓库id
	CoinType       string  `gorm:"column:coin_type"`       // 币种信息
}

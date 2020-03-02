package models

//交易唯一标示Id表
type TxIdList struct {
	Id          int     `gorm:"column:id;primary_key" json:"id"`
	OrderState  bool    `gorm:"column:order_state" json:"order_state"` //任务完成状态，false=未完成，true=完成
	WalletState bool    `gorm:"column:wallet_state" json:"wallet_state"`
	TxId        string  `gorm:"column:tx_id" json:"tx_id"`             //任务id
	UserId      string  `gorm:"column:user_id" json:"user_id"`         //任务id
	Comment     string  `gorm:"column:comment" json:"comment"`         // 解释注释
	CreateTime  string  `gorm:"column:expenditure" json:"expenditure"` //任务id
	Expenditure float64 `gorm:"column:expenditure" json:"expenditure"` //任务id
	InCome      float64 `gorm:"column:in_come" json:"in_come"`         //任务id
}

// peer a_bouns list
type PeerListABouns struct {
	Items []PeerAbouns `json:"items"` //数据列表
	Page  Page         `json:"page"`  //分页信息
}

type PeerAbouns struct {
	Id       int     `json:"id"`
	UserName string  `json:"user_name"`
	Level    string  `json:"level"`
	Tfors    float64 `json:"tfors"`
	Time     string  `json:"time"`
}

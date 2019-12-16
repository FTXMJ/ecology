package models

//交易唯一标示Id表
type TxIdList struct {
	Id          int     `orm:"column(id);pk;auto"`
	State       string  `orm:column(state)`       //任务完成状态，false=未完成，true=完成
	TxId        string  `orm:column(tx_id)`       //任务id
	UserId      string  `orm:column(user_id)`     //任务id
	CreateTime  string  `orm:column(create_time)` //任务id
	Expenditure float64 `orm:column(expenditure)` //任务id
	InCome      float64 `orm:column(in_come)`     //任务id
}

func (this *TxIdList) TableName() string {
	return "tx_id_list"
}

func (this *TxIdList) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *TxIdList) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

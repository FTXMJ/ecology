package models

type SuperPeerTable struct {
	Id         int     `orm:"column(id);pk;auto"`
	UserId     string  `orm:column(user_id)`
	CoinNumber float64 `orm:column(coin_number)`
}

func (this *SuperPeerTable) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

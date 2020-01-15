package models

type SuperPeerTable struct {
	Id         int     `orm:"column(id);pk;auto"`
	UserId     string  `orm:"column(user_id)"`
	CoinNumber float64 `orm:"column(coin_number)"`
}

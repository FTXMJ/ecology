package models

type SuperPeerTable struct {
	Id         int     `gorm:"column:id;primary_key" json:"id"`
	UserId     string  `gorm:"column:user_id" json:"user_id"`
	CoinNumber float64 `gorm:"column:coin_number" json:"coin_number"`
}

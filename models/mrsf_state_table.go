package models

type MrsfStateTable struct {
	Id           int     `gorm:"column:id;primary_key" json:"id"`
	UserId       string  `gorm:"column:user_id" json:"user_id"`
	UserName     string  `gorm:"column:user_name" json:"user_name"`
	State        bool    `gorm:"column:state" json:"state"`
	Time         string  `gorm:"column:time" json:"time"`
	OrderId      string  `gorm:"column:order_id" json:"order_id"`
	Date         string  `gorm:"column:date" json:"date"`
	ZiYouABouns  float64 `gorm:"column:ziyou_a_bouns" json:"ziyou_a_bouns"`
	ZhiTuiABouns float64 `gorm:"column:zhitui_a_bouns" json:"zhitui_a_bouns"`
	TeamABouns   float64 `gorm:"column:team_a_bouns" json:"team_a_bouns"`
	PeerABouns   float64 `gorm:"column:peer_a_bouns" json:"peer_a_bouns"`
}

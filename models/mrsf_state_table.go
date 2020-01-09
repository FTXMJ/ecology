package models

type MrsfStateTable struct {
	Id           int     `orm:"column(id);pk;auto" json:"id"`
	UserId       string  `orm:"column(user_id)" json:"user_id"`
	UserName     string  `orm:"column(user_name)" json:"user_name"`
	State        bool    `orm:"column(state)" json:"state"`
	Time         string  `orm:"column(time)" json:"time"`
	OrderId      string  `orm:"column(order_id)" json:"order_id"`
	Date         string  `orm:"column(date)" json:"date"`
	ZiYouABouns  float64 `orm:"column(ziyou_a_bouns)"json:"ziyou_a_bouns"`
	ZhiTuiABouns float64 `orm:"column(zhitui_a_bouns)" json:"zhitui_a_bouns"`
	TeamABouns   float64 `orm:"column(team_a_bouns)" json:"team_a_bouns"`
	PeerABouns   float64 `orm:"column(peer_a_bouns)" json:"peer_a_bouns"`
}

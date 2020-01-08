package models

type MrsfStateTable struct {
	Id       int    `orm:"column(id);pk;auto" json:"id"`
	UserId   string `orm:"column(user_id)" json:"user_id"`
	UserName string `orm:"column(user_name)" json:"user_name"`
	State    bool   `orm:"column(state)" json:"state"`
	Time     string `orm:"column(time)" json:"time"`
	OrderId  string `orm:"column(order_id)" json:"order_id"`
	Date     string `orm:"column(date)" json:"date"`
}

package models

// Forces Table
type ForceTable struct {
	Id                  int     `orm:"column(id);pk;auto" json:"id"`
	Level               string  `orm:"column(level)" json:"level"`
	LowHold             int     `orm:"column(low_hold)" json:"low_hold"`                           //低位
	HighHold            int     `orm:"column(high_hold)" json:"high_hold"`                         //高位
	ReturnMultiple      float64 `orm:"column(return_multiple)" json:"return_multiple"`             //杠杆
	HoldReturnRate      float64 `orm:"column(hold_return_rate)" json:"hold_return_rate"`           //本金自由算力
	RecommendReturnRate float64 `orm:"column(recommend_return_rate)" json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `orm:"column(team_return_rate)" json:"team_return_rate"`           //动态算力
	PictureUrl          string  `orm:"column(picture_url)" json:"picture_url"`
}

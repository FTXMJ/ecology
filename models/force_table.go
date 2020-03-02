package models

// Forces Table
type ForceTable struct {
	Id                  int     `gorm:"column:id;primary_key" json:"id"`
	Level               string  `gorm:"column:level" json:"level"`
	LowHold             int     `gorm:"column:low_hold" json:"low_hold"`                           //低位
	HighHold            int     `gorm:"column:high_hold" json:"high_hold"`                         //高位
	ReturnMultiple      float64 `gorm:"column:return_multiple" json:"return_multiple"`             //杠杆
	HoldReturnRate      float64 `gorm:"column:hold_return_rate" json:"hold_return_rate"`           //本金自由算力
	RecommendReturnRate float64 `gorm:"column:recommend_return_rate" json:"recommend_return_rate"` //直推算力
	TeamReturnRate      float64 `gorm:"column:team_return_rate" json:"team_return_rate"`           //动态算力
	PictureUrl          string  `gorm:"column:picture_url" json:"picture_url"`
}

package models

// 超级节点算力表
type SuperForceTable struct {
	Id               int     `gorm:"column:id;primary_key" json:"id"`
	Level            string  `gorm:"column:level" json:"level"`
	CoinNumberRule   float64 `gorm:"column:coin_number_rule" json:"coin_number_rule"`
	BonusCalculation float64 `gorm:"column:bonus_calculation" json:"bonus_calculation"`
}

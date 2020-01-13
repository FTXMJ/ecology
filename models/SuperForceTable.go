package models

// 超级节点算力表
type SuperForceTable struct {
	Id               int     `orm:"column(id);pk;auto" json:"id"`
	Level            string  `orm:"column(level)" json:"level" json:"level"`
	CoinNumberRule   float64 `orm:"column(coin_number_rule)" json:"coin_number_rule"`
	BonusCalculation float64 `orm:"column(bonus_calculation)" json:"bonus_calculation"`
}

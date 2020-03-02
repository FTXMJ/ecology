package models

//算力表
type CalculationPower struct {
	Id                   int     `gorm:"column:id;primary_key"`
	DateTime             string  `gorm:"column:date_time"`             //日期
	PrincipalCalculation float64 `gorm:"column:principal_calculation"` //本金算力
	DirectCalculation    float64 `gorm:"column:direct_calculation"`    //直推算力
	DynamicCalculation   float64 `gorm:"column:dynamic_calculation"`   //动态算力
}

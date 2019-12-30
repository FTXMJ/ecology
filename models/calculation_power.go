package models

//算力表
type CalculationPower struct {
	Id                   int     `orm:"column(id);pk;auto"`
	DateTime             string  `orm:"column(date_time)"`             //日期
	PrincipalCalculation float64 `orm:"column(principal_calculation)"` //本金算力
	DirectCalculation    float64 `orm:"column(direct_calculation)"`    //直推算力
	DynamicCalculation   float64 `orm:"column(dynamic_calculation)"`   //动态算力
}

func (this *CalculationPower) TableName() string {
	return "calculation_power"
}

func (this *CalculationPower) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *CalculationPower) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

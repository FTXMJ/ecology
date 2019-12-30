package models

//公式表
type Formula struct {
	Id                  int     `orm:"column(id);pk;auto"`
	EcologyId           int     `orm:"column(ecology_id)"`
	Level               string  `orm:"column(level)"`
	LowHold             int     `orm:"column(low_hold)"`              //低位
	HighHold            int     `orm:"column(high_hold)"`             //高位
	ReturnMultiple      float64 `orm:"column(return_multiple)"`       //杠杆
	HoldReturnRate      float64 `orm:"column(hold_return_rate)"`      //本金自由算力
	RecommendReturnRate float64 `orm:"column(recommend_return_rate)"` //直推算力
	TeamReturnRate      float64 `orm:"column(team_return_rate)"`      //动态算力
}

func (this *Formula) TableName() string {
	return "formula"
}

func (this *Formula) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *Formula) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

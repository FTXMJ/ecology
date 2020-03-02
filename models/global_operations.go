package models

// 全局收益控制表
type GlobalOperations struct {
	Id        int    `gorm:"column:id;primary_key"`
	Operation string `gorm:"column:operation" json:"operation"`
	State     bool   `gorm:"column:state" json:"state"`
}

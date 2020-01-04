package models

// 全局收益控制表
type GlobalOperations struct {
	Id        int    `orm:"column(id);pk;auto" json:"id"`
	Operation string `orm:"column(operation)" json:"operation"`
	State     bool   `orm:"column(state)" json:"state"`
}

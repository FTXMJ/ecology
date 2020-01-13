package models

// Page 分页参数  ---  历史信息
type HostryPageInfo struct {
	Items []HostryValues `json:"items"` //数据列表
	Page  Page           `json:"page"`  //分页信息
}

type Page struct {
	TotalPage   int `json:"totalPage"`   //总页数
	CurrentPage int `json:"currentPage"` //当前页数
	PageSize    int `json:"pageSize"`    //每页数据条数
	Count       int `json:"count"`       //总数据量
}

type HostryValues struct {
	Id             int     `json:"id"`
	UserId         string  `json:"user_id"`
	CurrentRevenue float64 `json:"current_revenue"` //上期支出
	CurrentOutlay  float64 `json:"current_outlay"`  //本期支出
	OpeningBalance float64 `json:"opening_balance"` //上期余额
	CurrentBalance float64 `json:"current_balance"` //本期余额
	CreateDate     string  `json:"create_date"`     //创建时间
	Comment        string  `json:"comment"`         //评论
	TxId           string  `json:"tx_id"`           //任务id
	Account        int     `json:"account"`         //生态仓库id
}

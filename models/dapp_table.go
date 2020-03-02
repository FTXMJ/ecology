package models

type DappTable struct {
	Id              int    `gorm:"column:id;primary_key" json:"id"`
	Name            string `gorm:"column:name" json:"name"`                         //  名字
	AgreementType   string `gorm:"column:agreement_type" json:"agreement_type"`     // DAPP类型
	State           bool   `gorm:"column:state" json:"state"`                       // 状态  禁用 -- 启用
	TheLinkAddress  string `gorm:"column:the_link_address" json:"the_link_address"` //链接地址
	ContractAddress string `gorm:"column:contract_address" json:"contract_address"` //链接地址
	Image           string `gorm:"column:image" json:"image"`                       //图片链接
	CreateTime      string `gorm:"column:create_time" json:"create_time"`           //创建时间
	UpdateTime      string `gorm:"column:update_time" json:"update_time"`           //更新时间
}

// user ecology information
type DappList struct {
	Items []DappTable `json:"items"` //数据列表
	Page  Page        `json:"page"`  //分页信息
}

type DappGroupList struct {
	Items []List `json:"items"` //数据列表
}

type List struct {
	Values []DappTable `json:"values"`
	Title  string      `json:"title"`
}

package models

/*{"code":200,"data":[{"t":"1579160167272","s":"TFORUSDT","c":"0.2578","h":"0.258","l":"0.2314","o":"0.2512","v":"1647998.05","qv":"407175.919004","m":"0.0263","e":662}]}*/

type RealTimePrice struct {
	Id        int    `orm:"column(id);pk;auto"`
	TimeStamp string `orm:"column(time_stamp)" json:"time_stamp"`
	Symbol    string `orm:"column(symbol)" json:"symbol"`
	Close     string `orm:"column(close)" json:"close"`
	High      string `orm:"column(high)" json:"high"`
	Low       string `orm:"column(low)" json:"low"`
	Open      string `orm:"column(open)" json:"open"`
	Volume    string `orm:"column(volume)" json:"volume"`
	Quantity  string `orm:"column(quantity)" json:"quantity"`
}

type RealTimePriceHistory struct {
	Id        int    `orm:"column(id);pk;auto"`
	TimeStamp string `orm:"column(time_stamp)" json:"time_stamp"`
	Symbol    string `orm:"column(symbol)" json:"symbol"`
	Close     string `orm:"column(close)" json:"close"`
	High      string `orm:"column(high)" json:"high"`
	Low       string `orm:"column(low)" json:"low"`
	Open      string `orm:"column(open)" json:"open"`
	Volume    string `orm:"column(volume)" json:"volume"`
	Quantity  string `orm:"column(quantity)" json:"quantity"`
}

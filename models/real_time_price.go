package models

/*{"code":200,"data":[{"t":"1579160167272","s":"TFORUSDT","c":"0.2578","h":"0.258","l":"0.2314","o":"0.2512","v":"1647998.05","qv":"407175.919004","m":"0.0263","e":662}]}*/

type QuoteTicker struct {
	Id            int    `gorm:"column:id;primary_key" json:"id"`
	TimeStamp     string `gorm:"column:time_stamp" json:"time_stamp"`
	Symbol        string `gorm:"column:symbol" json:"symbol"`
	Close         string `gorm:"column:close" json:"close"`
	High          string `gorm:"column:high" json:"high"`
	Low           string `gorm:"column:low" json:"low"`
	Open          string `gorm:"column:open" json:"open"`
	Volume        string `gorm:"column:volume" json:"volume"`
	Quantity      string `gorm:"column:quantity" json:"quantity"`
	BaseCurrency  string `gorm:"column:base_currency" json:"base_currency"`
	QuoteCurrency string `gorm:"column:quote_currency" json:"quote_currency"`
}

type QuoteTickerHistory struct {
	Id            int    `gorm:"column:id;primary_key" json:"id"`
	TimeStamp     string `gorm:"column:time_stamp" json:"time_stamp"`
	Symbol        string `gorm:"column:symbol" json:"symbol"`
	Close         string `gorm:"column:close" json:"close"`
	High          string `gorm:"column:high" json:"high"`
	Low           string `gorm:"column:low" json:"low"`
	Open          string `gorm:"column:open" json:"open"`
	Volume        string `gorm:"column:volume" json:"volume"`
	Quantity      string `gorm:"column:quantity" json:"quantity"`
	BaseCurrency  string `gorm:"column:base_currency" json:"base_currency"`
	QuoteCurrency string `gorm:"column:quote_currency" json:"quote_currency"`
	SymbolId      string `gorm:"column:symbol_id" json:"symbol_id"`
}

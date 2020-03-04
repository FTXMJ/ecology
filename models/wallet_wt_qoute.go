package models

import "github.com/jinzhu/gorm"

//
//type WtQuote struct {
//	gorm.Model
//	Code          string  `gorm:"column:code" json:"code"`
//	BaseCurrency  string  `gorm:"column:base_currency" json:"base_currency"`
//	QuoteCurrency string  `gorm:"column:quote_currency" json:"quote_currency"`
//	Price         float64 `gorm:"column:price" json:"price"`
//}

type WtQuote struct {
	gorm.Model
	Code          string  `gorm:"index;column:code;unique;comment:'标示';"`           // 标示
	BaseCurrency  string  `gorm:"index;column:base_currency;comment:'基础货币/大写字母';"`  // 基础货币
	QuoteCurrency string  `gorm:"index;column:quote_currency;comment:'报价货币/大写字母';"` // 基础货币
	Price         float64 `gorm:"column:price;default:0;comment:'金额';"`             // 金额
}

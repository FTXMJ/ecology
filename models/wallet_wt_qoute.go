package models

import (
	"time"
)

type WtQuote struct {
	Id            int       `gorm:"column:id;primary_key" json:"id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Code          string    `gorm:"column:code" json:"code"`
	BaseCurrency  string    `gorm:"column:base_currency" json:"base_currency"`
	QuoteCurrency string    `gorm:"column:quote_currency" json:"quote_currency"`
	Price         float64   `gorm:"column:price" json:"price"`
}

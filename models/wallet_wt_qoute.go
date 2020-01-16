package models

import "time"

type WtQuote struct {
	Id            int       `orm:"column(id);pk;auto"`
	CreatedAt     time.Time `orm:"column(created_at)"`
	UpdatedAt     time.Time `orm:"column(updated_at)"`
	DeletedAt     time.Time `orm:"column(deleted_at)"`
	Code          string    `orm:"column(code)"`
	BaseCurrency  string    `orm:"column(base_currency)"`
	QuoteCurrency string    `orm:"column(quote_currency)"`
	Price         float64   `orm:"column(price)"`
}

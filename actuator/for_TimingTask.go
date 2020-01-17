package actuator

import (
	db "ecology/db"
	"ecology/logs"
	"ecology/models"

	"github.com/astaxie/beego/orm"

	"strconv"
	"time"
)

// 定时获取 交易行情
func Second5s() {
	symbols := []Symbol{}
	s1 := Symbol{
		BaseCurrency:  "TFOR",
		QuoteCurrency: "USDT",
	}
	symbols = append(symbols, s1)
	s2 := Symbol{
		BaseCurrency:  "BTC",
		QuoteCurrency: "USDT",
	}
	symbols = append(symbols, s2)
	s3 := Symbol{
		BaseCurrency:  "ETH",
		QuoteCurrency: "USDT",
	}
	symbols = append(symbols, s3)
	price := 0.0

	for _, v := range symbols {
		p := UpdateOrInsert(v.BaseCurrency, v.QuoteCurrency)
		if p != 0 {
			price = p
		}
	}
	UpdateCoinsPrice(price)
}

func UpdateOrInsert(baseCurrency, quoteCurrency string) (price float64) {
	value, err := GetQuote(baseCurrency, quoteCurrency)

	//   更新 本地 数据
	o_ec := db.NewEcologyOrm()
	o_wa := db.NewWalletOrm()
	o_ec.Begin()
	o_wa.Begin()
	state := "成功"
	symbol := baseCurrency + "-" + quoteCurrency

	err = EcologyH(o_ec, symbol, baseCurrency, quoteCurrency, value)
	if err != nil {
		o_ec.Rollback()
		o_wa.Rollback()
		state = "失败"
		logs.Log.Info("更新"+symbol+" 时间: ", time.Now().Format("2006-01-02 15:04:05")+" 操作: "+state)
		return
	}

	if symbol == "TFOR-USDT" {
		err = WalletH(o_wa, symbol, baseCurrency, quoteCurrency, value)
		if err != nil {
			o_ec.Rollback()
			o_wa.Rollback()
			state = "失败"
			logs.Log.Info("更新"+symbol+" 时间: ", time.Now().Format("2006-01-02 15:04:05")+" 操作: "+state)
			return
		}
		price, _ = strconv.ParseFloat(value.Date[0].C, 64)
	}

	o_ec.Commit()
	o_wa.Commit()
	logs.Log.Info("更新"+symbol+" 时间: ", time.Now().Format("2006-01-02 15:04:05")+" 操作: "+state)

	return
}

func EcologyH(o orm.Ormer, symbol, baseCurrency, quoteCurrency string, value Data_r) error {
	var err error
	r_t_p := models.QuoteTicker{
		TimeStamp:     value.Date[0].T,
		Symbol:        symbol,
		Close:         value.Date[0].C,
		High:          value.Date[0].H,
		Low:           value.Date[0].L,
		Open:          value.Date[0].O,
		Volume:        value.Date[0].V,
		Quantity:      value.Date[0].Qv,
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
	}
	if count, _ := o.QueryTable("quote_ticker").Filter("symbol", symbol).Count(); count == 0 {
		_, err = o.Insert(&r_t_p)
		if err != nil {
			return err
		}
	} else {
		_, err = o.Raw("update quote_ticker set close=?,high=?,low=?,open=?,volume=?,quantity=?,time_stamp=? where symbol=?",
			r_t_p.Close,
			r_t_p.High,
			r_t_p.Low,
			r_t_p.Open,
			r_t_p.Volume,
			r_t_p.Quantity,
			r_t_p.TimeStamp,
			r_t_p.Symbol,
		).Exec()
		if err != nil {
			return err
		}
	}
	r_h := models.QuoteTickerHistory{
		Symbol:        symbol,
		Close:         value.Date[0].C,
		High:          value.Date[0].H,
		Low:           value.Date[0].L,
		Open:          value.Date[0].O,
		Volume:        value.Date[0].V,
		Quantity:      value.Date[0].Qv,
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
		SymbolId:      symbol + "-" + value.Date[0].T,
	}
	t, _ := strconv.Atoi(value.Date[0].T)
	r_h.TimeStamp = time.Unix(int64(t)/1000, 0).Format("2006-01-02 15:04:05")
	_, err = o.Insert(&r_h)
	if err != nil {
		return err
	}
	return nil
}

func WalletH(o orm.Ormer, symbol, baseCurrency, quoteCurrency string, value Data_r) error {
	var err error
	price, _ := strconv.ParseFloat(value.Date[0].C, 64)
	a := []models.WtQuote{}
	b := models.WtQuote{}
	o.Raw("select * from wt_quote where code=?", symbol).QueryRow(&b)
	o.Raw("select * from wt_quote", symbol).QueryRows(&a)
	if count, _ := o.QueryTable("wt_quote").Filter("code", symbol).Count(); count == 0 {
		w_q := models.WtQuote{
			CreatedAt:     time.Now(),
			Code:          symbol,
			BaseCurrency:  baseCurrency,
			QuoteCurrency: quoteCurrency,
			Price:         price,
		}
		_, err = o.Insert(&w_q)
		if err != nil {
			return err
		}
	} else {
		_, err = o.Raw("update wt_quote set price=? , updated_at=? where code=?", price, time.Now(), symbol).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateCoinsPrice(price float64) {
	o := db.NewWalletOrm()
	w_q := []models.WtQuote{}
	o.Raw("select * from wt_quote").QueryRows(&w_q)
	for _, v := range w_q {
		pp := price
		switch v.Code {
		case "USDD-TROF":
			o.Raw("update wt_quote set price=? where id=?", float64(1)/pp, v.Id).Exec()
		case "TFOR-USDD":
			o.Raw("update wt_quote set price=? where id=?", price, v.Id).Exec()
		case "USDT-TFOR":
			o.Raw("update wt_quote set price=? where id=?", float64(1)/pp, v.Id).Exec()
		}
	}
}

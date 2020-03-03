package actuator

import (
	db "ecology/db"
	"ecology/logs"
	"ecology/models"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"

	"strconv"
	"time"
)

// 定时获取 交易行情
func Second5s() {
	symbols := make([]Symbol, 0)
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
	var price float64 = 0

	//   更新 本地 数据
	if timing_orm.UpdateTime+3600 <= time.Now().Unix() || timing_orm.UpdateTime < 1 || timing_orm.EcologyConn == nil || timing_orm.WalletConn == nil {
		UpdateTimingOrm()
	}

	for _, v := range symbols {
		p := UpdateOrInsert(timing_orm.EcologyConn, timing_orm.WalletConn, v.BaseCurrency, v.QuoteCurrency)
		if p != 0 {
			price = p
		}
	}
	if price > 0 {
		UpdateCoinsPrice(timing_orm.WalletConn, price)
	}

}

type TimingOrm struct {
	EcologyConn *gorm.DB
	WalletConn  *gorm.DB
	UpdateTime  int64
}

var timing_orm TimingOrm

func UpdateTimingOrm() {
	o_ec := db.NewEcologyOrm()
	o_wa := db.NewWalletOrm()
	u_t := time.Now().Unix()
	timing_orm.UpdateTime = u_t
	timing_orm.EcologyConn = o_ec
	timing_orm.WalletConn = o_wa
}

func UpdateOrInsert(o_ec, o_wa *gorm.DB, baseCurrency, quoteCurrency string) (price float64) {
	value, err := GetQuote(baseCurrency, quoteCurrency)
	if value.Code == 0 || len(value.Date) == 0 {
		return 0
	}

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

func EcologyH(o *gorm.DB, symbol, baseCurrency, quoteCurrency string, value Data_r) error {
	var err error
	r_t_p := models.QuoteTicker{
		//TimeStamp:     value.Date[0].T,
		TimeStamp:     strconv.Itoa(value.Date[0].T),
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
	quote_ticker := []models.QuoteTicker{}
	if o.Table("quote_ticker").Where("symbol = ?", symbol).Find(&quote_ticker); len(quote_ticker) == 0 {
		er := o.Create(&r_t_p)
		if er.Error != nil {
			return err
		}
	} else {
		er := o.Model(&models.QuoteTicker{}).Where("symbol = ?", r_t_p.Symbol).Update(map[string]interface{}{
			"close":      r_t_p.Close,
			"high":       r_t_p.High,
			"low":        r_t_p.Low,
			"open":       r_t_p.Open,
			"volume":     r_t_p.Volume,
			"quantity":   r_t_p.Quantity,
			"time_stamp": r_t_p.TimeStamp,
		})
		if er.Error != nil {
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
		SymbolId:      symbol + "-" + strconv.Itoa(value.Date[0].T),
	}
	t := value.Date[0].T
	r_h.TimeStamp = time.Unix(int64(t)/1000, 0).Format("2006-01-02 15:04:05")
	er := o.Create(&r_h)
	if er.Error != nil {
		if err.Error() == "Error 1062: Duplicate entry '"+r_h.SymbolId+"' for key 'symbol_id'" {
			o.Update(&r_h)
		} else {
			return err
		}
	}
	return nil
}

func WalletH(o *gorm.DB, symbol, baseCurrency, quoteCurrency string, value Data_r) error {
	var err error
	price, _ := strconv.ParseFloat(value.Date[0].C, 64)
	a := make([]models.WtQuote, 0)
	b := models.WtQuote{}
	o.Raw("select * from wt_quote where code=?", symbol).First(&b)
	o.Raw("select * from wt_quote", symbol).Find(&a)
	wt := []models.WtQuote{}
	if o.Table("wt_quote").Where("code = ?", symbol).Find(&wt); len(wt) == 0 {
		w_q := models.WtQuote{
			CreatedAt:     time.Now(),
			Code:          symbol,
			BaseCurrency:  baseCurrency,
			QuoteCurrency: quoteCurrency,
			Price:         price,
		}
		er := o.Create(&w_q)
		if er.Error != nil {
			return err
		}
	} else {
		er := o.Model(&models.QuoteTicker{}).Where("code = ?", symbol).Update(map[string]interface{}{
			"price":      price,
			"updated_at": time.Now(),
		})
		if er.Error != nil {
			return err
		}
	}
	return nil
}

func UpdateCoinsPrice(o *gorm.DB, price float64) {
	w_q := make([]models.WtQuote, 0)
	var count float64 = 1
	o.Raw("select * from wt_quote").Find(&w_q)
	if len(w_q) > 0 {
		items := make([]models.WtQuote, 0)
		p := div(count, price)
		for _, v := range w_q {
			switch v.Code {
			case "USDD-TFOR":
				v.Price, _ = p.Float64()
			case "TFOR-USDD":
				v.Price = price
			case "USDT-TFOR":
				v.Price, _ = p.Float64()
			default:
				continue
			}
			items = append(items, v)
		}
		if len(items) > 0 {
			for _, v := range items {
				o.Save(&v)
			}
		}
	}
}

// 除法
func div(d1, d2 float64) decimal.Decimal {
	d11 := decimal.NewFromFloat(d1)
	d22 := decimal.NewFromFloat(d2)
	return d11.Div(d22)
}

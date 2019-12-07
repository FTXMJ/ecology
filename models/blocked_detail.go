package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

//铸币表
type BlockedDetail struct {
	Id             int       `orm:"column(id);pk;auto"`
	UserId         string    `orm:column(user_id)`
	CurrentRevenue float64   `orm:column(current_revenue)` //上期支出
	CurrentOutlay  float64   `orm:column(current_outlay)`  //本期支出
	OpeningBalance float64   `orm:column(opening_balance)` //上期余额
	CurrentBalance float64   `orm:column(current_balance)` //本期余额
	CreateDate     time.Time `orm:column(create_date)`     //创建时间
	Comment        string    `orm:column(comment)`         //评论
	TxId           string    `orm:column(tx_id)`           //任务id
	Account        int       `orm:column(account)`         //生态仓库id
}

func (this *BlockedDetail) TableName() string {
	return "blocked_detail"
}

func (this *BlockedDetail) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *BlockedDetail) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

func FindLimitOneAndSaveBlo_d(o orm.Ormer,user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	fmt.Println("account_id :",account_id)
	blocked_old := BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
		OrderBy("-create_date").
		Limit(1).
		One(&blocked_old)
	if blocked_old.Id == 0 {
		blocked_old.CurrentBalance = 0
	}
	for_mula := Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id",account_id).One(&for_mula)
	if err_for!=nil{
		return err_for
	}
	blocked_new := BlockedDetail{
		Id:             0,
		UserId:         user_id,
		CurrentRevenue: coin_in,
		CurrentOutlay:  coin_out,
		OpeningBalance: blocked_old.CurrentBalance,
		CurrentBalance: blocked_old.CurrentBalance + coin_in - coin_out,
		CreateDate:     time.Now(),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
	}
	blocked_new.CurrentBalance = blocked_old.CurrentBalance+coin_in*for_mula.ReturnMultiple-coin_out
	_,err:=o.Insert(&blocked_new)
	if err != nil {
		return err
	}

	_,err_txid := o.QueryTable("tx_id_list").Filter("tx_id",tx_id).Update(orm.Params{"state":"true"})
	if err_txid!=nil{
		return err_txid
	}

	_,err_up :=o.QueryTable("account").Filter("id",account_id).Update(orm.Params{"bocked_balance":blocked_new.CurrentBalance})
	if err_up!=nil{
		return err_up
	}
	return nil
}

//　把所有算力的值加起来
func AddFormula(coin_in,ganggan,ziyou,jiasu,dongtai float64) float64 {
	va := coin_in *ganggan
	va += coin_in * ziyou
	va += coin_in * jiasu
	va += coin_in * dongtai
	return va
}
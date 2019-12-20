package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//铸币表
type BlockedDetail struct {
	Id             int     `orm:"column(id);pk;auto" json:"id"`
	UserId         string  `orm:"column(user_id)" json:"user_id"`
	CurrentRevenue float64 `orm:"column(current_revenue)" json:"current_revenue"` //上期支出
	CurrentOutlay  float64 `orm:"column(current_outlay)" json:"current_outlay"`   //本期支出
	OpeningBalance float64 `orm:"column(opening_balance)" json:"opening_balance"` //上期余额
	CurrentBalance float64 `orm:"column(current_balance)" json:"current_balance"` //本期余额
	CreateDate     string  `orm:"column(create_date)" json:"create_date"`         //创建时间
	Comment        string  `orm:"column(comment)" json:"comment"`                 //评论
	TxId           string  `orm:"column(tx_id)" json:"tx_id"`                     //任务id
	Account        int     `orm:"column(account)" json:"account"`                 //生态仓库id
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

func FindLimitOneAndSaveBlo_d(o orm.Ormer, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	blocked_old := BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
		Filter("account", account_id).
		OrderBy("-create_date").
		Limit(1).
		One(&blocked_old)
	if blocked_old.Id == 0 {
		blocked_old.CurrentBalance = 0
	}
	for_mula := Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id", account_id).One(&for_mula)
	if err_for != nil {
		o.Rollback()
		return err_for
	}
	blocked_new := BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: coin_in * for_mula.ReturnMultiple,
		CurrentOutlay:  coin_out,
		OpeningBalance: blocked_old.CurrentBalance,
		CurrentBalance: blocked_old.CurrentBalance + coin_in*for_mula.ReturnMultiple - coin_out,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
	}
	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		o.Rollback()
		return err
	}

	// 更新任务完成状态
	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"state": "true"})
	if err_txid != nil {
		o.Rollback()
		return err_txid
	}

	//更新生态仓库属性
	_, err_up := o.QueryTable("account").Filter("id", account_id).Update(orm.Params{"bocked_balance": blocked_new.CurrentBalance})
	if err_up != nil {
		o.Rollback()
		return err_up
	}

	// 超级节点表生成与更新
	super_peer_table := SuperPeerTable{}
	err_super := o.QueryTable("super_peer_table").Filter("user_id", user_id).One(&super_peer_table)
	if err_super != nil {
		o.Rollback()
		return err_super
	}
	coin := super_peer_table.CoinNumber + (coin_in * for_mula.ReturnMultiple) - coin_out
	if coin < 0 {
		coin = 0
	}
	_, err_super_up := o.QueryTable("super_peer_table").Filter("user_id", user_id).Update(orm.Params{"coin_number": coin})
	if err_super_up != nil {
		o.Rollback()
		return err_super_up
	}

	//  直推收益
	user := User{}
	erruser := o.QueryTable("user").Filter("user_id", user_id).One(&user)
	if erruser != nil {
		o.Rollback()
		return erruser
	}
	errrr := ForAddCoin(o, user.FatherId, coin_in, 0.1)
	if errrr != nil {
		o.Rollback()
		return nil
	}

	return nil
}

func NewCreateAndSaveBlo_d(o orm.Ormer, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	for_mula := Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id", account_id).One(&for_mula)
	if err_for != nil {
		o.Rollback()
		return err_for
	}
	blocked_new := BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: coin_in * for_mula.ReturnMultiple,
		CurrentOutlay:  coin_out,
		OpeningBalance: 0,
		CurrentBalance: coin_in * for_mula.ReturnMultiple,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		o.Rollback()
		return err
	}

	// 更新任务完成状态
	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"state": "true"})
	if err_txid != nil {
		o.Rollback()
		return err_txid
	}

	//更新生态仓库属性
	_, err_up := o.QueryTable("account").Filter("id", account_id).Update(orm.Params{"bocked_balance": blocked_new.CurrentBalance})
	if err_up != nil {
		o.Rollback()
		return err_up
	}

	// 超级节点表生成与更新
	super_peer_table := SuperPeerTable{}
	err_super := o.QueryTable("super_peer_table").Filter("user_id", user_id).One(&super_peer_table)
	if err_super != nil {
		o.Rollback()
		return err_super
	}
	coin := super_peer_table.CoinNumber + (coin_in * for_mula.ReturnMultiple) - coin_out
	if coin < 0 {
		coin = 0
	}
	_, err_super_up := o.QueryTable("super_peer_table").Filter("user_id", user_id).Update(orm.Params{"coin_number": coin})
	if err_super_up != nil {
		o.Rollback()
		return err_super_up
	}

	//  直推收益
	user := User{}
	erruser := o.QueryTable("user").Filter("user_id", user_id).One(&user)
	if erruser != nil {
		o.Rollback()
		return erruser
	}
	if user.FatherId != "" {
		errrr := ForAddCoin(o, user.FatherId, coin_in, 0.1)
		if errrr != nil {
			o.Rollback()
			return errrr
		}
	}
	return nil
}

//　把所有算力的值加起来
func ForAddCoin(o orm.Ormer, father_id string, coin float64, proportion float64) error {
	user := User{}
	if err := o.QueryTable("user").Filter("user_id", father_id).One(&user); err != nil {
		return err
	}
	account := Account{}
	erraccount := o.QueryTable("account").Filter("user_id", father_id).One(&account)
	if erraccount != nil {
		return erraccount
	}
	new_coin := account.BockedBalance + (coin * proportion)
	_, err_up := o.QueryTable("account").Filter("user_id", father_id).Update(orm.Params{"bocked_balance": new_coin})
	if err_up != nil {
		return err_up
	}

	if (coin * proportion * 0.1) > 1 {
		ForAddCoin(o, user.FatherId, (coin * proportion), proportion*0.1)
	}
	return nil
}

// 查询所用对象
type FindObj struct {
	UserId    string
	TxId      string
	StartTime string
	EndTime   string
}

// Page 分页参数  ---  历史信息
type HostryFindInfo struct {
	Items []BlockedDetail `json:"items"` //数据列表
	Page  Page            `json:"page"`  //分页信息
}

/*
	条件查询
	对象包含的则视为条件
*/
func SelectPondMachinemsg(p FindObj, page Page, table_name string) ([]BlockedDetail, Page, error) {
	var list []BlockedDetail
	level := ""
	var err error
	s_ql := "select * from " + table_name + " where "
	if p.UserId != "" {
		level += "1"
	}
	if p.TxId != "" {
		level += "2"
	}
	if p.StartTime != "" && p.EndTime != "" {
		level += "3"
	}
	if level == "1" {
		s_ql = s_ql + "user_id=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.UserId).QueryRows(&list)
		err = er
	} else if level == "12" {
		s_ql = s_ql + "user_id=? and tx_id=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.UserId, p.TxId).QueryRows(&list)
		err = er
	} else if level == "123" {
		s_ql = s_ql + "user_id=? and tx_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.UserId, p.TxId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "13" {
		s_ql = s_ql + "user_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.UserId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "23" {
		s_ql = s_ql + "tx_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.TxId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "3" {
		s_ql = s_ql + "create_date > ? and create_date < ? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	}

	if err != nil {
		return []BlockedDetail{}, Page{}, err
	}

	page.Count = len(list)
	if page.PageSize < 10 {
		page.PageSize = 10
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	listle := []BlockedDetail{}
	if end > len(list) {
		for _, v := range list[start:] {
			listle = append(listle, v)
		}
	} else {
		for _, v := range list[start:end] {
			listle = append(listle, v)
		}
	}
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	return list, page, nil
}

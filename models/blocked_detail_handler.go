package models

import (
	"ecology/logs"
	"ecology/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

// 更新借贷表
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
	NetIncome += coin_in * for_mula.ReturnMultiple
	if user.FatherId != "" {
		ForAddCoin(o, user.FatherId, coin_in, 0.1)
	}
	return nil
}

// 创建第一条借贷记录
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

//　把所有算力的值加起来  -- 更新静态｀动态的验证
func ForAddCoin(o orm.Ormer, father_id string, coin float64, proportion float64) error {
	user := User{}
	if err_user := o.QueryTable("user").Filter("user_id", father_id).One(&user); err_user != nil {
		return err_user
	}
	account := Account{}
	erraccount := o.QueryTable("account").Filter("user_id", father_id).One(&account)
	if erraccount != nil {
		return erraccount
	}
	if account.DynamicRevenue == true {
		new_coin := account.BockedBalance + (coin * proportion)
		_, err_up := o.QueryTable("account").Filter("user_id", father_id).Update(orm.Params{"bocked_balance": new_coin})
		if err_up != nil {
			return err_up
		}
		//任务表 USDD  铸币记录
		tx_id_blo_d := utils.Shengchengstr("直推收益", father_id, "USDD")
		blo_txid_dcmt := TxIdList{
			TxId:        tx_id_blo_d,
			State:       "true",
			UserId:      father_id,
			CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
			Expenditure: 0,
			InCome:      (coin * proportion),
		}
		_, errtxid_blo := o.Insert(&blo_txid_dcmt)
		if errtxid_blo != nil {
			logs.Log.Error("直推算力累加错误", errtxid_blo)
			return errtxid_blo
		}
		blocked_old := BlockedDetail{}
		o.QueryTable("blocked_detail").
			Filter("user_id", father_id).
			Filter("account", account.Id).
			OrderBy("-create_date").
			Limit(1).
			One(&blocked_old)
		if blocked_old.Id == 0 {
			blocked_old.CurrentBalance = 0
		}

		blocked_new := BlockedDetail{
			UserId:         father_id,
			CurrentRevenue: (coin * proportion),
			CurrentOutlay:  0,
			OpeningBalance: blocked_old.CurrentBalance,
			CurrentBalance: blocked_old.CurrentBalance + (coin * proportion),
			CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
			Comment:        "直推收益",
			TxId:           tx_id_blo_d,
			Account:        account.Id,
		}
		if blocked_new.CurrentBalance < 0 {
			blocked_new.CurrentBalance = 0
		}
		_, err := o.Insert(&blocked_new)
		if err != nil {
			return err
		}
	}

	if coin*proportion > 1 {
		ForAddCoin(o, user.FatherId, (coin * proportion), proportion)
	}
	return nil
}

/*
	条件查询
	对象包含的则视为条件
*/
func SelectPondMachinemsg(p FindObj, page Page, table_name string) ([]BlockedDetail, Page, error) {
	list, err := SqlCreateValues1(p, table_name)
	if err != nil {
		return []BlockedDetail{}, Page{}, err
	}
	page.Count = len(list)
	if page.PageSize < 5 {
		page.PageSize = 5
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
	return listle, page, nil
}

// 释放流水查询－－处理
func SelectFlows(p FindObj, page Page, table_name string) ([]Flow, Page, error) {
	var list []BlockedDetail
	level := ""
	var err error
	s_ql := "select * from " + table_name + " where comment=? "
	if p.UserId != "" {
		level += "1"
	}
	if p.StartTime != "" && p.EndTime != "" {
		level += "2"
	}
	if level == "1" {
		s_ql = s_ql + "and user_id=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, "每日释放", p.UserId).QueryRows(&list)
		err = er
	} else if level == "12" {
		s_ql = s_ql + "and user_id=? and create_date>=? and create_date<=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, "每日释放", p.UserId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "2" {
		s_ql = s_ql + "and create_date>=? and create_date<=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, "每日释放", p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else {
		s_ql = s_ql + "and id>0 order by create_date desc"
		_, er := NewOrm().Raw(s_ql, "每日释放", p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	}
	if err != nil {
		return []Flow{}, Page{}, err
	}
	page.Count = len(list)
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	if end > len(list) {
		end = len(list)
	}
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	flows := []Flow{}
	for _, v := range list[start:end] {
		flow := Flow{}
		t, _ := time.Parse("2006-01-02 15:04:05", v.CreateDate)
		time_start := t.AddDate(0, 0, -1).Format("2006-01-02") + " 00:00:00"
		time_end := t.AddDate(0, 0, -1).Format("2006-01-02") + " 23:59:59"
		//　直推算力
		zhitui, err_zt := RecommendReturnRateEveryDay(v.UserId, time_start, time_end)
		if err_zt != nil {
			return []Flow{}, Page{}, err_zt
		}
		//　算力
		formula := Formula{
			EcologyId: v.Account,
		}
		o := NewOrm()
		err_read := o.Read(&formula, "ecology_id")
		if err_read != nil {
			return []Flow{}, Page{}, err_read
		}
		// 	已近释放
		blo := []BlockedDetail{}
		_, err_raw := o.Raw("select * from blocked_detail where user_id=? and comment=? and create_date<=? ", v.UserId, "每日释放", time_start).QueryRows(&blo)
		if err_raw != nil {
			return []Flow{}, Page{}, err_raw
		}
		coin := 0.0
		for _, v := range blo {
			coin += v.CurrentOutlay
		}
		flow.UserId = v.UserId
		flow.HoldReturnRate = formula.HoldReturnRate * v.CurrentBalance // 本金自由算力
		flow.RecommendReturnRate = zhitui
		flow.TeamReturnRate = formula.TeamReturnRate
		flow.Released = coin
		flow.UpdateTime = v.CreateDate
		flows = append(flows, flow)
	}

	return flows, page, nil
}

// 直推算力的计算　　　－－　　　当天
func RecommendReturnRate(user_id, time string) (float64, error) {
	blo := []BlockedDetail{}
	sql_str := "SELECT * from blocked_detail where user_id=? and create_date>=? and comment=? "
	_, err := NewOrm().Raw(sql_str, user_id, time, "直推收益").QueryRows(&blo)
	if err != nil {
		return 0, err
	}
	zhitui := 0.0
	for _, v := range blo {
		zhitui += v.CurrentRevenue
	}
	zhit, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", zhitui), 64)
	return zhit, nil
}

// 直推算力的计算　　　－－　　　任意天
func RecommendReturnRateEveryDay(user_id, time_start, time_end string) (float64, error) {
	blo := []BlockedDetail{}
	sql_str := "SELECT * from blocked_detail where user_id=? and create_date>=? and create_date<=? and comment=? "
	_, err := NewOrm().Raw(sql_str, user_id, time_start, time_end, "直推收益").QueryRows(&blo)
	if err != nil {
		return 0, err
	}
	zhitui := 0.0
	for _, v := range blo {
		zhitui += v.CurrentRevenue
	}
	return zhitui, nil
}

// Find user ecology information
func FindU_E_OBJ(page Page, user_id string) ([]U_E_OBJ, Page) {
	o := NewOrm()
	users := []User{}
	if user_id != "" {
		o.Raw("select * from user where user_id=? ", user_id).QueryRows(&users)
	} else {
		o.Raw("select * from user order by id").QueryRows(&users)
	}
	user_e_objs := []U_E_OBJ{}
	for _, v := range users {
		user_e_obj := U_E_OBJ{}
		account := Account{}
		formula := Formula{}
		blos := []BlockedDetail{}
		o.Raw("select * from account where user_id=? ", v.UserId).QueryRow(&account)
		o.Raw("select * from formula where ecology_id=? ", account.Id).QueryRow(&formula)
		user_e_obj.UserId = v.UserId
		user_e_obj.Level = account.Level
		user_e_obj.ReturnMultiple = formula.ReturnMultiple
		user_e_obj.CoinAll = account.Balance
		user_e_obj.ToBeReleased = account.BockedBalance
		o.Raw("select * from blocked_detail where user_id=? and comment=?", v.UserId, "每日释放").QueryRows(&blos)
		zhichu := 0.0
		for _, v := range blos {
			zhichu += v.CurrentOutlay
		}
		user_e_obj.Released = zhichu
		user_e_obj.HoldReturnRate = formula.HoldReturnRate * account.Balance
		zhitui, _ := RecommendReturnRate(v.UserId, time.Now().Format("2006-01-02")+" 00:00:00")
		user_e_obj.RecommendReturnRate = zhitui
		user_e_obj.TeamReturnRate = formula.TeamReturnRate

		user_e_objs = append(user_e_objs, user_e_obj)
	}
	page.Count = len(user_e_objs)
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	if end > len(user_e_objs) {
		return user_e_objs[start:], page
	} else {
		return user_e_objs[start:end], page
	}
}

// Find user ecology information
func FindUserAccountOFF(page Page, obj FindObj) ([]AccountOFF, Page, error) {
	accounts, err := SqlCreateValues2(obj, "account")
	if err != nil {
		return []AccountOFF{}, Page{}, err
	}
	user_accounts := []AccountOFF{}
	for _, v := range accounts {
		user_account := AccountOFF{
			UserId:         v.UserId,
			Account:        v.Id,
			DynamicRevenue: v.DynamicRevenue,
			StaticReturn:   v.StaticReturn,
			CreateDate:     v.CreateDate,
		}
		user_accounts = append(user_accounts, user_account)
	}
	page.Count = len(user_accounts)
	if page.PageSize < 5 {
		page.PageSize = 5
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	if end > len(user_accounts) {
		return user_accounts[start:], page, nil
	} else {
		return user_accounts[start:end], page, nil
	}
}

///*
//	条件查询
//	对象包含的则视为条件     sql 生成并　查询
//*/
func SqlCreateValues1(p FindObj, table_name string) ([]BlockedDetail, error) {
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
	} else {
		s_ql = s_ql + "id>0 order by create_date desc"
		_, er := NewOrm().Raw(s_ql, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	}
	if err != nil {
		return []BlockedDetail{}, err
	}
	return list, nil
}

//   ---   Find user ecology information　　　sql 生成并　查询
func SqlCreateValues2(obj FindObj, table_name string) ([]Account, error) {
	var list []Account
	level := ""
	var err error
	s_ql := "select * from " + table_name + " where "
	if obj.UserId != "" {
		level += "1"
	}
	if obj.TxId != "" {
		level += "2"
	}
	if obj.StartTime != "" && obj.EndTime != "" {
		level += "3"
	}
	if level == "1" {
		s_ql = s_ql + "user_id=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.UserId).QueryRows(&list)
		err = er
	} else if level == "12" {
		s_ql = s_ql + "user_id=? and id=? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.UserId, obj.TxId).QueryRows(&list)
		err = er
	} else if level == "123" {
		s_ql = s_ql + "user_id=? and id=? and create_date>? and create_date<? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.UserId, obj.TxId, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else if level == "13" {
		s_ql = s_ql + "user_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.UserId, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else if level == "23" {
		s_ql = s_ql + "id=? and create_date>? and create_date<? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.TxId, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else if level == "3" {
		s_ql = s_ql + "create_date > ? and create_date < ? order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else {
		s_ql = s_ql + "id>0 order by create_date desc"
		_, er := NewOrm().Raw(s_ql, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	}
	if err != nil {
		return []Account{}, err
	}
	return list, nil
}

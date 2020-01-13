package actuator

import (
	"ecology/logs"
	"ecology/models"
	"ecology/utils"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

// 更新借贷表
func FindLimitOneAndSaveBlo_d(o orm.Ormer, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	blocked_olds := []models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
		Filter("account", account_id).
		OrderBy("-create_date").
		Limit(3).
		All(&blocked_olds)
	var blocked_old models.BlockedDetail
	if len(blocked_olds) != 0 {
		for i := 0; i < len(blocked_olds)-1; i++ {
			for j := i + 1; j < len(blocked_olds); j++ {
				if blocked_olds[i].Id > blocked_olds[j].Id {
					blocked_olds[i], blocked_olds[j] = blocked_olds[j], blocked_olds[i]
				}
			}
		}

		blocked_old = blocked_olds[len(blocked_olds)-1]
		if blocked_old.Id == 0 {
			blocked_old.CurrentBalance = 0
		}
	}
	for_mula := models.Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id", account_id).One(&for_mula)
	if err_for != nil {
		o.Rollback()
		return err_for
	}
	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: coin_in * for_mula.ReturnMultiple,
		CurrentOutlay:  coin_out,
		OpeningBalance: blocked_old.CurrentBalance,
		CurrentBalance: blocked_old.CurrentBalance + coin_in*for_mula.ReturnMultiple - coin_out,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
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
	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"order_state": true})
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

	//  直推收益
	user := models.User{}
	erruser := o.QueryTable("user").Filter("user_id", user_id).One(&user)
	if erruser != nil {
		o.Rollback()
		return erruser
	}
	if user.FatherId != "" && coin_in >= 10 {
		ForAddCoin(o, user.FatherId, coin_in, 0.1)
	}
	return nil
}

// 创建第一条借贷记录
func NewCreateAndSaveBlo_d(o orm.Ormer, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	for_mula := models.Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id", account_id).One(&for_mula)
	if err_for != nil {
		o.Rollback()
		return err_for
	}
	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: coin_in * for_mula.ReturnMultiple,
		CurrentOutlay:  coin_out,
		OpeningBalance: 0,
		CurrentBalance: coin_in * for_mula.ReturnMultiple,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		o.Rollback()
		return err
	}

	// 更新任务完成状态
	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"order_state": true})
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

	//  直推收益
	user := models.User{}
	erruser := o.QueryTable("user").Filter("user_id", user_id).One(&user)
	if erruser != nil {
		o.Rollback()
		return erruser
	}
	if user.FatherId != "" && coin_in > 10 {
		errrr := ForAddCoin(o, user.FatherId, coin_in, 0.1)
		if errrr != nil {
			o.Rollback()
			return errrr
		}
	}
	o.Commit()
	return nil
}

//　把所有算力的值加起来  -- 更新静态｀动态的验证          bo 是是否调过　钱包的增加接口
func ForAddCoin(o orm.Ormer, father_id string, coin float64, proportion float64) error {
	user := models.User{}
	o.QueryTable("user").Filter("user_id", father_id).One(&user)

	account := models.Account{}
	o.QueryTable("account").Filter("user_id", father_id).One(&account)
	//任务表 USDD  铸币记录
	order_id := utils.TimeUUID()
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  true,
		WalletState: true,
		UserId:      father_id,
		Comment:     "直推收益",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: (coin * proportion),
		InCome:      0,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		logs.Log.Error("直推算力累加错误", errtxid_blo)
		return errtxid_blo
	}
	blocked_olds := []models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user.UserId).
		Filter("account", account.Id).
		OrderBy("-create_date").
		Limit(3).
		All(&blocked_olds)
	var blocked_old models.BlockedDetail
	if len(blocked_olds) != 0 {
		for i := 0; i < len(blocked_olds)-1; i++ {
			for j := i + 1; j < len(blocked_olds); j++ {
				if blocked_olds[i].Id > blocked_olds[j].Id {
					blocked_olds[i], blocked_olds[j] = blocked_olds[j], blocked_olds[i]
				}
			}
		}

		blocked_old = blocked_olds[len(blocked_olds)-1]
		if blocked_old.Id == 0 {
			blocked_old.CurrentBalance = 0
		}
	}

	blocked_new := models.BlockedDetail{
		UserId:         father_id,
		CurrentRevenue: 0,
		CurrentOutlay:  (coin * proportion),
		OpeningBalance: blocked_old.CurrentBalance,
		CurrentBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "直推收益",
		TxId:           order_id,
		Account:        account.Id,
		CoinType:       "USDD",
	}
	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		ForAddCoin(o, father_id, coin, proportion)
	}
	if coin*proportion*proportion >= 1 && user.FatherId != "" {
		ForAddCoin(o, user.FatherId, (coin * proportion), proportion)
	}
	return nil
}

/*
	条件查询
	对象包含的则视为条件
*/
func SelectPondMachinemsg(o orm.Ormer, p models.FindObj, page models.Page, table_name string) ([]models.BlockedDetailIndex, models.Page, error) {
	list, err := SqlCreateValues1(p, table_name)
	if err != nil {
		return []models.BlockedDetailIndex{}, models.Page{}, err
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
	listle := []models.BlockedDetail{}
	if end > len(list) && start < len(list) {
		for _, v := range list[start:] {
			listle = append(listle, v)
		}
	} else if start > len(list) {

	} else if end < len(list) && start < len(list) {
		for _, v := range list[start:end] {
			listle = append(listle, v)
		}
	}
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	lists := []models.BlockedDetailIndex{}
	for _, v := range listle {
		var u models.User
		u.UserId = v.UserId
		o.Read(&u, "user_id")
		blo := models.BlockedDetailIndex{
			Id:             v.Id,
			UserId:         v.UserId,
			UserName:       u.UserName,
			CurrentRevenue: v.CurrentRevenue,
			CurrentOutlay:  v.CurrentOutlay,
			OpeningBalance: v.OpeningBalance,
			CurrentBalance: v.CurrentBalance,
			CreateDate:     v.CreateDate,
			Comment:        v.Comment,
			TxId:           v.TxId,
			Account:        v.Account,
			CoinType:       v.CoinType,
		}
		lists = append(lists, blo)
	}
	return lists, page, nil
}

// 释放流水查询－－处理
func SelectFlows(o orm.Ormer, p models.FindObj, page models.Page, table_name string) ([]models.Flow, models.Page, error) {
	us := []models.User{}
	blos := []models.BlockedDetail{}
	if p.UserName != "" && p.UserId == "" {
		o.Raw("select * from user where user_name=?", p.UserName).QueryRows(&us)
		if len(us) == 0 {
			return []models.Flow{}, models.Page{}, errors.New("没有符合条件的用户!")
		}
	}

	if p.UserName != "" && p.UserId != "" {
		o.Raw("select * from user where user_name=? and user_id=?", p.UserName, p.UserId).QueryRows(&us)
		if len(us) == 0 {
			return []models.Flow{}, models.Page{}, errors.New("没有符合条件的用户!")
		}
	}

	if p.UserName == "" && p.UserId != "" {
		o.Raw("select * from user where user_id=?", p.UserId).QueryRows(&us)
		if len(us) == 0 {
			return []models.Flow{}, models.Page{}, errors.New("没有符合条件的用户!")
		}
	}

	if p.StartTime != "" && p.EndTime != "" && p.UserName == "" && p.UserId == "" {
		list := []models.BlockedDetail{}
		s_ql := "select * from " + table_name + " where comment=? and create_date>=? and create_date<=? order by create_date desc"
		o.Raw(s_ql, "每日释放收益", p.StartTime, p.EndTime).QueryRows(&list)
		//                                                                                  处理分页
		page.Count = len(list)
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
		listle := []models.BlockedDetail{}
		if end > len(list) && start < len(list) {
			for _, v := range list[start:] {
				listle = append(listle, v)
			}
		} else if start > len(list) {

		} else if end < len(list) && start < len(list) {
			for _, v := range list[start:end] {
				listle = append(listle, v)
			}
		}
		//                                                                                  拼接数据
		flows := []models.Flow{}
		for _, v := range listle {
			flow := models.Flow{}
			zhitui := models.BlockedDetail{}
			o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日直推收益").QueryRow(&zhitui)
			tuandui := models.BlockedDetail{}
			o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日团队收益").QueryRow(&tuandui)

			var u models.User
			u.UserId = v.UserId
			o.Read(&u, "user_id")
			flow.UserId = v.UserId
			flow.UserName = u.UserName
			flow.HoldReturnRate = v.CurrentBalance // 本金自由算力
			flow.RecommendReturnRate = zhitui.CurrentOutlay
			flow.TeamReturnRate = tuandui.CurrentOutlay
			flow.Released = zhitui.CurrentOutlay + tuandui.CurrentOutlay + v.CurrentOutlay
			flow.UpdateTime = v.CreateDate
			flows = append(flows, flow)
		}
		return flows, page, nil
	}

	if p.StartTime == "" && p.EndTime == "" && p.UserName == "" && p.UserId == "" {
		list := []models.BlockedDetail{}
		s_ql := "select * from " + table_name + " where comment=? order by create_date desc"
		o.Raw(s_ql, "每日释放收益").QueryRows(&list)
		//                                                                                  处理分页
		page.Count = len(list)
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
		listle := []models.BlockedDetail{}
		if end > len(list) && start < len(list) {
			for _, v := range list[start:] {
				listle = append(listle, v)
			}
		} else if start > len(list) {

		} else if end < len(list) && start < len(list) {
			for _, v := range list[start:end] {
				listle = append(listle, v)
			}
		}
		//                                                                                  拼接数据
		flows := []models.Flow{}
		for _, v := range listle {
			flow := models.Flow{}
			zhitui := models.BlockedDetail{}
			o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日直推收益").QueryRow(&zhitui)
			tuandui := models.BlockedDetail{}
			o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日团队收益").QueryRow(&tuandui)

			var u models.User
			u.UserId = v.UserId
			o.Read(&u, "user_id")
			flow.UserId = v.UserId
			flow.UserName = u.UserName
			flow.HoldReturnRate = v.CurrentOutlay // 本金自由算力
			flow.RecommendReturnRate = zhitui.CurrentOutlay
			flow.TeamReturnRate = tuandui.CurrentOutlay
			flow.Released = zhitui.CurrentOutlay + tuandui.CurrentOutlay + v.CurrentOutlay
			flow.UpdateTime = v.CreateDate
			flows = append(flows, flow)
		}
		return flows, page, nil
	}

	for _, v := range us {
		if p.StartTime != "" && p.EndTime != "" {
			list := []models.BlockedDetail{}
			o.Raw("select * from "+table_name+" where user_id=? and comment=? and create_date>=? and create_date<=? order by create_date desc", v.UserId, "每日释放收益", p.StartTime, p.EndTime).QueryRows(&list)
			for _, v := range list {
				blos = append(blos, v)
			}
		} else {
			list := []models.BlockedDetail{}
			o.Raw("select * from "+table_name+" where user_id=? and comment=? order by create_date desc", v.UserId, "每日释放收益").QueryRows(&list)
			for _, v := range list {
				blos = append(blos, v)
			}
		}
	}
	if len(blos) > 1 {
		QuickSortBlockedDetail(blos, 0, len(blos)-1)
	}
	//                                                                                  处理分页
	page.Count = len(blos)
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
	//                                                                                  拼接数据
	listle := []models.BlockedDetail{}
	if end > len(blos) && start < len(blos) {
		for _, v := range blos[start:] {
			listle = append(listle, v)
		}
	} else if start > len(blos) {

	} else if end < len(blos) && start < len(blos) {
		for _, v := range blos[start:end] {
			listle = append(listle, v)
		}
	}
	flows := []models.Flow{}
	for _, v := range listle {
		flow := models.Flow{}
		zhitui := models.BlockedDetail{}
		o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日直推收益").QueryRow(&zhitui)
		tuandui := models.BlockedDetail{}
		o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日团队收益").QueryRow(&tuandui)

		var u models.User
		u.UserId = v.UserId
		o.Read(&u, "user_id")
		flow.UserId = v.UserId
		flow.UserName = u.UserName
		hold, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.CurrentOutlay), 64)
		flow.HoldReturnRate = hold // 本金自由算力
		reco, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", zhitui.CurrentOutlay), 64)
		flow.RecommendReturnRate = reco
		team, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", tuandui.CurrentOutlay), 64)
		flow.TeamReturnRate = team
		rele, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", zhitui.CurrentOutlay+tuandui.CurrentOutlay+v.CurrentOutlay), 64)
		flow.Released = rele
		flow.UpdateTime = v.CreateDate
		flows = append(flows, flow)
	}
	return flows, page, nil
}

// 直推算力的计算　　　－－　　　当天
func RecommendReturnRate(user_id, time string) (float64, error) {
	blo := orm.ParamsList{}
	sql_str := "SELECT sum(current_outlay) from blocked_detail where user_id=? and create_date>=? and comment=? "
	_, err := models.NewOrm().Raw(sql_str, user_id, time, "直推收益").ValuesFlat(&blo)
	if err != nil {
		return 0, err
	}
	var zhitui float64
	if len(blo) > 0 && blo[0] != nil {
		z, _ := strconv.ParseFloat(blo[0].(string), 64)
		zhitui = z
	}
	zhit, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", zhitui), 64)
	return zhit, nil
}

// 直推算力的计算　　　－－　　　任意天
func RecommendReturnRateEveryDay(user_id, time_start, time_end string) (float64, error) {
	blo := orm.ParamsList{}
	sql_str := "SELECT sum(current_outlay) from blocked_detail where user_id=? and create_date>=? and create_date<=? and comment=? "
	_, err := models.NewOrm().Raw(sql_str, user_id, time_start, time_end, "直推收益").ValuesFlat(&blo)
	if err != nil {
		return 0, err
	}
	var zhitui float64
	if len(blo) > 0 && blo[0] != nil {
		z, _ := strconv.ParseFloat(blo[0].(string), 64)
		zhitui = z
	}
	return zhitui, nil
}

// Find user ecology information
func FindU_E_OBJ(o orm.Ormer, page models.Page, user_id, user_name string) ([]models.U_E_OBJ, models.Page) {
	users := []models.User{}
	if user_id != "" && user_name == "" {
		o.Raw("select * from user where user_id=? order by id", user_id).QueryRows(&users)
	} else if user_id != "" && user_name != "" {
		o.Raw("select * from user where user_name=? and user_id=? order by id", user_name, user_id).QueryRows(&users)
	} else if user_id == "" && user_name != "" {
		o.Raw("select * from user where user_name=? order by id", user_name).QueryRows(&users)
	} else {
		o.Raw("select * from user order by id").QueryRows(&users)
	}
	user_e_objs := []models.U_E_OBJ{}
	for _, v := range users {
		user_e_obj := models.U_E_OBJ{}
		account := models.Account{}
		formula := models.Formula{}
		blos := orm.ParamsList{}
		o.Raw("select * from account where user_id=? ", v.UserId).QueryRow(&account)
		o.Raw("select * from formula where ecology_id=? ", account.Id).QueryRow(&formula)
		user_e_obj.UserId = v.UserId
		user_e_obj.UserName = v.UserName
		user_e_obj.Level = account.Level
		user_e_obj.ReturnMultiple = formula.ReturnMultiple
		user_e_obj.CoinAll = account.Balance
		user_e_obj.ToBeReleased = account.BockedBalance
		o.Raw("select sum(current_outlay) from blocked_detail where user_id=? and comment=?", v.UserId, "每日释放").ValuesFlat(&blos)
		var zhichu float64
		if len(blos) > 0 && blos[0] != nil {
			z, _ := strconv.ParseFloat(blos[0].(string), 64)
			zhichu = z
		}
		user_e_obj.Released = zhichu
		user_e_obj.HoldReturnRate = formula.HoldReturnRate * account.Balance
		zhitui, _ := RecommendReturnRate(v.UserId, time.Now().Format("2006-01-02")+" 00:00:00")
		user_e_obj.RecommendReturnRate = zhitui
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
	list := []models.U_E_OBJ{}
	if end > len(user_e_objs) && start < len(user_e_objs) {
		list = user_e_objs[start:]
	} else if start > len(user_e_objs) {
		list = []models.U_E_OBJ{}
	} else if end < len(user_e_objs) && start < len(user_e_objs) {
		list = user_e_objs[start:end]
	}
	if len(list) == 0 {
		return []models.U_E_OBJ{}, page
	}
	for i := 0; i < len(list); i++ {
		team, _ := IndexTeamABouns(o, list[i].UserId)
		acc := models.Account{UserId: list[i].UserId}
		o.Read(&acc, "user_id")
		for_m := models.Formula{EcologyId: acc.Id}
		o.Read(&for_m, "ecology_id")
		list[i].TeamReturnRate = team * for_m.TeamReturnRate
	}
	u_e_objs := []models.U_E_OBJ{}
	for _, v := range list {
		u_e_obj := models.U_E_OBJ{}
		hold, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.HoldReturnRate), 64)
		u_e_obj.HoldReturnRate = hold
		reco, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.RecommendReturnRate), 64)
		u_e_obj.RecommendReturnRate = reco
		rele, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.Released), 64)
		u_e_obj.Released = rele
		team, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.TeamReturnRate), 64)
		u_e_obj.TeamReturnRate = team
		tobe, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.ToBeReleased), 64)
		u_e_obj.ToBeReleased = tobe
		coin, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.CoinAll), 64)
		u_e_obj.CoinAll = coin
		u_e_obj.Level = v.Level
		u_e_obj.UserId = v.UserId
		u_e_obj.UserName = v.UserName
		u_e_obj.ReturnMultiple = v.ReturnMultiple
		u_e_objs = append(u_e_objs, u_e_obj)
	}
	return u_e_objs, page
}

// Find user ecology information
func FindFalseUser(o orm.Ormer, page models.Page, user_id, user_name string) ([]models.FalseUser, models.Page) {
	users := []models.User{}
	if user_id != "" && user_name == "" {
		o.Raw("select * from user where user_id=? order by id", user_id).QueryRows(&users)
	} else if user_id != "" && user_name != "" {
		o.Raw("select * from user where user_name=? and user_id=? order by id", user_name, user_id).QueryRows(&users)
	} else if user_id == "" && user_name != "" {
		o.Raw("select * from user where user_name=? order by id", user_name).QueryRows(&users)
	} else {
		o.Raw("select * from user order by id").QueryRows(&users)
	}
	f_u_s := []models.FalseUser{}
	for _, v := range users {
		account := models.Account{}
		f_u := models.FalseUser{}
		o.Raw("select * from account where user_id=? and (dynamic_revenue=? or static_return=?)", v.UserId, false, false).QueryRow(&account)
		if account.Id > 0 {
			f_u.UserName = v.UserName
			f_u.UserId = v.UserId
			f_u.UpdateTime = account.UpdateDate
			f_u.Dongtai = account.DynamicRevenue
			f_u.Jintai = account.StaticReturn
			f_u.AccountId = account.Id
			f_u_s = append(f_u_s, f_u)
		}
	}
	page.Count = len(f_u_s)
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

	if end > len(f_u_s) && start < len(f_u_s) {
		return f_u_s[start:], page
	} else if start > len(f_u_s) {
		return []models.FalseUser{}, page
	} else if end < len(f_u_s) && start < len(f_u_s) {
		return f_u_s[start:end], page
	}
	return []models.FalseUser{}, page
}

// Find user ecology information
func FindUserAccountOFF(o orm.Ormer, page models.Page, obj models.FindObj) ([]models.AccountOFF, models.Page, error) {
	accounts, err := SqlCreateValues2(o, obj, "account")
	if err != nil {
		return []models.AccountOFF{}, models.Page{}, err
	}
	user_accounts := []models.AccountOFF{}
	g := []models.GlobalOperations{}
	o.Raw("select * from global_operations").QueryRows(&g)
	m := make(map[string]bool)
	for _, v := range g {
		m[v.Operation] = v.State
	}
	for _, v := range accounts {
		user_account := models.AccountOFF{
			UserId:     v.UserId,
			Account:    v.Id,
			CreateDate: v.CreateDate,
		}
		var dynamic_revenue bool = v.DynamicRevenue
		var static_return bool = v.StaticReturn
		var peer_state bool = v.PeerState
		if m["全局动态收益控制"] == false {
			dynamic_revenue = false
		}
		if m["全局静态收益控制"] == false {
			static_return = false
		}
		if m["全局节点分红控制"] == false {
			peer_state = false
		}
		user_account.DynamicRevenue = dynamic_revenue
		user_account.StaticReturn = static_return
		user_account.PeerState = peer_state
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
		for i := start; i < len(user_accounts); i++ {
			var u models.User
			u.UserId = user_accounts[i].UserId
			o.Read(&u, "user_id")
			user_accounts[i].UserName = u.UserName
		}
		return user_accounts[start:], page, nil
	} else {
		for i := start; i < len(user_accounts); i++ {
			var u models.User
			u.UserId = user_accounts[i].UserId
			o.Read(&u, "user_id")
			user_accounts[i].UserName = u.UserName
		}
		return user_accounts[start:end], page, nil
	}
	if end > len(user_accounts) && start < len(user_accounts) {
		for i := start; i < len(user_accounts); i++ {
			var u models.User
			u.UserId = user_accounts[i].UserId
			o.Read(&u, "user_id")
			user_accounts[i].UserName = u.UserName
		}
		return user_accounts[start:], page, nil
	} else if start > len(user_accounts) {
		return []models.AccountOFF{}, page, nil
	} else {
		for i := start; i < len(user_accounts); i++ {
			var u models.User
			u.UserId = user_accounts[i].UserId
			o.Read(&u, "user_id")
			user_accounts[i].UserName = u.UserName
		}
		return user_accounts[start:end], page, nil
	}
}

///*
//	条件查询
//	对象包含的则视为条件     sql 生成并　查询
//*/
func SqlCreateValues1(p models.FindObj, table_name string) ([]models.BlockedDetail, error) {
	var list []models.BlockedDetail
	o := models.NewOrm()
	level := ""
	name := ""
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
	if p.UserName != "" {
		name = "4"
	}
	if level == "1" {
		s_ql = s_ql + "user_id=? order by create_date desc"
		_, er := o.Raw(s_ql, p.UserId).QueryRows(&list)
		err = er
	} else if level == "12" {
		s_ql = s_ql + "user_id=? and tx_id=? order by create_date desc"
		_, er := o.Raw(s_ql, p.UserId, p.TxId).QueryRows(&list)
		err = er
	} else if level == "123" {
		s_ql = s_ql + "user_id=? and tx_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := o.Raw(s_ql, p.UserId, p.TxId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "13" {
		s_ql = s_ql + "user_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := o.Raw(s_ql, p.UserId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "23" {
		s_ql = s_ql + "tx_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := o.Raw(s_ql, p.TxId, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else if level == "3" {
		s_ql = s_ql + "create_date > ? and create_date < ? order by create_date desc"
		_, er := o.Raw(s_ql, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	} else {
		s_ql = s_ql + "id>0 order by create_date desc"
		_, er := o.Raw(s_ql, p.StartTime, p.EndTime).QueryRows(&list)
		err = er
	}
	if err != nil {
		return []models.BlockedDetail{}, err
	}
	list_last := []models.BlockedDetail{}
	if name != "" {
		for _, v := range list {
			u := models.User{}
			o.Raw("select * from user where user_id=?", v.UserId).QueryRow(&u)
			if u.UserName == p.UserName {
				list_last = append(list_last, v)
			}
		}
		return list_last, nil
	}
	return list, nil
}

//   ---   Find user ecology information　　　sql 生成并　查询
func SqlCreateValues2(o orm.Ormer, obj models.FindObj, table_name string) ([]models.Account, error) {
	var list []models.Account
	level := ""
	name := ""
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
	if obj.UserName != "" {
		name = "4"
	}
	if level == "1" {
		s_ql = s_ql + "user_id=? order by create_date desc"
		_, er := o.Raw(s_ql, obj.UserId).QueryRows(&list)
		err = er
	} else if level == "12" {
		s_ql = s_ql + "user_id=? and id=? order by create_date desc"
		_, er := o.Raw(s_ql, obj.UserId, obj.TxId).QueryRows(&list)
		err = er
	} else if level == "123" {
		s_ql = s_ql + "user_id=? and id=? and create_date>? and create_date<? order by create_date desc"
		_, er := o.Raw(s_ql, obj.UserId, obj.TxId, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else if level == "13" {
		s_ql = s_ql + "user_id=? and create_date>? and create_date<? order by create_date desc"
		_, er := o.Raw(s_ql, obj.UserId, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else if level == "23" {
		s_ql = s_ql + "id=? and create_date>? and create_date<? order by create_date desc"
		_, er := o.Raw(s_ql, obj.TxId, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else if level == "3" {
		s_ql = s_ql + "create_date > ? and create_date < ? order by create_date desc"
		_, er := o.Raw(s_ql, obj.StartTime, obj.EndTime).QueryRows(&list)
		err = er
	} else {
		s_ql = s_ql + "id>0 order by create_date desc"
		_, er := o.Raw(s_ql).QueryRows(&list)
		err = er
	}
	if err != nil {
		return []models.Account{}, err
	}
	list_last := []models.Account{}
	if name != "" {
		for _, v := range list {
			u := models.User{}
			o.Raw("select * from user where user_id=?", v.UserId).QueryRow(&u)
			if u.UserName == obj.UserName {
				list_last = append(list_last, v)
			}
		}
		return list_last, nil
	}

	return list, nil
}

func ShowMrsfTable(o orm.Ormer, page models.Page, user_name, user_id, date string, state bool) ([]models.MrsfStateTable, models.Page, error) {
	list := []models.MrsfStateTable{}
	acc := models.Account{UserId: user_id}
	o.Read(&acc, "user_id")
	level := ""
	var err error
	s_ql := "select * from mrsf_state_table where "
	if user_id != "" {
		level += "1"
	}
	if user_name != "" {
		level += "2"
	}
	if date != "" {
		level += "3"
	}
	if level == "1" {
		s_ql = s_ql + "user_id=? and state=? order by time desc"
		_, er := o.Raw(s_ql, user_id, state).QueryRows(&list)
		err = er
	} else if level == "12" {
		s_ql = s_ql + "user_id=? and user_name=? and state=? order by time desc"
		_, er := o.Raw(s_ql, user_id, user_name, state).QueryRows(&list)
		err = er
	} else if level == "123" {
		s_ql = s_ql + "user_id=? and user_name=? and order_id=? and state=? and date=? order by time desc"
		_, er := o.Raw(s_ql, user_id, user_name, strconv.Itoa(acc.Id)+date, state, date).QueryRows(&list)
		err = er
	} else if level == "13" {
		s_ql = s_ql + "user_id=? and order_id=? and state=? and date=? order by time desc"
		_, er := o.Raw(s_ql, user_id, strconv.Itoa(acc.Id)+date, state, date).QueryRows(&list)
		err = er
	} else if level == "23" {
		s_ql = s_ql + "user_name=? and state=? and date=? order by time desc"
		_, er := o.Raw(s_ql, user_name, state, date).QueryRows(&list)
		err = er
	} else if level == "3" {
		s_ql = s_ql + "state=? and date=? order by time desc"
		_, er := o.Raw(s_ql, state, date).QueryRows(&list)
		err = er
	} else if level == "2" {
		s_ql = s_ql + "state=? and user_name=? order by time desc"
		_, er := o.Raw(s_ql, state, user_name).QueryRows(&list)
		err = er
	} else {
		s_ql = s_ql + "id>0 and state=? order by time desc"
		_, er := o.Raw(s_ql, state).QueryRows(&list)
		err = er
	}
	if err != nil {
		return []models.MrsfStateTable{}, page, err
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
	page.TotalPage = (page.Count / page.PageSize) + 1 //总页数
	if page.Count <= 5 {
		page.CurrentPage = 1
	}
	if end > len(list) && start < len(list) {
		return list[start:], page, nil
	} else if start > len(list) {
		return []models.MrsfStateTable{}, page, nil
	} else if end <= len(list) && start <= len(list) {
		return list[start:end], page, nil
	}
	return []models.MrsfStateTable{}, page, nil
}

// 查看节点收益流水
func SelectPeerABounsList(o orm.Ormer, page models.Page, user_name, start_time, end_time string) ([]models.PeerAbouns, models.Page, error) {
	peer_a_bouns := []models.TxIdList{}
	time := ""
	if start_time != "" {
		time = "and create_time>=" + "'" + start_time + "'" + " and create_time<=" + "'" + end_time + "'"
	}
	switch user_name {
	case "":
		o.Raw("select * from tx_id_list where comment=? "+time+" order by create_time desc", "节点分红").QueryRows(&peer_a_bouns)
	default:
		users := []models.User{}
		_, err_1 := o.Raw("select * from user where user_name=?", user_name).QueryRows(&users)
		if err_1 != nil || len(users) < 1 {
			return []models.PeerAbouns{}, page, err_1
		}
		for _, v := range users {
			_, err := o.Raw("select * from tx_id_list where user_id=? and comment=? "+time+" order by create_time desc", v.UserId, "节点分红").QueryRows(&peer_a_bouns)
			if err != nil || len(peer_a_bouns) < 1 {
				return []models.PeerAbouns{}, page, err
			}
		}
	}
	if len(peer_a_bouns) < 1 {
		return []models.PeerAbouns{}, page, errors.New("没有相关数据!")
	}
	models.QuickSortPeerABouns(peer_a_bouns, 0, len(peer_a_bouns)-1)

	page.Count = len(peer_a_bouns)
	if page.PageSize < 10 {
		page.PageSize = 10
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize

	listle := []models.PeerAbouns{}

	if start > len(peer_a_bouns) {
		return []models.PeerAbouns{}, page, nil
	} else if end > len(peer_a_bouns) {
		end = len(peer_a_bouns)
	}
	if start == 0 && end == 0 {
		return []models.PeerAbouns{}, page, nil
	}
	if len(peer_a_bouns[start:end]) > 0 {
		for _, v := range peer_a_bouns[start:end] {
			u := models.User{
				UserId: v.UserId,
			}
			o.Read(&u, "user_id")
			_, level, _, err_tfor := ReturnSuperPeerLevel(v.UserId)
			if err_tfor != nil {
				return []models.PeerAbouns{}, page, err_tfor
			}
			p := models.PeerAbouns{
				Id:       v.Id,
				UserName: u.UserName,
				Level:    level,
				Tfors:    v.Expenditure,
				Time:     v.CreateTime,
			}
			listle = append(listle, p)
		}
	}
	return listle, page, nil
}

// 查看用户有史以来所有的收益
func AddAllSum(o orm.Ormer, user_id string) float64 {
	var blos orm.ParamsList
	//o.Raw("select * from blocked_detail where user_id=? and comment!=?",user_id,"直推收益").QueryRows(&blos)
	o.Raw("select sum(current_outlay) from blocked_detail where user_id=? and comment!=?", user_id, "直推收益").ValuesFlat(&blos)
	var zhitui float64
	if len(blos) > 0 && blos[0] != nil {
		z, _ := strconv.ParseFloat(blos[0].(string), 64)
		zhitui = z
	}
	return zhitui
}

// 处理器，计算所有用户的收益  并发布任务和 分红记录
func HandlerOperation(users []string, user_id string) (float64, error) {
	o := models.NewOrm()
	var coin_abouns float64
	for _, v := range users {
		// 拿到生态项目实例
		account := models.Account{}
		err_acc := o.QueryTable("account").Filter("user_id", v).One(&account)
		if err_acc != nil {
			if err_acc.Error() != "<QuerySeter> no row found" {
				return 0, err_acc
			}
		}
		// 拿到生态项目对应的算力表
		formula := models.Formula{}
		err_for := o.QueryTable("formula").Filter("ecology_id", account.Id).One(&formula)
		if err_for != nil {
			if err_for.Error() != "<QuerySeter> no row found" {
				return 0, err_for
			}
		}
		coin_abouns += formula.HoldReturnRate * account.Balance
	}
	return coin_abouns, nil
}

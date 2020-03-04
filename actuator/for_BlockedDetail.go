package actuator

import (
	"ecology/logs"
	"ecology/models"
	"ecology/utils"
	"github.com/jinzhu/gorm"

	"errors"
	"fmt"
	"strconv"
	"time"
)

// 更新借贷表
func FindLimitOneAndSaveBlo_d(o *gorm.DB, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account_id).First(&blocked_old)
	for_mula := models.Formula{}
	o.Table("formula").Where("ecology_id = ?", account_id).First(&for_mula)

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
	err := o.Create(&blocked_new)
	if err.Error != nil {
		return err.Error
	}

	// 更新任务完成状态
	err_txid := o.Model(&models.TxIdList{}).Where("tx_id = ?", tx_id).Update("order_state", true)
	if err_txid.Error != nil {
		return err_txid.Error
	}

	//更新生态仓库属性
	err_up := o.Model(&models.Account{}).Where("id = ?", account_id).Update("bocked_balance", blocked_new.CurrentBalance)
	if err_up.Error != nil {
		return err_up.Error
	}

	//  直推收益
	user := models.User{}
	o.Table("user").Where("user_id = ?", user_id).First(&user)
	if user.FatherId != "" && coin_in >= 10 {
		ForAddCoin(o, user.FatherId, coin_in, 0.1)
	}
	return nil
}

// 创建第一条借贷记录
func NewCreateAndSaveBlo_d(o *gorm.DB, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	for_mula := models.Formula{}
	o.Table("formula").Where("ecology_id = ?", account_id).First(&for_mula)

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
	err := o.Create(&blocked_new)
	if err.Error != nil {
		return err.Error
	}

	// 更新任务完成状态
	err_txid := o.Model(&models.TxIdList{}).Where("tx_id = ?", tx_id).Update("order_state", true)
	if err_txid.Error != nil {
		return err_txid.Error
	}

	//更新生态仓库属性
	err_up := o.Model(&models.Account{}).Where("id = ?", account_id).Update("bocked_balance", blocked_new.CurrentBalance)
	if err_up.Error != nil {
		return err_up.Error
	}

	//  直推收益
	user := models.User{UserId: user_id}
	o.Table("user").Where("user_id = ?", user_id).First(&user)
	if user.FatherId != "" && coin_in > 10 {
		errrr := ForAddCoin(o, user.FatherId, coin_in, 0.1)
		if errrr != nil {
			return errrr
		}
	}
	return nil
}

//　把所有算力的值加起来  -- 更新静态｀动态的验证          bo 是是否调过　钱包的增加接口
func ForAddCoin(o *gorm.DB, father_id string, coin float64, proportion float64) error {
	user := models.User{}
	o.Table("user").Where("user_id = ?", father_id).First(&user)

	account := models.Account{}
	o.Table("account").Where("user_id = ?", father_id).First(&account)

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
	errtxid_blo := o.Create(&blo_txid_dcmt)
	if errtxid_blo.Error != nil {
		logs.Log.Error("直推算力累加错误", errtxid_blo.Error)
		return errtxid_blo.Error
	}

	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user.UserId, account.Id).First(&blocked_old)

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
	err := o.Create(&blocked_new)
	if err.Error != nil {
		ForAddCoin(o, father_id, coin, proportion)
	} else if coin*proportion*proportion >= 1 && user.FatherId != "" {
		ForAddCoin(o, user.FatherId, (coin * proportion), proportion)
	}
	return nil
}

//条件查询对象包含的则视为条件
func SelectPondMachinemsg(o *gorm.DB, p models.FindObj, page models.Page) ([]models.BlockedDetailIndex, models.Page, error) {
	list, err := SqlCreateValues(o, p, "account_detail")
	if err != nil {
		return []models.BlockedDetailIndex{}, models.Page{}, err
	}

	start, end := InitPage(&page, len(list))

	listle := ListLimit(list, start, end)

	lists := []models.BlockedDetailIndex{}
	for _, v := range listle {
		value, _ := v.(models.AccountDetail)
		var u models.User
		o.Table("user").Where("user_id = ?", value.UserId).First(&u)
		blo := models.BlockedDetailIndex{
			Id:                value.Id,
			UserId:            value.UserId,
			UserName:          u.UserName,
			AccCurrentRevenue: value.CurrentRevenue,
			CreateDate:        value.CreateDate,
			Comment:           value.Comment,
			TxId:              value.TxId,
			Account:           value.Account,
			CoinType:          value.CoinType,
		}
		lists = append(lists, blo)
	}
	for i := 0; i < len(lists); i++ {
		var blo models.BlockedDetail
		o.Table("blocked_detail").Where("tx_id = ?", lists[i].TxId).First(&blo)
		lists[i].BloCurrentRevenue = blo.CurrentRevenue
		rm := blo.CurrentRevenue / lists[i].AccCurrentRevenue
		lists[i].ReturnMultiple = int(rm)
	}
	return lists, page, nil
}

// 释放流水查询－－处理
func SelectFlows(o *gorm.DB, p models.FindObj, page models.Page, table_name string) ([]models.Flow, models.Page, error) {
	us := []models.User{}
	blos := []models.BlockedDetail{}

	q_user := o.Table("user")

	if p.UserName != "" {
		q_user = q_user.Where("user_name = ?", p.UserName)
	}
	if p.UserId != "" {
		q_user = q_user.Where("user_id = ?", p.UserId)
	}
	q_user.Find(&us)
	user_ids := []string{}
	for _, v := range us {
		user_ids = append(user_ids, v.UserId)
	}

	q_blos := o.Table(table_name).Where("user_id in (?)", user_ids)

	if p.StartTime != "" && p.EndTime != "" {
		q_blos = q_blos.Where("create_date >= ?", p.StartTime).Where("create_date <= ?", p.EndTime)
	}
	q_blos.Order("create_date desc").Find(&blos)

	QuickSortBlockedDetail(blos, 0, len(blos)-1)

	start, end := InitPage(&page, len(blos))

	value_list := []interface{}{}
	for _, v := range blos {
		value_list = append(value_list, v)
	}

	list := ListLimit(value_list, start, end)

	flows := make([]models.Flow, 0)
	for _, value := range list {
		v := value.(models.BlockedDetail)
		flow := models.Flow{}
		zhitui := models.BlockedDetail{}
		o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日直推收益").First(&zhitui)
		tuandui := models.BlockedDetail{}
		o.Raw("select * from "+table_name+" where create_date=? and user_id=? and comment=?", v.CreateDate, v.UserId, "每日团队收益").First(&tuandui)

		var u models.User
		o.Table("user").Where("user_id = ?", v.UserId).First(&u)
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

// Find user ecology information
func FindU_E_OBJ(o *gorm.DB, page models.Page, user_id, user_name string) ([]models.U_E_OBJ, models.Page) {
	users := []models.User{}

	q_user := o.Table("user")
	if user_name != "" {
		q_user = q_user.Where("user_name = ?", user_name)
	}
	if user_id != "" {
		q_user = q_user.Where("user_id = ?", user_id)
	}
	q_user.Find(&users)
	user_ids := []string{}
	for _, v := range users {
		user_ids = append(user_ids, v.UserId)
	}

	user_e_objs := []models.U_E_OBJ{}
	for _, v := range users {
		user_e_obj := models.U_E_OBJ{}
		account := models.Account{}
		formula := models.Formula{}
		blos := []models.BlockedDetail{}

		o.Raw("select * from account where user_id=? ", v.UserId).First(&account)
		o.Raw("select * from formula where ecology_id=? ", account.Id).First(&formula)

		user_e_obj.UserId = v.UserId
		user_e_obj.UserName = v.UserName
		user_e_obj.Level = account.Level
		user_e_obj.ReturnMultiple = formula.ReturnMultiple
		user_e_obj.CoinAll = account.Balance
		user_e_obj.ToBeReleased = account.BockedBalance

		o.Raw("select * from blocked_detail where user_id=? and comment=?", v.UserId, "每日释放").Find(&blos)

		var zhichu float64
		if len(blos) > 0 {
			for _, v := range blos {
				zhichu += v.CurrentOutlay
			}
		}
		user_e_obj.Released = zhichu
		user_e_obj.HoldReturnRate = formula.HoldReturnRate * account.Balance
		zhitui, _ := RecommendReturnRate(v.UserId, time.Now().Format("2006-01-02")+" 00:00:00")
		user_e_obj.RecommendReturnRate = zhitui
		user_e_objs = append(user_e_objs, user_e_obj)
	}

	start, end := InitPage(&page, len(user_e_objs))
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
		acc := models.Account{}
		o.Table("account").Where("user_id = ?", list[i].UserId).First(&acc)

		for_m := models.Formula{}
		o.Table("account").Where("ecology_id = ?", acc.Id).First(&for_m)

		list[i].TeamReturnRate = team * for_m.TeamReturnRate
	}
	u_e_objs := make([]models.U_E_OBJ, 0)
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
func FindFalseUser(o *gorm.DB, page models.Page, user_id, user_name string) ([]models.FalseUser, models.Page) {
	users := []models.User{}

	q_user := o.Table("user")
	if user_name != "" {
		q_user = q_user.Where("user_name = ?", user_name)
	}
	if user_id != "" {
		q_user = q_user.Where("user_id = ?", user_id)
	}
	q_user.Find(&users)

	f_u_s := []interface{}{}
	for _, v := range users {
		account := models.Account{}
		f_u := models.FalseUser{}
		o.Raw("select * from account where user_id=? and (dynamic_revenue=? or static_return=?)", v.UserId, false, false).First(&account)
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

	start, end := InitPage(&page, len(f_u_s))
	list := []interface{}{}
	if list = ListLimit(f_u_s, start, end); len(list) == 0 {
		return []models.FalseUser{}, page
	}

	last_list := make([]models.FalseUser, 0)
	for _, v := range list {
		last_list = append(last_list, v.(models.FalseUser))
	}

	return last_list, page
}

// Find user ecology information
func FindUserAccountOFF(o *gorm.DB, page models.Page, obj models.FindObj) ([]models.AccountOFF, models.Page, error) {
	accounts, err := SqlCreateValues(o, obj, "account")
	if err != nil {
		return []models.AccountOFF{}, models.Page{}, err
	}

	user_accounts := make([]models.AccountOFF, 0)
	g := []models.GlobalOperations{}
	o.Raw("select * from global_operations").Find(&g)
	m := make(map[string]bool)
	for _, v := range g {
		m[v.Operation] = v.State
	}
	for _, velue := range accounts {
		v := velue.(models.Account)
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

	start, end := InitPage(&page, len(user_accounts))

	if end > len(user_accounts) {
		for i := start; i < len(user_accounts); i++ {
			var u models.User
			o.Table("user").Where("user_id = ?", user_accounts[i].UserId).First(&u)
			user_accounts[i].UserName = u.UserName
		}
		return user_accounts[start:], page, nil
	} else {
		for i := start; i < len(user_accounts); i++ {
			var u models.User
			o.Table("user").Where("user_id = ?", user_accounts[i].UserId).First(&u)
			user_accounts[i].UserName = u.UserName
		}
		return user_accounts[start:end], page, nil
	}
}

//	对象包含的则视为条件     sql 生成并　查询
func SqlCreateValues(o *gorm.DB, p models.FindObj, table_name string) ([]interface{}, error) {
	// 根据条件  进行数据插叙
	if list, err := GeneratedSQLAndExec(o, table_name, p); err != nil {
		return []interface{}{}, err
	} else {
		return list, nil
	}
}

func ShowMrsfTable(o *gorm.DB, page models.Page, user_name, user_id, date string, state bool) ([]models.MrsfStateTable, models.Page, error) {
	list := make([]models.MrsfStateTable, 0)
	acc := models.Account{}
	o.Table("account").Where("user_id = ?", user_id).First(&acc)

	q := o.Table("mrsf_state_table")
	if user_name != "" {
		q = q.Where("user_name = ?", user_name)
	}
	if user_id != "" {
		q = q.Where("user_id = ?", user_id)
	}
	if date != "" {
		q = q.Where("date = ?", date)
	}
	q.Where("state", state).Order("time desc").Find(&list)

	start, end := InitPage(&page, len(list))
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
func SelectPeerABounsList(o *gorm.DB, page models.Page, user_name, start_time, end_time string) ([]models.PeerAbouns, models.Page, error) {
	peer_a_bouns := []models.TxIdList{}
	time := ""
	if start_time != "" {
		time = "and create_time>=" + "'" + start_time + "'" + " and create_time<=" + "'" + end_time + "'"
	}
	switch user_name {
	case "":
		o.Raw("select * from tx_id_list where comment=? "+time+" order by create_time desc", "节点分红").Find(&peer_a_bouns)
	default:
		users := []models.User{}
		err_1 := o.Raw("select * from user where user_name=?", user_name).Find(&users)
		if err_1.Error != nil || len(users) < 1 {
			return []models.PeerAbouns{}, page, err_1.Error
		}
		for _, v := range users {
			err := o.Raw("select * from tx_id_list where user_id=? and comment=? "+time+" order by create_time desc", v.UserId, "节点分红").Find(&peer_a_bouns)
			if err.Error != nil || len(peer_a_bouns) < 1 {
				return []models.PeerAbouns{}, page, err.Error
			}
		}
	}
	if len(peer_a_bouns) < 1 {
		return []models.PeerAbouns{}, page, errors.New("没有相关数据!")
	}
	QuickSortPeerABouns(peer_a_bouns, 0, len(peer_a_bouns)-1)

	start, end := InitPage(&page, len(peer_a_bouns))

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
			u := models.User{}
			o.Table("user").Where("user_id = ?", v.UserId).First(&u)
			_, level, _, err_tfor := ReturnSuperPeerLevel(o, v.UserId)
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

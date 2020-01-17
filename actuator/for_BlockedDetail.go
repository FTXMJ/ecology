package actuator

import (
	"ecology/logs"
	"ecology/models"
	"ecology/utils"

	"github.com/astaxie/beego/orm"

	"errors"
	"fmt"
	"strconv"
	"time"
)

// 更新借贷表
func FindLimitOneAndSaveBlo_d(o orm.Ormer, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account_id).QueryRow(&blocked_old)
	for_mula := models.Formula{EcologyId: account_id}
	o.Read(&for_mula, "ecology_id")

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
		return err
	}

	// 更新任务完成状态
	_, err_txid := o.Raw("update tx_id_list set order_state=? where tx_id=?", true, tx_id).Exec()
	if err_txid != nil {
		return err_txid
	}

	//更新生态仓库属性
	_, err_up := o.Raw("update account set bocked_balance=? where id=?", blocked_new.CurrentBalance, account_id).Exec()
	if err_up != nil {
		return err_up
	}

	//  直推收益
	user := models.User{UserId: user_id}
	o.Read(&user, "user_id")
	if user.FatherId != "" && coin_in >= 10 {
		ForAddCoin(o, user.FatherId, coin_in, 0.1)
	}
	return nil
}

// 创建第一条借贷记录
func NewCreateAndSaveBlo_d(o orm.Ormer, user_id, comment, tx_id string, coin_out, coin_in float64, account_id int) error {
	for_mula := models.Formula{EcologyId: account_id}
	o.Read(&for_mula, "ecology_id")

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
		return err
	}

	// 更新任务完成状态
	_, err_txid := o.Raw("update tx_id_list set order_state=? where tx_id=?", true, tx_id).Exec()
	if err_txid != nil {
		return err_txid
	}

	//更新生态仓库属性
	_, err_up := o.Raw("update account set bocked_balance=? where id=?", blocked_new.CurrentBalance, account_id).Exec()
	if err_up != nil {
		return err_up
	}

	//  直推收益
	user := models.User{UserId: user_id}
	o.Read(&user, "user_id")
	if user.FatherId != "" && coin_in > 10 {
		errrr := ForAddCoin(o, user.FatherId, coin_in, 0.1)
		if errrr != nil {
			return errrr
		}
	}
	return nil
}

//　把所有算力的值加起来  -- 更新静态｀动态的验证          bo 是是否调过　钱包的增加接口
func ForAddCoin(o orm.Ormer, father_id string, coin float64, proportion float64) error {
	user := models.User{UserId: father_id}
	o.Read(&user, "user_id")

	account := models.Account{UserId: father_id}
	o.Read(&account, "user_id")

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

	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user.UserId, account.Id).QueryRow(&blocked_old)

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
	} else if coin*proportion*proportion >= 1 && user.FatherId != "" {
		ForAddCoin(o, user.FatherId, (coin * proportion), proportion)
	}
	return nil
}

//条件查询对象包含的则视为条件
func SelectPondMachinemsgForAcc(o orm.Ormer, p models.FindObj, page models.Page, table_name string) ([]models.BlockedDetailIndex, models.Page, error) {
	list, err := SqlCreateValues(o, p, table_name)
	if err != nil {
		return []models.BlockedDetailIndex{}, models.Page{}, err
	}

	start, end := InitPage(&page, len(list))

	listle := ListLimit(list, start, end)

	lists := []models.BlockedDetailIndex{}
	for _, v := range listle {
		value, ok := v.(models.AccountDetail)
		fmt.Println(ok)
		var u models.User
		u.UserId = value.UserId
		o.Read(&u, "user_id")
		blo := models.BlockedDetailIndex{
			Id:             value.Id,
			UserId:         value.UserId,
			UserName:       u.UserName,
			CurrentRevenue: value.CurrentRevenue,
			CurrentOutlay:  value.CurrentOutlay,
			OpeningBalance: value.OpeningBalance,
			CurrentBalance: value.CurrentBalance,
			CreateDate:     value.CreateDate,
			Comment:        value.Comment,
			TxId:           value.TxId,
			Account:        value.Account,
			CoinType:       value.CoinType,
		}
		lists = append(lists, blo)
	}
	return lists, page, nil
}

func SelectPondMachinemsgForBlo(o orm.Ormer, p models.FindObj, page models.Page, table_name string) ([]models.BlockedDetailIndex, models.Page, error) {
	list, err := SqlCreateValues(o, p, table_name)
	if err != nil {
		return []models.BlockedDetailIndex{}, models.Page{}, err
	}

	start, end := InitPage(&page, len(list))

	listle := ListLimit(list, start, end)

	lists := []models.BlockedDetailIndex{}
	for _, v := range listle {
		value, ok := v.(models.BlockedDetail)
		fmt.Println(ok)
		var u models.User
		u.UserId = value.UserId
		o.Read(&u, "user_id")
		blo := models.BlockedDetailIndex{
			Id:             value.Id,
			UserId:         value.UserId,
			UserName:       u.UserName,
			CurrentRevenue: value.CurrentRevenue,
			CurrentOutlay:  value.CurrentOutlay,
			OpeningBalance: value.OpeningBalance,
			CurrentBalance: value.CurrentBalance,
			CreateDate:     value.CreateDate,
			Comment:        value.Comment,
			TxId:           value.TxId,
			Account:        value.Account,
			CoinType:       value.CoinType,
		}
		lists = append(lists, blo)
	}
	return lists, page, nil
}

// 释放流水查询－－处理
func SelectFlows(o orm.Ormer, p models.FindObj, page models.Page, table_name string) ([]models.Flow, models.Page, error) {
	us := []models.User{}
	blos := []models.BlockedDetail{}

	q_user := o.QueryTable("user")

	if p.UserName != "" {
		q_user = q_user.Filter("user_name", p.UserName)
	}
	if p.UserId != "" {
		q_user = q_user.Filter("user_id", p.UserId)
	}
	q_user.All(&us)
	user_ids := []string{}
	for _, v := range us {
		user_ids = append(user_ids, v.UserId)
	}

	q_blos := o.QueryTable(table_name).Filter("user_id__in", user_ids)

	if p.StartTime != "" && p.EndTime != "" {
		q_blos = q_blos.Filter("create_date__gte", p.StartTime).Filter("create_date__lte", p.EndTime)
	}
	q_blos.OrderBy("-create_date").All(&blos)

	QuickSortBlockedDetail(blos, 0, len(blos)-1)

	start, end := InitPage(&page, len(blos))

	value_list := []interface{}{}
	for _, v := range blos {
		value_list = append(value_list, v)
	}

	list := ListLimit(value_list, start, end)

	flows := []models.Flow{}
	for _, value := range list {
		v := value.(models.BlockedDetail)
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

// Find user ecology information
func FindU_E_OBJ(o orm.Ormer, page models.Page, user_id, user_name string) ([]models.U_E_OBJ, models.Page) {
	users := []models.User{}

	q_user := o.QueryTable("user")
	if user_name != "" {
		q_user = q_user.Filter("user_name", user_name)
	}
	if user_id != "" {
		q_user = q_user.Filter("user_id", user_id)
	}
	q_user.All(&users)
	user_ids := []string{}
	for _, v := range users {
		user_ids = append(user_ids, v.UserId)
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

	q_user := o.QueryTable("user")
	if user_name != "" {
		q_user = q_user.Filter("user_name", user_name)
	}
	if user_id != "" {
		q_user = q_user.Filter("user_id", user_id)
	}
	q_user.All(&users)

	f_u_s := []interface{}{}
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

	start, end := InitPage(&page, len(f_u_s))
	list := []interface{}{}
	if list = ListLimit(f_u_s, start, end); len(list) == 0 {
		return []models.FalseUser{}, page
	}

	last_list := []models.FalseUser{}
	for _, v := range list {
		last_list = append(last_list, v.(models.FalseUser))
	}

	return last_list, page
}

// Find user ecology information
func FindUserAccountOFF(o orm.Ormer, page models.Page, obj models.FindObj) ([]models.AccountOFF, models.Page, error) {
	accounts, err := SqlCreateValues(o, obj, "account")
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
}

//	对象包含的则视为条件     sql 生成并　查询
func SqlCreateValues(o orm.Ormer, p models.FindObj, table_name string) ([]interface{}, error) {
	// 根据条件  进行数据插叙
	if list, err := GeneratedSQLAndExec(o, table_name, p); err != nil {
		return []interface{}{}, err
	} else {
		return list, nil
	}
}

func ShowMrsfTable(o orm.Ormer, page models.Page, user_name, user_id, date string, state bool) ([]models.MrsfStateTable, models.Page, error) {
	list := []models.MrsfStateTable{}

	acc := models.Account{UserId: user_id}
	o.Read(&acc, "user_id")

	q := o.QueryTable("mrsf_state_table")
	if user_name != "" {
		q = q.Filter("user_name", user_name)
	}
	if user_id != "" {
		q = q.Filter("user_id", user_id)
	}
	if date != "" {
		q = q.Filter("date", date)
	}
	q.Filter("state", state).OrderBy("-time").All(&list)

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

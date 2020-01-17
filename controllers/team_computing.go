package controllers

import (
	"ecology/actuator"
	db "ecology/db"
	"ecology/models"
	"ecology/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"errors"
	"strconv"
	"time"
)

type Test struct {
	beego.Controller
}

// 用户每日任务数值列表
type UserDayTx struct {
	UserId string
	BenJin float64
	Team   float64
	ZhiTui float64
}

type operatio_n struct {
	Jintai  bool
	Dongtai bool
	Peer    bool
}

type info struct {
	peer_a_bouns float64
	one          int
	two          int
	three        int
}

// @Tags 测试每日释放
// @Accept  json
// @Produce json
// @Success 200
// @router /admin/test_mrsf [GET]
func (this *Test) DailyDividendAndReleaseTest() {
	o := db.NewEcologyOrm()
	user := []models.User{}
	o.QueryTable("user").All(&user)

	//    每日释放___and___团队收益___and___直推收益
	error_users := ProducerEcology(user, "") // 返回错误的用户名单

	//    给失败的用户　添加失败的任务记录表
	CreateErrUserTxList(error_users)

	// 超级节点的分红
	in_fo := info{}
	err_peer := ProducerPeer(user, &in_fo, "")
	if err_peer == nil {
		perr_h := models.PeerHistory{
			Time:             time.Now().Format("2006-01-02 15:04:05"),
			WholeNetworkTfor: db.NetIncome,
			PeerABouns:       in_fo.peer_a_bouns,
			DiamondsPeer:     in_fo.one,
			SuperPeer:        in_fo.two,
			CreationPeer:     in_fo.three,
		}
		o.Insert(&perr_h)
	}

	// 让收益回归今日
	blo := []models.BlockedDetail{}
	o.Raw("select * form blocked_detail where create_date>=?", time.Now().Format("2006-01-02")+" 00:00:00").QueryRows(&blo)
	shouyi := 0.0
	if len(blo) >= 1 {
		for _, v := range blo {
			shouyi += v.CurrentOutlay
			shouyi += v.CurrentRevenue
		}
	}
	db.NetIncome = shouyi
}

func DailyDividendAndRelease() {
	o := db.NewEcologyOrm()
	user := []models.User{}
	o.QueryTable("user").All(&user)

	//    每日释放___and___团队收益___and___直推收益
	error_users := ProducerEcology(user, "") // 返回错误的用户名单

	//    给失败的用户　添加失败的任务记录表
	CreateErrUserTxList(error_users)

	// 超级节点的分红
	in_fo := info{}
	err_peer := ProducerPeer(user, &in_fo, "")
	if err_peer == nil {
		perr_h := models.PeerHistory{
			Time:             time.Now().Format("2006-01-02 15:04:05"),
			WholeNetworkTfor: db.NetIncome,
			PeerABouns:       in_fo.peer_a_bouns,
			DiamondsPeer:     in_fo.one,
			SuperPeer:        in_fo.two,
			CreationPeer:     in_fo.three,
		}
		o.Insert(&perr_h)
	}

	// 让收益回归今日
	blo := []models.BlockedDetail{}
	o.Raw("select * form blocked_detail where create_date>=?", time.Now().Format("2006-01-02")+" 00:00:00").QueryRows(&blo)
	shouyi := 0.0
	for _, v := range blo {
		shouyi += v.CurrentOutlay
		shouyi += v.CurrentRevenue
	}
	db.NetIncome = shouyi
}

func DailyDividendAndReleaseToSomeOne(user []string, order_id string) {
	o := db.NewEcologyOrm()
	users := []models.User{}
	for _, v := range user {
		u := models.User{UserId: v}
		o.Read(&u)
		users = append(users, u)
	}

	//    每日释放___and___团队收益___and___直推收益
	error_users := ProducerEcology(users, order_id) // 返回错误的用户名单

	//    给失败的用户　添加失败的任务记录表
	CreateErrUserTxList(error_users)

	// 超级节点的分红
	in_fo := info{}
	err_peer := ProducerPeer(users, &in_fo, order_id)
	if err_peer == nil {
		perr_h := models.PeerHistory{
			Time:             time.Now().Format("2006-01-02 15:04:05"),
			WholeNetworkTfor: db.NetIncome,
			PeerABouns:       in_fo.peer_a_bouns,
			DiamondsPeer:     in_fo.one,
			SuperPeer:        in_fo.two,
			CreationPeer:     in_fo.three,
		}
		o.Insert(&perr_h)
	}

	// 让收益回归今日
	blo := []models.BlockedDetail{}
	o.Raw("select * form blocked_detail where create_date>=?", time.Now().Format("2006-01-02")+" 00:00:00").QueryRows(&blo)
	shouyi := 0.0
	for _, v := range blo {
		shouyi += v.CurrentOutlay
		shouyi += v.CurrentRevenue
	}
	db.NetIncome = shouyi
}

// 设置 全局状态
func OperationSet(o orm.Ormer, account *models.Account) {
	g_o := []models.GlobalOperations{}
	o.Raw("select * from global_operations").QueryRows(&g_o)
	for _, v := range g_o {
		switch v.Operation {
		case "全局静态收益控制":
			if account.StaticReturn == true {
				account.StaticReturn = v.State
			}
		case "全局动态收益控制":
			if account.DynamicRevenue == true {
				account.DynamicRevenue = v.State
			}
		case "全局节点分红控制":
			if account.PeerState == true {
				account.PeerState = v.State
			}
		}
	}
}

//生态仓库释放　－　团队收益  --  直推收益
func ProducerEcology(users []models.User, order_id string) []models.User {
	error_users := []models.User{}
	for _, v := range users {
		if err := Worker(v, order_id); err != nil {
			error_users = append(error_users, v)
		}
	}
	return error_users
}

//超级节点　的　释放
func ProducerPeer(users []models.User, in_fo *info, order_id string) error {
	o := db.NewEcologyOrm()
	g_o := models.GlobalOperations{Operation: "全局节点分红控制"}
	o.Read(&g_o, "operation")
	if g_o.State == false && g_o.Id > 0 {
		return errors.New("err")
	}
	error_users := []models.User{}
	m := make(map[string][]string)
	for _, v := range users {
		_, level, _, err := actuator.ReturnSuperPeerLevel(v.UserId)
		if err != nil {
			error_users = append(error_users, v)
		} else if level == "" && err == nil {
			// 没有出错，但是不符合超级节点的规则
		} else if level != "" && err == nil {
			m[level] = append(m[level], v.UserId)
		}
	}
	if len(error_users) > 0 {
		ProducerPeer(error_users, in_fo, order_id)
	}
	HandlerMap(o, m, in_fo, order_id)
	return nil
}

// 工作　函数
func Worker(user models.User, order_id string) error {
	o := db.NewEcologyOrm()
	team_a_bouns := 0.0
	ziyou_a_bouns := 0.0
	zhitui_a_bouns := 0.0
	o.Begin()
	account := models.Account{
		UserId: user.UserId,
	}
	o.Read(&account, "user_id")

	OperationSet(o, &account)

	if account.DynamicRevenue != true && account.StaticReturn != true {
		JintaiBuShiFang(o, user.UserId)
		DongtaiBuShiFang(o, user.UserId)
		TeamBuShiFang(o, user.UserId)
		o.Commit()
		CreateMrsfTable(o, user, account, true, order_id, 0, 0, 0)
		return nil
	} else if account.DynamicRevenue == true && account.StaticReturn != true { // 动态可以，静态禁止
		err_jin := JintaiBuShiFang(o, user.UserId)
		if err_jin != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_jin
		}
		t, err_team := Team(o, user)
		team_a_bouns = t
		if err_team != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_team
		}
		z, err_zhitui := ZhiTui(o, user.UserId)
		zhitui_a_bouns = z
		if err_zhitui != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_zhitui
		}

	} else if account.StaticReturn == true && account.DynamicRevenue != true { //静态可以，动态禁止
		err_dong := DongtaiBuShiFang(o, user.UserId)
		if err_dong != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_dong
		}
		err_team := TeamBuShiFang(o, user.UserId)
		if err_team != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_dong
		}
		j, err_jintai := Jintai(o, user)
		ziyou_a_bouns = j
		if err_jintai != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_jintai
		}
	} else { // 都可以
		t, err_team := Team(o, user)
		team_a_bouns = t
		if err_team != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_team
		}
		z, err_zhitui := ZhiTui(o, user.UserId)
		zhitui_a_bouns = z
		if err_zhitui != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_zhitui
		}
		j, err_jintai := Jintai(o, user)
		ziyou_a_bouns = j
		if err_jintai != nil {
			o.Rollback()
			CreateMrsfTable(o, user, account, false, order_id, 0, 0, 0)
			return err_jintai
		}
	}
	CreateMrsfTable(o, user, account, true, order_id, ziyou_a_bouns, zhitui_a_bouns, team_a_bouns)
	o.Commit()
	return nil
}

func CreateMrsfTable(o orm.Ormer, user models.User, account models.Account, s_tate bool, order_id string, ziyou, zhitui, team float64) {
	switch order_id {
	case "":
		mrsf_table := models.MrsfStateTable{
			UserId:       user.UserId,
			UserName:     user.UserName,
			State:        s_tate,
			Time:         time.Now().Format("2006-01-02 15:04:05"),
			OrderId:      strconv.Itoa(account.Id) + time.Now().Format("2006-01-02"),
			Date:         time.Now().Format("2006-01-02"),
			ZiYouABouns:  ziyou,
			ZhiTuiABouns: zhitui,
			TeamABouns:   team,
		}
		_, err := o.Insert(&mrsf_table)
		if err != nil {
			CreateMrsfTable(o, user, account, s_tate, order_id, ziyou, zhitui, team)
		}
	default:
		m := models.MrsfStateTable{OrderId: order_id}
		o.Read(&m, "order_id")
		m.State = s_tate
		m.Time = time.Now().Format("2006-01-02 15:04:05")
		m.ZiYouABouns = ziyou
		m.ZhiTuiABouns = zhitui
		m.TeamABouns = team
		_, err := o.Update(&m)
		if err != nil {
			CreateMrsfTable(o, user, account, s_tate, order_id, ziyou, zhitui, team)
		}
	}
}

func Team(o orm.Ormer, user models.User) (float64, error) {
	coins := []float64{}
	user_current_layer := []models.User{}
	// 团队收益　开始
	o.QueryTable("user").Filter("father_id", user.UserId).All(&user_current_layer)
	if len(user_current_layer) > 0 {
		for _, v := range user_current_layer {
			if user.UserId != v.UserId {
				// 获取用户teams
				team_user, err := actuator.GetTeams(v)
				if err != nil {
					if err.Error() != "用户未激活或被拉入黑名单" {
						return 0, err
					}
				}
				if len(team_user) > 0 {
					// 去处理这些数据 // 处理器，计算所有用户的收益  并发布任务和 分红记录
					coin, err_handler := actuator.HandlerOperation(team_user, v.UserId)
					if err_handler != nil {
						return 0, err_handler
					}
					coins = append(coins, coin)
				}
			}
		}
	}
	value, err_sort_a_r := SortABonusRelease(o, coins, user.UserId)
	if err_sort_a_r != nil {
		return 0, err_sort_a_r
	}
	// 团队收益　结束
	return value, nil
}

func Jintai(o orm.Ormer, user models.User) (float64, error) {
	z, err := DailyRelease(o, user.UserId)
	if err != nil {
		return 0, err
	}
	return z, nil
}

// 去掉最大的 团队收益
func SortABonusRelease(o orm.Ormer, coins []float64, user_id string) (float64, error) {
	for i := 0; i < len(coins)-1; i++ {
		for j := i + 1; j < len(coins); j++ {
			if coins[i] > coins[j] {
				coins[i], coins[j] = coins[j], coins[i]
			}
		}
	}
	value := 0.0
	for i := 0; i < len(coins)-1; i++ {
		value += coins[i]
	}

	if value == 0 {
		return 0, nil
	}

	acc := models.Account{
		UserId: user_id,
	}
	o.Read(&acc, "user_id")
	for_m := models.Formula{
		EcologyId: acc.Id,
	}
	o.Read(&for_m, "ecology_id")

	value = value * for_m.TeamReturnRate

	var account = models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	var formula = models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")

	if value > account.BockedBalance {
		value = account.BockedBalance
	}

	if value == 0 {
		return 0, nil
	}

	//任务表 USDD  铸币记录
	order_id := utils.TimeUUID()
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  true,
		WalletState: false,
		UserId:      user_id,
		Comment:     "每日团队收益",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: value,
		InCome:      0,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		return 0, errtxid_blo
	}

	//找最近的数据记录表
	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account.Id).QueryRow(&blocked_old)

	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: 0,
		CurrentOutlay:  value,
		CurrentBalance: blocked_old.CurrentBalance - value + 0,
		OpeningBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "每日团队收益",
		TxId:           order_id,
		Account:        account.Id,
		CoinType:       "USDD",
	}

	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		return 0, err
	}

	//更新生态仓库属性
	_, err_up := o.Raw("update account set bocked_balance=? where id=?", blocked_new.CurrentBalance, account.Id).Exec()
	if err_up != nil {
		return 0, err_up
	}

	err_ping_shifang := actuator.PingAddWalletCoin(user_id, value)
	if err_ping_shifang != nil {
		return 0, err_ping_shifang
	}

	_, err_up_tx := o.Raw("update tx_id_list set wallet_state=? where tx_id=?", true, order_id).Exec()
	if err_up_tx != nil {
		return 0, err_up_tx
	}

	db.NetIncome += value
	return value, nil
}

// 超级节点的分红
func AddFormulaABonus(user_id string, abonus float64) {
	o := db.NewEcologyOrm()
	o.Begin()

	//任务表 USDD  铸币记录
	order_id := utils.TimeUUID()
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  true,
		WalletState: true,
		UserId:      user_id,
		Comment:     "节点分红",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: abonus,
		InCome:      0,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		o.Rollback()
		AddFormulaABonus(user_id, abonus)
		return
	}
	o.Commit()
	return
}

// 每日释放
func DailyRelease(o orm.Ormer, user_id string) (float64, error) {
	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	formula := models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")
	blocked_yestoday := models.AccountDetail{}
	date_time := time.Now().AddDate(0, 0, -1).Format("2006-01-02 ") + "23:59:59"
	err_raw := o.Raw(
		"select * from account_detail where user_id=? and create_date<=? order by create_date desc,id desc limit 1", user_id, date_time).QueryRow(&blocked_yestoday)
	if err_raw != nil {
		if err_raw.Error() != "<QuerySeter> no row found" {
			return 0, err_raw
		}
	}
	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account.Id).QueryRow(&blocked_old)
	abonus := formula.HoldReturnRate * blocked_yestoday.CurrentBalance
	aabonus := blocked_old.CurrentBalance - abonus
	if aabonus < 0 {
		aabonus = 0
		abonus = blocked_old.CurrentBalance
	}
	if abonus == 0 {
		return 0, nil
	}
	//任务表 USDD  铸币记录
	order_id := utils.TimeUUID()
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  true,
		WalletState: false,
		UserId:      user_id,
		Comment:     "每日释放收益",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: abonus,
		InCome:      0,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		return 0, errtxid_blo
	}

	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: 0,
		CurrentOutlay:  abonus,
		CurrentBalance: aabonus,
		OpeningBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "每日释放收益",
		TxId:           order_id,
		Account:        account.Id,
		CoinType:       "USDD",
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		return 0, err
	}

	//更新生态仓库属性
	account.BockedBalance = aabonus
	_, err_up := o.Raw("update account set bocked_balance=? where id=?", aabonus, account.Id).Exec()
	if err_up != nil {
		return 0, err_up
	}

	// 钱包　数据　修改
	err_ping := actuator.PingAddWalletCoin(user_id, abonus)
	if err_ping != nil {
		return 0, err_ping
	}
	_, err_up_tx := o.Raw("update tx_id_list set wallet_state=? where tx_id=?", true, order_id).Exec()
	if err_up_tx != nil {
		return 0, err_up_tx
	}
	db.NetIncome += abonus
	return abonus, nil
}

//　直推收益
func ZhiTui(o orm.Ormer, user_id string) (float64, error) {
	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")

	formula := models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")

	blos := []models.BlockedDetail{}
	time_start := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 00:00:00"
	time_end := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 23:59:59"
	_, err := o.Raw("select * from blocked_detail where user_id=? and create_date>=? and create_date<=? and comment=?", user_id, time_start, time_end, "直推收益").QueryRows(&blos)
	if err != nil {
		if err.Error() != "<QuerySeter> no row found" {
			return 0, err
		}
	}
	shouyi := 0.0
	for _, v := range blos {
		shouyi += v.CurrentOutlay
	}

	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account.Id).QueryRow(&blocked_old)
	shouyia := blocked_old.CurrentBalance - shouyi
	if shouyia < 0 {
		shouyia = 0
		shouyi = blocked_old.CurrentBalance
	}

	if shouyi == 0 {
		return 0, nil
	}

	//任务表 USDD  铸币记录
	order_id := utils.TimeUUID()
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  true,
		WalletState: false,
		UserId:      user_id,
		Comment:     "每日直推收益",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: shouyi,
		InCome:      0,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		return 0, errtxid_blo
	}

	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: 0,
		CurrentOutlay:  shouyi,
		CurrentBalance: shouyia,
		OpeningBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "每日直推收益",
		TxId:           order_id,
		Account:        account.Id,
		CoinType:       "USDD",
	}

	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	_, err_in := o.Insert(&blocked_new)
	if err_in != nil {
		return 0, err_in
	}

	account.BockedBalance = blocked_new.CurrentBalance
	_, err_update := o.Raw("update account set bocked_balance=? where id=?", blocked_new.CurrentBalance, account.Id).Exec()
	if err_update != nil {
		return 0, err_update
	}

	// 钱包　数据　修改
	err_ping := actuator.PingAddWalletCoin(user_id, shouyi)
	if err_ping != nil {
		return 0, err_ping
	}
	_, err_up_tx := o.Raw("update tx_id_list set wallet_state=? where tx_id=?", true, order_id).Exec()
	if err_up_tx != nil {
		return 0, err_up_tx
	}

	db.NetIncome += shouyi
	return shouyi, nil
}

// 创建用于超级节点　等级记录的　map 每个　values 第一个元素都是　等级标示
func ReturnMap(m map[string][]string) {
	s_f_t := []models.SuperForceTable{}
	db.NewEcologyOrm().QueryTable("super_force_table").All(&s_f_t)
	for _, v := range s_f_t {
		if m[v.Level] == nil {
			m[v.Level] = append(m[v.Level], v.Level)
		}
	}
}

// 处理map数据并给定收益
func HandlerMap(o orm.Ormer, m map[string][]string, in_fo *info, order_id string) {
	err_m := make(map[string][]string)
	for k_level, vv := range m {
		s_f_t := models.SuperForceTable{
			Level: k_level,
		}
		o.Read(&s_f_t, "level")
		tfor_some := db.NetIncome * s_f_t.BonusCalculation
		for _, v := range vv {
			acc := models.Account{UserId: v}
			o.Read(&acc, "user_id")
			if acc.PeerState == true {
				err := actuator.PingAddWalletCoin(v, tfor_some/float64(len(vv)))
				if err != nil {
					err_m[k_level] = append(err_m[k_level], v)
				} else {
					if tfor_some/float64(len(vv)) != 0 {
						AddFormulaABonus(v, tfor_some/float64(len(vv)))
						in_fo.peer_a_bouns += tfor_some / float64(len(vv))
						if order_id != "" {
							o.Raw("update mrsf_state_table set peer_a_bouns=? where order_id=? and user_id=?", tfor_some/float64(len(vv)), order_id, v).Exec()
						} else {
							o.Raw(
								"update mrsf_state_table set peer_a_bouns=? where order_id=? and user_id=?",
								tfor_some/float64(len(vv)),
								time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
								v).Exec()
						}
					}
					if k_level == "钻石节点" {
						in_fo.one++
					} else if k_level == "超级节点" {
						in_fo.two++
					} else if k_level == "创世节点" {
						in_fo.three++
					}
				}
			}
		}
	}
	if len(err_m) != 0 {
		HandlerMap(o, err_m, in_fo, order_id)
	}
}

// 给失败的用户　添加失败的任务记录表
func CreateErrUserTxList(users []models.User) {
	o := db.NewEcologyOrm()
	err_users := []models.User{}
	for _, v := range users {
		//任务表 USDD  铸币记录
		order_id := utils.TimeUUID()
		blo_txid_dcmt := models.TxIdList{
			TxId:        order_id,
			OrderState:  false,
			WalletState: false,
			UserId:      v.UserId,
			Comment:     "每日任务失败用户",
			CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
			Expenditure: 0,
			InCome:      0,
		}
		_, errtxid_blo := o.Insert(&blo_txid_dcmt)
		if errtxid_blo != nil {
			err_users = append(err_users, v)
		}
	}
	if len(err_users) != 0 {
		CreateErrUserTxList(err_users)
	}
}

//静态的收益累加
func JintaiBuShiFang(o orm.Ormer, user_id string) error {
	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	formula := models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")
	blocked_yestoday := models.AccountDetail{}
	err_raw := o.Raw(
		"select * from account_detail where user_id=? and create_date<=? order by create_date desc,id desc limit 1",
		user_id,
		time.Now().AddDate(0, 0, -1).Format("2006-01-02 ")+"23:59:59").
		QueryRow(&blocked_yestoday)
	if err_raw != nil {
		if err_raw.Error() != "<QuerySeter> no row found" {
			return err_raw
		}
	}
	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account.Id).QueryRow(&blocked_old)
	abonus := formula.HoldReturnRate * blocked_yestoday.CurrentBalance
	aabonus := blocked_old.CurrentBalance - abonus
	if aabonus < 0 {
		aabonus = 0
		abonus = blocked_old.CurrentBalance
	}
	if abonus == 0 {
		return nil
	}

	db.NetIncome += abonus
	return nil
}

//直推累加
func DongtaiBuShiFang(o orm.Ormer, user_id string) error {
	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")

	formula := models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")

	blos := []models.BlockedDetail{}
	time_start := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 00:00:00"
	time_end := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + " 23:59:59"
	_, err := o.Raw("select * from blocked_detail where user_id=? and create_date>=? and create_date<=? and comment=?", user_id, time_start, time_end, "直推收益").QueryRows(&blos)
	if err != nil {
		if err.Error() != "<QuerySeter> no row found" {
			return err
		}
	}
	shouyi := 0.0
	for _, v := range blos {
		shouyi += v.CurrentOutlay
	}

	blocked_old := models.BlockedDetail{}
	o.Raw("select * from blocked_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account.Id).QueryRow(&blocked_old)

	shouyia := blocked_old.CurrentBalance - shouyi
	if shouyia < 0 {
		shouyia = 0
		shouyi = blocked_old.CurrentBalance
	}

	if shouyi == 0 {
		return nil
	}
	db.NetIncome += shouyi
	return nil
}

// 团队累加
func TeamBuShiFang(o orm.Ormer, user_id string) error {
	coins := []float64{}
	user_current_layer := []models.User{}
	// 团队收益　开始
	o.Raw("select * from user where father_id=?", user_id).QueryRows(&user_current_layer)
	if len(user_current_layer) > 0 {
		for _, v := range user_current_layer {
			if user_id != v.UserId {
				// 获取用户teams
				team_user, err := actuator.GetTeams(v)
				if err != nil {
					if err.Error() != "用户未激活或被拉入黑名单" {
						return err
					}
				}
				if len(team_user) > 0 {
					// 去处理这些数据 // 处理器，计算所有用户的收益  并发布任务和 分红记录
					coin, err_handler := actuator.HandlerOperation(team_user, v.UserId)
					if err_handler != nil {
						return err_handler
					}
					coins = append(coins, coin)
				}
			}
		}
	}
	err_sort_a_r := TeamWork(o, coins, user_id)
	if err_sort_a_r != nil {
		return err_sort_a_r
	}
	// 团队收益　结束
	return nil
}

//　计算团队的收益,但是不给定
func TeamWork(o orm.Ormer, coins []float64, user_id string) error {
	for i := 0; i < len(coins)-1; i++ {
		for j := i + 1; j < len(coins); j++ {
			if coins[i] > coins[j] {
				coins[i], coins[j] = coins[j], coins[i]
			}
		}
	}
	value := 0.0
	for i := 0; i < len(coins)-1; i++ {
		value += coins[i]
	}

	if value == 0 {
		return nil
	}

	acc := models.Account{
		UserId: user_id,
	}
	o.Read(&acc, "user_id")
	for_m := models.Formula{
		EcologyId: acc.Id,
	}
	o.Read(&for_m, "ecology_id")

	value = value * for_m.TeamReturnRate

	var account = models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	var formula = models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")
	value = value * formula.TeamReturnRate

	if value > account.BockedBalance {
		value = account.BockedBalance
	}

	if value == 0 {
		return nil
	}
	db.NetIncome += value
	return nil
}

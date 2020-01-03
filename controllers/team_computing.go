package controllers

import (
	"ecology/consul"
	"ecology/logs"
	"ecology/models"
	"ecology/utils"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// @Tags 测试每日释放
// @Accept  json
// @Produce json
// @Success 200
// @router /test_mrsf [GET]
func (this *Test) DailyDividendAndReleaseTest() {
	logs.Log.Info("开始")
	o := models.NewOrm()
	user := []models.User{}
	o.QueryTable("user").All(&user)

	//    每日释放___and___团队收益___and___直推收益
	error_users := ProducerEcology(user) // 返回错误的用户名单

	//    给失败的用户　添加失败的任务记录表
	CreateErrUserTxList(error_users)

	// 超级节点的分红
	peer_a_bouns := 0.0
	one := 0
	two := 0
	three := 0
	ProducerPeer(user, peer_a_bouns, one, two, three)
	perr_h := models.PeerHistory{
		Time:             time.Now().Format("2006-01-02 15:04:05"),
		WholeNetworkTfor: models.NetIncome,
		PeerABouns:       peer_a_bouns,
		DiamondsPeer:     one,
		SuperPeer:        two,
		CreationPeer:     three,
	}
	o.Insert(&perr_h)

	// 让收益回归今日
	blo := []models.BlockedDetail{}
	o.Raw("select * form blocked_detail where create_date>=?", time.Now().Format("2006-01-02")+" 00:00:00").QueryRows(&blo)
	shouyi := 0.0
	for _, v := range blo {
		shouyi += v.CurrentOutlay
		shouyi += v.CurrentRevenue
	}
	models.NetIncome = shouyi
	logs.Log.Info("结束")
}

func DailyDividendAndRelease() {
	logs.Log.Info("开始")
	o := models.NewOrm()
	user := []models.User{}
	o.QueryTable("user").All(&user)

	//    每日释放___and___团队收益___and___直推收益
	error_users := ProducerEcology(user) // 返回错误的用户名单

	//    给失败的用户　添加失败的任务记录表
	CreateErrUserTxList(error_users)

	// 超级节点的分红
	peer_a_bouns := 0.0
	one := 0
	two := 0
	three := 0
	ProducerPeer(user, peer_a_bouns, one, two, three)
	perr_h := models.PeerHistory{
		Time:             time.Now().Format("2006-01-02 15:04:05"),
		WholeNetworkTfor: models.NetIncome,
		PeerABouns:       peer_a_bouns,
		DiamondsPeer:     one,
		SuperPeer:        two,
		CreationPeer:     three,
	}
	o.Insert(&perr_h)

	// 让收益回归今日
	blo := []models.BlockedDetail{}
	o.Raw("select * form blocked_detail where create_date>=?", time.Now().Format("2006-01-02")+" 00:00:00").QueryRows(&blo)
	shouyi := 0.0
	for _, v := range blo {
		shouyi += v.CurrentOutlay
		shouyi += v.CurrentRevenue
	}
	models.NetIncome = shouyi
	logs.Log.Info("结束")
}

//生态仓库释放　－　团队收益  --  直推收益
func ProducerEcology(users []models.User) []models.User {
	error_users := []models.User{}
	for _, v := range users {
		if err := Worker(v); err != nil {
			error_users = append(error_users, v)
		}
	}
	return error_users
}

//超级节点　的　释放
func ProducerPeer(users []models.User, peer_a_bouns float64, one, two, three int) {
	error_users := []models.User{}
	m := make(map[string][]string)
	for _, v := range users {
		_, level, _, err := ReturnSuperPeerLevel(v.UserId)
		if err != nil {
			error_users = append(error_users, v)
		} else if level == "" && err == nil {
			// 没有出错，但是不符合超级节点的规则
		} else if level != "" && err == nil {
			m[level] = append(m[level], v.UserId)
		}
	}
	if len(error_users) > 0 {
		ProducerPeer(error_users, peer_a_bouns, one, two, three)
	}
	HandlerMap(m, peer_a_bouns, one, two, three)
}

// 工作　函数
func Worker(user models.User) error {
	o := models.NewOrm()
	o.Begin()
	account := models.Account{
		UserId: user.UserId,
	}
	o.Read(&account, "user_id")

	if account.DynamicRevenue != true && account.StaticReturn != true {
		o.Commit()
		return nil
	} else if account.DynamicRevenue == true && account.StaticReturn != true { // 动态可以，静态禁止
		err_team := Team(o, user)
		if err_team != nil {
			o.Rollback()
			return err_team
		}
		err_zhitui := ZhiTui(o, user.UserId)
		if err_zhitui != nil {
			o.Rollback()
			return err_zhitui
		}

	} else if account.StaticReturn == true && account.DynamicRevenue != true { //静态可以，动态禁止
		err_jintai := Jintai(o, user)
		if err_jintai != nil {
			o.Rollback()
			return err_jintai
		}
	} else { // 都可以
		err_team := Team(o, user)
		if err_team != nil {
			o.Rollback()
			return err_team
		}
		err_zhitui := ZhiTui(o, user.UserId)
		if err_zhitui != nil {
			o.Rollback()
			return err_zhitui
		}
		err_jintai := Jintai(o, user)
		if err_jintai != nil {
			o.Rollback()
			return err_jintai
		}
	}
	o.Commit()
	return nil
}

func Team(o orm.Ormer, user models.User) error {
	coins := []float64{}
	user_current_layer := []models.User{}
	// 团队收益　开始
	o.QueryTable("user").Filter("father_id", user.UserId).All(&user_current_layer)
	if len(user_current_layer) > 0 {
		for _, v := range user_current_layer {
			if user.UserId != v.UserId {
				// 获取用户teams
				team_user, err := GetTeams(v)
				if err != nil {
					if err.Error() != "用户未激活或被拉入黑名单" {
						return err
					}
				}
				if len(team_user) > 0 {
					// 去处理这些数据 // 处理器，计算所有用户的收益  并发布任务和 分红记录
					coin, err_handler := HandlerOperation(team_user, user.UserId)
					if err_handler != nil {
						return err_handler
					}
					coins = append(coins, coin)
				}
			}
		}
	}
	err_sort_a_r := SortABonusRelease(o, coins, user.UserId)
	if err_sort_a_r != nil {
		return err_sort_a_r
	}
	// 团队收益　结束
	return nil
}

func Jintai(o orm.Ormer, user models.User) error {
	err := DailyRelease(o, user.UserId)
	if err != nil {
		return err
	}
	return nil
}

type data_users struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}

// 从晓东那里获取团队 成员  直推
func GetTeams(user models.User) ([]string, error) {
	client := &http.Client{}
	//生成要访问的url
	url := consul.GetUserApi + beego.AppConfig.String("api::apiurl_get_team")
	//提交请求
	reqest, errnr := http.NewRequest("GET", url, nil)

	b, token := generateToken(user)
	if b != true {
		return nil, errors.New("err")
	}

	//增加header选项
	reqest.Header.Add("Authorization", token)
	if errnr != nil {
		return nil, errnr
	}

	//处理返回结果
	response, errdo := client.Do(reqest)
	defer response.Body.Close()
	if errdo != nil {
		return nil, errdo
	}
	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return nil, err_read
	}
	values := data_users{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return values.Data, errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return values.Data, errors.New(values.Msg)
	}
	users := []string{}
	for _, v := range values.Data {
		users = append(users, v)
	}
	return users, nil
}

// 处理器，计算所有用户的收益  并发布任务和 分红记录
func HandlerOperation(users []string, user_id string) (float64, error) {
	o := models.NewOrm()
	var coin_abouns float64
	for _, v := range users {
		if user_id != v {
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
	}
	return coin_abouns, nil
}

// 去掉最大的 团队收益
func SortABonusRelease(o orm.Ormer, coins []float64, user_id string) error {
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

	acc := models.Account{
		UserId: user_id,
	}
	o.Read(&acc, "user_id")
	for_m := models.Formula{
		EcologyId: acc.Id,
	}
	o.Read(&for_m, "ecology_id")

	value += acc.BockedBalance * for_m.HoldReturnRate

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
		return errtxid_blo
	}

	//找最近的数据记录表
	blocked_olds := []models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
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
		return err
	}

	//更新生态仓库属性
	_, err_up := o.QueryTable("account").Filter("id", account.Id).Update(orm.Params{"bocked_balance": blocked_new.CurrentBalance})
	if err_up != nil {
		return err_up
	}

	err_ping_shifang := PingAddWalletCoin(user_id, value)
	if err_ping_shifang != nil {
		return err_ping_shifang
	}
	order := models.TxIdList{
		TxId: order_id,
	}
	err_rea := o.Read(&order, "tx_id")
	if err_rea != nil {
		return err_rea
	}
	order.WalletState = true
	_, err_up_tx := o.Update(&order, "wallet_state")
	if err_up_tx != nil {
		return err_up_tx
	}

	models.NetIncome += value
	return nil
}

// 超级节点的分红
func AddFormulaABonus(user_id string, abonus float64) {
	o := models.NewOrm()
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
func DailyRelease(o orm.Ormer, user_id string) error {
	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	formula := models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")
	blocked_yestoday := []models.AccountDetail{}
	_, err_raw := o.Raw(
		"select * from account_detail where user_id=? and create_date<=? order by create_date desc limit 3",
		user_id,
		time.Now().AddDate(0, 0, -1).Format("2006-01-02 ")+"23:59:59").
		QueryRows(&blocked_yestoday)
	if err_raw != nil {
		if err_raw.Error() != "<QuerySeter> no row found" {
			return err_raw
		}
	}
	var blocked_old1 models.AccountDetail
	if len(blocked_yestoday) != 0 {
		for i := 0; i < len(blocked_yestoday)-1; i++ {
			for j := i + 1; j < len(blocked_yestoday); j++ {
				if blocked_yestoday[i].Id > blocked_yestoday[j].Id {
					blocked_yestoday[i], blocked_yestoday[j] = blocked_yestoday[j], blocked_yestoday[i]
				}
			}
		}

		blocked_old1 = blocked_yestoday[len(blocked_yestoday)-1]
		if blocked_old1.Id == 0 {
			blocked_old1.CurrentBalance = 0
		}
	}
	blocked_olds := []models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
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
	abonus := formula.HoldReturnRate * blocked_old1.CurrentBalance
	aabonus := blocked_old.CurrentBalance - abonus
	if aabonus < 0 {
		aabonus = 0
		abonus = blocked_old.CurrentBalance
	}
	if abonus == 0 {
		return nil
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
		return errtxid_blo
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
		return err
	}

	//更新生态仓库属性
	account.BockedBalance = aabonus
	_, err_up := o.Update(&account, "bocked_balance")
	if err_up != nil {
		return err_up
	}

	// 钱包　数据　修改
	err_ping := PingAddWalletCoin(user_id, abonus)
	if err_ping != nil {
		return err_ping
	}
	order := models.TxIdList{
		TxId: order_id,
	}
	err_rea := o.Read(&order, "tx_id")
	if err_rea != nil {
		return err_rea
	}
	order.WalletState = true
	_, err_up_tx := o.Update(&order)
	if err_up_tx != nil {
		return err_up_tx
	}

	models.NetIncome += abonus
	return nil
}

//　直推收益
func ZhiTui(o orm.Ormer, user_id string) error {
	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")

	formula := models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula)

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

	blocked_olds := []models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
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
	shouyia := blocked_old.CurrentBalance - shouyi
	if shouyia < 0 {
		shouyia = 0
		shouyi = blocked_old.CurrentBalance
	}

	if shouyi == 0 {
		return nil
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
		return errtxid_blo
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
		return err_in
	}

	account.BockedBalance = blocked_new.CurrentBalance
	_, err_update := o.Update(&account, "bocked_balance")
	if err_update != nil {
		return err_update
	}

	// 钱包　数据　修改
	err_ping := PingAddWalletCoin(user_id, shouyi)
	if err_ping != nil {
		return err_ping
	}
	order := models.TxIdList{
		TxId: order_id,
	}
	err_rea := o.Read(&order, "tx_id")
	if err_rea != nil {
		return err_rea
	}
	order.WalletState = true
	_, err_up_tx := o.Update(&order)
	if err_up_tx != nil {
		return err_up_tx
	}

	models.NetIncome += shouyi
	return nil
}

// 远端连接  -  给定分红收益  释放通用
func PingAddWalletCoin(user_id string, abonus float64) error {
	if abonus == 0 {
		return nil
	}
	user := models.User{
		UserId: user_id,
	}
	b, token := generateToken(user)
	if b != true {
		return errors.New("err")
	}
	//生成要访问的url
	apiurl := consul.GetWalletApi
	resoure := beego.AppConfig.String("api::apiurl_share_bonus")
	data := url.Values{}
	data.Set("money", strconv.FormatFloat(abonus, 'f', -1, 64))
	data.Set("symbol", "USDD")

	u, _ := url.ParseRequestURI(apiurl)
	u.Path = resoure
	urlStr := u.String()

	client := &http.Client{}
	req, err1 := http.NewRequest(`POST`, urlStr, strings.NewReader(data.Encode()))
	if err1 != nil {
		return err1
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	//处理返回结果
	response, errdo := client.Do(req)
	if errdo != nil {
		return errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return err_read
	}
	values := models.Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return errors.New(values.Msg)
	}
	response.Body.Close()
	return nil
}

// TFOR 数量查询
func PingSelectTforNumber(user_id string) (string, float64, error) {
	user := models.User{
		UserId: user_id,
	}
	b, token := generateToken(user)
	if b != true {
		return "", 0.0, errors.New("err")
	}
	//生成要访问的url
	apiurl := consul.GetWalletApi
	resoure := beego.AppConfig.String("api::apiurl_tfor_info")
	data := url.Values{}

	u, _ := url.ParseRequestURI(apiurl)
	u.Path = resoure
	urlStr := u.String()

	client := &http.Client{}
	req, err1 := http.NewRequest(`GET`, urlStr, strings.NewReader(data.Encode()))
	if err1 != nil {
		return "", 0.0, err1
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	//处理返回结果
	response, errdo := client.Do(req)
	if errdo != nil {
		return "", 0.0, errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return "", 0.0, err_read
	}
	values := models.Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return "", 0.0, errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return "", 0.0, errors.New(values.Msg)
	}
	response.Body.Close()
	aa := values.Data["balance"].(string)
	time := values.Data["updated_at"].(string)
	bb, err := strconv.ParseFloat(aa, 64)
	return time, bb, nil
}

// 返回超级节点的等级
func ReturnSuperPeerLevel(user_id string) (time, level string, tfor float64, err error) {
	s_f_t := []models.SuperForceTable{}
	models.NewOrm().QueryTable("super_force_table").All(&s_f_t)
	up_time, tfor_number, err_tfor := PingSelectTforNumber(user_id)
	if err_tfor != nil {
		return "", "", 0.0, err_tfor
	}

	for i := 0; i < len(s_f_t); i++ {
		for j := 1; j < len(s_f_t)-1; j++ {
			if s_f_t[i].CoinNumberRule > s_f_t[j].CoinNumberRule {
				s_f_t[i].CoinNumberRule, s_f_t[j].CoinNumberRule = s_f_t[j].CoinNumberRule, s_f_t[i].CoinNumberRule
			}
		}
	}
	index := []int{}
	for i, v := range s_f_t {
		if tfor_number > float64(v.CoinNumberRule) {
			index = append(index, i)
		}
	}
	if len(index) > 0 {
		return up_time, s_f_t[index[len(index)-1]].Level, tfor_number, nil
	}
	return "", "", 0.0, err_tfor
}

// 创建用于超级节点　等级记录的　map 每个　values 第一个元素都是　等级标示
func ReturnMap(m map[string][]string) {
	s_f_t := []models.SuperForceTable{}
	models.NewOrm().QueryTable("super_force_table").All(&s_f_t)
	for _, v := range s_f_t {
		if m[v.Level] == nil {
			m[v.Level] = append(m[v.Level], v.Level)
		}
	}
}

// 处理map数据并给定收益
func HandlerMap(m map[string][]string, peer_a_bouns float64, one, two, three int) {
	err_m := make(map[string][]string)
	for k_level, vv := range m {
		s_f_t := models.SuperForceTable{
			Level: k_level,
		}
		models.NewOrm().Read(&s_f_t, "level")
		tfor_some := models.NetIncome * s_f_t.BonusCalculation
		for _, v := range vv {
			err := PingAddWalletCoin(v, tfor_some/float64(len(vv)))
			if err != nil {
				err_m[k_level] = append(err_m[k_level], v)
			} else {
				if tfor_some/float64(len(vv)) != 0 {
					AddFormulaABonus(v, tfor_some/float64(len(vv)))
					peer_a_bouns += tfor_some / float64(len(vv))
				}
				if k_level == "钻石节点" {
					one++
				} else if k_level == "超级节点" {
					two++
				} else if k_level == "创世节点" {
					three++
				}
			}
		}
	}
	if len(err_m) != 0 {
		HandlerMap(err_m, peer_a_bouns, one, two, three)
	}
}

// 给失败的用户　添加失败的任务记录表
func CreateErrUserTxList(users []models.User) {
	o := models.NewOrm()
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

/*
			每日任务－－－思路
	users := []UserDayTx{}// 准备好的　obj　用来接受　最终处理数据
	// 计算　所有　用户　当日　发放　的值
	for _, v := range u {
		if err := UsersPermissionsFiltering(v,users); err != nil {
			error_users = append(error_users, v)
		}
	}

	// 错误的用户,生成　未完成　的用户表　　　　等待"客服人员"操作
	//TODO

	// 集体性的　钱包数据更新　　以及生态数据库的更新
*/
/*
// 执行赋值操作
func UsersPermissionsFiltering(u models.User,users []UserDayTx) error {
	o := models.NewOrm()
	account := models.Account{
		UserId: u.UserId,
	}
	o.Read(&account, "user_id")

	user := UserDayTx{
		UserId: u.UserId,
		BenJin: 0,
		Team:   0,
		ZhiTui: 0,
	}
	var err error
	if account.DynamicRevenue != true && account.StaticReturn != true {

	} else if account.DynamicRevenue == true && account.StaticReturn != true { // 动态可以，静态禁止
		err = TeamValue(user)
		err = ZhiTuiValue(user)
	} else if account.StaticReturn == true && account.DynamicRevenue != true { //静态可以，动态禁止
		err = BenJinValue(user)
	} else { // 都可以
		err = TeamValue(user)
		err = ZhiTuiValue(user)
		err = BenJinValue(user)
	}
	return nil
}

// 给所有人的　每日本金自由算力收益　　　赋值
func BenJinValue(user UserDayTx) error {

	return nil
}

// 给所有人的　每日团队收益　　			赋值
func TeamValue(user UserDayTx) error {

	return nil
}

// 给所有人的　前一天的直推收益　　　	赋值
func ZhiTuiValue(user UserDayTx) error {

	return nil
}
*/

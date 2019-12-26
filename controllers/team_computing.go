package controllers

import (
	"ecology/consul"
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

// @Tags 测试每日释放
// @Accept  json
// @Produce json
// @Success 200
// @router /test_mrsf [GET]
func (this *Test) DailyDividendAndRelease() {
	o := models.NewOrm()
	users := []models.User{}
	o.QueryTable("user").All(&users)
	//每日释放　－－　团队收益
	ProducerEcology(users)

	// 超级节点的分红
	ProducerPeer(users)
	models.NetIncome = 0.0
}

//生态仓库释放　－　团队收益
func ProducerEcology(users []models.User) {
	error_users := []models.User{}
	for _, v := range users {
		if err := Worker(v); err != nil {
			error_users = append(error_users, v) //  385ce4bc5a5141a790de7c009cc89e92
			//TODO  写日志,告知管理者失败的原因
		}
	}
	if len(error_users) > 0 {
		ProducerEcology(error_users)
	}
}

//超级节点　的　释放
func ProducerPeer(users []models.User) {
	error_users := []models.User{}
	m := make(map[string][]string)
	for _, v := range users {
		level, err := ReturnSuperPeerLevel(v.UserId)
		if err != nil {
			error_users = append(error_users, v)
		} else if level == "" && err == nil {
			// 没有出错，但是不符合超级节点的规则
		} else if level != "" && err == nil {
			m[level] = append(m[level], v.UserId)
		}
	}
	if len(error_users) > 0 {
		ProducerPeer(error_users)
	}
	HandlerMap(m)
}

func Worker(user models.User) error {
	o := models.NewOrm()
	o.Begin()
	user_current_layer := []models.User{}
	account := models.Account{
		UserId: user.UserId,
	}
	coins := []float64{}
	o.Read(&account, "user_id")

	if account.DynamicRevenue != true && account.StaticReturn != true {
		o.Commit()
		return nil
	} else if account.DynamicRevenue == true && account.StaticReturn != true { // 动态可以，静态禁止
		err_team := Team(o, user_current_layer, user, coins)
		if err_team != nil {
			o.Rollback()
			return err_team
		}
	} else if account.StaticReturn == true && account.DynamicRevenue != true { //静态可以，动态禁止
		err_jintai := Jintai(o, user)
		if err_jintai != nil {
			o.Rollback()
			return err_jintai
		}
	} else { // 都可以
		err_team := Team(o, user_current_layer, user, coins)
		if err_team != nil {
			o.Rollback()
			return err_team
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

func Team(o orm.Ormer, user_current_layer []models.User, user models.User, coins []float64) error {
	// 团队收益　开始
	o.QueryTable("user").Filter("father_id", user.UserId).All(&user_current_layer)
	if len(user_current_layer) > 0 {
		for _, v := range user_current_layer {
			// 获取用户teams
			team_user, err := GetTeams(v)
			if err != nil {
				return err
			}
			// 去处理这些数据 // 处理器，计算所有用户的收益  并发布任务和 分红记录
			coin, err_handler := HandlerOperation(team_user)
			if err_handler != nil {
				return err_handler
			}
			coins = append(coins, coin)
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
	/*// 超级节点分红
	err_super_peer := AddFormulaABonus(o, user.UserId)
	if err_super_peer != nil {
		return err_super_peer
	}*/
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

// 从晓东那里获取团队 成员
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
		return nil, err
	}
	return values.Data, nil
}

// 处理器，计算所有用户的收益  并发布任务和 分红记录
func HandlerOperation(users []string) (float64, error) {
	o := models.NewOrm()
	coin_abouns := 0.0
	for _, v := range users {
		// 拿到生态项目实例
		account := models.Account{}
		err_acc := o.QueryTable("account").Filter("user_id", v).One(&account)
		if err_acc != nil {
			if err_acc.Error() == "<QuerySeter> no row found" {
				return 0, nil
			} else {
				return 0, err_acc
			}
		}
		// 拿到生态项目对应的算力表
		formula := models.Formula{}
		err_for := o.QueryTable("formula").Filter("ecology_id", account.Id).One(&formula)
		if err_for != nil {
			return 0, err_for
		}
		coin_abouns += (formula.HoldReturnRate * account.Balance)
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

	//任务表 USDD  铸币记录
	tx_id_blo_d := utils.Shengchengstr("每日团队收益", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		State: "true",
		TxId:  tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		o.Rollback()
		return errtxid_blo
	}

	//找最近的数据记录表
	account := models.Account{}
	o.QueryTable("account").Filter("user_id", user_id).One(&account)
	blocked_old := models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
		Filter("account", account.Id).
		OrderBy("-create_date").
		Limit(1).
		One(&blocked_old)
	if blocked_old.Id == 0 {
		blocked_old.CurrentBalance = 0
	}
	for_mula := models.Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id", account.Id).One(&for_mula)
	if err_for != nil {
		o.Rollback()
		return err_for
	}
	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: value * for_mula.TeamReturnRate,
		CurrentOutlay:  0,
		CurrentBalance: blocked_old.CurrentBalance + (value * for_mula.TeamReturnRate) - 0,
		OpeningBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "每日团队收益",
		TxId:           tx_id_blo_d,
		Account:        account.Id,
	}

	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
		o.Rollback()
		return err
	}

	//更新生态仓库属性
	_, err_up := o.QueryTable("account").Filter("id", account.Id).Update(orm.Params{"bocked_balance": blocked_new.CurrentBalance})
	if err_up != nil {
		o.Rollback()
		return err_up
	}

	// 超级节点表生成与更新
	super_peer_table := models.SuperPeerTable{}
	err_super := o.QueryTable("super_peer_table").Filter("user_id", user_id).One(&super_peer_table)
	if err_super != nil {
		o.Rollback()
		return err_super
	}
	coin := super_peer_table.CoinNumber + (value * for_mula.TeamReturnRate) - 0
	if coin < 0 {
		coin = 0
	}
	_, err_super_up := o.QueryTable("super_peer_table").Filter("user_id", user_id).Update(orm.Params{"coin_number": coin})
	if err_super_up != nil {
		o.Rollback()
		return err_super_up
	}
	models.NetIncome += value
	return nil
}

// 超级节点的分红
func AddFormulaABonus(user_id string, abonus float64) {
	o := models.NewOrm()
	o.Begin()

	//任务表 USDD  铸币记录
	tx_id_blo_d := utils.Shengchengstr("超级节点分红", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		State: "true",
		TxId:  tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		o.Rollback()
		AddFormulaABonus(user_id, abonus)
		return
	}

	account := models.Account{}
	o.QueryTable("account").Filter("user_id", user_id).One(&account)
	blocked_old := models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
		Filter("account", account.Id).
		OrderBy("-create_date").
		Limit(1).
		One(&blocked_old)
	if blocked_old.Id == 0 {
		blocked_old.CurrentBalance = 0
	}
	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: 0,
		CurrentOutlay:  abonus,
		CurrentBalance: blocked_old.CurrentBalance - abonus,
		OpeningBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "超级节点分红",
		TxId:           tx_id_blo_d,
		Account:        account.Id,
	}

	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	_, err := o.Insert(&blocked_new)
	if err != nil {
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
	o.Read(&formula)
	blocked_yestoday := models.BlockedDetail{}
	err_raw := o.Raw(
		"select * from blocked_detail where user_id=? and create_date<=? order by create_date desc limit 1",
		user_id,
		time.Now().AddDate(0, 0, -1).Format("2006-01-02 ")+"23:59:59").
		QueryRow(&blocked_yestoday)
	if err_raw != nil {
		if err_raw.Error() == "<QuerySeter> no row found" {
			return nil
		} else {
			return err_raw
		}
	}
	if blocked_yestoday.Id == 0 {
		blocked_yestoday.CurrentBalance = 0
	}
	abonus := formula.HoldReturnRate * blocked_yestoday.CurrentBalance

	//任务表 USDD  铸币记录
	tx_id_blo_d := utils.Shengchengstr("每日释放", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		State: "true",
		TxId:  tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		return errtxid_blo
	}

	blocked_old := models.BlockedDetail{}
	o.QueryTable("blocked_detail").
		Filter("user_id", user_id).
		Filter("account", account.Id).
		OrderBy("-create_date").
		Limit(1).
		One(&blocked_old)
	if blocked_old.Id == 0 {
		blocked_old.CurrentBalance = 0
	}
	blocked_new := models.BlockedDetail{
		UserId:         user_id,
		CurrentRevenue: 0,
		CurrentOutlay:  abonus,
		CurrentBalance: blocked_old.CurrentBalance - abonus,
		OpeningBalance: blocked_old.CurrentBalance,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        "每日释放",
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
	models.NetIncome += abonus
	return nil
}

// 远端连接  -  给定分红收益  释放通用
func PingAddWalletCoin(user_id string, abonus float64) error {
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
func PingSelectTforNumber(user_id string) (float64, error) {
	user := models.User{
		UserId: user_id,
	}
	b, token := generateToken(user)
	if b != true {
		return 0.0, errors.New("err")
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
		return 0.0, err1
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", token)

	//处理返回结果
	response, errdo := client.Do(req)
	if errdo != nil {
		return 0.0, errdo
	}

	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return 0.0, err_read
	}
	values := models.Data_wallet{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return 0.0, errors.New("钱包金额操作失败!")
	} else if values.Code != 200 {
		return 0.0, errors.New(values.Msg)
	}
	response.Body.Close()
	return values.Data["balance"].(float64), nil
}

// 返回超级节点的等级
func ReturnSuperPeerLevel(user_id string) (string, error) {
	s_f_t := []models.SuperForceTable{}
	models.NewOrm().QueryTable("super_force_table").All(&s_f_t)
	tfor_number, err_tfor := PingSelectTforNumber(user_id)
	if err_tfor != nil {
		return "", err_tfor
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
		return s_f_t[index[len(index)-1]].Level, nil
	}
	return "", nil
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
func HandlerMap(m map[string][]string) {
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
				AddFormulaABonus(v, tfor_some/float64(len(vv)))
			}
		}
	}
	if len(err_m) != 0 {
		HandlerMap(err_m)
	}
}

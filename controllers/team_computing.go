package controllers

import (
	"ecology/common"
	"ecology/consul"
	"ecology/models"
	"ecology/utils"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func DailyDividendAndRelease() {
	o := models.NewOrm()
	users := []models.User{}
	o.QueryTable("user").All(&users)
	Producer(users)
	//TODO logs 打印成功信息
}

func Producer(users []models.User) {
	error_users := []models.User{}
	for _, v := range users {
		if err := Worker(v); err != nil {
			error_users = append(error_users, v)
			//TODO  写日志,告知管理者失败的原因
		}
	}
	if len(error_users) > 0 {
		Producer(error_users)
	}
}

func Worker(user models.User) error {
	o := models.NewOrm()
	o.Begin()
	user_current_layer := []models.User{}
	coins := []float64{}
	o.QueryTable("user").Filter("father_id", user.UserId).All(&user_current_layer)
	if len(user_current_layer) > 0 {
		for _, v := range user_current_layer {
			// 获取用户teams
			team_user, err := GetTeams(v)
			if err != nil {
				return err
			}
			// 去处理这些数据
			coin, err_handler := HandlerOperation(team_user)
			if err_handler != nil {
				return err_handler
			}
			coins = append(coins, coin)
		}
	}
	err_sort_a_r := SortABonusRelease(o, coins, user.UserId)
	if err_sort_a_r != nil {
		o.Rollback()
		return err_sort_a_r
	}
	o.Commit()
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
			return 0, err_acc
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

// 去掉最大的 并更新分红和释放
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
		CurrentRevenue: blocked_old.CurrentOutlay,
		CurrentOutlay:  value * for_mula.TeamReturnRate,
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
	return nil
}

// 超级节点的分红
func AddFormulaABonus(user_id string) error {
	s_f_t := []models.SuperForceTable{}
	models.NewOrm().QueryTable("super_force_table").All(&s_f_t)
	s_p_t := models.SuperPeerTable{}
	models.NewOrm().QueryTable("super_peer_table").Filter("user_id", user_id).One(&s_p_t)

	for i := 0; i < len(s_f_t); i++ {
		for j := 1; j < len(s_f_t)-1; j++ {
			if s_f_t[i].CoinNumberRule > s_f_t[j].CoinNumberRule {
				s_f_t[i].CoinNumberRule, s_f_t[j].CoinNumberRule = s_f_t[j].CoinNumberRule, s_f_t[i].CoinNumberRule
			}
		}
	}
	index := []int{}
	for i, v := range s_f_t {
		if s_p_t.CoinNumber > float64(v.CoinNumberRule) {
			index = append(index, i)
		}
	}
	if len(index) > 0 {
		abonus := s_f_t[index[len(index)-1]].BonusCalculation * s_p_t.CoinNumber
		// 调老罗接口
		if err := PingAddWalletCoin(user_id, abonus); err != nil {
			return err
		}
	}
	return nil
}

// 远端连接  -  给定分红收益
func PingAddWalletCoin(user_id string, abonus float64) error {
	user := models.User{
		UserId: user_id,
	}
	b, user_str := generateToken(user)
	if b != true {
		return errors.New("err")
	}
	client := &http.Client{}
	//生成要访问的url
	url := beego.AppConfig.String("api::apiurl_abonus")
	//提交请求
	reqest, errnr := http.NewRequest("POST", url, nil)

	//增加header选项
	reqest.Header.Add("Authorization", user_str)
	reqest.Form.Add("coin", strconv.FormatFloat(abonus, 'f', -1, 64))
	reqest.Form.Add("coin_type", "TFOR")

	if errnr != nil {
		return errnr
	}
	//处理返回结果
	response, errdo := client.Do(reqest)
	defer response.Body.Close()
	if errdo != nil {
		return errdo
	}
	bys, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		return err_read
	}
	values := common.ResponseData{}
	err := json.Unmarshal(bys, &values)
	if err != nil {
		return err_read
	}
	if values.Code != 200 {
		return errors.New("err")
	}
	return nil
}

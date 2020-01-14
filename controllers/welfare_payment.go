package controllers

import (
	"ecology/models"
	"ecology/utils"

	"github.com/astaxie/beego/orm"
	"github.com/robfig/cron"

	"time"
)

var c = cron.New()

func WelfarePayment() {
	o := models.NewOrm()
	o.Begin()
	task := models.DailyDividendTasks{
		Time:  time.Now().Format("2006-01-02"),
		State: "false",
	}
	if _, err := o.Insert(&task); err != nil {
		//TODO logs
		c.AddFunc("0 0/5 * * * ? *", WelfarePayment)
		c.Start()
		return
	}
	users := []models.User{}
	_, err_read_user := o.QueryTable("user").All(&users)
	if err_read_user != nil {
		o.Rollback()
		c.AddFunc("0 0/5 * * * ? *", WelfarePayment)
		c.Start()
		return
	}
	for _, vuser := range users {
		accounts := []models.Account{}
		int, err := o.QueryTable("account").Filter("user_id", vuser.UserId).All(&accounts)
		if err != nil {
			//TODO logs
			o.Rollback()
			c.AddFunc("0 0/5 * * * ? *", WelfarePayment)
			c.Start()
			return
		}
		if int > 0 {
			for _, vaccount := range accounts {
				if err := ABonus(o, vuser.UserId, vaccount); err != nil {
					//TODO logs
					o.Rollback()
					c.AddFunc("0 0/5 * * * ? *", WelfarePayment)
					c.Start()
					return
				}
			}
			AddFormulaABonus(vuser.UserId, 0.0)
		}
	}

	o.Commit()
}

func ABonus(o orm.Ormer, user_id string, account models.Account) error {
	formula := models.Formula{}
	err_formula := o.QueryTable("formula").Filter("ecology_id", account.Id).One(&formula)
	if err_formula != nil {
		//TODO logs
		return err_formula
	}

	//任务表 USDD  铸币记录
	tx_id_blo_d := utils.Shengchengstr("释放记录", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		OrderState: false,
		TxId:       tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		//TODO logs
		return errtxid_blo
	}

	//铸币交易记录
	err_blo_d := FindLimitOneAndSaveBlo_dAbonus(o, user_id, "每日释放奖励", tx_id_blo_d, account.Id)
	if err_blo_d != nil {
		//TODO logs
		return err_blo_d
	}

	return nil
}

// 铸币表释放 以及 超级节点表的 同步数据
func FindLimitOneAndSaveBlo_dAbonus(o orm.Ormer, user_id, comment, tx_id string, account_id int) error {
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
		return err_for
	}
	shifang := (blocked_old.CurrentBalance * for_mula.HoldReturnRate) + (blocked_old.CurrentBalance * for_mula.RecommendReturnRate) + (blocked_old.CurrentBalance * for_mula.TeamReturnRate)
	blocked_new := models.BlockedDetail{
		Id:             0,
		UserId:         user_id,
		CurrentRevenue: 0,
		CurrentOutlay:  shifang,
		OpeningBalance: blocked_old.CurrentBalance,
		CurrentBalance: blocked_old.CurrentBalance - shifang,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
	}
	if blocked_new.CurrentBalance < 0 {
		blocked_new.CurrentBalance = 0
	}
	blocked_new.CurrentBalance = blocked_old.CurrentBalance + 0*for_mula.ReturnMultiple - shifang
	_, err := o.Insert(&blocked_new)
	if err != nil {
		return err
	}

	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"state": "true"})
	if err_txid != nil {
		return err_txid
	}

	_, err_up := o.QueryTable("account").Filter("id", account_id).Update(orm.Params{"bocked_balance": blocked_new.CurrentBalance})
	if err_up != nil {
		return err_up
	}

	super_peer_table := models.SuperPeerTable{}
	err_super := o.QueryTable("super_peer_table").Filter("user_id", user_id).One(&super_peer_table)
	if err_super != nil {
		return err_super
	}
	coin := super_peer_table.CoinNumber + (0 * for_mula.ReturnMultiple) - shifang
	if coin < 0 {
		coin = 0
	}
	_, err_super_up := o.QueryTable("super_peer_table").Filter("user_id", user_id).Update(orm.Params{"coin_number": coin})
	if err_super_up != nil {
		return err_super_up
	}
	return nil
}

//// 超级节点的分红
//func AddFormulaABonus(user_id string) error {
//	s_f_t := []models.SuperForceTable{}
//	models.NewOrm().QueryTable("super_force_table").All(&s_f_t)
//	s_p_t := models.SuperPeerTable{}
//	models.NewOrm().QueryTable("super_peer_table").Filter("user_id",user_id).One(&s_p_t)
//
//	for i := 0; i < len(s_f_t); i++ {
//		for j := 1; j < len(s_f_t)-1; j++ {
//			if s_f_t[i].CoinNumberRule > s_f_t[j].CoinNumberRule {
//				s_f_t[i].CoinNumberRule, s_f_t[j].CoinNumberRule = s_f_t[j].CoinNumberRule, s_f_t[i].CoinNumberRule
//			}
//		}
//	}
//	index := []int{}
//	for i, v := range s_f_t {
//		if s_p_t.CoinNumber > float64(v.CoinNumberRule) {
//			index = append(index, i)
//		}
//	}
//	if len(index) > 0 {
//		abonus := s_f_t[index[len(index)-1]].BonusCalculation * s_p_t.CoinNumber
//		// 调老罗接口
//		if err := PingAddWalletCoin(user_id,abonus);err!=nil{
//			return err
//		}
//	}
//	return nil
//}
//
//// 远端连接  -  给定分红收益
//func PingAddWalletCoin(user_id string,abonus float64) error{
//	user := models.User{
//		UserId:   user_id,
//	}
//	b,user_str := generateToken(user)
//	if b != true{
//		return errors.New("err")
//	}
//	client := &http.Client{}
//	//生成要访问的url
//	url := beego.AppConfig.String("api::apiurl_abonus")
//	//提交请求
//	reqest, errnr := http.NewRequest("POST", url, nil)
//
//	//增加header选项
//	reqest.Header.Add("Authorization", user_str)
//	reqest.Form.Add("coin",strconv.FormatFloat(abonus, 'E', -1, 64))
//	reqest.Form.Add("coin_type","TFOR")
//
//	if errnr != nil {
//		return errnr
//	}
//	//处理返回结果
//	response, errdo := client.Do(reqest)
//	defer response.Body.Close()
//	if errdo != nil {
//		return errdo
//	}
//	bys,err_read := ioutil.ReadAll(response.Body)
//	if err_read!=nil {
//		return err_read
//	}
//	values := common.ResponseData{}
//	err := json.Unmarshal(bys,&values)
//	if err!=nil {
//		return err_read
//	}
//	if values.Code != 200{
//		return errors.New("err")
//	}
//	return nil
//}

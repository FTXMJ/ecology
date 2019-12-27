package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

// 根据等级 进行算力的更新
func JudgeLevel(o orm.Ormer, user_id, level string, formula *Formula) error {
	//if JudgeLevelFor_wh_mx(o, user_id, level) == nil {
	//	force := ForceTable{}
	//	err := NewOrm().QueryTable("force_table").Filter("level", level).One(&force)
	//	if err != nil {
	//		return err
	//	}
	//	formula.Level = force.Level
	//	formula.LowHold = force.LowHold
	//	formula.HighHold = force.HighHold
	//	formula.ReturnMultiple = force.ReturnMultiple
	//	formula.HoldReturnRate = force.HoldReturnRate
	//	formula.RecommendReturnRate = force.RecommendReturnRate
	//	formula.TeamReturnRate = force.TeamReturnRate
	//	return nil
	//}
	if PanDuanLevel(user_id, level) == true {
		force := ForceTable{}
		err := NewOrm().QueryTable("force_table").Filter("level", level).One(&force)
		if err != nil {
			return err
		}
		formula.Level = force.Level
		formula.LowHold = force.LowHold
		formula.HighHold = force.HighHold
		formula.ReturnMultiple = force.ReturnMultiple
		formula.HoldReturnRate = force.HoldReturnRate
		formula.RecommendReturnRate = force.RecommendReturnRate
		formula.TeamReturnRate = force.TeamReturnRate
		return nil
	}
	return errors.New("升级条件不满足!")
}

// 验证当前用户的父亲用户是否达标  可以升级的条件
func JudgeLevelFor_wh_mx(o orm.Ormer, user_id, level string) error {
	if level == "代言人" {
		var u_ser User
		o.QueryTable("user").Filter("user_id", user_id).One(&u_ser)
		var u_sers []User
		o.QueryTable("user").Filter("father_id", u_ser.FatherId).All(&u_sers)
		count := 0
		if len(u_sers) < 1 {
			return nil
		}
		for _, v := range u_sers {
			account := Account{}
			o.QueryTable("account").Filter("user_id", v.UserId).Filter("level", "代言人").One(&account)
			if account.Id > 0 {
				count += 1
			}
			if count >= 2 {
				return UpdateLevel(o, u_ser.FatherId, "网红")
			}
		}
		return nil
	} else if level == "网红" {
		var u_ser User
		o.QueryTable("user").Filter("user_id", user_id).One(&u_ser)
		var u_sers []User
		o.QueryTable("user").Filter("father_id", u_ser.FatherId).All(&u_sers)
		count := 0
		if len(u_sers) < 1 {
			return nil
		}
		for _, v := range u_sers {
			account := Account{}
			o.QueryTable("account").Filter("user_id", v.UserId).Filter("level", "网红").One(&account)
			if account.Id > 0 {
				count += 1
			}
			if count >= 3 {
				return UpdateLevel(o, u_ser.FatherId, "明星")
			}
		}
		return nil
	}
	return nil
}

// 升级父亲用户的 生态等级
func UpdateLevel(o orm.Ormer, father_account_id, level string) error {
	if level == "网红" {
		force := ForceTable{}
		o.QueryTable("force_table").Filter("level", level).One(&force)
		if _, err := o.QueryTable("account").Filter("user_id", father_account_id).Filter("bocked_balance__gte", force.LowHold).Update(orm.Params{"level": level}); err == nil {
			account := Account{}
			o.QueryTable("account").Filter("user_id", father_account_id).One(&account)
			_, err := NewOrm().QueryTable("formula").Filter("ecology_id", account.Id).Update(
				orm.Params{
					"level":                 force.Level,
					"low_hold":              force.LowHold,
					"high_hold":             force.HighHold,
					"return_multiple":       force.ReturnMultiple,
					"hold_return_rate":      force.HoldReturnRate,
					"recommend_return_rate": force.RecommendReturnRate,
					"team_return_rate":      force.TeamReturnRate,
				})
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	} else if level == "明星" {
		force := ForceTable{}
		o.QueryTable("force_table").Filter("level", level).One(&force)
		if _, err := o.QueryTable("account").Filter("user_id", father_account_id).Filter("bocked_balance__gte", force.LowHold).Update(orm.Params{"level": level}); err == nil {
			account := Account{}
			o.QueryTable("account").Filter("user_id", father_account_id).One(&account)
			_, err := NewOrm().QueryTable("formula").Filter("ecology_id", account.Id).Update(
				orm.Params{
					"level":                 force.Level,
					"low_hold":              force.LowHold,
					"high_hold":             force.HighHold,
					"return_multiple":       force.ReturnMultiple,
					"hold_return_rate":      force.HoldReturnRate,
					"recommend_return_rate": force.RecommendReturnRate,
					"team_return_rate":      force.TeamReturnRate,
				})
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}
	return nil
}

// 如果升级等级是　？？　就需要判断是否　符合升级条件                 ------          每个人只有一改生态仓库
func PanDuanLevel(user_id, level string) bool {
	if level == "侯爵" || level == "公爵" {
		o := NewOrm()
		sun_users := []User{}
		sun_accounts := []Account{}
		o.Raw("select * from user where father_id=?", user_id).QueryRows(&sun_users)
		for _, v := range sun_users {
			sun_account := Account{}
			o.Raw("select * from account where user_id=?", v.UserId).QueryRow(&sun_account)
			sun_accounts = append(sun_accounts, sun_account)
		}
		switch level {
		case "侯爵":
			l := 0
			for _, v := range sun_accounts {
				if v.Level == "伯爵" {
					l++
				}
			}
			if l >= 2 {
				return true
			}
			return false
		case "公爵":
			l := 0
			for _, v := range sun_accounts {
				if v.Level == "侯爵" {
					l++
				}
			}
			if l >= 3 {
				return true
			}
			return false
		}
	}
	return true
}

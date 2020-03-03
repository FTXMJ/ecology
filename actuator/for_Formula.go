package actuator

import (
	"ecology/models"
	"github.com/jinzhu/gorm"

	"errors"
)

// 根据等级 进行算力的更新
func JudgeLevel(o *gorm.DB, user_id, level string, formula *models.Formula) error {
	if PanDuanLevel(o, user_id, level) == true {
		force := models.ForceTable{}
		err := o.Table("force_table").Where("level = ?", level).Find(&force)
		if err.Error != nil {
			return err.Error
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
func JudgeLevelFor_wh_mx(o *gorm.DB, user_id, level string) error {
	if level == "代言人" {
		var u_ser models.User
		o.Table("user").Where("user_id = ?", user_id).First(&u_ser)
		var u_sers = make([]models.User, 0)
		o.Table("user").Where("father_id = ?", u_ser.FatherId).Find(&u_sers)
		count := 0
		if len(u_sers) < 1 {
			return nil
		}
		for _, v := range u_sers {
			account := models.Account{}
			o.Table("account").Where("user_id = ?", v.UserId).Where("level = ?", "代言人").First(&account)
			if account.Id > 0 {
				count += 1
			}
			if count >= 2 {
				return UpdateLevel(o, u_ser.FatherId, "网红")
			}
		}
		return nil
	} else if level == "网红" {
		var u_ser models.User
		o.Table("user").Where("user_id = ?", user_id).First(&u_ser)
		var u_sers = make([]models.User, 0)
		o.Table("user").Where("father_id = ?", u_ser.FatherId).Find(&u_sers)
		count := 0
		if len(u_sers) < 1 {
			return nil
		}
		for _, v := range u_sers {
			account := models.Account{}
			o.Table("account").Where("user_id = ?", v.UserId).Where("level = ?", "网红").First(&account)
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
func UpdateLevel(o *gorm.DB, father_account_id, level string) error {
	if level == "网红" {
		force := models.ForceTable{}
		o.Table("force_table").Where("level = ?", level).First(&force)
		if er := o.Model(&models.Account{}).Where("user_id = ?", father_account_id).Where("bocked_balance__gte = ?", force.LowHold).Update("level", level); er.Error == nil {
			account := models.Account{}
			o.Table("account").Where("user_id = ?", father_account_id).First(&account)
			err := o.Model(&models.Formula{}).Where("ecology_id = ?", account.Id).Update(
				map[string]interface{}{
					"level":                 force.Level,
					"low_hold":              force.LowHold,
					"high_hold":             force.HighHold,
					"return_multiple":       force.ReturnMultiple,
					"hold_return_rate":      force.HoldReturnRate,
					"recommend_return_rate": force.RecommendReturnRate,
					"team_return_rate":      force.TeamReturnRate,
				})
			if err.Error != nil {
				return err.Error
			}
			return nil
		}
		return nil
	} else if level == "明星" {
		force := models.ForceTable{}
		o.Table("force_table").Where("level = ?", level).First(&force)
		if er := o.Table("account").Where("user_id = ?", father_account_id).Where("bocked_balance__gte = ?", force.LowHold).Update("level", level); er.Error == nil {
			account := models.Account{}
			o.Table("account").Where("user_id", father_account_id).First(&account)
			err := o.Table("formula").Where("ecology_id = ?", account.Id).Update(
				map[string]interface{}{
					"level":                 force.Level,
					"low_hold":              force.LowHold,
					"high_hold":             force.HighHold,
					"return_multiple":       force.ReturnMultiple,
					"hold_return_rate":      force.HoldReturnRate,
					"recommend_return_rate": force.RecommendReturnRate,
					"team_return_rate":      force.TeamReturnRate,
				})
			if err.Error != nil {
				return err.Error
			}
			return nil
		}
		return nil
	}
	return nil
}

// 如果升级等级是　？？　就需要判断是否　符合升级条件                 ------          每个人只有一改生态仓库
func PanDuanLevel(o *gorm.DB, user_id, level string) bool {
	if level == "侯爵" || level == "公爵" {
		sun_users := make([]models.User, 0)
		sun_accounts := make([]models.Account, 0)
		o.Raw("select * from user where father_id=?", user_id).Find(&sun_users)
		for _, v := range sun_users {
			sun_account := models.Account{}
			o.Raw("select * from account where user_id=?", v.UserId).Find(&sun_account)
			sun_accounts = append(sun_accounts, sun_account)
		}
		switch level {
		case "侯爵":
			l := 0
			for _, v := range sun_accounts {
				if v.Level == "伯爵" {
					l++
				}
				if v.Level == "公爵" {
					l++
				}
				if v.Level == "侯爵" {
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
				if v.Level == "公爵" {
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

// 返回超级节点的等级
func ReturnSuperPeerLevel(o *gorm.DB, user_id string) (time, level string, tfor float64, err error) {
	s_f_t := make([]models.SuperForceTable, 0)
	o.Table("super_force_table").Find(&s_f_t)
	up_time, tfor_number, err_tfor := PingSelectTforNumber(user_id)
	if err_tfor != nil {
		return "", "", 0.0, err_tfor
	}

	for i := 0; i < len(s_f_t); i++ {
		for j := i + 1; j < len(s_f_t)-1; j++ {
			if s_f_t[i].CoinNumberRule > s_f_t[j].CoinNumberRule {
				s_f_t[i], s_f_t[j] = s_f_t[j], s_f_t[i]
			}
		}
	}
	index := make([]int, 0)
	for i, v := range s_f_t {
		if tfor_number >= float64(v.CoinNumberRule) {
			index = append(index, i)
		}
	}
	if len(index) > 0 {
		return up_time, s_f_t[index[len(index)-1]].Level, tfor_number, nil
	}
	return "", "", 0.0, err_tfor
}

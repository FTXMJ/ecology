package actuator

import (
	"ecology/logs"
	"ecology/models"
	"github.com/astaxie/beego/orm"
	"time"
)

func TheWheel(o orm.Ormer, user_id string, acc models.Account, index *models.Ecology_index_obj) error {
	formula_index := models.Formula{EcologyId: acc.Id}
	errfor := o.Read(&formula_index, "ecology_id")
	if errfor != nil {
		logs.Log.Error(errfor)
		return errfor
	}
	f := models.Formulaindex{
		Id:             acc.Id,
		Level:          acc.Level,
		BockedBalance:  acc.BockedBalance,
		Balance:        acc.Balance,
		LowHold:        formula_index.LowHold,
		HighHold:       formula_index.HighHold,
		ReturnMultiple: formula_index.ReturnMultiple,
		HoldReturnRate: formula_index.HoldReturnRate * acc.Balance,
	}
	t := time.Now().Format("2006-01-02") + " 00:00:00"
	zhitui, err := RecommendReturnRate(user_id, t)
	if err != nil {
		logs.Log.Error("计算用户当前直推收益出错!")
		return err
	}
	f.RecommendReturnRate = zhitui
	team_coins, err_team := IndexTeamABouns(o, user_id)
	if err_team != nil {
		logs.Log.Error(err_team.Error())
		return err_team
	}
	f.TeamReturnRate = team_coins * formula_index.TeamReturnRate
	to_day_rate := zhitui + f.TeamReturnRate + f.HoldReturnRate
	f.ToDayRate = to_day_rate
	index.Ecological_poject = append(index.Ecological_poject, f)
	return nil
}

// 查看用户团队收益 首页查看
func IndexTeamABouns(o orm.Ormer, user_id string) (float64, error) {
	coins := []float64{}
	user_current_layer := []models.User{}
	// 团队收益　开始
	o.QueryTable("user").Filter("father_id", user_id).All(&user_current_layer)
	if len(user_current_layer) > 0 {
		for _, v := range user_current_layer {
			if user_id != v.UserId {
				// 获取用户teams
				team_user, err := GetTeams(v)
				if err != nil {
					if err.Error() != "用户未激活或被拉入黑名单" {
						return 0, err
					}
				}
				if len(team_user) > 0 {
					// 去处理这些数据 // 处理器，计算所有用户的收益  并发布任务和 分红记录
					coin, err_handler := HandlerOperation(team_user, user_id)
					if err_handler != nil {
						return 0, err_handler
					}
					coins = append(coins, coin)
				}
			}
		}
	}
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
	// 团队收益　结束
	return value, nil
}

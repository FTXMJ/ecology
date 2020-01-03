package controllers

import (
	"ecology/common"
	"ecology/models"
	"errors"
	"github.com/astaxie/beego"
)

type FirstController struct {
	beego.Controller
}

// @Tags 心跳检测
// @Accept  json
// @Produce json
// @Success 200
// @router /check [GET]
func (this *FirstController) Check() {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	data = common.NewResponse(nil)
	return
}

// 计算团队收益
func SumTeamProfit(user_id string) (float64, error) {
	o := models.NewOrm()
	user_team := []models.User{}
	_, err_raw := o.Raw("select * from user where father_id=?", user_id).QueryRows(&user_team)
	if err_raw != nil {
		if err_raw.Error() != "<QuerySeter> no row found" {
			return 0.0, err_raw
		}
	}
	coins := []float64{}
	coin_number := 0.0
	if len(user_team) > 0 {
		for _, v := range user_team {
			// 获取用户teams
			team_user, err := GetTeams(v)
			if err != nil {
				return 0.0, errors.New("查询用户团队直推成员时出错,请重试")
			}
			// 去处理这些数据 // 处理器，计算所有用户的收益  并发布任务和 分红记录
			coin, err_handler := HandlerOperation(team_user, user_id)
			if err_handler != nil {
				return 0.0, errors.New("查询用户团队成员时出错,请重试")
			}
			coins = append(coins, coin)
		}
		for i := 0; i < len(coins)-1; i++ {
			for j := i + 1; j < len(coins); j++ {
				if coins[i] > coins[j] {
					coins[i], coins[j] = coins[j], coins[i]
				}
			}
		}
		for i := 0; i < len(coins)-1; i++ {
			coin_number += coins[i]
		}
	}
	var account = models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	var formula = models.Formula{
		EcologyId: account.Id,
	}
	o.Read(&formula, "ecology_id")
	coin_number = coin_number * formula.TeamReturnRate
	return coin_number, nil
}

package controllers

import (
	"ecology/models"
	"errors"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// 查看节点收益流水
func SelectPeerABounsList(page models.Page, user_name string) ([]models.PeerAbouns, models.Page, error) {
	peer_a_bouns := []models.TxIdList{}
	o := models.NewOrm()
	switch user_name {
	case "":
		o.Raw("select * from tx_id_list where comment=? order by create_time", "节点分红").QueryRows(&peer_a_bouns)
	default:
		users := []models.User{}
		_, err_1 := o.Raw("select * from user where user_name=?", user_name).QueryRows(&users)
		if err_1 != nil || len(users) < 1 {
			return []models.PeerAbouns{}, page, err_1
		}
		for _, v := range users {
			_, err := o.Raw("select * from tx_id_list where user_id=? and where comment=?", v.UserId, "节点分红").QueryRows(&peer_a_bouns)
			if err != nil || len(peer_a_bouns) < 1 {
				return []models.PeerAbouns{}, page, err
			}
		}
	}
	if len(peer_a_bouns) < 1 {
		return []models.PeerAbouns{}, page, errors.New("没有相关数据!")
	}
	models.QuickSortPeerABouns(peer_a_bouns, 0, len(peer_a_bouns)-1)

	page.Count = len(peer_a_bouns)
	if page.PageSize < 10 {
		page.PageSize = 10
	}
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	start := (page.CurrentPage - 1) * page.PageSize
	end := start + page.PageSize

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
			o.Read(&u)
			_, level, tfors, err_tfor := ReturnSuperPeerLevel(v.UserId)
			if err_tfor != nil {
				return []models.PeerAbouns{}, page, err_tfor
			}
			p := models.PeerAbouns{
				Id:       v.Id,
				UserName: u.UserName,
				Level:    level,
				Tfors:    tfors,
				Time:     v.CreateTime,
			}
			listle = append(listle, p)
		}
	}
	return listle, page, nil
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

// 查看用户有史以来所有的收益
func AddAllSum(o orm.Ormer, user_id string) float64 {
	var blos orm.ParamsList
	//o.Raw("select * from blocked_detail where user_id=? and comment!=?",user_id,"直推收益").QueryRows(&blos)
	o.Raw("select sum(current_outlay) from blocked_detail where user_id=? and comment!=?", user_id, "直推收益").ValuesFlat(&blos)
	var zhitui float64
	if len(blos) > 0 && blos[0] != nil {
		z, _ := strconv.ParseFloat(blos[0].(string), 64)
		zhitui = z
	}
	return zhitui
}

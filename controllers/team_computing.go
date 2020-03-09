package controllers

import (
	db "ecology/db"
	"ecology/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

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

var Ecology_orm orm.Ormer

// @Tags 测试每日释放
// @Accept  json
// @Produce json
// @Success 200
// @router /admin/test_mrsf [GET]
func (this *Test) Mrsf() {
	if Ecology_orm == nil {
		Ecology_orm = db.NewEcologyOrm()
	}
	users := []models.User{}
	Ecology_orm.QueryTable("user").All(&users)

	for _, v := range users {
		DailyDividendAndReleaseTest(v)
	}
	PeerABounsHandler(users)
}

// 释放 - 团队 - 直推   收益的给定
func DailyDividendAndReleaseTest(user models.User) {
	if Ecology_orm == nil {
		Ecology_orm = db.NewEcologyOrm()
	}
	o := Ecology_orm

	//    每日释放___and___团队收益___and___直推收益
	ProducerEcology(o, user, "") // 返回错误的用户名单

	//    给失败的用户　添加失败的任务记录表
	CreateErrUserTxList(user)
}

// 节点用户的收益分发 - 全网总收益回归正常
func PeerABounsHandler(users []models.User) {
	if Ecology_orm == nil {
		Ecology_orm = db.NewEcologyOrm()
	}
	o := Ecology_orm
	// 超级节点的分红
	in_fo := info{}
	err_peer := ProducerPeer(users, &in_fo, "")
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

// 执行错误的释放用户
func ErrorUserMrsf(users_str []string, order_id string) {
	if Ecology_orm == nil {
		Ecology_orm = db.NewEcologyOrm()
	}
	o := Ecology_orm
	users := []models.User{}
	for _, v := range users_str {
		u := models.User{UserId: v}
		o.Read(&u)
		users = append(users, u)
	}

	for _, v := range users {
		//    每日释放___and___团队收益___and___直推收益
		ProducerEcology(o, v, order_id) // 返回错误的用户名单

		//    给失败的用户　添加失败的任务记录表
		CreateErrUserTxList(v)
	}
}

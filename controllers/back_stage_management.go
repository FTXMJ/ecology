package controllers

import (
	"ecology1/common"
	"ecology1/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

// 后台管理
type BackStageManagement struct {
	beego.Controller
}

// @Tags 算力表显示   后台操作　or 用户查看　都可
// @Accept  json
// @Produce json
// @Success 200
// @router /show_formula_list____算力表显示___后台操作__or__用户查看__都可 [Post]
func (this *BackStageManagement) ShowFormulaList() {
	var (
		data       *common.ResponseData
		o          = models.NewOrm()
		force_list []models.ForceTable
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	o.QueryTable("force_table").All(&force_list)
	models.QuickSortForce(force_list, 0, len(force_list)-1)

	data = common.NewResponse(force_list)
	return
}

// @Tags 算力表信息修改
// @Accept  json
// @Produce json
// @Param force_id query string true "算力信息的id  如果是删除多条，　用　, 逗号隔开"
// @Param action query string true "具体的操作　　delete=删除　update=更新　insert=新增    更新必须传这条数据的全部数据回来，id 除外，和新增一样，删除只要ｉｄ就够"
// @Param levelstr query string true "等级名字"
// @Param low_hold query string true "低位"
// @Param high_hold query string true "高位"
// @Param return_multiple query string true "杠杆"
// @Param hold_return_rate query string true "本金自由算力"
// @Param recommend_return_rate query string true "直推算力"
// @Param team_return_rate query string true "动态算力"
// @Success 200
// @router /operation_formula_list____算力表信息修改 [Post]
func (this *BackStageManagement) OperationFormulaList() {
	var (
		data                      *common.ResponseData
		o                         = models.NewOrm()
		force_id                  = this.GetString("force_id")
		action                    = this.GetString("action")
		levelstr                  = this.GetString("levelstr")
		low_hold, _               = this.GetInt("low_hold")
		high_hold, _              = this.GetInt("high_hold")
		return_multiple_str       = this.GetString("return_multiple")
		hold_return_rate_str      = this.GetString("hold_return_rate")
		recommend_return_rate_str = this.GetString("recommend_return_rate")
		team_return_rate_str      = this.GetString("team_return_rate")
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	return_multiple, _ := strconv.ParseFloat(return_multiple_str, 64)
	hold_return_rate, _ := strconv.ParseFloat(hold_return_rate_str, 64)
	recommend_return_rate, _ := strconv.ParseFloat(recommend_return_rate_str, 64)
	team_return_rate, _ := strconv.ParseFloat(team_return_rate_str, 64)

	switch action {
	case "delete":
		id_strs := strings.Split(force_id, ",")
		for _, v := range id_strs {
			id, _ := strconv.Atoi(v)
			_, err := o.QueryTable("force_table").Filter("id", id).Delete()
			if err != nil {
				data = common.NewErrorResponse(500)
				return
			}
			data = common.NewResponse(nil)
			return
		}
	case "update":
		id, _ := strconv.Atoi(force_id)
		_, err := o.QueryTable("force_table").
			Filter("id", id).
			Update(orm.
			Params{"level": levelstr, "low_hold": low_hold, "high_hold": high_hold, "return_multiple": return_multiple, "hold_return_rate": hold_return_rate, "recommend_return_rate": recommend_return_rate, "": team_return_rate})
		if err != nil {
			data = common.NewErrorResponse(500)
			return
		}
		data = common.NewResponse(nil)
		return
	case "insert":
		force := models.ForceTable{
			Level:               levelstr,
			LowHold:             low_hold,
			HighHold:            high_hold,
			ReturnMultiple:      return_multiple,
			HoldReturnRate:      hold_return_rate,
			RecommendReturnRate: recommend_return_rate,
			TeamReturnRate:      team_return_rate,
		}
		_, err := o.Insert(&force)
		if err != nil {
			data = common.NewErrorResponse(500)
			return
		}
		data = common.NewResponse(nil)
		return
	default:
		data = common.NewErrorResponse(500)
		return
	}
}

// @Tags 超级节点算力表显示   后台操作　or 用户查看　都可
// @Accept  json
// @Produce json
// @Success 200
// @router /show_super_formula_list____超级节点算力表显示___后台操作__or_用户查看__都可 [Post]
func (this *BackStageManagement) ShowSuperFormulaList() {
	var (
		data       *common.ResponseData
		o          = models.NewOrm()
		force_list []models.SuperForceTable
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	o.QueryTable("super_force_table").All(&force_list)
	models.QuickSortSuperForce(force_list, 0, len(force_list)-1)

	data = common.NewResponse(force_list)
	return
}

// @Tags 超级节点算力表信息修改
// @Accept  json
// @Produce json
// @Param super_force_id query string true "超级节点的算力信息的id  如果是删除多条，　用　, 逗号隔开"
// @Param action query string true "具体的操作　　delete=删除　update=更新　insert=新增    更新必须传这条数据的全部数据回来，id 除外，和新增一样，删除只要id就够"
// @Param levelstr query string true "等级名字"
// @Param coin_number query string true "要求的持币数量"
// @Param force query string true "算力　　要以小数的格式返回　　如 : 15% = 0.15 "
// @Success 200
// @router /operation_super_formula_list___超级节点算力表信息修改 [Post]
func (this *BackStageManagement) OperationSuperFormulaList() {
	var (
		data           *common.ResponseData
		o              = models.NewOrm()
		super_force_id = this.GetString("super_force_id")
		action         = this.GetString("action")
		levelstr       = this.GetString("levelstr")
		coin_number, _ = this.GetInt("coin_number")
		force_str          = this.GetString("force")
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	force, _ := strconv.ParseFloat(force_str, 64)

	switch action {
	case "delete":
		id_strs := strings.Split(super_force_id, ",")
		for _, v := range id_strs {
			id, _ := strconv.Atoi(v)
			_, err := o.QueryTable("super_force_table").Filter("id", id).Delete()
			if err != nil {
				data = common.NewErrorResponse(500)
				return
			}
			data = common.NewResponse(nil)
			return
		}
	case "update":
		id, _ := strconv.Atoi(super_force_id)
		_, err := o.QueryTable("super_force_table").
			Filter("id", id).
			Update(orm.
			Params{"level":levelstr,"coin_number_rule":coin_number,"bonus_calculation":force})
		if err != nil {
			data = common.NewErrorResponse(500)
			return
		}
		data = common.NewResponse(nil)
		return
	case "insert":
		super_force := models.SuperForceTable{
			Level:            levelstr,
			CoinNumberRule:   coin_number,
			BonusCalculation: force,
		}
		_, err := o.Insert(&super_force)
		if err != nil {
			data = common.NewErrorResponse(500)
			return
		}
		data = common.NewResponse(nil)
		return
	default:
		data = common.NewErrorResponse(500)
		return
	}
}

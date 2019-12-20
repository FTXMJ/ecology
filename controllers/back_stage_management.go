package controllers

import (
	"ecology/common"
	"ecology/logs"
	"ecology/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

// 后台管理
type BackStageManagement struct {
	beego.Controller
}

// @Tags 算力表显示   后台操作　or 用户查看　都可
// @Accept  json
// @Produce json
// @Success 200__算力表显示后台操or用户查_都可 {object} models.ForceTable_test
// @router /show_formula_list [GET]
func (this *BackStageManagement) ShowFormulaList() {
	var (
		data       *common.ResponseData
		o          = models.NewOrm()
		force_list []models.ForceTable
		api_url    = this.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	_, err := o.QueryTable("force_table").All(&force_list)
	if err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewErrorResponse(500, "算力数据获取失败!")
		return
	}
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
// @Success 200___算力表信息修改
// @router /operation_formula_list [POST]
func (this *BackStageManagement) OperationFormulaList() {
	var (
		data                      *common.ResponseData
		o                         = models.NewOrm()
		api_url                   = this.Ctx.Request.RequestURI
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
				logs.Log.Error(api_url, err)
				data = common.NewErrorResponse(500, "算力表　删除失败!")
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
				Params{"level": levelstr, "low_hold": low_hold, "high_hold": high_hold, "return_multiple": return_multiple, "hold_return_rate": hold_return_rate, "recommend_return_rate": recommend_return_rate, "team_return_rate": team_return_rate})
		if err != nil {
			logs.Log.Error(api_url, err)
			data = common.NewErrorResponse(500, "算力表　更新失败!")
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
			logs.Log.Error(api_url, err)
			data = common.NewErrorResponse(500, "算力表新增失败!")
			return
		}
		data = common.NewResponse(nil)
		return
	default:
		logs.Log.Error(api_url, "未知操作!")
		data = common.NewErrorResponse(500, "未知操作!")
		return
	}
}

// @Tags 超级节点算力表显示   后台操作　or 用户查看　都可
// @Accept  json
// @Produce json
// @Success 200__超级节点算力表显示后台操作or用户查看都可以  {object} models.SuperForceTable_test
// @router /show_super_formula_list [GET]
func (this *BackStageManagement) ShowSuperFormulaList() {
	var (
		data       *common.ResponseData
		o          = models.NewOrm()
		api_url    = this.Ctx.Request.RequestURI
		force_list []models.SuperForceTable
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	_, err := o.QueryTable("super_force_table").All(&force_list)
	if err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewErrorResponse(500, "超级节点算力数据获取失败!")
		return
	}
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
// @Success 200___超级节点算力表信息修改
// @router /operation_super_formula_list [POST]
func (this *BackStageManagement) OperationSuperFormulaList() {
	var (
		data           *common.ResponseData
		o              = models.NewOrm()
		api_url        = this.Ctx.Request.RequestURI
		super_force_id = this.GetString("super_force_id")
		action         = this.GetString("action")
		levelstr       = this.GetString("levelstr")
		coin_number, _ = this.GetInt("coin_number")
		force_str      = this.GetString("force")
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
				logs.Log.Error(api_url, err)
				data = common.NewErrorResponse(500, "超级节点算力表 删除失败!")
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
				Params{"level": levelstr, "coin_number_rule": coin_number, "bonus_calculation": force})
		if err != nil {
			logs.Log.Error(api_url, err)
			data = common.NewErrorResponse(500, "超级节点算力表 更新失败!")
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
			logs.Log.Error(api_url, err)
			data = common.NewErrorResponse(500, "超级节点算力表 新增失败!")
			return
		}
		data = common.NewResponse(nil)
		return
	default:
		logs.Log.Error(api_url, "未知操作!")
		data = common.NewErrorResponse(500, "未知操作!")
		return
	}
}

// @Tags root-历史
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /return_page_hostry_root [GET]
func (this *BackStageManagement) ReturnPageHostryRoot() {
	var (
		data            *common.ResponseData
		current_page, _ = this.GetInt("page")
		page_size, _    = this.GetInt("pageSize")
		api_url         = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	page := models.Page{
		CurrentPage: current_page,
		PageSize:    page_size,
	}
	values, p, err := models.SelectHosteryRoot(page)
	if err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewResponse(models.HostryPageInfo{})
		return
	}
	hostory_list := models.HostryPageInfo{
		Items: values,
		Page:  p,
	}
	data = common.NewResponse(hostory_list)
	return
}

// @Tags root-历史-筛选
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param type query string true "查询的数据类型　　 =铸币表　account_detail=充值表"
// @Param user_id query string true "用户id"
// @Param tx_id query string true "订单id"
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /filter_history_info [POST]
func (this *BackStageManagement) FilterHistoryInfo() {
	var (
		data              *common.ResponseData
		current_page, _   = this.GetInt("page")
		page_size, _      = this.GetInt("pageSize")
		table_name        = this.GetString("type")
		user_id           = this.GetString("user_id")
		tx_id             = this.GetString("tx_id")
		start_time_int, _ = this.GetInt64("start_time")
		end_time_int, _   = this.GetInt64("end_time")
		api_url           = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	start_time := time.Unix(start_time_int, 0).Format("2006-01-02 15:04:05")
	end_time := time.Unix(end_time_int, 0).Format("2006-01-02 15:04:05")

	find_obj := models.FindObj{
		UserId:    user_id,
		TxId:      tx_id,
		StartTime: start_time,
		EndTime:   end_time,
	}
	p := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	list, page, err := models.SelectPondMachinemsg(find_obj, p, table_name)
	if err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewResponse(models.HostryFindInfo{})
		return
	}
	hostory_list := models.HostryFindInfo{
		Items: list,
		Page:  page,
	}

	data = common.NewResponse(hostory_list)
	return
}

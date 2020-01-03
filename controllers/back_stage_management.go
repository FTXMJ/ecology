package controllers

import (
	"ecology/common"
	"ecology/logs"
	"ecology/models"
	"fmt"
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
		data = common.NewErrorResponse(500, "算力数据获取失败!", []models.ForceTable{})
		return
	}
	models.QuickSortForce(force_list, 0, len(force_list)-1)
	data = common.NewResponse(force_list)
	return
}

// @Tags 算力等级详情
// @Accept  json
// @Produce json
// @Success 200__算力表显示后台操or用户查_都可 {object} models.ForceTable_test
// @Success 200____算力等级详情 {object} models.ForceTable_test_yq
// @router /show_user_formula [GET]
func (this *BackStageManagement) ShowUserFormula() {
	var (
		data *common.ResponseData
		o    = models.NewOrm()
		//api_url    = this.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	jwtValues := GetJwtValues(this.Ctx)
	user_id := jwtValues.UserID
	account := models.Account{UserId: user_id}
	o.Read(&account, "user_id")
	for_mula := models.Formula{EcologyId: account.Id}
	o.Read(&for_mula, "ecology_id")
	for_mula_table := models.ForceTable{Level: for_mula.Level}
	o.Read(&for_mula_table, "level")
	data = common.NewResponse(for_mula_table)
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
// @Param picture_url query string true "图片url"
// @Success 200___算力表信息修改
// @router /admin/operation_formula_list [POST]
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
		picture_url               = this.GetString("picture_url")
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
				data = common.NewErrorResponse(500, "算力表　删除失败!", nil)
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
				Params{"level": levelstr, "low_hold": low_hold, "high_hold": high_hold, "return_multiple": return_multiple, "hold_return_rate": hold_return_rate, "recommend_return_rate": recommend_return_rate, "team_return_rate": team_return_rate, "picture_url": picture_url})
		if err != nil {
			logs.Log.Error(api_url, err)
			data = common.NewErrorResponse(500, "算力表　更新失败!", nil)
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
			PictureUrl:          picture_url,
		}
		_, err := o.Insert(&force)
		if err != nil {
			logs.Log.Error(api_url, err)
			data = common.NewErrorResponse(500, "算力表新增失败!", nil)
			return
		}
		data = common.NewResponse(nil)
		return
	default:
		logs.Log.Error(api_url, "未知操作!")
		data = common.NewErrorResponse(500, "未知操作!", nil)
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
		data = common.NewErrorResponse(500, "节点算力数据获取失败!", []models.SuperForceTable{})
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
// @router /admin/operation_super_formula_list [POST]
func (this *BackStageManagement) OperationSuperFormulaList() {
	var (
		data            *common.ResponseData
		o               = models.NewOrm()
		api_url         = this.Ctx.Request.RequestURI
		super_force_id  = this.GetString("super_force_id")
		action          = this.GetString("action")
		levelstr        = this.GetString("levelstr")
		coin_number_str = this.GetString("coin_number")
		force_str       = this.GetString("force")
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	force, _ := strconv.ParseFloat(force_str, 64)
	coin_number, _ := strconv.ParseFloat(coin_number_str, 64)

	switch action {
	case "delete":
		id_strs := strings.Split(super_force_id, ",")
		for _, v := range id_strs {
			id, _ := strconv.Atoi(v)
			_, err := o.QueryTable("super_force_table").Filter("id", id).Delete()
			if err != nil {
				logs.Log.Error(api_url, err)
				data = common.NewErrorResponse(500, "节点算力表 删除失败!", nil)
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
			data = common.NewErrorResponse(500, "节点算力表 更新失败!", nil)
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
			data = common.NewErrorResponse(500, "节点算力表 新增失败!", nil)
			return
		}
		data = common.NewResponse(nil)
		return
	default:
		logs.Log.Error(api_url, "未知操作!")
		data = common.NewErrorResponse(500, "未知操作!", nil)
		return
	}
}

// @Tags root-历史
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /admin/return_page_hostry_root [GET]
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
		logs.Log.Error(api_url+"   更新状态失败,数据库错误", err)
		data = common.NewErrorResponse(500, "更新状态失败,数据库错误", models.HostryPageInfo{})
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
// @Param type query string true "查询的s数据类型　　 blocked_detail=铸币表　account_detail=充值表"
// @Param user_id query string true "用户id"
// @Param tx_id query string true "订单id"
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /admin/filter_history_info [GET]
func (this *BackStageManagement) FilterHistoryInfo() {
	var (
		data              *common.ResponseData
		current_page, _   = this.GetInt("page")
		page_size, _      = this.GetInt("pageSize")
		table_name        = this.GetString("type")
		user_id           = this.GetString("user_id")
		user_name         = this.GetString("user_name")
		tx_id             = this.GetString("tx_id")
		start_time_int, _ = this.GetInt64("start_time")
		end_time_int, _   = this.GetInt64("end_time")
		api_url           = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	start_time := ""
	end_time := ""
	if start_time_int == 0 || end_time_int == 0 {
		start_time = "2006-01-02 15:04:05"
		end_time = time.Now().Format("2006-01-02 15:04:05")
	} else {
		start_time = time.Unix(start_time_int, 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(end_time_int, 0).Format("2006-01-02 15:04:05")
	}

	find_obj := models.FindObj{
		UserId:    user_id,
		TxId:      tx_id,
		UserName:  user_name,
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
		logs.Log.Error(api_url+"   更新状态失败,数据库错误", err)
		data = common.NewErrorResponse(500, "更新状态失败,数据库错误", []models.HostryFindInfo{})
		return
	}
	hostory_list := models.HostryFindInfo{
		Items: list,
		Page:  page,
	}

	data = common.NewResponse(hostory_list)
	return
}

// @Tags root-用户生态列表
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param user_id query string true "用户id  不搜就传空，搜索就传user_id"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____用户生态列表 {object} models.UEOBJList_test
// @router /admin/user_ecology_list [GET]
func (this *BackStageManagement) UserEcologyList() {
	var (
		data            *common.ResponseData
		current_page, _ = this.GetInt("page")
		page_size, _    = this.GetInt("pageSize")
		user_id         = this.GetString("user_id")
		user_name       = this.GetString("user_name")
		//api_url         = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	u_e_obj_list, p := models.FindU_E_OBJ(page, user_id, user_name)
	u_e_objs := []models.U_E_OBJ{}
	for _, v := range u_e_obj_list {
		u_e_obj := models.U_E_OBJ{}
		hold, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.HoldReturnRate), 64)
		u_e_obj.HoldReturnRate = hold
		reco, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.RecommendReturnRate), 64)
		u_e_obj.RecommendReturnRate = reco
		rele, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.Released), 64)
		u_e_obj.Released = rele
		team, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.TeamReturnRate), 64)
		u_e_obj.TeamReturnRate = team
		tobe, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.ToBeReleased), 64)
		u_e_obj.ToBeReleased = tobe
		coin, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.CoinAll), 64)
		u_e_obj.CoinAll = coin
		u_e_obj.Level = v.Level
		u_e_obj.UserId = v.UserId
		u_e_obj.UserName = v.UserName
		u_e_obj.ReturnMultiple = v.ReturnMultiple
		u_e_objs = append(u_e_objs, u_e_obj)
	}
	list := models.UEOBJList{
		Items: u_e_objs,
		Page:  p,
	}
	data = common.NewResponse(list)
	return
}

// @Tags root-用户生态禁止列表
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param user_id query string true "用户id  不搜就传空，搜索就传user_id"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____用户生态列表 {object} models.UEOBJList_test
// @router /admin/user_ecology_false_list [GET]
func (this *BackStageManagement) UserEcologyFalseList() {
	var (
		data            *common.ResponseData
		current_page, _ = this.GetInt("page")
		page_size, _    = this.GetInt("pageSize")
		user_id         = this.GetString("user_id")
		user_name       = this.GetString("user_name")
		//api_url         = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	u_e_obj_list, p := models.FindFalseUser(page, user_id, user_name)
	list := models.UserFalse{
		Items: u_e_obj_list,
		Page:  p,
	}
	data = common.NewResponse(list)
	return

}

// @Tags root-算力流水表
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param user_id query string true "用户id"
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____算力流水表 {object} models.FlowList_test
// @router /admin/computational_flow [GET]
func (this *BackStageManagement) ComputationalFlow() {
	var (
		data              *common.ResponseData
		current_page, _   = this.GetInt("page")
		page_size, _      = this.GetInt("pageSize")
		user_id           = this.GetString("user_id")
		user_name         = this.GetString("user_name")
		start_time_int, _ = this.GetInt64("start_time")
		end_time_int, _   = this.GetInt64("end_time")
		api_url           = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	start_time := ""
	end_time := ""
	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(start_time_int, 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(end_time_int, 0).Format("2006-01-02 15:04:05")
	}

	find_obj := models.FindObj{
		UserId:    user_id,
		UserName:  user_name,
		TxId:      "",
		StartTime: start_time,
		EndTime:   end_time,
	}
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	flows, p, err := models.SelectFlows(find_obj, page, "blocked_detail")
	if err != nil {
		logs.Log.Error(api_url+"    更新状态失败,数据库错误", err)
		data = common.NewErrorResponse(500, err.Error(), models.FlowList{})
		return
	}
	var flowss []models.Flow
	for _, v := range flows {
		flow := models.Flow{}
		hold, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.HoldReturnRate), 64)
		flow.HoldReturnRate = hold
		reco, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.RecommendReturnRate), 64)
		flow.RecommendReturnRate = reco
		rele, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.Released), 64)
		flow.Released = rele
		team, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", v.TeamReturnRate), 64)
		flow.TeamReturnRate = team
		flow.UserId = v.UserId
		flow.UserName = v.UserName
		flow.UpdateTime = v.UpdateTime
		flowss = append(flowss, flow)
	}
	user_SF_information := models.FlowList{
		Items: flowss,
		Page:  p,
	}
	data = common.NewResponse(user_SF_information)
	return
}

// @Tags root-用户收益控制＿＿展示
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param user_id query string true "用户id"
// @Param account_id query string true "生态仓库id"
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____用户收益控制＿＿展示
// @router /admin/ecological_income_control [GET]
func (this *BackStageManagement) EcologicalIncomeControl() {
	var (
		data              *common.ResponseData
		current_page, _   = this.GetInt("page")
		page_size, _      = this.GetInt("pageSize")
		user_id           = this.GetString("user_id")
		account_id        = this.GetString("account_id")
		user_name         = this.GetString("user_name")
		start_time_int, _ = this.GetInt64("start_time")
		end_time_int, _   = this.GetInt64("end_time")
		api_url           = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	start_time := ""
	end_time := ""
	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(start_time_int, 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(end_time_int, 0).Format("2006-01-02 15:04:05")
	}

	p := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	find_obj := models.FindObj{
		UserId:    user_id,
		TxId:      account_id,
		UserName:  user_name,
		StartTime: start_time,
		EndTime:   end_time,
	}
	account_off, page, err := models.FindUserAccountOFF(p, find_obj)
	if err != nil {
		logs.Log.Error(api_url+"    数据库错误,数据查询失败", err)
		data = common.NewErrorResponse(500, "数据库错误", models.UserAccountOFF{})
		return
	}
	user_SF_information := models.UserAccountOFF{
		Items: account_off,
		Page:  page,
	}
	data = common.NewResponse(user_SF_information)
	return
}

// @Tags root-用户收益控制＿＿修改
// @Accept  json
// @Produce json
// @Param account_id query string true "生态仓库id"
// @Param profit_type query string true "静态=1  动态=2 节点=3"
// @Param profit_start query string true "启用=1  禁用=2"
// @Success 200____用户收益控制＿＿修改
// @router /admin/ecological_income_control_update [POST]
func (this *BackStageManagement) EcologicalIncomeControlUpdate() {
	var (
		data                *common.ResponseData
		profit_type_int, _  = this.GetInt("profit_type")
		profit_start_int, _ = this.GetInt("profit_start")
		strs                = this.GetString("account_id")
		//api_url             = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	profit_start := false
	if profit_start_int == 1 {
		profit_start = true
	}

	str := strings.Split(strs, ",")

	o := models.NewOrm()
	err_user := ""
	for _, v := range str {
		id_int, _ := strconv.Atoi(v)
		acc := models.Account{
			Id: id_int,
		}
		o.Read(&acc)
		switch profit_type_int {
		case 2:
			acc.DynamicRevenue = profit_start
		case 1:
			acc.StaticReturn = profit_start
		case 3:
			acc.PeerState = profit_start
		}
		acc.UpdateDate = time.Now().Format("2006-01-02 15:04:05")
		_, err := o.Update(&acc)
		if err != nil {
			if len(err_user) == 0 {
				err_user += v
			}
			err_user += "," + v
		}
	}
	if len(err_user) != 0 {
		data = common.NewErrorResponse(500, "这些用户更新失败:"+err_user, nil)
		return
	}
	data = common.NewResponse(nil)
	return
}

// @Tags 节点用户列表
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____节点用户列表
// @router /admin/peer_user_list [GET]
func (this *BackStageManagement) PeerUserList() {
	var (
		data            *common.ResponseData
		current_page, _ = this.GetInt("page")
		page_size, _    = this.GetInt("pageSize")
		user_name       = this.GetString("user_name")
		//api_url             = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	p_u_s := []models.PeerUser{}
	user := []models.User{}
	switch user_name {
	case "":
		models.NewOrm().Raw("select * from user").QueryRows(&user)
	default:
		models.NewOrm().Raw("select * from user where user_name=?", user_name).QueryRows(&user)
	}
	for _, v := range user {
		p_u := models.PeerUser{}
		update_date, level, tfor, _ := ReturnSuperPeerLevel(v.UserId)
		if level != "" {
			acc := models.Account{UserId: v.UserId}
			models.NewOrm().Read(&acc, "user_id")
			p_u.AccountId = acc.Id
			p_u.UserId = v.UserId
			p_u.UserName = v.UserName
			p_u.Level = level
			p_u.State = acc.PeerState
			p_u.Number = tfor
			p_u.UpdateTime = update_date
			p_u_s = append(p_u_s, p_u)
		}
	}
	peer_users, p := models.PageS(p_u_s, page)
	peer_user_list := models.PeerUserFalse{
		Items: peer_users,
		Page:  p,
	}
	data = common.NewResponse(peer_user_list)
	return
}

// @Tags 节点历史记录
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Success 200____节点历史记录  {object} models.PeerHistoryList_test
// @router /admin/peer_a_bouns_list [GET]
func (this *BackStageManagement) PeerABounsList() {
	var (
		data              *common.ResponseData
		current_page, _   = this.GetInt("page")
		page_size, _      = this.GetInt("pageSize")
		start_time_int, _ = this.GetInt64("start_time")
		end_time_int, _   = this.GetInt64("end_time")
		//api_url             = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	start_time := ""
	end_time := ""
	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(start_time_int, 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(end_time_int, 0).Format("2006-01-02 15:04:05")
	}
	peer_history := []models.PeerHistory{}
	switch start_time {
	case "":
		_, err := models.NewOrm().Raw("select * from peer_history where time>=? and time<=?", start_time, end_time).QueryRows(&peer_history)
		if err != nil {
			data = common.NewResponse(models.PeerHistoryList{})
			return
		}
	default:
		_, err := models.NewOrm().Raw("select * from peer_history").QueryRows(&peer_history)
		if err != nil {
			data = common.NewResponse(models.PeerHistoryList{})
			return
		}
	}
	peer_users, p := models.PageHistory(peer_history, page)
	peer_user_list := models.PeerHistoryList{
		Items: peer_users,
		Page:  p,
	}
	data = common.NewResponse(peer_user_list)
	return
}

// @Tags 节点收益记录流水
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Success 200____节点收益记录流水 {object} models.PeerHistoryList_test
// @router /admin/peer_a_bouns_history_list [GET]
func (this *BackStageManagement) PeerABounsHistoryList() {
	var (
		data            *common.ResponseData
		current_page, _ = this.GetInt("page")
		page_size, _    = this.GetInt("pageSize")
		user_name       = this.GetString("user_name")
		//api_url             = this.Controller.Ctx.Request.RequestURI
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	list, p, err := SelectPeerABounsList(page, user_name)
	if err != nil {
		data = common.NewErrorResponse(500, "数据库操作失败", models.PeerListABouns{})
		return
	}
	one := models.PeerListABouns{
		Items: list,
		Page:  p,
	}
	data = common.NewResponse(one)
	return
}

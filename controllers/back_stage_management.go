package controllers

import (
	"ecology/actuator"
	"ecology/common"
	db "ecology/db"
	"ecology/filter"
	"ecology/logs"
	"ecology/models"
	"github.com/gin-gonic/gin"

	"errors"
	"strconv"
	"strings"
	"time"
)

// 后台管理

// @Tags 算力表显示   后台操作　or 用户查看　都可
// @Accept  json
// @Produce json
// @Success 200__算力表显示后台操or用户查_都可 {object} models.ForceTable_test
// @router /show_formula_list [GET]
func ShowFormulaList(c *gin.Context) {
	var (
		data       *common.ResponseData
		o          = db.NewEcologyOrm()
		force_list = make([]models.ForceTable, 0)
	)
	defer func() {
		c.JSON(200, data)
	}()
	err := o.Table("force_table").Find(&force_list)
	if err.Error != nil {
		logs.Log.Error(err)
		data = common.NewErrorResponse(500, "算力数据获取失败!", []models.ForceTable{})
		return
	}
	actuator.QuickSortForce(force_list, 0, len(force_list)-1)
	data = common.NewResponse(force_list[1:])
	return
}

// @Tags 算力等级详情
// @Accept  json
// @Produce json
// @Success 200__算力表显示后台操or用户查_都可 {object} models.ForceTable_test
// @Success 200____算力等级详情 {object} models.ForceTable_test_yq
// @router /show_user_formula [GET]
func ShowUserFormula(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()
	)
	defer func() {
		c.JSON(200, data)
	}()
	jwtValues := filter.GetJwtValues(c)
	user_id := jwtValues.UserID
	account := models.Account{}
	o.Table("account").Where("user_id = ?", user_id).First(&account)
	for_mula := models.Formula{}
	o.Table("formula").Where("ecology_id = ?", account.Id).First(&for_mula)
	for_mula_table := models.ForceTable{Level: for_mula.Level, ReturnMultiple: for_mula.ReturnMultiple}
	if for_mula.Level == "" {
		for_mula_table.ReturnMultiple = 1
	}
	o.Table("force_table").Where("return_multiple = ?", for_mula_table.ReturnMultiple).First(&for_mula_table)
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
func OperationFormulaList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()
		err  error

		force_id                  = c.PostForm("force_id")
		action                    = c.PostForm("action")
		levelstr                  = c.PostForm("levelstr")
		low_hold_int              = c.PostForm("low_hold")
		low_hold, _               = strconv.Atoi(low_hold_int)
		high_hold_int             = c.PostForm("high_hold")
		high_hold, _              = strconv.Atoi(high_hold_int)
		return_multiple_str       = c.PostForm("return_multiple")
		hold_return_rate_str      = c.PostForm("hold_return_rate")
		recommend_return_rate_str = c.PostForm("recommend_return_rate")
		team_return_rate_str      = c.PostForm("team_return_rate")
		picture_url               = c.PostForm("picture_url")

		return_multiple, _       = strconv.ParseFloat(return_multiple_str, 64)
		hold_return_rate, _      = strconv.ParseFloat(hold_return_rate_str, 64)
		recommend_return_rate, _ = strconv.ParseFloat(recommend_return_rate_str, 64)
		team_return_rate, _      = strconv.ParseFloat(team_return_rate_str, 64)
	)
	defer func() {
		if err != nil {
			o.Rollback()
			logs.Log.Error(err)
			data = common.NewErrorResponse(500, err.Error(), nil)
			c.JSON(200, data)
		} else {
			c.JSON(200, data)
		}
	}()
	o.Begin()

	switch action {
	case "delete":
		id_strs := strings.Split(force_id, ",")
		for _, v := range id_strs {
			id, _ := strconv.Atoi(v)
			if er := o.Delete(&models.ForceTable{}, "id = ?", strconv.Itoa(id)); er.Error != nil {
				err = errors.New("算力表　删除失败!")
				return
			}
		}
	case "update":
		id, _ := strconv.Atoi(force_id)
		force_table := models.ForceTable{}
		o.Table("force_tale").Where("id = ?", id).First(&force_table)
		force_table.Id = id
		force_table.Level = levelstr
		force_table.LowHold = low_hold
		force_table.HighHold = high_hold
		force_table.ReturnMultiple = return_multiple
		force_table.HoldReturnRate = hold_return_rate
		force_table.RecommendReturnRate = recommend_return_rate
		force_table.TeamReturnRate = team_return_rate
		force_table.PictureUrl = picture_url
		o.Update(&force_table)
		if er := o.Update(&force_table); er.Error != nil {
			err = errors.New("算力表　更新失败!")
			return
		}
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
		if er := o.Create(&force); er.Error != nil {
			err = errors.New("算力表新增失败!")
			return
		}
	default:
		err = errors.New("未知操作!")
	}
	o.Commit()

	data = common.NewResponse(nil)
	return
}

// @Tags 超级节点算力表显示   后台操作　or 用户查看　都可
// @Accept  json
// @Produce json
// @Success 200__超级节点算力表显示后台操作or用户查看都可以  {object} models.SuperForceTable_test
// @router /show_super_formula_list [GET]
func ShowSuperFormulaList(c *gin.Context) {
	var (
		data       *common.ResponseData
		o          = db.NewEcologyOrm()
		force_list = make([]models.SuperForceTable, 0)
	)
	defer func() {
		c.JSON(200, data)
	}()
	o.Raw("select * from super_force_table").Find(&force_list)
	actuator.QuickSortSuperForce(force_list, 0, len(force_list)-1)
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
func OperationSuperFormulaList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()
		err  error

		super_force_id  = c.PostForm("super_force_id")
		action          = c.PostForm("action")
		levelstr        = c.PostForm("levelstr")
		coin_number_str = c.PostForm("coin_number")
		force_str       = c.PostForm("force")

		force, _       = strconv.ParseFloat(force_str, 64)
		coin_number, _ = strconv.ParseFloat(coin_number_str, 64)
	)
	defer func() {
		if err != nil {
			o.Rollback()
			logs.Log.Error(err)
			data = common.NewErrorResponse(500, err.Error(), nil)
			c.JSON(200, data)
		} else {
			c.JSON(200, data)
		}
	}()

	o.Begin()

	switch action {
	case "delete":
		id_strs := strings.Split(super_force_id, ",")
		for _, v := range id_strs {
			id, _ := strconv.Atoi(v)
			o.Delete(&models.SuperForceTable{}, "id = ?", id)
			if er := o.Delete(&models.SuperForceTable{}, "id = ?", id); er.Error != nil {
				err = errors.New("节点算力表 删除失败!")
				return
			}
		}
	case "update":
		id, _ := strconv.Atoi(super_force_id)
		if er := o.Model(&models.SuperForceTable{}).Where("id = ?", id).Updates(map[string]interface{}{"level": levelstr, "coin_number_rule": coin_number, "bonus_calculation": force}); er.Error != nil {
			err = errors.New("节点算力表 更新失败!")
			return
		}
	case "insert":
		super_force := models.SuperForceTable{
			Level:            levelstr,
			CoinNumberRule:   coin_number,
			BonusCalculation: force,
		}
		if er := o.Create(&super_force); er.Error != nil {
			err = errors.New("节点算力表 新增失败!")
			return
		}
	default:
		err = errors.New("未知操作!")
	}
	o.Commit()

	data = common.NewResponse(nil)
	return
}

// @Tags root-历史
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /admin/return_page_hostry_root [GET]
func ReturnPageHostryRoot(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)
	)
	defer func() {
		c.JSON(200, data)
	}()
	page := models.Page{
		CurrentPage: current_page,
		PageSize:    page_size,
	}
	values, p, err := actuator.SelectHosteryRoot(o, page)
	if err != nil {
		logs.Log.Error("   更新状态失败,数据库错误", err)
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
func FilterHistoryInfo(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str   = c.PostForm("page")
		current_page, _    = strconv.Atoi(current_page_str)
		page_size_str      = c.PostForm("pageSize")
		page_size, _       = strconv.Atoi(page_size_str)
		table_name         = c.PostForm("type")
		user_id            = c.PostForm("user_id")
		user_name          = c.PostForm("user_name")
		tx_id              = c.PostForm("tx_id")
		start_time_int_str = c.PostForm("start_time")
		start_time_int, _  = strconv.Atoi(start_time_int_str)
		end_time_int_str   = c.PostForm("end_time")
		end_time_int, _    = strconv.Atoi(end_time_int_str)
		start_time         = ""
		end_time           = ""
	)
	defer func() {
		c.JSON(200, data)
	}()

	if start_time_int == 0 || end_time_int == 0 {
		start_time = "2006-01-02 15:04:05"
		end_time = time.Now().Format("2006-01-02 15:04:05")
	} else {
		start_time = time.Unix(int64(start_time_int), 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(int64(end_time_int), 0).Format("2006-01-02 15:04:05")
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

	list := make([]models.BlockedDetailIndex, 0)
	page := models.Page{}
	var err error
	if table_name == "account_detail" {
		list, page, err = actuator.SelectPondMachinemsgForAcc(o, find_obj, p, table_name)
	} else if table_name == "blocked_detail" {
		list, page, err = actuator.SelectPondMachinemsgForBlo(o, find_obj, p, table_name)
	}
	if err != nil {
		logs.Log.Error("   更新状态失败,数据库错误", err)
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
func UserEcologyList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)
		user_id          = c.PostForm("user_id")
		user_name        = c.PostForm("user_name")
	)
	defer func() {
		c.JSON(200, data)
	}()

	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	u_e_obj_list, p := actuator.FindU_E_OBJ(o, page, user_id, user_name)

	list := models.UEOBJList{
		Items: u_e_obj_list,
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
func UserEcologyFalseList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		user_id   = c.PostForm("user_id")
		user_name = c.PostForm("user_name")
	)
	defer func() {
		c.JSON(200, data)
	}()

	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	u_e_obj_list, p := actuator.FindFalseUser(o, page, user_id, user_name)
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
func ComputationalFlow(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		user_id   = c.PostForm("user_id")
		user_name = c.PostForm("user_name")

		start_time_int_str = c.PostForm("start_time")
		start_time_int, _  = strconv.Atoi(start_time_int_str)
		end_time_int_str   = c.PostForm("end_time")
		end_time_int, _    = strconv.Atoi(end_time_int_str)

		start_time = ""
		end_time   = ""
	)
	defer func() {
		c.JSON(200, data)
	}()

	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(int64(start_time_int), 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(int64(end_time_int), 0).Format("2006-01-02 15:04:05")
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

	flows, p, err := actuator.SelectFlows(o, find_obj, page, "blocked_detail")
	if err != nil {
		logs.Log.Error("    更新状态失败,数据库错误", err)
		data = common.NewErrorResponse(500, err.Error(), models.FlowList{})
		return
	}

	user_SF_information := models.FlowList{
		Items: flows,
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
func EcologicalIncomeControl(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		user_id    = c.PostForm("user_id")
		account_id = c.PostForm("account_id")
		user_name  = c.PostForm("user_name")

		start_time_int_str = c.PostForm("start_time")
		start_time_int, _  = strconv.Atoi(start_time_int_str)
		end_time_int_str   = c.PostForm("end_time")
		end_time_int, _    = strconv.Atoi(end_time_int_str)

		start_time = ""
		end_time   = ""
	)
	defer func() {
		c.JSON(200, data)
	}()

	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(int64(start_time_int), 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(int64(end_time_int), 0).Format("2006-01-02 15:04:05")
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
	account_off, page, err := actuator.FindUserAccountOFF(o, p, find_obj)
	if err != nil {
		logs.Log.Error("    数据库错误,数据查询失败", err)
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
func EcologicalIncomeControlUpdate(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		profit_type_int_str  = c.PostForm("profit_type")
		profit_type_int, _   = strconv.Atoi(profit_type_int_str)
		profit_start_int_str = c.PostForm("profit_start")
		profit_start_int, _  = strconv.Atoi(profit_start_int_str)

		strs = c.PostForm("account_id")
	)
	defer func() {
		c.JSON(200, data)
	}()
	profit_start := false
	if profit_start_int == 1 {
		profit_start = true
	}

	str := strings.Split(strs, ",")

	err_user := ""
	for _, v := range str {
		id_int, _ := strconv.Atoi(v)
		acc := models.Account{}
		o.Table("account").Where("id = ?", id_int).First(&acc)
		switch profit_type_int {
		case 2:
			acc.DynamicRevenue = profit_start
		case 1:
			acc.StaticReturn = profit_start
		case 3:
			acc.PeerState = profit_start
		}
		acc.UpdateDate = time.Now().Format("2006-01-02 15:04:05")
		o.Update(&acc)
		er := o.Update(&acc)
		if er.Error != nil {
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
func PeerUserList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		user_name = c.PostForm("user_name")

		p_u_s = []models.PeerUser{}
		user  = []models.User{}
	)
	defer func() {
		c.JSON(200, data)
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	switch user_name {
	case "":
		o.Raw("select * from user").Find(&user)
	default:
		o.Raw("select * from user where user_name=?", user_name).Find(&user)
	}
	g := models.GlobalOperations{}
	o.Raw("select * from global_operations where operation=?", "全局节点分红控制").First(&g)

	for _, v := range user {
		p_u := models.PeerUser{}
		update_date, level, tfor, _ := actuator.ReturnSuperPeerLevel(o, v.UserId)
		if level != "" {
			acc := models.Account{}
			o.Table("account").Where("user_id = ?", v.UserId).First(&acc)
			p_u.AccountId = acc.Id
			p_u.UserId = v.UserId
			p_u.UserName = v.UserName
			p_u.Level = level
			p_u.Number = tfor
			p_u.UpdateTime = update_date

			var peer_state bool = acc.PeerState
			if g.State == false {
				peer_state = false
			}
			p_u.State = peer_state
			p_u_s = append(p_u_s, p_u)
		}
	}
	peer_users, p := actuator.PageS(p_u_s, page)
	peer_user_list := models.PeerUserFalse{
		Items: peer_users,
		Page:  p,
	}
	if len(peer_users) == 0 {
		peer_user_list.Items = []models.PeerUser{}
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
func PeerABounsList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		start_time_int_str = c.PostForm("start_time")
		start_time_int, _  = strconv.Atoi(start_time_int_str)
		end_time_int_str   = c.PostForm("end_time")
		end_time_int, _    = strconv.Atoi(end_time_int_str)

		start_time   = ""
		end_time     = ""
		peer_history = make([]models.PeerHistory, 0)
	)
	defer func() {
		c.JSON(200, data)
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(int64(start_time_int), 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(int64(end_time_int), 0).Format("2006-01-02 15:04:05")
	}

	switch start_time {
	case "":
		er := o.Raw("select * from peer_history order by time desc").Find(&peer_history)
		if er.Error != nil {
			data = common.NewResponse(models.PeerHistoryList{})
			return
		}
	default:
		er := o.Raw("select * from peer_history where time>=? and time<=? order by time desc", start_time, end_time).Find(&peer_history)
		if er.Error != nil {
			data = common.NewResponse(models.PeerHistoryList{})
			return
		}
	}
	peer_users, p := actuator.PageHistory(peer_history, page)
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
// @Param start_time query string true "开始时间"
// @Param end_time query string true "结束时间"
// @Success 200____节点收益记录流水 {object} models.PeerListABouns_test
// @router /admin/peer_a_bouns_history_list [GET]
func PeerABounsHistoryList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		start_time_int_str = c.PostForm("start_time")
		start_time_int, _  = strconv.Atoi(start_time_int_str)
		end_time_int_str   = c.PostForm("end_time")
		end_time_int, _    = strconv.Atoi(end_time_int_str)

		user_name = c.PostForm("user_name")

		a          = models.PeerListABouns{}
		start_time = ""
		end_time   = ""
	)
	defer func() {
		c.JSON(200, data)
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	if start_time_int == 0 || end_time_int == 0 {
		start_time = ""
		end_time = ""
	} else {
		start_time = time.Unix(int64(start_time_int), 0).Format("2006-01-02 15:04:05")
		end_time = time.Unix(int64(end_time_int), 0).Format("2006-01-02 15:04:05")
	}
	list, p, err := actuator.SelectPeerABounsList(o, page, user_name, start_time, end_time)
	if err != nil {
		data = common.NewErrorResponse(500, "请重新尝试!", a)
		return
	}
	one := models.PeerListABouns{
		Items: list,
		Page:  p,
	}
	data = common.NewResponse(one)
	return
}

// @Tags 全局状态显示
// @Accept  json
// @Produce json
// @Success 200__全局状态显示
// @router /admin/show_global_operations [GET]
func ShowGlobalOperations(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()
	)
	defer func() {
		c.JSON(200, data)
	}()
	operation_list := make([]models.GlobalOperations, 0)
	er := o.Raw("select * from global_operations").Find(&operation_list)
	if er.Error != nil {
		logs.Log.Error(er.Error.Error())
		data = common.NewErrorResponse(500, "全局控制信息获取失败!", []models.GlobalOperations{})
		return
	}
	data = common.NewResponse(operation_list)
	return
}

// @Tags 全局状态修改
// @Accept  json
// @Produce json
// @Param operation_id query string true "操作_id"
// @Success 200__全局状态修改
// @router /admin/update_global_operations [POST]
func UpdateGlobalOperations(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		operation_id = c.PostForm("operation_id")

		err error
	)
	defer func() {
		if err != nil {
			o.Rollback()
			logs.Log.Error(err)
			data = common.NewErrorResponse(500, err.Error(), nil)
			c.JSON(200, data)
		}
		c.JSON(200, data)
	}()

	o.Begin()
	ids := strings.Split(operation_id, ";")
	for _, v := range ids {
		id := strings.Split(v, "-")
		switch id[1] {
		case "1": //UPDATE 表名称 SET 列名称 = 新值 WHERE 列名称 = 某值
			if er := o.Model(&models.GlobalOperations{}).Where("id = ?", id[0]).Update("state", true); er.Error != nil {
				err = errors.New("全局控制信息更新失败!")
				return
			}
		case "2":
			o.Model(&models.GlobalOperations{}).Where("id = ?", id[0]).Update("state", false)
			if er := o.Model(&models.GlobalOperations{}).Where("id = ?", id[0]).Update("state", false); er.Error != nil {
				err = errors.New("全局控制信息更新失败!")
				return
			}
		}
	}
	o.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 每日释放任务表
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param date_time query string true "开始时间"
// @Param user_name query string true "用户名字  不搜就传空，搜索就传user_name"
// @Param user_id query string true "用户id  不搜就传空，搜索就传user_id"
// @Param state query string true "状态　1=完成的 2=未完成的"
// @Success 200__每日释放任务表 {object} models.MrsfTable_test
// @router /admin/show_one_day_mrsf [GET]
func ShowOneDayMrsf(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)
		state_int_str    = c.PostForm("state")
		state_int, _     = strconv.Atoi(state_int_str)

		date_time_int_str = c.PostForm("date_time")
		date_time_int, _  = strconv.Atoi(date_time_int_str)

		user_name = c.PostForm("user_name")
		user_id   = c.PostForm("user_id")

		date_time = ""
	)
	defer func() {
		c.JSON(200, data)
	}()
	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}

	if date_time_int == 0 {
		date_time = ""
	} else {
		date_time = time.Unix(int64(date_time_int), 0).Format("2006-01-02")
	}
	state := true //1578412800
	if state_int == 2 {
		state = false
	}
	mrsf_list, p, err := actuator.ShowMrsfTable(o, page, user_name, user_id, date_time, state)
	if err != nil {
		logs.Log.Error(err)
		data = common.NewErrorResponse(500, "请尝试刷新!", []models.MrsfTable{})
		return
	}
	m_t := models.MrsfTable{
		Items: mrsf_list,
		Page:  p,
	}
	data = common.NewResponse(m_t)
	return
}

// @Tags 每日释放任务_手动释放错误用户
// @Accept  json
// @Produce json
// @Param user_id query string true "用户id+order_id   格式－   user_id-order_id;user_id-order_id;....."
// @Success 200__每日释放任务表
// @router /admin/the_release_of_err_users [POST]
func TheReleaseOfErrUsers(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		user_id = c.PostForm("user_id")

		order_id  = ""
		user_mrsf = []string{}
	)
	defer func() {
		c.JSON(200, data)
	}()
	users := strings.Split(user_id, ";")

	for i, v := range users {
		user := strings.Split(v, "_")
		if i == 0 {
			order_id = user[1]
		}
		m_s_t := models.MrsfStateTable{}
		o.Table("mrsf_state_table").Where("order_id = ?", user[1]).First(&m_s_t)
		if m_s_t.State == false && order_id == user[1] {
			user_mrsf = append(user_mrsf, user[0])
		} else {
			data = common.NewErrorResponse(500, "只能释放未释放没有释放的用户!", nil)
			return
		}
	}
	DailyDividendAndReleaseToSomeOne(users, order_id)
	data = common.NewResponse(nil)
	return
}

// @Tags 展示_DAPP_列表
// @Accept  json
// @Produce json
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Param dapp_name query string true "dapp的应用 名字"
// @Param dapp_id query string true "dapp的应用 id"
// @Param dapp_type query string true "类型"
// @Success 200____展示_DAPP_列表 {object} models.DAPPListTest
// @router /admin/show_dapp_list [GET]
func ShowDAPPList(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		current_page_str = c.PostForm("page")
		current_page, _  = strconv.Atoi(current_page_str)
		page_size_str    = c.PostForm("pageSize")
		page_size, _     = strconv.Atoi(page_size_str)

		dapp_name = c.PostForm("dapp_name")
		dapp_id   = c.PostForm("dapp_id")
		dapp_type = c.PostForm("dapp_type")
	)
	defer func() {
		c.JSON(200, data)
	}()

	page := models.Page{
		TotalPage:   0,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       0,
	}
	dapp_list, err := actuator.SelectDAPP(o, dapp_name, dapp_id, dapp_type, &page)
	if err != nil {
		l := []models.DappTable{}
		data = common.NewErrorResponse(500, "出现错误,请再次刷新!", models.DappList{Items: l, Page: page})
		return
	} else if len(dapp_list) == 0 {
		l := []models.DappTable{}
		data = common.NewResponse(models.DappList{Items: l, Page: page})
		return
	}

	dapp := models.DappList{
		Items: dapp_list,
		Page:  page,
	}
	data = common.NewResponse(dapp)
	return
}

// @Tags 插入_DAPP
// @Accept  json
// @Produce json
// @Param dapp_name query string true "名字"
// @Param image_url query string true "图片_url"
// @Param dapp_link_address query string true "dapp的链接地址"
// @Param dapp_contract_address query string true "dapp的合约地址"
// @Param dapp_type query string true "类型"
// @Success 200____插入_DAPP
// @router /admin/insert_dapp [POST]
func InsertDAPP(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		dapp_name             = c.PostForm("dapp_name")
		image_url             = c.PostForm("image_url")
		dapp_type             = c.PostForm("dapp_type")
		dapp_link_address     = c.PostForm("dapp_link_address")
		dapp_contract_address = c.PostForm("dapp_contract_address")
	)
	defer func() {
		c.JSON(200, data)
	}()

	dapp_table := models.DappTable{}
	o.Table("dapp_table").Where("name = ?", dapp_name).First(&dapp_table)
	if dapp_table.Id != 0 {
		data = common.NewErrorResponse(500, "该名字已存在!", nil)
		return
	}

	dapp := models.DappTable{
		Name:            dapp_name,
		AgreementType:   dapp_type,
		State:           true,
		TheLinkAddress:  dapp_link_address,
		ContractAddress: dapp_contract_address,
		Image:           image_url,
		CreateTime:      time.Now().Format("2006-01-02 15:04:05"),
	}
	er := o.Create(&dapp)
	if er.Error != nil {
		data = common.NewErrorResponse(500, "新增失败,请重试!", nil)
		return
	}
	data = common.NewResponse(nil)
	return
}

// @Tags 更新_DAPP
// @Accept  json
// @Produce json
// @Param dapp_id query string true "dapp id"
// @Param dapp_name query string true "名字"
// @Param image_url query string true "图片_url"
// @Param dapp_link_address query string true "dapp的链接地址"
// @Param dapp_contract_address query string true "dapp的合约地址"
// @Param dapp_type query string true "类型"
// @Success 200____更新_DAPP
// @router /admin/update_dapp [POST]
func UpdateDAPP(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		dapp_id_str = c.PostForm("dapp_id")
		dapp_id, _  = strconv.Atoi(dapp_id_str)

		dapp_name             = c.PostForm("dapp_name")
		image_url             = c.PostForm("image_url")
		dapp_type             = c.PostForm("dapp_type")
		dapp_link_address     = c.PostForm("dapp_link_address")
		dapp_contract_address = c.PostForm("dapp_contract_address")
	)
	defer func() {
		c.JSON(200, data)
	}()

	er := o.Model(&models.DappTable{}).Where("id = ?", dapp_id).Updates(map[string]interface{}{
		"name":             dapp_name,
		"agreement_type":   dapp_type,
		"the_link_address": dapp_link_address,
		"contract_address": dapp_contract_address,
		"image":            image_url,
		"update_time":      time.Now().Format("2006-01-02 15:04:05"),
	})
	if er.Error != nil {
		data = common.NewErrorResponse(500, "更新失败,请重试!", nil)
		return
	}
	data = common.NewResponse(nil)
	return
}

// @Tags 修改状态_DAPP
// @Accept  json
// @Produce json
// @Param dapp_id query string true "dapp id"
// @Param dapp_state query string true "状态 1=true(开启)  2=false(失败)"
// @Success 200____修改状态_DAPP
// @router /admin/update_dapp_state [POST]
func UpdateDAPPState(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		dapp_id_str = c.PostForm("dapp_id")
		dapp_id, _  = strconv.Atoi(dapp_id_str)

		dapp_state_str = c.PostForm("dapp_state")
		dapp_state, _  = strconv.Atoi(dapp_state_str)

		state = true
	)
	defer func() {
		c.JSON(200, data)
	}()

	if dapp_state == 2 {
		state = false
	}
	er := o.Model(&models.DappTable{}).Where("id = ?", dapp_id).Updates(map[string]interface{}{
		"state":       state,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	})

	if er.Error != nil {
		data = common.NewErrorResponse(500, "更改状态失败,请重试!", nil)
		return
	}

	data = common.NewResponse(nil)
	return
}

// @Tags 删除_DAPP
// @Accept  json
// @Produce json
// @Param dapp_id query string true "dapp id"
// @Success 200____删除_DAPP
// @router /admin/delete_dapp [POST]
func DeleteDAPP(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		dapp_id_str = c.PostForm("dapp_id")
		dapp_id, _  = strconv.Atoi(dapp_id_str)
	)
	defer func() {
		c.JSON(200, data)
	}()

	er := o.Delete(&models.DappTable{}, "id = ?", dapp_id)
	if er.Error != nil {
		data = common.NewErrorResponse(500, "删除应用失败,请重试!", nil)
		return
	}

	data = common.NewResponse(nil)
	return
}

// @Tags 分组_展示_DAPP_列表_to_the_app
// @Accept  json
// @Produce json
// @Success 200____分组_展示_DAPP_列表_to_the_app {object} models.DAPPListTest
// @router /show_group_by_type [GET]
func ShowGroupByType(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()

		list = []models.DappTable{}
	)
	defer func() {
		c.JSON(200, data)
	}()

	m := make(map[string][]models.DappTable)

	o.Raw("select * from dapp_table").Find(&list)

	for _, v := range list {
		m[v.AgreementType] = append(m[v.AgreementType], v)
	}

	values := models.DappGroupList{}
	for i, v := range m {
		a := models.List{
			Title: i,
		}
		for _, vv := range v {
			a.Values = append(a.Values, vv)
		}
		values.Items = append(values.Items, a)
	}
	if len(values.Items) == 0 {
		l := []models.List{}
		a := models.DappGroupList{Items: l}
		data = common.NewResponse(a)
		return
	}
	data = common.NewResponse(values)
	return
}

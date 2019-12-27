package controllers

import (
	"ecology/common"
	"ecology/logs"
	"ecology/models"
	"ecology/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type EcologyIndexController struct {
	beego.Controller
}

// @Tags 生态首页展示
// @Accept  json
// @Produce json
// @Success 200___生态首页展示 {object} models.Ecology_index_ob_test
// @router /show_ecology_index [GET]
func (this *EcologyIndexController) ShowEcologyIndex() {
	var (
		data          *common.ResponseData
		account_index []models.Account
		api_url       = this.Controller.Ctx.Request.RequestURI
	)

	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	token := GetJwtValues(this.Ctx)
	user_id := token.UserID

	o := models.NewOrm()
	i, erracc := o.QueryTable("account").Filter("user_id", user_id).All(&account_index)
	if erracc != nil {
		data = common.NewErrorResponse(500, "数据库操作失败!", models.Ecology_index_obj{})
		logs.Log.Error(api_url, erracc)
		return
	}
	indexValues := models.Ecology_index_obj{}
	if 1 > i {
		models.CouShu(&indexValues)
		indexValues.Ecological_poject_bool = false
		indexValues.Super_peer_bool = false
		logs.Log.Error(api_url, "没有查询到用户的生态仓库!")
		data = common.NewResponse(indexValues)
		return
	}
	indexValues.Ecological_poject_bool = true
	if len(account_index) > 0 {
		for _, v := range account_index {
			var formula_index []models.Formula
			_, errfor := o.QueryTable("formula").Filter("ecology_id", v.Id).All(&formula_index)
			if errfor != nil {
				data = common.NewErrorResponse(500, "数据库操作失败!", models.Ecology_index_obj{})
				logs.Log.Error(api_url, errfor)
				return
			}
			f := models.Formulaindex{
				Id:             v.Id,
				Level:          v.Level,
				BockedBalance:  v.BockedBalance,
				Balance:        v.Balance,
				ReturnMultiple: formula_index[0].ReturnMultiple,
				//ToDayRate:           formula_index[0].HoldReturnRate + formula_index[0].RecommendReturnRate + formula_index[0].TeamReturnRate,
				HoldReturnRate: formula_index[0].HoldReturnRate * v.Balance,
				//RecommendReturnRate: formula_index[0].RecommendReturnRate,
				//TeamReturnRate: formula_index[0].TeamReturnRate,
			}
			zhitui, err := models.RecommendReturnRate(user_id, time.Now().Format("2006-01-02")+" 00:00:00")
			if err != nil {
				models.CouShu(&indexValues)
				logs.Log.Error(api_url, "计算用户当前直推收益出错!")
				data = common.NewErrorResponse(500, "计算用户当前直推收益出错!", models.Ecology_index_obj{})
				return
			}
			to_day_rate := zhitui + f.TeamReturnRate + f.HoldReturnRate
			f.ToDayRate = to_day_rate
			f.RecommendReturnRate = zhitui
			team_coins, err_team := SumTeamProfit(user_id)
			if err_team != nil {
				models.CouShu(&indexValues)
				logs.Log.Error(api_url, err_team.Error())
				data = common.NewErrorResponse(500, err_team.Error(), models.Ecology_index_obj{})
				return
			}
			f.TeamReturnRate = team_coins
			indexValues.Ecological_poject = append(indexValues.Ecological_poject, f)
		}
	}
	tfors, err_tfor := PingSelectTforNumber(user_id)
	if err_tfor != nil {
		data = common.NewErrorResponse(500, "查看钱包　TFOR 数量时错误!", models.Ecology_index_obj{})
		return
	}
	models.SuperLevelSet(user_id, &indexValues, tfors)
	data = common.NewResponse(indexValues)
	return
}

// @Tags 新增生态仓库
// @Accept  json
// @Produce json
// @Param user_id query string true "当前用户的id   ---- 放在 header"
// @Param coin_number query string true "铸(发)币的数量"
// @Param levelstr query string true "等级数据"
// @Success 200____新增生态仓库
// @router /create_new_warehouse [POST]
func (this *EcologyIndexController) CreateNewWarehouse() {
	var (
		data    *common.ResponseData
		o       = models.NewOrm()
		api_url = this.Controller.Ctx.Request.RequestURI
	)

	coin_number_str := this.GetString("coin_number")
	coin_number, _ := strconv.ParseFloat(coin_number_str, 64)
	levelstr := this.GetString("levelstr")

	err_orm_begin := o.Begin()
	if err_orm_begin != nil {
		logs.Log.Error(api_url, err_orm_begin)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	jwtValues := GetJwtValues(this.Ctx)
	user_id := jwtValues.UserID

	//项目生成
	account := models.Account{
		UserId:        user_id,
		Balance:       0,
		Currency:      "USDD",
		BockedBalance: 0,
		Level:         levelstr,
	}
	_, err_acc := o.Insert(&account)

	if err_acc != nil {
		logs.Log.Error(api_url, err_acc)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	//生态项目的算力表
	formula := models.Formula{}
	err_level := models.JudgeLevel(o, user_id, levelstr, &formula)
	if err_level != nil {
		o.Rollback()
		logs.Log.Error(api_url, err_level)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	formula.EcologyId = account.Id
	_, err_insert := o.Insert(&formula)
	if err_insert != nil {
		o.Rollback()
		logs.Log.Error(api_url, err_insert)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	//任务表 TFOR
	order_id := utils.TimeUUID()
	acc_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  false,
		WalletState: false,
		UserId:      user_id,
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: 0,
		InCome:      coin_number,
	}
	_, errtxid_acc := o.Insert(&acc_txid_dcmt)
	if errtxid_acc != nil {
		o.Rollback()
		logs.Log.Error(api_url, errtxid_acc)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	o.Commit()

	token := this.Ctx.Request.Header.Get("Authorization")
	err := models.PingWalletAdd(token, coin_number)
	if err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewErrorResponse(500, err.Error(), nil)
		return
	}

	oo := models.NewOrm()
	oo.Begin()
	//TFOR交易记录 - 更新生态仓库的交易余额
	err_acc_d := models.NewCreateAndSaveAcc_d(oo, user_id, "普通转入", order_id, 0, coin_number, account.Id)
	if err_acc_d != nil {
		logs.Log.Error(api_url, err_acc_d)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	//铸币交易记录
	err_blo_d := models.NewCreateAndSaveBlo_d(oo, user_id, "转入铸币", order_id, 0, coin_number, account.Id)
	if err_blo_d != nil {
		logs.Log.Error(api_url, err_blo_d)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	oo.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 转USDD到生态仓库
// @Accept  json
// @Produce json
// @Param user_id query string true "用户id   ---- 放在 header"
// @Param ecology_id query string int "生态仓库id的id"
// @Param order_id query string int "交易的id"
// @Param coin_number query string true "铸(发)币的数量"
// @Success 200___转USDD到生态仓库
// @router /to_change_into_USDD [POST]
func (this *EcologyIndexController) ToChangeIntoUSDD() {
	var (
		data    *common.ResponseData
		o       = models.NewOrm()
		api_url = this.Controller.Ctx.Request.RequestURI
	)

	coin_number_str := this.GetString("coin_number")
	order_id := this.GetString("order_id")
	coin_number, _ := strconv.ParseFloat(coin_number_str, 64)
	ecology_id, _ := this.GetInt("ecology_id")
	jwtValues := GetJwtValues(this.Ctx)
	user_id := jwtValues.UserID

	order_if := models.TxIdList{
		TxId: order_id,
	}
	o.Read(&order_if)
	if order_if.Id != 0 {
		logs.Log.Error(api_url, "重复的交易－多次提交")
		data = common.NewErrorResponse(500, "订单以存在，请勿提交!!", nil)
		return
	}

	err_orm_begin := o.Begin()
	if err_orm_begin != nil {
		logs.Log.Error(api_url, err_orm_begin)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	//任务表 USDD  铸币记录
	formula := models.Formula{}
	err_for := o.QueryTable("formula").Filter("ecology_id", ecology_id).One(&formula)
	if err_for != nil {
		o.Rollback()
		logs.Log.Error(api_url, err_for)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  false,
		WalletState: false,
		UserId:      user_id,
		Comment:     "转入交易",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: 0,
		InCome:      formula.ReturnMultiple * coin_number,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		o.Rollback()
		logs.Log.Error(api_url, errtxid_blo)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	o.Commit()

	token := this.Ctx.Request.Header.Get("Authorization")
	if err := models.PingWalletAdd(token, coin_number); err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewErrorResponse(500, err.Error(), nil)
		return
	}
	order := models.TxIdList{
		TxId: order_id,
	}
	o.Read(&order, "tx_id")
	order.WalletState = true
	o.Update(&order)

	oo := models.NewOrm()
	oo.Begin()
	//TFOR交易记录 - 更新生态仓库的交易余额
	err_acc_d := models.FindLimitOneAndSaveAcc_d(oo, user_id, "普通转入", order_id, 0, coin_number, ecology_id)
	if err_acc_d != nil {
		logs.Log.Error(api_url, err_acc_d)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	//铸币交易记录
	err_blo_d := models.FindLimitOneAndSaveBlo_d(oo, user_id, "转入铸币", order_id, 0, coin_number, ecology_id)
	if err_blo_d != nil {
		logs.Log.Error(api_url, err_blo_d)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	oo.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 升级生态仓库
// @Accept  json
// @Produce json
// @Param user_id query string true "用户id   ---- 放在 header"
// @Param ecology_id query string int "生态仓库id的id"
// @Param order_id query string int "交易的id"
// @Param cion_number query string true "铸(发)币的数量"
// @Param levelstr query string true "升级后的等级"
// @Success 200____升级生态仓库
// @router /upgrade_warehouse [POST]
func (this *EcologyIndexController) UpgradeWarehouse() {
	var (
		data    *common.ResponseData
		o       = models.NewOrm()
		api_url = this.Controller.Ctx.Request.RequestURI
	)

	coin_number_str := this.GetString("cion_number")
	order_id := this.GetString("order_id")
	coin_number, _ := strconv.ParseFloat(coin_number_str, 64)
	ecology_id, _ := this.GetInt("ecology_id")
	levelstr := this.GetString("levelstr")
	jwtValues := GetJwtValues(this.Ctx)
	user_id := jwtValues.UserID

	order_if := models.TxIdList{
		TxId: order_id,
	}
	o.Read(&order_if)
	if order_if.Id != 0 {
		logs.Log.Error(api_url, "重复的交易－多次提交")
		data = common.NewErrorResponse(500, "订单以存在，请勿提交!!", nil)
		return
	}

	err_orm_begin := o.Begin()
	if err_orm_begin != nil {
		logs.Log.Error(api_url, err_orm_begin)
		data = common.NewErrorResponse(500, "数据库操作失败", nil)
		return
	}

	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	formula_table := models.ForceTable{
		Level: levelstr,
	}
	err_r := o.Read(&formula_table, "level")
	if err_r != nil || float64(formula_table.LowHold) > coin_number {
		logs.Log.Error(api_url, err_r)
		data = common.NewErrorResponse(500, "不满足升级条件,请填入规定的升级金额", nil)
		return
	}

	formula := models.Formula{EcologyId: ecology_id}
	err_read := o.Read(&formula, "ecology_id")
	if err_read != nil {
		o.Rollback()
		logs.Log.Error(api_url, err_read)
		data = common.NewErrorResponse(500, "数据库操作失败", nil)
		return
	}

	errJu := models.JudgeLevel(o, user_id, levelstr, &formula)
	if errJu != nil {
		o.Rollback()
		logs.Log.Error(api_url, errJu)
		data = common.NewErrorResponse(500, errJu.Error(), nil)
		return
	}

	_, err_up_acc := o.QueryTable("account").Filter("id", ecology_id).Update(orm.Params{"level": levelstr})
	if err_up_acc != nil {
		logs.Log.Error(api_url, err_up_acc)
		data = common.NewErrorResponse(500, "数据库操作失败", nil)
		return
	}

	if _, err_up_for := o.Update(&formula, "level", "low_hold", "high_hold", "return_multiple", "hold_return_rate", "recommend_return_rate", "team_return_rate"); err_up_for != nil {
		o.Rollback()
		logs.Log.Error(api_url, err_up_for)
		data = common.NewErrorResponse(500, "数据库操作失败", nil)
		return
	}
	//任务表 USDD  铸币记录
	blo_txid_dcmt := models.TxIdList{
		TxId:        order_id,
		OrderState:  false,
		WalletState: false,
		UserId:      user_id,
		Comment:     "升级交易",
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Expenditure: 0,
		InCome:      formula.ReturnMultiple * coin_number,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		o.Rollback()
		logs.Log.Error(api_url, errtxid_blo)
		data = common.NewErrorResponse(500, "数据库操作失败", nil)
		return
	}

	o.Commit()

	//钱包操作
	token := this.Ctx.Request.Header.Get("Authorization")
	if err := models.PingWalletAdd(token, coin_number); err != nil {
		logs.Log.Error(api_url, err)
		data = common.NewErrorResponse(500, err.Error(), nil)
		return
	}
	order := models.TxIdList{
		TxId: order_id,
	}
	err_rea := o.Read(&order, "tx_id")
	if err_rea != nil {
		logs.Log.Error(api_url, err_rea)
		data = common.NewErrorResponse(500, err_rea.Error(), nil)
		return
	}
	order.WalletState = true
	_, err_up := o.Update(&order)
	if err_up != nil {
		logs.Log.Error(api_url, err_up)
		data = common.NewErrorResponse(500, err_up.Error(), nil)
		return
	}

	oo := models.NewOrm()
	err_oo := oo.Begin()
	if err_oo != nil {
		logs.Log.Error(api_url, err_oo)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	//TFOR交易记录 - 更新生态仓库的交易余额
	err_acc_d := models.FindLimitOneAndSaveAcc_d(oo, user_id, "升级转入", order_id, 0, coin_number, ecology_id)
	if err_acc_d != nil {
		logs.Log.Error(api_url, err_acc_d)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}

	//铸币交易记录
	err_blo_d := models.FindLimitOneAndSaveBlo_d(oo, user_id, "升级铸币", order_id, 0, coin_number, ecology_id)
	if err_blo_d != nil {
		logs.Log.Error(api_url, err_blo_d)
		data = common.NewErrorResponse(500, "数据库操作失败!", nil)
		return
	}
	oo.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 交易的历史记录
// @Accept  json
// @Produce json
// @Param ecology_id query string int "生态仓库id的id"
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /return_page_list_hostry [GET]
func (this *EcologyIndexController) ReturnPageListHostry() {
	var (
		data            *common.ResponseData
		ecology_id, _   = this.GetInt("ecology_id")
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
	values, p, err := models.SelectHostery(ecology_id, page)
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

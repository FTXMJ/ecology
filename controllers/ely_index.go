package controllers

import (
	"ecology1/common"
	"ecology1/models"
	"ecology1/utils"
	"github.com/astaxie/beego"
)

type EcologyIndexController struct {
	beego.Controller
}

// @Tags 生态首页展示
// @Accept  json
// @Produce json
// @Param user_id query string true "当前用户的id"
// @Success 200  生态首页展示
// @router /show_ecology_index___生态首页展示 [Post]
func (this *EcologyIndexController) ShowEcologyIndex() {
	var (
		data          *common.ResponseData
		account_index []models.Account
		user_id       = this.GetString("user_id")
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	//TODO 调用老罗的钱包接口
	usdd := 0.0
	i, erracc := models.NewOrm().QueryTable("account").Filter("user_id", user_id).All(&account_index)
	if erracc != nil {
		data = common.NewErrorResponse(500)
		//TODO log
		return
	}
	indexValues := models.Ecology_index_obj{}
	if 1 > i {
		indexValues.Usdd = usdd
		indexValues.Ecological_poject_bool = false
		indexValues.Super_peer_bool = false
		data = common.NewResponse(indexValues)
		return
	}
	indexValues.Usdd = usdd
	indexValues.Ecological_poject_bool = true
	for o, v := range account_index {
		var formula_index models.Formula
		errfor := models.NewOrm().QueryTable("formula").Filter("ecology_id", v.Id).One(&formula_index)
		if errfor != nil {
			data = common.NewErrorResponse(500)
			//TODO log
			return
		}
		indexValues.Ecological_poject[o].Level = v.Level
		indexValues.Ecological_poject[o].BockedBalance = v.BockedBalance
		indexValues.Ecological_poject[o].ReturnMultiple = formula_index.ReturnMultiple
		indexValues.Ecological_poject[o].ToDayRate = formula_index.HoldReturnRate + formula_index.ReturnMultiple + formula_index.TeamReturnRate
		indexValues.Ecological_poject[o].HoldReturnRate = formula_index.HoldReturnRate
		indexValues.Ecological_poject[o].RecommendReturnRate = formula_index.RecommendReturnRate
		indexValues.Ecological_poject[o].TeamReturnRate = formula_index.TeamReturnRate
	}
	models.SuperLevelSet(usdd, &indexValues)
	data = common.NewResponse(indexValues)
	return
}

// @Tags 新增生态仓库
// @Accept  json
// @Produce json
// @Param user_id query string true "当前用户的id   ---- 放在 header"
// @Param coin_number query string true "铸(发)币的数量"
// @Param levelstr query string true "等级数据"
// @Param ecology_id query string true "生态仓库id"
// @Success 200
// @router /create_new_warehouse__新增生态仓库 [Post]
func (this *EcologyIndexController) CreateNewWarehouse() {
	var (
		data *common.ResponseData
		o    = models.NewOrm()
	)
	net_values := models.NetValueType{}
	this.ParseForm(&net_values)
	coin_number := net_values.CoinNumber
	levelstr := net_values.LevelStr

	err_orm_begin := o.Begin()
	if err_orm_begin != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
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
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//生态项目的算力表
	formula := models.Formula{}
	err_level := models.JudgeLevel(o,user_id,levelstr, &formula)
	if err_level != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}
	formula.EcologyId = account.Id
	_, err_insert := o.Insert(&formula)
	if err_insert != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//任务表 TFOR
	tx_id_acc_d := utils.Shengchengstr("转入记录", user_id, "USDD")
	acc_txid_dcmt := models.TxIdList{
		State: "false",
		TxId:  tx_id_acc_d,
	}
	_, errtxid_acc := o.Insert(&acc_txid_dcmt)
	if errtxid_acc != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//任务表 USDD
	tx_id_blo_d := utils.Shengchengstr("铸币记录", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		State: "false",
		TxId:  tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//TFOR交易记录 - 更新生态仓库的交易余额
	err_acc_d := models.FindLimitOneAndSaveAcc_d(o, user_id, "新增生态仓库转入-USDD", tx_id_acc_d, 0, coin_number, account.Id)
	if err_acc_d != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//铸币交易记录
	err_blo_d := models.FindLimitOneAndSaveBlo_d(o, user_id, "生态仓库铸币", tx_id_blo_d, 0, coin_number, account.Id)
	if err_blo_d != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}
	o.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 转USDD到生态仓库
// @Accept  json
// @Produce json
// @Param user_id query string true "用户id   ---- 放在 header"
// @Param ecology_id query string int "生态仓库id的id"
// @Param coin_number query string true "铸(发)币的数量"
// @Param levelstr query string true "等级数据"
// @Success 200
// @router /to_change_into_USDD__转USDD到生态仓库 [Post]
func (this *EcologyIndexController) ToChangeIntoUSDD() {
	var (
		data *common.ResponseData
		o    = models.NewOrm()
	)
	net_values := models.NetValueType{}
	this.ParseForm(&net_values)
	coin_number := net_values.CoinNumber
	ecology_id := net_values.EcologyId
	jwtValues := GetJwtValues(this.Ctx)
	user_id := jwtValues.UserID

	err_orm_begin := o.Begin()
	if err_orm_begin != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	//任务表 USDD   转入记录
	tx_id_acc_d := utils.Shengchengstr("转入记录", user_id, "USDD")
	acc_txid_dcmt := models.TxIdList{
		State: "false",
		TxId:  tx_id_acc_d,
	}
	_, errtxid_acc := o.Insert(&acc_txid_dcmt)
	if errtxid_acc != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//任务表 USDD  铸币记录
	tx_id_blo_d := utils.Shengchengstr("铸币记录", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		State: "false",
		TxId:  tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//TFOR交易记录 - 更新生态仓库的交易余额
	err_acc_d := models.FindLimitOneAndSaveAcc_d(o, user_id, "新增生态仓库转入-USDD", tx_id_acc_d, 0, coin_number, ecology_id)
	if err_acc_d != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//铸币交易记录
	err_blo_d := models.FindLimitOneAndSaveBlo_d(o, user_id, "生态仓库铸币", tx_id_blo_d, 0, coin_number, ecology_id)
	if err_blo_d != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}
	o.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 升级生态仓库
// @Accept  json
// @Produce json
// @Param user_id query string true "用户id   ---- 放在 header"
// @Param ecology_id query string int "生态仓库id的id"
// @Param cion_number query string true "铸(发)币的数量"
// @Param levelstr query string true "升级后的等级"
// @Success 200
// @router /upgrade_warehouse__升级生态仓库 [Post]
func (this *EcologyIndexController) UpgradeWarehouse() {
	var (
		data *common.ResponseData
		o    = models.NewOrm()
	)
	net_values := models.NetValueType{}
	this.ParseForm(&net_values)
	coin_number := net_values.CoinNumber
	ecology_id := net_values.EcologyId
	levelstr := net_values.LevelStr
	jwtValues := GetJwtValues(this.Ctx)
	user_id := jwtValues.UserID

	err_orm_begin := o.Begin()
	if err_orm_begin != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	formula := models.Formula{EcologyId: ecology_id}
	err_read := o.Read(&formula)
	if err_read != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}
	if _, err_up_for := o.Update(&formula, "level", "low_hold", "high_hold", "return_multiple", "hold_return_rate", "recommend_return_rate", "team_return_rate"); err_up_for != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	errJu := models.JudgeLevel(o,user_id,levelstr, &formula)
	if errJu!=nil{
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//任务表 USDD   转入记录
	tx_id_acc_d := utils.Shengchengstr("转入记录", user_id, "USDD")
	acc_txid_dcmt := models.TxIdList{
		State: "false",
		TxId:  tx_id_acc_d,
	}
	_, errtxid_acc := o.Insert(&acc_txid_dcmt)
	if errtxid_acc != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//任务表 USDD  铸币记录
	tx_id_blo_d := utils.Shengchengstr("铸币记录", user_id, "USDD")
	blo_txid_dcmt := models.TxIdList{
		State: "false",
		TxId:  tx_id_blo_d,
	}
	_, errtxid_blo := o.Insert(&blo_txid_dcmt)
	if errtxid_blo != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//TFOR交易记录 - 更新生态仓库的交易余额
	err_acc_d := models.FindLimitOneAndSaveAcc_d(o, user_id, "升级生态仓库　转入-USDD", tx_id_acc_d, 0, coin_number, ecology_id)
	if err_acc_d != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}

	//铸币交易记录
	err_blo_d := models.FindLimitOneAndSaveBlo_d(o, user_id, "生态仓库铸币", tx_id_blo_d, 0, coin_number, ecology_id)
	if err_blo_d != nil {
		//TODO logs
		data = common.NewErrorResponse(500)
		return
	}
	o.Commit()
	data = common.NewResponse(nil)
	return
}

// @Tags 交易的历史记录
// @Accept  json
// @Produce json
// @Param ecology_id query string int "生态仓库id的id"
// @Param page_count query string true "分页信息　－　数据总条数"
// @Param total_page query string true "分页信息　－　总页数"
// @Param current_page query string true "分页信息　－　当前页数"
// @Param page_size query string true "分页信息　－　每页数据量"
// @Success 200
// @router /upgrade_warehouse__交易的历史记录 [Post]
func (this *EcologyIndexController) ReturnPageListHostry() {
	var (
		data            *common.ResponseData
		ecology_id, _   = this.GetInt("ecology_id")
		page_count, _   = this.GetInt("page_count")
		total_page, _   = this.GetInt("total_page")
		current_page, _ = this.GetInt("current_page")
		page_size, _    = this.GetInt("page_size")
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	page := models.Page{
		TotalPage:   total_page,
		CurrentPage: current_page,
		PageSize:    page_size,
		Count:       page_count,
	}
	values, p, err := models.SelectHostery(ecology_id, page)
	if err != nil {
		data = common.NewErrorResponse(500)
		return
	}
	hostory_list := models.HostryPageInfo{
		Items: values,
		Page:  p,
	}
	data = common.NewResponse(hostory_list)
	return
}

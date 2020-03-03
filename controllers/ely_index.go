package controllers

import (
	"ecology/actuator"
	"ecology/common"
	db "ecology/db"
	"ecology/filter"
	"ecology/logs"
	"ecology/models"
	"ecology/utils"
	"github.com/gin-gonic/gin"

	"errors"
	"strconv"
	"time"
)

// @Tags 生态首页展示
// @Accept  json
// @Produce json
// @Success 200___生态首页展示 {object} models.Ecology_index_ob_test
// @router /show_ecology_index [GET]
func ShowEcologyIndex(c *gin.Context) {
	var (
		data        *common.ResponseData
		account     = models.Account{}
		indexValues = models.Ecology_index_obj{Ecological_poject_bool: true}
		err         error
		o           = db.NewEcologyOrm()

		token   = filter.GetJwtValues(c)
		user_id = token.UserID
	)
	defer func() {
		if err != nil {
			logs.Log.Error(err)
			data = common.NewErrorResponse(500, "数据库操作失败!", models.Ecology_index_obj{})
			c.JSON(200, data)
		} else {
			c.JSON(200, data)
		}
	}()

	if er := o.Raw("select * from account where user_id=?", user_id).Find(&account); er.Error != nil {
		err = er.Error
		return
	}

	if err = actuator.TheWheel(o, user_id, account, &indexValues); err != nil {
		return
	}

	tfors := 0.0
	if _, tfors, err = actuator.PingSelectTforNumber(user_id); err != nil {
		return
	}

	indexValues.Usdd = actuator.AddAllSum(o, user_id)
	actuator.SuperLevelSet(o, user_id, &indexValues, tfors)
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
func CreateNewWarehouse(c *gin.Context) {
	var (
		data            *common.ResponseData
		o               = db.NewEcologyOrm()
		coin_number_str = c.GetString("coin_number")
		coin_number, _  = strconv.ParseFloat(coin_number_str, 64)
		levelstr        = c.GetString("levelstr")
		jwtValues       = filter.GetJwtValues(c)
		user_id         = jwtValues.UserID
		err             error
	)

	defer func() {
		if err != nil {
			o.Rollback()
			logs.Log.Error(err)
			data = common.NewErrorResponse(500, "数据库操作失败!", nil)
			c.JSON(200, data)
		}
		c.JSON(200, data)
	}()

	o.Begin()
	//项目生成
	account := models.Account{
		UserId:   user_id,
		Currency: "USDD",
		Level:    levelstr,
	}
	if er := o.Create(&account); er.Error != nil {
		err = er.Error
		return
	}

	//生态项目的算力表
	formula := models.Formula{}
	if err = actuator.JudgeLevel(o, user_id, levelstr, &formula); err != nil {
		return
	}
	formula.EcologyId = account.Id
	if er := o.Create(&formula); er.Error != nil {
		err = er.Error
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
	if er := o.Create(&acc_txid_dcmt); er.Error != nil {
		err = er.Error
		return
	}
	o.Commit()

	token := c.Request.Header.Get("Authorization")
	if err = actuator.PingWalletAdd(token, coin_number); err != nil {
		return
	}

	o.Begin()
	//TFOR交易记录 - 更新生态仓库的交易余额
	if err = actuator.NewCreateAndSaveAcc_d(o, user_id, "普通转入", order_id, 0, coin_number, account.Id); err != nil {
		return
	}
	//铸币交易记录
	if err = actuator.NewCreateAndSaveBlo_d(o, user_id, "转入铸币", order_id, 0, coin_number, account.Id); err != nil {
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
// @Param order_id query string int "交易的id"
// @Param coin_number query string true "铸(发)币的数量"
// @Success 200___转USDD到生态仓库
// @router /to_change_into_USDD [POST]
func ToChangeIntoUSDD(c *gin.Context) {
	var (
		data            *common.ResponseData
		o               = db.NewEcologyOrm()
		coin_number_str = c.PostForm("coin_number")
		order_id        = c.PostForm("order_id")
		coin_number, _  = strconv.ParseFloat(coin_number_str, 64)
		ecology_id_str  = c.PostForm("ecology_id")
		ecology_id, _   = strconv.Atoi(ecology_id_str)
		jwtValues       = filter.GetJwtValues(c)
		user_id         = jwtValues.UserID
		err             error
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
	order_if := models.TxIdList{}

	if er := o.Table("tx_id_list").Where("tx_id = ?", order_id).First(&order_if); er.Error != nil {
		if er.Error.Error() != "<QuerySeter> no row found" {
			err = errors.New("订单以存在，请勿提交!!")
			return
		}
	} else if order_if.Id > 0 {
		err = errors.New("订单以存在，请勿提交!!")
		return
	}

	//任务表 USDD  铸币记录
	formula := models.Formula{}
	o.Raw("select * from formula where ecology_id=?", ecology_id).First(&formula)
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
	if er := o.Create(&blo_txid_dcmt); er.Error != nil {
		err = errors.New("数据库操作失败!")
		return
	}
	o.Commit()

	token := c.Request.Header.Get("Authorization")
	if err = actuator.PingWalletAdd(token, coin_number); err != nil {
		return
	}
	order := models.TxIdList{}
	o.Table("tx_id_list").Where("tx_id = ?", order_id).First(&order)
	order.WalletState = true
	o.Update(&order)

	o.Begin()
	//TFOR交易记录 - 更新生态仓库的交易余额
	if err = actuator.FindLimitOneAndSaveAcc_d(o, user_id, "普通转入", order_id, 0, coin_number, ecology_id); err != nil {
		err = errors.New("数据库操作失败!!")
		return
	}
	//铸币交易记录
	if err = actuator.FindLimitOneAndSaveBlo_d(o, user_id, "转入铸币", order_id, 0, coin_number, ecology_id); err != nil {
		err = errors.New("数据库操作失败!!")
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
// @Param order_id query string int "交易的id"
// @Param cion_number query string true "铸(发)币的数量"
// @Param levelstr query string true "升级后的等级"
// @Success 200____升级生态仓库
// @router /upgrade_warehouse [POST]
func UpgradeWarehouse(c *gin.Context) {
	var (
		data            *common.ResponseData
		o               = db.NewEcologyOrm()
		coin_number_str = c.GetString("cion_number")
		order_id        = c.GetString("order_id")
		coin_number, _  = strconv.ParseFloat(coin_number_str, 64)
		ecology_id      = c.GetInt("ecology_id")
		levelstr        = c.GetString("levelstr")
		jwtValues       = filter.GetJwtValues(c)
		user_id         = jwtValues.UserID
		err             error
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
	order_if := models.TxIdList{}
	if er := o.Table("tx_id_list").Where("tx_id = ?", order_id).First(&order_if); er.Error != nil {
		if er.Error.Error() != "<QuerySeter> no row found" {
			err = errors.New("订单以存在，请勿提交!!")
			return
		}
	} else if order_if.Id > 0 {
		err = errors.New("订单以存在，请勿提交!!")
		return
	}
	formula_table := models.ForceTable{}
	o.Table("formula_table").Where("level = ?", levelstr).First(&formula_table)
	if float64(formula_table.LowHold) > coin_number {
		err = errors.New("不满足升级条件,请填入规定的升级金额")
		return
	}

	formula := models.Formula{EcologyId: ecology_id}
	o.Table("formula").Where("ecology_id = ?", ecology_id).First(&formula)
	if err = actuator.JudgeLevel(o, user_id, levelstr, &formula); err != nil {
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
	if er := o.Create(&blo_txid_dcmt); er.Error != nil {
		err = errors.New("数据库操作失败!")
		return
	}
	o.Commit()

	//钱包操作
	token := c.Request.Header.Get("Authorization")
	if err := actuator.PingWalletAdd(token, coin_number); err != nil {
		return
	}

	if er := o.Model(&models.TxIdList{}).Where("tx_id = ?", order_id).Update("wallet_state", true); er.Error != nil {
		err = errors.New("数据库操作失败!")
		return
	}

	o.Begin()
	o.Model(&models.Account{}).Where("id = ?", ecology_id).Update("level", levelstr)
	if er := o.Model(&models.Account{}).Where("id = ?", ecology_id).Update("level", levelstr); er.Error != nil {
		err = errors.New("数据库操作失败!")
		return
	}
	o.Save(&formula)
	if er := o.Save(&formula); er.Error != nil {
		err = errors.New("数据库操作失败!")
		return
	}
	//TFOR交易记录 - 更新生态仓库的交易余额
	if err = actuator.FindLimitOneAndSaveAcc_d(o, user_id, "升级转入", order_id, 0, coin_number, ecology_id); err != nil {
		err = errors.New("数据库操作失败!")
		return
	}
	//铸币交易记录
	if err = actuator.FindLimitOneAndSaveBlo_d(o, user_id, "升级铸币", order_id, 0, coin_number, ecology_id); err != nil {
		err = errors.New("数据库操作失败!")
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
// @Param page query string true "分页信息　－　当前页数"
// @Param pageSize query string true "分页信息　－　每页数据量"
// @Success 200____交易的历史记录 {object} models.HostryPageInfo_test
// @router /return_page_list_hostry [GET]
func ReturnPageListHostry(c *gin.Context) {
	var (
		data         *common.ResponseData
		ecology_id   = c.GetInt("ecology_id")
		current_page = c.GetInt("page")
		page_size    = c.GetInt("pageSize")
	)
	defer func() {
		c.JSON(200, data)
	}()
	page := models.Page{
		CurrentPage: current_page,
		PageSize:    page_size,
	}
	values, p, err := actuator.SelectHostery(ecology_id, page)
	if err != nil {
		logs.Log.Error(err)
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

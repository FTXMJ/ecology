package controllers

import (
	"ecology/common"
	db "ecology/db"
	"ecology/models"
	"github.com/gin-gonic/gin"
)

// @Tags 心跳检测
// @Accept  json
// @Produce json
// @Success 200
// @router /check [GET]
func Check(c *gin.Context) {
	var data *common.ResponseData
	data = common.NewResponse(nil)

	c.JSON(200, data)
	return
}

// @Tags 给定实时数据
// @Accept  json
// @Produce json
// @Success 200____给定实时数据 {object} models.RealTimePriceTest
// @router /ticker [GET]
func ShowRealTimePrice(c *gin.Context) {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()
	)

	value := []models.QuoteTicker{}
	o.Raw("select * from quote_ticker").Find(&value)

	data = common.NewResponse(value)
	c.JSON(200, data)
	return
}

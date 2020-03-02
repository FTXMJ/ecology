package routers

import (
	"ecology/conf"
	"ecology/controllers"
	"ecology/filter"
	"ecology/logs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func Router() {
	router := gin.Default()
	router.GET("/api/v1/ecology/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/check", controllers.Check)
	router.Use(Cors())
	router.Use(filter.HTTPInterceptor())
	router.Use(logs.LoggerToFile())

	// 文档界面访问URL
	// http://127.0.0.1:2019/api/v1/article/swagger/index.html
	v1 := router.Group("/api/v1/ecology")

	AuthAPI(v1)

	http.Handle("/", router)

	gin.SetMode(gin.ReleaseMode)
	errrr := router.Run(":" + conf.ConfInfo.GinIpPort)
	if errrr != nil {
		logs.Log.Error(errrr, "程序启动失败  --- ecology")
	}
}

func AuthAPI(route *gin.RouterGroup) {
	// 生态页面展示－升级－充值－历史记录
	route.GET("/show_ecology_index", controllers.ShowEcologyIndex)
	route.POST("/to_change_into_USDD", controllers.ToChangeIntoUSDD)
	route.POST("/upgrade_warehouse", controllers.UpgradeWarehouse)
	route.POST("/return_page_list_hostry", controllers.ReturnPageListHostry)

	// 后台操作模块
	route.POST("/admin/operation_formula_list", controllers.OperationFormulaList)
	route.POST("/admin/operation_super_formula_list", controllers.OperationSuperFormulaList)
	route.GET("/admin/return_page_hostry_root", controllers.ReturnPageHostryRoot)
	route.GET("/admin/filter_history_info", controllers.FilterHistoryInfo)
	route.GET("/admin/user_ecology_list", controllers.UserEcologyList)
	route.GET("/admin/user_ecology_false_list", controllers.UserEcologyFalseList)
	route.GET("/admin/computational_flow", controllers.ComputationalFlow)
	route.GET("/admin/ecological_income_control", controllers.EcologicalIncomeControl)
	route.POST("/admin/ecological_income_control_update", controllers.EcologicalIncomeControlUpdate)
	route.GET("/admin/peer_user_list", controllers.PeerUserList)
	route.GET("/admin/peer_a_bouns_list", controllers.PeerABounsList)
	route.GET("/admin/peer_a_bouns_history_list", controllers.PeerABounsHistoryList)
	route.GET("/admin/show_global_operations", controllers.ShowGlobalOperations)
	route.POST("/admin/update_global_operations", controllers.UpdateGlobalOperations)
	route.GET("/admin/show_one_day_mrsf", controllers.ShowOneDayMrsf)
	route.POST("/admin/the_release_of_err_users", controllers.TheReleaseOfErrUsers)

	// dapp 操作相关
	route.GET("/admin/show_dapp_list", controllers.ShowDAPPList)
	route.POST("/admin/insert_dapp", controllers.InsertDAPP)
	route.POST("/admin/update_dapp", controllers.UpdateDAPP)
	route.POST("/admin/update_dapp_state", controllers.UpdateDAPPState)
	route.POST("/admin/delete_dapp", controllers.DeleteDAPP)
	route.POST("/admin/show_group_by_type", controllers.ShowGroupByType)
	route.POST("/admin/test_mrsf", controllers.DailyDividendAndReleaseTest)

	// 展示　普通算力表__个人等级info___超级节点算力表显示
	route.GET("/show_formula_list", controllers.ShowFormulaList)
	route.GET("/show_user_formula", controllers.ShowUserFormula)
	route.GET("/show_super_formula_list", controllers.ShowSuperFormulaList)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

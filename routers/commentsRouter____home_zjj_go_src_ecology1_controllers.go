package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"],
		beego.ControllerComments{
			Method:           "OperationFormulaList",
			Router:           `/operation_formula_list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"],
		beego.ControllerComments{
			Method:           "OperationSuperFormulaList",
			Router:           `/operation_super_formula_list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"],
		beego.ControllerComments{
			Method:           "ShowFormulaList",
			Router:           `/show_formula_list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology/controllers:BackStageManagement"],
		beego.ControllerComments{
			Method:           "ShowSuperFormulaList",
			Router:           `/show_super_formula_list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"],
		beego.ControllerComments{
			Method:           "CreateNewWarehouse",
			Router:           `/create_new_warehouse`,
			AllowHTTPMethods: []string{"Post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"],
		beego.ControllerComments{
			Method:           "ShowEcologyIndex",
			Router:           `/show_ecology_index`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"],
		beego.ControllerComments{
			Method:           "ToChangeIntoUSDD",
			Router:           `/to_change_into_USDD`,
			AllowHTTPMethods: []string{"Post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"],
		beego.ControllerComments{
			Method:           "UpgradeWarehouse",
			Router:           `/upgrade_warehouse`,
			AllowHTTPMethods: []string{"Post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology/controllers:EcologyIndexController"],
		beego.ControllerComments{
			Method:           "ReturnPageListHostry",
			Router:           `/upgrade_warehouse`,
			AllowHTTPMethods: []string{"Post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:FirstController"] = append(beego.GlobalControllerRouter["ecology/controllers:FirstController"],
		beego.ControllerComments{
			Method:           "Check",
			Router:           `/check`,
			AllowHTTPMethods: []string{"GET"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:FirstController"] = append(beego.GlobalControllerRouter["ecology/controllers:FirstController"],
		beego.ControllerComments{
			Method:           "CreateUserAbout",
			Router:           `/create_user`,
			AllowHTTPMethods: []string{"POST"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["ecology/controllers:FirstController"] = append(beego.GlobalControllerRouter["ecology/controllers:FirstController"],
		beego.ControllerComments{
			Method:           "DailyDividendAndRelease",
			Router:           `/yanshi`,
			AllowHTTPMethods: []string{"GET"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}

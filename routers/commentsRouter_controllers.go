package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "OperationFormulaList",
            Router: `/operation_formula_list____算力表信息修改`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "OperationSuperFormulaList",
            Router: `/operation_super_formula_list___超级节点算力表信息修改`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "ShowFormulaList",
            Router: `/show_formula_list____算力表显示___后台操作__or__用户查看__都可`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "ShowSuperFormulaList",
            Router: `/show_super_formula_list____超级节点算力表显示___后台操作__or_用户查看__都可`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "CreateNewWarehouse",
            Router: `/create_new_warehouse__新增生态仓库`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "ShowEcologyIndex",
            Router: `/show_ecology_index___生态首页展示`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "ToChangeIntoUSDD",
            Router: `/to_change_into_USDD__转USDD到生态仓库`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "ReturnPageListHostry",
            Router: `/upgrade_warehouse__交易的历史记录`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "UpgradeWarehouse",
            Router: `/upgrade_warehouse__升级生态仓库`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

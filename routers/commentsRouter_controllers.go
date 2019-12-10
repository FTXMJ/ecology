package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "OperationFormulaList",
            Router: `/operation_formula_list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "OperationSuperFormulaList",
            Router: `/operation_super_formula_list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "ShowFormulaList",
            Router: `/show_formula_list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"] = append(beego.GlobalControllerRouter["ecology1/controllers:BackStageManagement"],
        beego.ControllerComments{
            Method: "ShowSuperFormulaList",
            Router: `/show_super_formula_list`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "CreateNewWarehouse",
            Router: `/create_new_warehouse`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "ShowEcologyIndex",
            Router: `/show_ecology_index`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "ToChangeIntoUSDD",
            Router: `/to_change_into_USDD`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "UpgradeWarehouse",
            Router: `/upgrade_warehouse`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"] = append(beego.GlobalControllerRouter["ecology1/controllers:EcologyIndexController"],
        beego.ControllerComments{
            Method: "ReturnPageListHostry",
            Router: `/upgrade_warehouse`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ecology1/controllers:FirstController"] = append(beego.GlobalControllerRouter["ecology1/controllers:FirstController"],
        beego.ControllerComments{
            Method: "CreateUserAbout",
            Router: `/test_read_all`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

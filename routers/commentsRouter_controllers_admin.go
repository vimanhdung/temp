package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"],
        beego.ControllerComments{
            Method: "CreateAdminAccount",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"],
        beego.ControllerComments{
            Method: "GetListAcount",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"],
        beego.ControllerComments{
            Method: "Detail",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"],
        beego.ControllerComments{
            Method: "UpdateAccount",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/admin:AdminAccountsController"],
        beego.ControllerComments{
            Method: "UpdateMyAccount",
            Router: `/update/myaccount`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

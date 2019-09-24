package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/user_app:UserAppAccountsController"] = append(beego.GlobalControllerRouter["indetail/controllers/user_app:UserAppAccountsController"],
        beego.ControllerComments{
            Method: "CreateUserAppAccount",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

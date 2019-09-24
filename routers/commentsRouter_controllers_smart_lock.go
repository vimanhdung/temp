package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"] = append(beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"] = append(beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"] = append(beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"] = append(beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"] = append(beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"] = append(beego.GlobalControllerRouter["indetail/controllers/smart_lock:SmartLocksController"],
        beego.ControllerComments{
            Method: "ChangeState",
            Router: `/changeState/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

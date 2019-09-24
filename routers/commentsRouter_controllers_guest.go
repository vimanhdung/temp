package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"] = append(beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"],
        beego.ControllerComments{
            Method: "CreateGuest",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"] = append(beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"],
        beego.ControllerComments{
            Method: "GetListGuest",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"] = append(beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"],
        beego.ControllerComments{
            Method: "UpdateGuest",
            Router: `/:guestId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"] = append(beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"],
        beego.ControllerComments{
            Method: "DetailGuest",
            Router: `/:guestId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"] = append(beego.GlobalControllerRouter["indetail/controllers/guest:GuestsController"],
        beego.ControllerComments{
            Method: "DeletedGuest",
            Router: `/:guestId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

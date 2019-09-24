package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"] = append(beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"],
        beego.ControllerComments{
            Method: "CreateRoom",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"] = append(beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"],
        beego.ControllerComments{
            Method: "GetListRoom",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"] = append(beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"],
        beego.ControllerComments{
            Method: "Detail",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"] = append(beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"],
        beego.ControllerComments{
            Method: "UpdateRoom",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"] = append(beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"] = append(beego.GlobalControllerRouter["indetail/controllers/room:RoomsController"],
        beego.ControllerComments{
            Method: "GetListRoomAvailable",
            Router: `/getListRoomAvailable/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

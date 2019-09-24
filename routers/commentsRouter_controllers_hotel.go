package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"] = append(beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"],
        beego.ControllerComments{
            Method: "CreateHotel",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"] = append(beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"],
        beego.ControllerComments{
            Method: "UpdateHotel",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"] = append(beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"] = append(beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"],
        beego.ControllerComments{
            Method: "GetAllHotels",
            Router: `/GetAllHotels/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"] = append(beego.GlobalControllerRouter["indetail/controllers/hotel:HotelsController"],
        beego.ControllerComments{
            Method: "DetailHotel",
            Router: `/detailHotel/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

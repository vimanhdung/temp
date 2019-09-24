package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers:BaseController"] = append(beego.GlobalControllerRouter["indetail/controllers:BaseController"],
        beego.ControllerComments{
            Method: "GetConfig",
            Router: `/configs/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers:FileController"] = append(beego.GlobalControllerRouter["indetail/controllers:FileController"],
        beego.ControllerComments{
            Method: "UploadImage",
            Router: `/file/upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

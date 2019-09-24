package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "CreateBooking",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "GetListBooking",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "Detail",
            Router: `/:bookingId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put","post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "BookingVerifyPassport",
            Router: `/checkin/verifyPassport/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "ConfirmCheckin",
            Router: `/confirmCheckin/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "ConfirmCheckout",
            Router: `/confirmCheckout/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "Reject",
            Router: `/reject/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"] = append(beego.GlobalControllerRouter["indetail/controllers/booking:BookingsController"],
        beego.ControllerComments{
            Method: "BookingCharge",
            Router: `/reject/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

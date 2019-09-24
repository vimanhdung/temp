package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"],
        beego.ControllerComments{
            Method: "ChangePassword",
            Router: `/changePassword`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"],
        beego.ControllerComments{
            Method: "KioskLogout",
            Router: `/kiosk/logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"],
        beego.ControllerComments{
            Method: "Me",
            Router: `/me`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthController"],
        beego.ControllerComments{
            Method: "RefreshToken",
            Router: `/refreshToken`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthGuestController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthGuestController"],
        beego.ControllerComments{
            Method: "GuestChangePassword",
            Router: `/guest/changePassword`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:AuthUserAppAcountController"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:AuthUserAppAcountController"],
        beego.ControllerComments{
            Method: "UserAppAccountChangePassword",
            Router: `/userAppAccount/changePassword`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"],
        beego.ControllerComments{
            Method: "ConfirmCode",
            Router: `/forgotPassword/confirmCode`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"],
        beego.ControllerComments{
            Method: "InputNewPassword",
            Router: `/forgotPassword/inputNewPassword`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"],
        beego.ControllerComments{
            Method: "SendCode",
            Router: `/forgotPassword/sendCode`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"],
        beego.ControllerComments{
            Method: "GuestConfirmCode",
            Router: `/guest/forgotPassword/confirmCode`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"],
        beego.ControllerComments{
            Method: "InputNewPasswordForGuest",
            Router: `/guest/forgotPassword/inputNewPassword`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:ForgotPassword"],
        beego.ControllerComments{
            Method: "ForgotPassword",
            Router: `/guest/forgotPassword/sendCode`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:Login"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:Login"],
        beego.ControllerComments{
            Method: "LoginCMS",
            Router: `/cms/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:Login"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:Login"],
        beego.ControllerComments{
            Method: "LoginGuest",
            Router: `/guest/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:Login"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:Login"],
        beego.ControllerComments{
            Method: "LoginHotel",
            Router: `/hotel/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["indetail/controllers/auth:Login"] = append(beego.GlobalControllerRouter["indetail/controllers/auth:Login"],
        beego.ControllerComments{
            Method: "LoginUserAppAccount",
            Router: `/userAppAccount/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

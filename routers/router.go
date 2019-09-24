// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"indetail/controllers"
	"indetail/controllers/admin"
	"indetail/controllers/auth"
	"indetail/controllers/booking"
	"indetail/controllers/guest"
	"indetail/controllers/hotel"
	"indetail/controllers/room"
	"indetail/controllers/smart_lock"
	"indetail/controllers/user_app"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&auth.AuthController{},
				&auth.AuthGuestController{},
				&auth.AuthUserAppAcountController{},
				&auth.ForgotPassword{},
				&auth.Login{},
			),
		),

		beego.NSNamespace("/account",
			beego.NSInclude(
				&admin.AdminAccountsController{},
			),
		),

		beego.NSNamespace("/userApp",
			beego.NSInclude(
				&user_app.UserAppAccountsController{},
			),
		),

		beego.NSNamespace("/bookings",
			beego.NSInclude(
				&booking.BookingsController{},
			),
		),

		beego.NSNamespace("/hotels",
			beego.NSInclude(
				&hotel.HotelsController{},
			),
		),

		beego.NSNamespace("/rooms",
			beego.NSInclude(
				&room.RoomsController{},
			),
		),

		beego.NSNamespace("/guest",
			beego.NSInclude(
				&guest.GuestsController{},
			),
		),

		beego.NSNamespace("/smart_lock",
			beego.NSInclude(
				&smart_lock.SmartLocksController{},
			),
		),

		beego.NSNamespace("/common",
			beego.NSInclude(
				&controllers.BaseController{},
				&controllers.FileController{},
			),
		),

		beego.NSRouter("/configs", &controllers.BaseController{}, "post:GetConfig"),
	)
	beego.AddNamespace(ns)

	beego.SetStaticPath("/storage", "storage")
}

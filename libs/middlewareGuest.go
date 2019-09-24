package libs

import (
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/models"
)

type MiddlewareGuest struct {
	beego.Controller
}

var GlobalGuest *models.Guests

//var PermissionsCode []string

// Check Token valid & check user exists
func (c *MiddlewareGuest) Prepare() {
	token := ParseToken(c.Ctx.Request.Header.Get("Authorization"))
	et := EasyToken{}
	valid, email, jwt_id, _ := et.ValidateToken(token)
	guest, _ := models.GetGuestByEmail(email)
	b := models.CheckJWTExists(jwt_id)
	if !valid || guest == nil || !b {
		c.Data["json"] = ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Invalid token",
			map[string]interface{}{"Token": conf.TOKEN_INVALID},
		)
		c.ServeJSON()
		return
	}
	GlobalGuest = guest
	return
}

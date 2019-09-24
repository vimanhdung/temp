package libs

import (
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/models"
)

type MiddlewareUserAppAccount struct {
	beego.Controller
}

var GlobalUserAppAccount *models.UserAppAccounts

//var PermissionsCode []string

// Check Token valid & check user exists
func (c *MiddlewareUserAppAccount) Prepare() {
	token := ParseToken(c.Ctx.Request.Header.Get("Authorization"))
	et := EasyToken{}
	valid, userName, jwt_id, _ := et.ValidateToken(token)
	userAppAccount, _ := models.GetUserAppAccountByLoginName(userName)
	b := models.CheckJWTExists(jwt_id)
	if !valid || userAppAccount == nil || !b {
		c.Data["json"] = ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Invalid token",
			map[string]interface{}{"Token": conf.TOKEN_INVALID},
		)
		c.ServeJSON()
		return
	}
	GlobalUserAppAccount = userAppAccount
	return
}


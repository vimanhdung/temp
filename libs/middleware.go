package libs

import (
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/models"
)

type Middleware struct {
	beego.Controller
}

var GlobalUser models.GlobalUsers

var PermissionsCode []string

// Check Token valid & check user exists
func (c *Middleware) Prepare() {
	token := ParseToken(c.Ctx.Request.Header.Get("Authorization"))
	et := EasyToken{}
	valid, email, jwtId, _ := et.ValidateToken(token)
	var user models.GlobalUsers

	b, jwtToken, _ := models.GetJWT(jwtId)
	if !valid || !b {
		c.Data["json"] = ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Invalid token",
			map[string]interface{}{"Token": conf.TOKEN_INVALID},
		)
		c.ServeJSON()
		return
	}

	switch jwtToken.Type {
	case conf.TYPE_USER:
		user, _ = models.PermissionGetUserAppByLoginName(email)
		user.Type = conf.TYPE_USER
	case conf.TYPE_GUEST:
		user, _ = models.PermissionGetGuestByLoginName(email)
		user.Type = conf.TYPE_GUEST
	default:
		user, _ = models.PermissionGetAdminByEmail(email)
		user.Type = conf.TYPE_ADMIN
	}

	GlobalUser = user
	return
}

// Permission Denied
// @Param permissionCode string
func (c *Middleware) PermissionDenied(code string) bool {
	if !c.HasPermission(code) {
		c.Data["json"] = ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Permission denied",
			map[string]interface{}{"Auth": conf.PERMISSION_DENY},
		)
		c.ServeJSON()
		return false
	}
	return true
}

// Check Has Permission
// @Param code string "permission code"
// @return bool
func (c *Middleware) HasPermission(code string) bool {
	if len(PermissionsCode) == 0 {
		return false
	}
	for _, value := range PermissionsCode {
		if code == value {
			return true
		}
	}
	return false
}

package auth

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
)

// AuthController operations for Auth
type AuthController struct {
	libs.Middleware
}

// Post ...
// @Title Refresh Token
// @Description Refresh Token
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} libs.ResponseJson Refresh Token Success
// @Failure 403 104 : Invalid Token
// @router /refreshToken [post]
func (c *AuthController) RefreshToken() {
	tokenString := libs.ParseToken(c.Ctx.Request.Header.Get("Authorization"))
	userGlobal := libs.GlobalUser
	var user *models.AdminAccounts
	user.Id = userGlobal.Id
	et := libs.EasyToken{}
	token, _ := et.RefreshToken(tokenString, user)
	c.Data["json"] = libs.ResultJson(
		map[string]interface{}{"token": token},
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Refresh Token Success!",
		nil,
	)
	c.ServeJSON()
}

// Post
// @Title Change Password
// @Description Change Password
// @Param Authorization header string true "Bearer token"
// @Param body body models.ChangePassFields true "Body for change password"
// @Success 200 {object} libs.ResponseJson Change Password Success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 221 : Parse json false <br> 111 : Incorrect password <br> 215 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number and new password must be different from your previous password) <br> 217 : Fields required
// @router /changePassword [post]
func (c *AuthController) ChangePassword() {
	// Check Permission
	if !c.PermissionDenied(conf.AUTH_CHANGE_PASSWORD) {
		return
	}
	var ob models.ChangePassFields
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err == nil {
		b, errCode := validation.CheckValidate(ob)
		if !b {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Validate False", errCode)
			c.ServeJSON()
			return
		}
		if ob.NewPass == ob.OldPass {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Incorrect password format",
				map[string]interface{}{"Password": conf.PASSWORD_FORMAT_INVALID},
			)
			c.ServeJSON()
			return
		}
		userGlobal := libs.GlobalUser
		var user *models.AdminAccounts
		user.Id = userGlobal.Id
		user.Password = userGlobal.Password
		if !libs.CheckHash(ob.OldPass, user.Password) {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Incorrect password ",
				map[string]interface{}{"Password": conf.INCORRECT},
			)
			c.ServeJSON()
			return
		}
		user.Password = libs.GetHashPassword(ob.NewPass)
		// Update Password for User Account
		// Destroy token of this user
		if models.UpdateAdminAccountsAndDestroyToken(user) {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.SUCCESS_STATUS),
				"Change Password Success",
				nil,
			)
		}
	} else {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
	}
	c.ServeJSON()
}

// Post
// @Title Logout
// @Description Logout
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} libs.ResponseJson Logout Success
// @Failure 403 104 : Invalid Token <br> 105 : Logout false
// @router /logout [post]
func (c *AuthController) Logout() {
	tokenString := libs.ParseToken(c.Ctx.Request.Header.Get("Authorization"))
	et := libs.EasyToken{}
	_, _, jti, _ := et.ValidateToken(tokenString)
	if err := models.DeleteJwtTokens(orm.NewOrm(), jti); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Logout false",
			map[string]interface{}{"Logout": conf.LOGOUT_FALSE},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Logout Success",
		nil,
	)
	c.ServeJSON()
}

// Post
// @Title Logout
// @Description Logout
// @Param Authorization header string true "Bearer token"
// @Param body body models.KioskLogout true "Password for logout"
// @Success 200 {object} libs.ResponseJson Logout Success
// @Failure 403 104 : Invalid Token <br> 217 : Fields required <br> 111 : Incorrect password <br> 105 : Logout false <br> 221 : Parse json false
// @router /kiosk/logout [post]
func (c *AuthController) KioskLogout() {
	var ob models.KioskLogout
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err == nil {
		b, errCode := validation.CheckValidate(ob)
		if !b {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Validate False",
				errCode,
			)
			c.ServeJSON()
			return
		}
		user := libs.GlobalUser;
		if !libs.CheckHash(ob.Password, user.Password) {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Password incorrect",
				map[string]interface{}{"Password": conf.INCORRECT},
			)
			c.ServeJSON()
			return
		}
		tokenString := libs.ParseToken(c.Ctx.Request.Header.Get("Authorization"))
		et := libs.EasyToken{}
		_, _, jti, _ := et.ValidateToken(tokenString)
		if err := models.DeleteJwtTokens(orm.NewOrm(), jti); err != nil {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Logout false",
				map[string]interface{}{"Logout": conf.LOGOUT_FALSE},
			)
			c.ServeJSON()
			return
		}
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.SUCCESS_STATUS),
			"Logout Success",
			nil,
		)
	} else {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
	}
	c.ServeJSON()
}

// Post
// @Title Get my account infomation
// @Description Get my account infomation
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} libs.ResponseJson Success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token
// @router /me [get]
func (c *AuthController) Me() {
	// Check Permission
	globalUser := libs.GlobalUser

	var user interface{}
	switch globalUser.Type {
	case conf.TYPE_USER:
		user, _ = models.GetUserAppAccountsById(globalUser.Id)
	case conf.TYPE_GUEST:
		user, _ = models.GetGuestById(globalUser.Id)
	default:
		user, _ = models.GetAdminAccountsById(globalUser.Id)
	}

	c.Data["json"] = libs.ResultJson(
		user,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Get Me Success",
		nil,
	)
	c.ServeJSON()
}

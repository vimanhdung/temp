package auth

import (
	"encoding/json"
	"fmt"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
)

// AuthUserAppAcountController operations for AuthUserAppAcount
type AuthUserAppAcountController struct {
	libs.Middleware
}

// UserAppAccountChangePassword
// @Title Change Password
// @Description Change Password
// @Param Authorization header string true "Bearer token"
// @Param body body models.ChangePassFields true "Body for change password"
// @Success 200 {object} models.SwaggerDefault Change Password Success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 221 : Parse json false <br> 111 : Incorrect password <br> 215 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number and new password must be different from your previous password) <br> 217 : Fields required
// @router /userAppAccount/changePassword [post]
func (c *AuthUserAppAcountController) UserAppAccountChangePassword() {
	var ob models.ChangePassFields
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
		c.ServeJSON()
		return
	}

	isValidate, errCode := validation.CheckValidate(ob)
	if !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Validate False", errCode)
		c.ServeJSON()
		return
	}

	userGlobal := libs.GlobalUser
	var userAppAccount *models.UserAppAccounts
	userAppAccount.Id = userGlobal.Id
	userAppAccount.Password = userGlobal.Password

	if !libs.CheckHash(ob.OldPass, userAppAccount.Password) {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Incorrect password ",
			map[string]interface{}{"Password": conf.INCORRECT},
		)
		c.ServeJSON()
		return
	}

	userAppAccount.Password = libs.GetHashPassword(ob.NewPass)
	// Update Password for User Account
	// Destroy token of this user
	if !models.ChangePasswordUserAppAndLogout(userAppAccount) {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.SUCCESS_STATUS),
			"Change Password False",
			nil,
		)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Change Password Success",
		nil,
	)
	c.ServeJSON()
	return
}

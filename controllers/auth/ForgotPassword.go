package auth

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/mail"
	"indetail/libs/validation"
	"indetail/models"
)

type ForgotPassword struct {
	beego.Controller
}

// ForgotPassword
// @Title Reset Password
// @Description Send link reset password to email
// @Param body body models.ForgotSendCode true "Input Email for forgot password"
// @Success 200 {object} models.SwaggerDefault Send mail success
// @Failure 403 217 : Fields required <br> 213 : Invalid email format <br> 302 : Email not found <br> 221 : Parse json false <br> 303 : gen code false
// @router /guest/forgotPassword/sendCode [post]
func (c *ForgotPassword) ForgotPassword() {
	var ob models.ForgotSendCode
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
			"Validate false",
			errCode,
		)
		c.ServeJSON()
		return
	}

	guest, err := models.GetGuestsByEmail(ob.Email)
	if err != nil || guest == nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.SUCCESS_STATUS),
			"Email not exists",
			nil,
		)
		c.ServeJSON()
		return
	}

	resetPassword := models.GetStructGuestInsert(guest)

	var code string
	if err, code = models.CreateOrUpdatePasswordReset(&resetPassword); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	//Send mail
	mail.SendMail(ob.Email, guest.FullName, code)
	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Send mail success",
		nil,
	)
	c.ServeJSON()
	return
}

// GuestConfirmCode
// @Title Confirm Code
// @Description Input code verify
// @Param body body models.ForgotConfirmCode true "Body for Input new password"
// @Success 200 {object} models.SwaggerDefault Code exactly
// @Failure 403 217 : Fields required <br> 111 : Incorrect code <br> 221 : Parse json false
// @router /guest/forgotPassword/confirmCode [post]
func (c *ForgotPassword) GuestConfirmCode() {
	GetTokenResetPassword(c)
}

// InputNewPasswordForGuest
// @Title Input New Password
// @Description Input code verify & new password
// @Param body body models.ForgotNewPassword true "Body for Input new password"
// @Success 200 {object} models.SwaggerDefault A new password has been created
// @Failure 403 217 : Fields required <br> 111 : Token reset incorrect <br> 225 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number) <br> 303 : Update false <br> 221 : Parse json false
// @router /guest/forgotPassword/inputNewPassword [post]
func (c *ForgotPassword) InputNewPasswordForGuest() {
	var ob models.ForgotNewPassword
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
			"Valide False",
			errCode,
		)
		c.ServeJSON()
		return
	}

	var rp models.PasswordResets
	if isValidate, rp = models.CheckTokenResetPassword(ob); !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Token reset password invalid or expired",
			map[string]interface{}{"TokenReset": conf.INCORRECT},
		)
		c.ServeJSON()
		return
	}

	// Change Token
	guest, _ := models.GetGuestById(rp.AccountId)
	guest.Password = libs.GetHashPassword(ob.NewPassword)

	if !models.UpdatePasswordGuestAndDeleteToken(guest) {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Update Password false",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"A new password has been created",
		nil,
	)
	c.ServeJSON()
	return
}

// SendCode
// @Title Send Code to Confirm Reset
// @Description Send code verify to email
// @Param body body models.ForgotSendCode true "Input Email for forgot password"
// @Success 200 {object} models.SwaggerDefault Send mail success
// @Failure 403 217 : Fields required <br> 213 : Invalid email format <br> 302 : Email not found <br> 221 : Parse json false <br> 303 : gen code false
// @router /forgotPassword/sendCode [post]
func (c *ForgotPassword) SendCode() {
	var ob models.ForgotSendCode
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
			"Validate false",
			errCode,
		)
		c.ServeJSON()
		return
	}
	user, err := models.GetUserByEmail(ob.Email)
	if err != nil || user == nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Email not exists",
			map[string]interface{}{"Email": conf.RECORD_NOT_FOUND})
		c.ServeJSON()
		return
	}
	// Insert password_resets table
	rp := models.GetStructInsert(user)
	var code string
	if err, code = models.CreateOrUpdatePasswordReset(&rp); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	//Send mail
	mail.SendMail(ob.Email, user.FullName, code)
	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Send mail success",
		nil,
	)
	c.ServeJSON()
	return
}

// ConfirmCode
// @Title Confirm Code
// @Description Input code verify
// @Param body body models.ForgotConfirmCode true "Body for Input new password"
// @Success 200 {object} models.SwaggerDefault Code exactly
// @Failure 403 217 : Fields required <br> 111 : Incorrect code <br> 221 : Parse json false
// @router /forgotPassword/confirmCode [post]
func (c *ForgotPassword) ConfirmCode() {
	GetTokenResetPassword(c)
}

func GetTokenResetPassword(c *ForgotPassword)  {
	var ob models.ForgotConfirmCode
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
			"Validate False",
			errCode,
		)
		c.ServeJSON()
		return
	}

	var rp models.PasswordResets
	if isValidate, rp = models.CheckCodeExists(ob); !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Code invalid or expired",
			map[string]interface{}{"Code": conf.INCORRECT},
		)
		c.ServeJSON()
		return
	}

	// Update Token reset password
	rp.Token = libs.GenerateMd5String()
	if err := models.UpdatePasswordResetsById(&rp); err == nil {
		c.Data["json"] = libs.ResultJson(
			map[string]interface{}{"token": rp.Token},
			fmt.Sprint(conf.SUCCESS_STATUS),
			"Code exactly",
			nil,
		)
		c.ServeJSON()
		return
	}
}

// InputNewPassword
// @Title Input New Password
// @Description Input code verify & new password
// @Param body body models.ForgotNewPassword true "Body for Input new password"
// @Success 200 {object} models.SwaggerDefault A new password has been created
// @Failure 403 217 : Fields required <br> 111 : Token reset incorrect <br> 215 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number) <br> 303 : Update false <br> 221 : Parse json false
// @router /forgotPassword/inputNewPassword [post]
func (c *ForgotPassword) InputNewPassword() {
	var ob models.ForgotNewPassword
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
	}

	isValidate, errCode := validation.CheckValidate(ob)
	if !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Valide False",
			errCode,
		)
		c.ServeJSON()
		return
	}

	var rp models.PasswordResets
	if isValidate, rp = models.CheckTokenResetPassword(ob); !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Token reset password invalid or expired",
			map[string]interface{}{"TokenReset": conf.INCORRECT},
		)
		c.ServeJSON()
		return
	}

	// Update password
	// Delete Password Reset
	user, _ := models.GetAdminAccountsById(rp.AccountId)
	user.Password = libs.GetHashPassword(ob.NewPassword)
	if !models.UpdateAccountAndDeleteResetRow(user) {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Update Password false",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"A new password has been created",
		nil,
	)
	c.ServeJSON()
	return
}

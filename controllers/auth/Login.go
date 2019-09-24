package auth

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
)

type Login struct {
	beego.Controller
}

// Post ...
// @Title CMS Login
// @Description Auth Login For CMS
// @Param body body models.LoginFields true "Body for Login CMS"
// @Success 200 {object} libs.ResponseJson Login Success
// @Failure 403 217 : Fields required <br> 213 : Email format invalid <br> 101 : Permission denied <br> 111 : Password incorrect <br> 221 : Parse json false
// @router /cms/login [post]
func (c *Login) LoginCMS() {
	isValidate, user := c.ParseJsonAndValidate(conf.AUTH_LOGIN_CMS)
	if !isValidate {
		return
	}

	userLogin := models.GetJwtTokenLogin(user.Id, conf.TYPE_ADMIN, user.Email, )

	c.Data["json"] = libs.GetToken(&userLogin, conf.TokenExpires)
	c.ServeJSON()
}

// Post ...
// @Title Hotel Login
// @Description Auth Login For Hotel
// @Param body body models.LoginFields true "Body for Login Hotel"
// @Success 200 {object} libs.ResponseJson Login Success
// @Failure 403 217 : Fields required <br> 213 : Email format invalid <br> 101 : Permission denied <br> 111 : Password incorrect <br> 221 : Parse json false
// @router /hotel/login [post]
func (c *Login) LoginHotel() {
	isValidate, user := c.ParseJsonAndValidate(conf.AUTH_LOGIN_HOTEL)
	if !isValidate {
		return
	}

	userLogin := models.GetJwtTokenLogin(user.Id, conf.TYPE_ADMIN, user.Email, )

	c.Data["json"] = libs.GetToken(&userLogin, conf.TokenExpires)
	c.ServeJSON()
}

// LoginGuest
// @Title Customer Login
// @Description Auth Login For Customer
// @Param body body models.LoginFields true "Body for Login Customer"
// @Success 200 {object} libs.ResponseJson Login Success
// @Failure 403 217 : Fields required <br> 213 : Email format invalid <br> 101 : Permission denied <br> 111 : Password incorrect <br> 221 : Parse json false
// @router /guest/login [post]
func (c *Login) LoginGuest() {
	isValidate, guest := CheckGuestLogin(c)
	if !isValidate {
		return
	}
	userLogin := models.GetJwtTokenLogin(guest.Id, conf.TYPE_GUEST, guest.Email, )

	c.Data["json"] = libs.GetToken(&userLogin, conf.TokenExpires)
	c.ServeJSON()
}

// Post ...
// @Title Customer Login
// @Description Auth Login For Customer
// @Param body body models.LoginUserAppFields true "Body for Login Customer"
// @Success 200 {object} libs.ResponseJson Login Success
// @Failure 403 217 : Fields required <br> 213 : Email format invalid <br> 101 : Permission denied <br> 111 : Password incorrect <br> 221 : Parse json false
// @router /userAppAccount/login [post]
func (c *Login) LoginUserAppAccount() {
	isValidate, userAppAccount := CheckUserAppLogin(c)
	if !isValidate {
		return
	}
	userLogin := models.GetJwtTokenLogin(userAppAccount.Id, conf.TYPE_USER, userAppAccount.LoginName)

	c.Data["json"] = libs.GetToken(&userLogin, conf.TokenExpires)
	c.ServeJSON()
}

// Check user app login
func CheckUserAppLogin(c *Login) (b bool, userApp *models.UserAppAccounts) {
	var ob models.LoginUserAppFields
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
		c.ServeJSON()
		return false, userApp
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
		return false, userApp
	}

	isValidate, userApp = libs.CheckAuthUserAppAccount(ob.LoginName, ob.Password)

	if !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Login Name or Password is incorrect ! Please try again !",
			map[string]interface{}{"Password": conf.INCORRECT},
		)
		c.ServeJSON()
		return false, userApp
	}

	return true, userApp
}

// Parse Json And Validate
func CheckGuestLogin(c *Login) (b bool, guest *models.Guests) {
	var ob models.LoginFields
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
		c.ServeJSON()
		return false, guest
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
		return false, guest
	}

	isValidate, guest = libs.CheckAuthGuest(ob.Email, ob.Password)

	if !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Email or Password is incorrect ! Please try again !",
			map[string]interface{}{"Password": conf.INCORRECT},
		)
		c.ServeJSON()
		return false, guest
	}

	return true, guest
}

// Parse Json And Validate
func (c *Login) ParseJsonAndValidate(permission string) (b bool, user *models.AdminAccounts) {
	var ob models.LoginFields
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
		c.ServeJSON()
		return false, user
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
		return false, user
	}

	isValidate, user = libs.CheckAuth(ob.Email, ob.Password)
	if !isValidate {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Email or Password is incorrect ! Please try again !",
			map[string]interface{}{"Password": conf.INCORRECT},
		)
		c.ServeJSON()
		return false, user
	}
	return true, user
}

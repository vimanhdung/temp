package user_app

import (
	"encoding/json"
	"fmt"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
)

// UserAppAccountsController operations for AdminAccounts
type UserAppAccountsController struct {
	libs.Middleware
}

// URLMapping ...
func (c *UserAppAccountsController) URLMapping() {
	c.Mapping("CreateUserAppAccount", c.CreateUserAppAccount)
}

// CreateAdminAccount ...
// @Title Create Admin Account
// @Description create AdminAccounts
// @Param Authorization header string true "Bearer token"
// @Param	body		body 	models.InsertUserAppFields	true	"body for UserAppAccount content"
// @Success 200 {object} models.SwaggerDetailAccount Add new account success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 224 : Fields value not exists <br> 221 : Parse json body false <br> 215 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number) <br> 217 : Field required <br> 213 : Invalid Email format <br> 210 : Field not numeric <br> 301 : Account Exists <br> 303 : Save Database false
// @router / [post]
func (c *UserAppAccountsController) CreateUserAppAccount() {
	// Check Permission
	var userAppModel models.InsertUserAppFields
	user := libs.GlobalUser
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &userAppModel); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE},
		)
		c.ServeJSON()
		return
	}

	// Check Validate
	if b, errCode := validation.CheckValidate(userAppModel); !b {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Validate false",
			errCode,
		)
		c.ServeJSON()
		return
	}

	// Check Login Name Exists
	if models.CheckLoginNameExists(conf.NEW_ACCOUNT, user.Id, userAppModel.LoginName) {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Account Exists",
			map[string]interface{}{"Account": conf.RECORD_EXISTS},
		)
		c.ServeJSON()
		return
	}
	// Parse Data
	var entity models.UserAppAccounts
	json.Unmarshal(c.Ctx.Input.RequestBody, &entity)
	entity.CreatedUser = user.Id
	entity.Password = libs.GetHashPassword(entity.Password)
	// Insert Account
	if _, err := models.AddUserAppAccounts(&entity); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	entity.Password = ""
	c.Data["json"] = libs.ResultJson(
		entity,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Add new account success",
		nil,
	)
	c.ServeJSON()
}



package admin

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
	"strconv"
)

// AdminAccountsController operations for AdminAccounts
type AdminAccountsController struct {
	libs.Middleware
}

// URLMapping ...
func (c *AdminAccountsController) URLMapping() {
	c.Mapping("CreateAdminAccount", c.CreateAdminAccount)
	c.Mapping("Detail", c.Detail)
	c.Mapping("GetListAcount", c.GetListAcount)
	c.Mapping("UpdateAccount", c.UpdateAccount)
	c.Mapping("Delete", c.Delete)
}

// CreateAdminAccount ...
// @Title Create Admin Account
// @Description create AdminAccounts
// @Param Authorization header string true "Bearer token"
// @Param	body		body 	models.InsertAccountFields	true		"body for AdminAccounts content"
// @Success 200 {object} models.SwaggerDetailAccount Add new account success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 224 : Fields value not exists <br> 221 : Parse json body false <br> 215 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number) <br> 217 : Field required <br> 213 : Invalid Email format <br> 210 : Field not numeric <br> 301 : Account Exists <br> 303 : Save Database false
// @router / [post]
func (c *AdminAccountsController) CreateAdminAccount() {
	// Check Permission
	if !c.PermissionDenied(conf.ACCOUNT_CREATE) {
		return
	}
	var v models.InsertAccountFields
	user := libs.GlobalUser
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
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
	if b, errCode := validation.CheckValidate(v); !b {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Validate false",
			errCode,
		)
		c.ServeJSON()
		return
	}
	// Check Email Exists
	if models.CheckEmailExists(conf.NEW_ACCOUNT, user.Id, v.Email) {
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
	var entity models.AdminAccounts
	json.Unmarshal(c.Ctx.Input.RequestBody, &entity)
	entity.CreatedUser = user.Id
	entity.Password = libs.GetHashPassword(entity.Password)
	// Insert Account
	if _, err := models.AddAdminAccounts(&entity); err != nil {
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

// Detail ...
// @Title Detail
// @Description get AdminAccounts by id
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SwaggerDetailAccount Get Success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Account not found
// @router /:id [get]
func (c *AdminAccountsController) Detail() {
	// Check Permission
	if !c.PermissionDenied(conf.ACCOUNT_DETAIL) {
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetAdminAccountsById(id)
	if err != nil || v == nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Id": conf.RECORD_NOT_FOUND},
		)
	} else {
		v.Password = ""
		c.Data["json"] = libs.ResultJson(
			v,
			fmt.Sprint(conf.SUCCESS_STATUS),
			"Get Success",
			nil,
		)
	}
	c.ServeJSON()
}

// GetListAcount ...
// @Title Get List Acount
// @Description get AdminAccounts
// @Param	Authorization	header	string	true	"token"
// @Param	hotelId			query	int		false	"hotel id"
// @Param	fullName		query	string	false	"account full name"
// @Param	email			query	string	false	"account email"
// @Param	status			query	string	false	"account status"
// @Param	fields	query	string	false	"Default "*".Fields returned. e.g. col1,col2 ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	page	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SwaggerListAccount
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br>
// @router / [get]
func (c *AdminAccountsController) GetListAcount() {
	var fields = "admin_accounts.admin_account_id, admin_accounts.role_id, admin_accounts.email, admin_accounts.status, admin_accounts.full_name," +
		" admin_accounts.created_user, admin_accounts.updated_user, admin_accounts.created_at, admin_accounts.updated_at"
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64
	var page int64 = 1

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = v
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("page"); err == nil {
		page = v
		if page < 1 {
			page = 1
		}
	}
	// hotelId
	if v, err := c.GetInt64("hotelId"); err == nil {
		query["hotel_administrators.hotel_id"] = fmt.Sprint(v)
	}
	// admin name
	if v := c.GetString("fullName"); v != "" {
		query["full_name__searchLike"] = fmt.Sprint(v)
	}
	// admin email
	if v := c.GetString("email"); v != "" {
		query["admin_accounts.email__searchLike"] = fmt.Sprint(v)
	}
	// admin status
	if v, err := c.GetInt64("status"); err == nil {
		query["admin_accounts.status"] = fmt.Sprint(v)
	}

	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var errorMsg = ""
	total, error := models.CountAdminAccounts(query)
	if total < 1 || error != nil {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultPagingJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, limit, page, total)
		c.ServeJSON()
		return
	}

	if limit <= 0 {
		limit = total
	}
	offset = (page - 1) * limit
	l, err := models.GetAllAdminAccounts(query, fields, sortby, order, offset, limit)

	// set status code
	c.Data["json"] = libs.ResultPagingJson(l, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, limit, page, total)
	if err != nil {
		c.Data["json"] = libs.ResultPagingJson(nil, errorCode, errorMsg, limit, page, total)
	}
	c.ServeJSON()
}

// UpdateAccount ...
// @Title Update Account
// @Description update the AdminAccounts
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.UpdateAccountFields	true		"body for AdminAccounts content"
// @Success 200 {object} models.AddAdminAccountStruct Update account success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 224 : Fields value not exists <br> 221 : Parse json body false <br> 215 : Incorrect password format (Password must contain at least eight characters to 60 characters and uppercase and number) <br> 217 : Field required <br> 213 : Valid Email format <br> 210 : Field not numeric <br> 301 : Email Exists <br> 303 : Save Database false
// @router /:id [put]
func (c *AdminAccountsController) UpdateAccount() {
	// Check Permission
	if !c.PermissionDenied(conf.ACCOUNT_UPDATE) {
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	var v models.UpdateAccountFields
	user := libs.GlobalUser
	entity, err := models.GetAdminAccountsById(id)
	// Check isset account
	if entity == nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			err.Error(),
			map[string]interface{}{"Id": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
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
	if b, errCode := validation.CheckValidate(v); !b {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Validate false",
			errCode,
		)
		c.ServeJSON()
		return
	}
	// Check Email Exists
	if models.CheckEmailExists(conf.EDIT_ACCOUNT, id, v.Email) {
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
	password := entity.Password
	json.Unmarshal(c.Ctx.Input.RequestBody, &entity)
	entity.UpdatedUser = user.Id
	transaction := false
	if v.Password != "" {
		transaction = true
		entity.Password = libs.GetHashPassword(v.Password)
	} else {
		entity.Password = password
	}
	// Update Account with transaction
	if transaction {
		if !models.UpdateAdminAccountsAndDestroyToken(entity) {
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Update False",
				map[string]interface{}{"Database": conf.SAVE_FAILURES},
			)
			c.ServeJSON()
			return
		}
	}
	// Update account without transaction
	if err := models.UpdateAdminAccountsById(orm.NewOrm(), entity); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Update False",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	entity.Password = ""
	c.Data["json"] = libs.ResultJson(
		entity,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Update account success",
		nil,
	)
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the AdminAccounts
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} libs.ResponseJson delete success!
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Account Not Found <br> 303 : Delete False
// @router /:id [delete]
func (c *AdminAccountsController) Delete() {
	// Check Permission
	c.HasPermission(conf.ACCOUNT_DELETE)
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	exists, ok := models.SoftDeleteAccount(id)
	if !exists {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Account not exists",
			map[string]interface{}{"Account": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}
	if !ok {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Delete False",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Delete Success",
		nil,
	)
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description Update My Account
// @Param Authorization header string true "Bearer token"
// @Param body body models.UpdateMyAccount true "body for update my account infomation"
// @Success 200 {object} models.JsonResponse Update Success
// @Failure 403 104 : Invalid Token <br> 221 : Parse json false <br> 210 : Phone is not numeric <br> 303 : Update false
// @router /update/myaccount [put]
func (c *AdminAccountsController) UpdateMyAccount() {
	var ob models.UpdateMyAccount
	// set status code
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
	if b, errCode := validation.CheckValidate(ob); !b {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Fields invalid",
			errCode,
		)
		c.ServeJSON()
		return
	}
	user := libs.GlobalUser
	updateInfo, _ := models.GetAdminAccountsById(user.Id)
	// Update Account
	updateInfo.FullName = ob.FullName
	if err := models.UpdateAdminAccountsById(orm.NewOrm(), updateInfo); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Update account false",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(
		nil,
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Update success",
		nil,
	)
	c.ServeJSON()
}

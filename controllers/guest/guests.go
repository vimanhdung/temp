package guest

import (
	"encoding/json"
	"fmt"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
	"strconv"
	"strings"
	"time"
)

// CustomersController operations for Customers
type GuestsController struct {
	libs.Middleware
}

// URLMapping ...
func (c *GuestsController) URLMapping() {
	c.Mapping("CreateGuest", c.CreateGuest)
	c.Mapping("UpdateGuest", c.UpdateGuest)
	c.Mapping("GetListGuest", c.GetListGuest)
}

// CreateGuest ...
// @Title CreateGuest
// @Description create Guest
// @Param Authorization header string true "Bearer token"
// @Param	body		body 	models.AddGuests	true		"body for Guest content"
// @Success 201 {int} models.Guests.Id
// @Failure 101 Permission deny <br> 201 Missing param <br> 217 Variable require <br> 219 Variable is not positive integer <br> 220 Variable is not date <br> 211 Json fail <br> 226 Limit record
// @router / [post]
func (c *GuestsController) CreateGuest() {
	var customerStruct models.Guests

	c.Data["json"] = SaveGuestProcess(customerStruct, c.Ctx.Input.RequestBody, false)
	c.ServeJSON()
}

// UpdateGuest ...
// @Title UpdateGuest
// @Description update the Guest
// @Param Authorization header string true "Bearer token"
// @Param	guestId		path 	string				true	"The guest id you want to update"
// @Param	body		body 	models.UpdateGuests	true	"body for guest content"
// @Success 200 {object} models.SwaggerGuest
// @Failure 101 Permission deny <br> 201 Missing param <br> 217 Variable require <br> 219 Variable is not positive integer <br> 220 Variable is not date <br> 211 Json fail <br> 226 Limit record
// @router /:guestId [put]
func (c *GuestsController) UpdateGuest() {
	var guestStruct = models.Guests{}

	idStr := c.Ctx.Input.Param(":guestId")
	id, _ := strconv.Atoi(idStr)
	guestStruct.Id = id

	jsonBody := SaveGuestProcess(guestStruct, c.Ctx.Input.RequestBody, true)
	c.Data["json"] = jsonBody
	c.ServeJSON()
}

func SaveGuestProcess(guestStruct models.Guests, requestBody []byte, isUpdate bool) (jsonBody interface{}) {
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var detailErrorCode = make(map[string]interface{})
	var guestTmpInterface map[string]interface{}
	if errTmp := json.Unmarshal(requestBody, &guestTmpInterface); errTmp != nil {
		errorMsg = "Wrong format"
		detailErrorCode["GuestData"] = conf.VARIABLE_IS_NOT_JSON
		jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		return
	}

	var files interface{}
	var bookingId int = -1
	var bookingIdErr error
	var passportNumber string
	for index, value := range guestTmpInterface {
		if index == "Files" && value != "" {
			files = value
			guestTmpInterface[index] = ""
		} else if index == "Status" {
			if len(fmt.Sprint(value)) < 1 {
				errorMsg = "Missing param"
				detailErrorCode["Status"] = conf.VARIABLE_REQUIRED
				jsonBody = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
				return
			}
		} else if index == "BookingId" && len(value.(string)) > 0 {
			if bookingId, bookingIdErr = strconv.Atoi(fmt.Sprint(value)); bookingIdErr != nil {
				errorMsg = "Param error"
				detailErrorCode["BookingId"] = conf.VARIABLE_IS_NOT_NUMERIC
				jsonBody = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
				return
			}
		} else if index == "BirthDay" {
			var err error
			guestTmpInterface[index], err = time.Parse(conf.RFC3339, fmt.Sprint(value))
			if err != nil {
				errorMsg = "Param error"
				detailErrorCode["BirthDay"] = conf.VARIABLE_IS_NOT_DATE
				jsonBody = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
				return
			}
		} else if index == "PassportExpired" {
			var err error
			guestTmpInterface[index], err = time.Parse(conf.RFC3339, fmt.Sprint(value))
			if err != nil {
				errorMsg = "Param error"
				detailErrorCode["PassportExpired"] = conf.VARIABLE_IS_NOT_DATE
				jsonBody = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
				return
			}
		} else if index == "PassportNumber" {
			passportNumber = fmt.Sprint(value)
		}
	}
	if models.CheckExitPassport(passportNumber, fmt.Sprint(guestStruct.Id)) {
		errorMsg = "Param error"
		detailErrorCode["PassportNumber"] = conf.RECORD_EXISTS
		jsonBody = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
		return
	}

	oldTmpGuestData, errOldGuest := models.GetGuestById(guestStruct.Id)

	//remove BookingId
	delete(guestTmpInterface, "BookingId")
	if !isUpdate {
		delete(guestTmpInterface, "Id")
	} else {
		if errOldGuest != nil || oldTmpGuestData == nil {
			errorMsg = "Record not found"
			detailErrorCode["GuestId"] = conf.RECORD_NOT_FOUND
			jsonBody = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
			return
		}
	}
	//validate file
	var guestFilesInsert = "[]"
	var rebuildValidate = make(map[string]interface{})
	if files != nil {
		guestFilesInsert, rebuildValidate = models.ValidateGuestFiles(files, map[string]interface{}{})
		if guestFilesInsert == "" {
			errorMsg = "Wrong format"
			detailErrorCode["Files"] = conf.VARIABLE_IS_NOT_JSON
			jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
			return
		}
	}
	reMakeGuestInfo, _ := json.Marshal(guestTmpInterface)
	if err := json.Unmarshal(reMakeGuestInfo, &guestStruct); err != nil {
		errorMsg = "Wrong format"
		jsonBody = libs.ResultJson(nil, errorCode, errorMsg, rebuildValidate)
		return
	}
	if isUpdate {
		guestStruct.Id = oldTmpGuestData.Id
		guestStruct.CreatedAt = oldTmpGuestData.CreatedAt
		guestStruct.UpdatedAt = time.Now()
		guestStruct.UpdatedUser = libs.GlobalUser.Id
		guestStruct.Password = oldTmpGuestData.Password
		guestStruct.LoginName = oldTmpGuestData.LoginName
	} else {
		guestStruct.CreatedUser = libs.GlobalUser.Id
		guestStruct.Password = libs.GetHashPassword(conf.DEFAULT_PASSWORD)
	}
	guestStruct.Files = guestFilesInsert
	var isPass bool
	jsonBody, isPass = libs.ValidateStatus(conf.GUEST_STATUS_LIST, int(guestStruct.Status))
	if !isPass {
		return
	}
	isValid, listErrorCode := validation.CheckValidate(guestStruct)
	if !isValid {
		for key, value := range listErrorCode {
			rebuildValidate[key] = value
		}
		errorMsg = "Wrong format"
		jsonBody = libs.ResultJson(nil, errorCode, errorMsg, rebuildValidate)
		return
	}

	if isUpdate {
		if errCode := models.UpdateGuestById(&guestStruct); errCode != nil {
			errorCode = strconv.Itoa(conf.SAVE_FAILURES)
			errorMsg = errCode.Error()
		}
	} else {
		if len(guestStruct.Email) > 0 {
			guestStruct.LoginName = guestStruct.Email[:strings.Index(guestStruct.Email, "@")]
		}
		if err := models.AddGuestWithTransaction(&guestStruct, bookingId); err != "" {
			rebuildValidate["Guest"] = conf.LIMIT_RECORD
			errorCode = strconv.Itoa(conf.SAVE_FAILURES)
			errorMsg = err
		}
	}

	jsonBody = libs.ResultJson(nil, errorCode, errorMsg, rebuildValidate)
	//override out put
	if errorMsg == "" {
		guestStruct.Password = ""
		jsonBody = libs.ResultJson(guestStruct, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, rebuildValidate)
	}
	return
}

// GetListGuest ...
// @Title Get List Guest
// @Description get list guest
// @Param Authorization header string 	true "Bearer token"
// @Param	bookingId	query	int		true	"booking id"
// @Param	fullName	query	string	false	"guest full name. search like"
// @Param	status		query	int		false	"guest status"
// @Param	ageFrom		query	string	false	"guest age"
// @Param	ageTo		query	string	false	"guest age"
// @Param	fields		query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	limit		query	string	false	"Limit the size of result set. Must be an integer"
// @Param	page		query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SwaggerListGuest
// @Failure 403 101 Permission deny <br> 201 Missing param <br> 219 Variable is not positive integer <br> 302 Record not found
// @router / [get]
func (c *GuestsController) GetListGuest() {
	//process booking
	var fields = "guests.*"
	var sortBy []string
	var order []string
	var ageGte string
	var ageLte string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64
	var page int64 = 1
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var detailErrorCode = make(map[string]interface{})

	//permission
	globalUser := libs.GlobalUser
	if globalUser.Type == conf.TYPE_USER {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Permission denied",
			map[string]interface{}{"Auth": conf.PERMISSION_DENY},
		)
		c.ServeJSON()
		return
	}

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

	if v, err := c.GetInt64("bookingId"); err == nil {
		query["BookingId"] = fmt.Sprint(v)
	}

	if v := c.GetString("fullName"); v != "" {
		query["FullName"] = v
	}

	if v, err := c.GetInt64("status"); err == nil {
		query["Status"] = fmt.Sprint(v)
	}

	if v, err := c.GetInt64("ageFrom"); err == nil {
		ageGte = fmt.Sprint(v)
	}
	if v, err := c.GetInt64("ageTo"); err == nil {
		ageLte = fmt.Sprint(v)
	}

	//validate age
	if ageGte == "" && ageLte != "" {
		detailErrorCode["Age"] = conf.MISSNG_PARAM
		c.Data["json"] = libs.ResultJson("", errorCode, "Missing param AgeTo", detailErrorCode)
		c.ServeJSON()
		return
	}
	if ageGte == "" && ageLte != "" {
		detailErrorCode["Age"] = conf.MISSNG_PARAM
		c.Data["json"] = libs.ResultJson("", errorCode, "Missing param AgeTo", detailErrorCode)
		c.ServeJSON()
		return
	}
	//end validate age

	//count total
	totalRecord, _ := models.CountTotalRecord(query, fields)
	if totalRecord == 0 {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultJson("", fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, nil)
		c.ServeJSON()
		return
	}

	if limit <= 0 {
		limit = totalRecord
	}
	//set offset
	offset = (page - 1) * limit
	listGuest, msgErr := models.GetAllGuestCondition(query, fields, sortBy, order, offset, limit)

	if msgErr != "" || listGuest == nil {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultJson("", fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultPagingJson(listGuest, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, limit, page, totalRecord)
	c.ServeJSON()
}

// Detail Guest ...
// @Title DetailGuest
// @Description detail the Guest
// @Param Authorization header string true "Bearer token"
// @Param	guestId		path 	string	true	"Id of guest"
// @Success 200 {object} models.SwaggerGuest
// @Failure 403 :id is not int
// @router /:guestId [get]
func (c *GuestsController) DetailGuest() {
	guestIdStr := c.Ctx.Input.Param(":guestId")
	guestId, _ := strconv.Atoi(guestIdStr)

	guest, err := models.GetGuestById(guestId)

	if err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Guest is not exists",
			map[string]interface{}{"Guest": conf.RECORD_NOT_FOUND,
			})
		c.ServeJSON()
		return
	}
	guest.Password = ""
	c.Data["json"] = libs.ResultJson(guest, fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Guest
// @Param Authorization header string true "Bearer token"
// @Param	guestId		path 	string	true	"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty <br> 302 Record not found <br> 303 Save data failures
// @router /:guestId [delete]
func (c *GuestsController) DeletedGuest() {
	guestIdStr := c.Ctx.Input.Param(":guestId")
	guestId, _ := strconv.Atoi(guestIdStr)

	err := models.DeleteGuests(guestId)

	if err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Guest is not exists",
			map[string]interface{}{"Guest": conf.RECORD_NOT_FOUND,
			})
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson("", fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

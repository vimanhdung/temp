package smart_lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"indetail/conf"
	"indetail/libs"
	"indetail/models"
	"strconv"
	"strings"
)

// SmartLocksController operations for SmartLocks
type SmartLocksController struct {
	libs.Middleware
}

// URLMapping ...
func (c *SmartLocksController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create SmartLocks
// @Param Authorization header string true "Bearer token"
// @Param	body		body 	models.SmartLocks	true		"body for SmartLocks content"
// @Success 201 {int} models.SmartLocks
// @Failure 403 body is empty
// @router / [post]
func (c *SmartLocksController) Post() {
	var smartLock models.SmartLocks
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &smartLock); err != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), err.Error(), map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE})
		c.ServeJSON()
		return
	}

	if _, err := models.AddSmartLocks(&smartLock); err != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), err.Error(), map[string]interface{}{"SmartLock": conf.SAVE_FAILURES})
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(smartLock, fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get SmartLocks by id
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SmartLocks
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SmartLocksController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, intErr := strconv.Atoi(idStr)
	if intErr != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), intErr.Error(), map[string]interface{}{"Id": conf.VARIABLE_IS_NOT_POSITIVE_INTEGER})
		c.ServeJSON()
		return
	}
	smartLock, err := models.GetSmartLocksById(id)
	if err != nil || smartLock == nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Data not found", map[string]interface{}{"SmartLock": conf.RECORD_NOT_FOUND})
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(smartLock, fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get SmartLocks
// @Param Authorization header string true "Bearer token"
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SmartLocks
// @Failure 403
// @router / [get]
func (c *SmartLocksController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllSmartLocks(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the SmartLocks
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.SmartLocks	true		"body for SmartLocks content"
// @Success 200 {object} models.SmartLocks
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SmartLocksController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, intErr := strconv.Atoi(idStr)
	if intErr != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), intErr.Error(), map[string]interface{}{"Id": conf.VARIABLE_IS_NOT_POSITIVE_INTEGER})
		c.ServeJSON()
		return
	}

	smartLock := models.SmartLocks{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &smartLock); err != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), err.Error(), map[string]interface{}{"Body": conf.PARSE_JSON_BODY_FALSE})
		c.ServeJSON()
		return
	}

	queryCount := map[string]string{
		"Id":         strconv.Itoa(id),
		"Deleted_at": fmt.Sprint(conf.NOT_DELETED),
	}
	if total, errTotal := models.CountSmartLocks(queryCount); total < 1 || errTotal != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Data not found", map[string]interface{}{"SmartLock": conf.RECORD_NOT_FOUND})
		c.ServeJSON()
		return
	}

	if err := models.UpdateSmartLocksById(&smartLock); err != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), err.Error(), map[string]interface{}{"SmartLock": conf.SAVE_FAILURES})
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(smartLock, fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the SmartLocks
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SmartLocksController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, intErr := strconv.Atoi(idStr)
	if intErr != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), intErr.Error(), map[string]interface{}{"Id": conf.VARIABLE_IS_NOT_POSITIVE_INTEGER})
		c.ServeJSON()
		return
	}

	smartLock := models.SmartLocks{
		Id:        id,
		DeletedAt: conf.IS_DELETED,
	}
	if err := models.UpdateSmartLocksById(&smartLock); err != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), err.Error(), map[string]interface{}{"SmartLock": conf.SAVE_FAILURES})
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(smartLock, fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

// Lock/Unlock ...
// @Title Lock/unlock
// @Description lock or unlock smart lock
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The id of booking you want to change state smart lock"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /changeState/:id [post]
func (c *SmartLocksController) ChangeState() {
	idStr := c.Ctx.Input.Param(":id")
	bookingId, intErr := strconv.Atoi(idStr)
	if intErr != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), intErr.Error(),
			map[string]interface{}{"Id": conf.VARIABLE_IS_NOT_POSITIVE_INTEGER})
		c.ServeJSON()
		return
	}

	//get sesame key
	queryDevice := map[string]string{
		"bookingId": fmt.Sprint(bookingId),
	}
	if libs.GlobalUser.Type == conf.TYPE_GUEST {
		queryDevice["guestId"] = fmt.Sprint(libs.GlobalUser.Id)
	} else if libs.GlobalUser.Type == conf.TYPE_GUEST {
		queryDevice["userAppId"] = fmt.Sprint(libs.GlobalUser.Id)
	}
	deviceId, err := models.GetSmartLockDeviceIdByBookingId(queryDevice)
	if deviceId == "" || err != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Smart lock not found",
			map[string]interface{}{"SmartLock": conf.RECORD_NOT_FOUND})
		c.ServeJSON()
		return
	}

	//check key status
	request := httplib.Get(conf.CANDY_HOUSE_API_BASE_URL + "/" + deviceId)
	request.Header("Authorization", conf.CANDY_HOUSE_API_TOKEN)
	response, _ := request.String()

	var statusResponse models.GetStatusSmartLock
	if errStatus := json.Unmarshal([]byte(response), &statusResponse); errStatus != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), errStatus.Error(),
			map[string]interface{}{"SmartLock": conf.RECORD_NOT_FOUND})
		c.ServeJSON()
		return
	}
	// change state
	state := "lock"
	if statusResponse.Locked {
		state = "unlock"
	}
	requestState := httplib.Post(conf.CANDY_HOUSE_API_BASE_URL + "/" + deviceId)
	requestState.Header("Authorization", conf.CANDY_HOUSE_API_TOKEN)
	command := map[string]string{
		"command": state,
	}
	requestState.JSONBody(command)
	requestState.String()
	c.Data["json"] = libs.ResultJson(map[string]string{"State": state}, fmt.Sprint(conf.SUCCESS_STATUS), "Change smart lock state success", nil)
	c.ServeJSON()
}

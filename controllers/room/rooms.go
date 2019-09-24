package room

import (
	"encoding/json"
	"fmt"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
	"strconv"
	"strings"
)

// RoomsController operations for Rooms
type RoomsController struct {
	libs.Middleware
}

// URLMapping ...
func (c *RoomsController) URLMapping() {
	c.Mapping("CreateRoom", c.CreateRoom)
	c.Mapping("Detail", c.Detail)
	c.Mapping("GetListRoom", c.GetListRoom)
	c.Mapping("UpdateRoom", c.UpdateRoom)
	c.Mapping("Delete", c.Delete)
}

// CreateRoom ...
// @Title Create Room
// @Description create Rooms
// @Param	Authorization	header	string	true	"token"
// @Param	body		body 	models.InsertRooms	true		"body for Rooms content"
// @Success 201 {int} models.Rooms.Id
// @Failure 403 body is empty <br> 211 Variable is not json type <br> 217 Variable required <br> 219 Variable is not positive integer
// <br> 220 Variable is not date type <br> 303 Save data failures
// @router / [post]
func (c *RoomsController) CreateRoom() {
	var v models.Rooms

	jsonBody := SaveRoomProcess(v, c.Ctx.Input.RequestBody, false)
	c.Data["json"] = jsonBody
	c.ServeJSON()
}

// Detail ...
// @Title Detail
// @Description get Rooms by id
// @Param	Authorization	header	string	true	"token"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SwaggerDetailRoom
// @Failure 403 :id is empty <br> 201 Missing param <br> 219 Variable is not positive integer <br> 302 Record not found
// @router /:id [get]
func (c *RoomsController) Detail() {
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var detailErrorCode = make(map[string]interface{})

	idStr := c.Ctx.Input.Param(":id")

	if idStr == "" {
		errorMsg = "Missing room id"
		detailErrorCode["Id"] = conf.MISSNG_PARAM
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	id, idError := strconv.Atoi(idStr)

	if idError != nil || id < 0 {
		errorMsg = "Room id must be positive integer"
		detailErrorCode["Id"] = conf.VARIABLE_IS_NOT_POSITIVE_INTEGER
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	roomDetail, err := models.GetRoomDataById(id)
	if err {
		errorMsg = "No record found"
		detailErrorCode["Room"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(roomDetail, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
	c.ServeJSON()
}

// GetListRoom ...
// @Title Get List Room
// @Description get Rooms
// @Param	Authorization	header	string	true	"token"
// @Param	hotelId			query	int		true	"hotel id"
// @Param	roomName		query	string	false	"search like room name"
// @Param	status			query	string	false	"room status"
// @Param	fields			query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	limit			query	string	false	"Limit the size of result set. Must be an integer"
// @Param	page			query	string	false	"Start page of result set. Must be an integer"
// @Success 200 {object} models.SwaggerListRoom
// @Failure 403 <br> 219 Variable is not positive integer <br> 302 Record not found
// @router / [get]
func (c *RoomsController) GetListRoom() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64
	var page int64 = 1

	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var errorMsg string
	var detailErrorCode = make(map[string]interface{})

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// page: 0 (default is 0)
	if v, err := c.GetInt64("page"); err == nil {
		page = v
		if page < 1 {
			page = 1
		}
	}
	//status
	if v := c.GetString("status"); v != "" {
		if _, err := strconv.Atoi(v); err != nil {
			errorMsg = "Room status must be integer"
			detailErrorCode["Status"] = conf.VARIABLE_IS_NOT_POSITIVE_INTEGER
			c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
			c.ServeJSON()
			return
		}
		query["status"] = v
	}
	// hotelId
	if v, err := c.GetInt64("hotelId"); err == nil {
		query["HotelId"] = fmt.Sprint(v)
	} else {
		errorMsg = "hotelId must be integer"
		detailErrorCode["HotelId"] = conf.VARIABLE_IS_NOT_POSITIVE_INTEGER
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	//roomName
	if v := c.GetString("roomName"); v != "" {
		query["RoomName.in"] = v
	}

	queryParam := map[string]interface{}{
		"hotel_id": query["HotelId"],
	}

	if isExists, _ := models.FindHotel(queryParam); !isExists {
		detailErrorCode["Hotel"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, errorCode, "Hotel is not exists", nil)
		c.ServeJSON()
		return
	}

	totalRecord, _ := models.CountRooms(query)
	if totalRecord < 1 {
		errorMsg = "No record found"
		detailErrorCode["Room"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	if limit <= 0 {
		limit = totalRecord
	}
	offset = (page - 1) * limit

	listRoom, err := models.GetAllRooms(query, fields, sortby, order, offset, limit)
	if err != nil || listRoom == nil {
		errorMsg = "No record found"
		detailErrorCode["Room"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultPagingJson(listRoom, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, limit, page, totalRecord)
	c.ServeJSON()
}

// GetListRoom ...
// @Title Get List Room
// @Description get Rooms
// @Param	Authorization	header	string	true	"token"
// @Param	hotelId			query	int		true	"hotel id"
// @Param	checkinDate		query	string	false	"check in date. e.g 2018/12/11"
// @Param	checkoutDate	query	string	false	"check out date. e.g 2018/12/13"
// @Param	limit			query	int	false	"limit"
// @Success 200 {object} models.SwaggerListRoom
// @Failure 403 <br> 219 Variable is not positive integer <br> 302 Record not found
// @router /getListRoomAvailable/ [get]
func (c *RoomsController) GetListRoomAvailable() {
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64 = 0

	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var errorMsg string
	var detailErrorCode = make(map[string]interface{})
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}

	if v := c.GetString("checkinDate"); v != "" {
		query["checkinDate"] = v
	}
	if v := c.GetString("checkoutDate"); v != "" {
		query["checkoutDate"] = v
	}
	// hotelId
	if v, err := c.GetInt64("hotelId"); err == nil {
		query["HotelId"] = fmt.Sprint(v)
	} else {
		errorMsg = "hotelId must be integer"
		detailErrorCode["HotelId"] = conf.VARIABLE_IS_NOT_POSITIVE_INTEGER
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	if libs.GlobalUser.Email != "line@gmail.com" {
		queryParam := map[string]interface{}{
			"hotel_id":         query["HotelId"],
			"admin_account_id": libs.GlobalUser.Id,
		}
		if isExists, _ := models.CheckHotelExists(queryParam); !isExists {
			detailErrorCode["Hotel"] = conf.RECORD_NOT_FOUND
			c.Data["json"] = libs.ResultJson(nil, errorCode, "Hotel is not exists", nil)
			c.ServeJSON()
			return
		}
	} else {
		queryParam := map[string]interface{}{
			"hotel_id": query["HotelId"],
		}
		if isExists, _ := models.FindHotel(queryParam); !isExists {
			detailErrorCode["Hotel"] = conf.RECORD_NOT_FOUND
			c.Data["json"] = libs.ResultJson(nil, errorCode, "Hotel is not exists", nil)
			c.ServeJSON()
			return
		}
	}

	totalRecord, _ := models.CountRoomsAvailable(query)
	if totalRecord < 1 {
		errorMsg = "No record found"
		detailErrorCode["Room"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}
	if limit < 1 {
		limit = totalRecord
	}
	listRoom, err := models.GetRoomsAvailable(query, int(offset), int(limit))
	if err != nil || listRoom == nil {
		errorMsg = "No record found"
		detailErrorCode["Room"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(listRoom, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
	c.ServeJSON()
}

// UpdateRoom ...
// @Title UpdateRoom
// @Description update the Rooms
// @Param	Authorization	header	string	true	"token"
// @Param	id		path 	int	true		"The id you want to update"
// @Param	body		body 	models.InsertRooms	true		"body for Rooms content"
// @Success 200 {object} models.SwaggerCreateRoom
// @Failure 403 :id is not int <br> 211 Variable is not json type <br> 217 Variable required <br> 219 Variable is not positive integer
// <br> 220 Variable is not date type <br> 303 Save data failures
// @router /:id [put]
func (c *RoomsController) UpdateRoom() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Rooms{Id: id}

	jsonBody := SaveRoomProcess(v, c.Ctx.Input.RequestBody, true)
	c.Data["json"] = jsonBody
	c.ServeJSON()
}

func SaveRoomProcess(roomStruct models.Rooms, requestBody []byte, isUpdate bool) (jsonBody interface{}) {
	var isExistRoom bool
	var oldRoom models.Rooms
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var detailErrorCode = make(map[string]interface{})

	if errTmp := json.Unmarshal(requestBody, &roomStruct); errTmp != nil {
		errorMsg = "Wrong format"
		detailErrorCode["RoomData"] = conf.VARIABLE_IS_NOT_JSON
		jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		return
	}

	var isValid bool
	if isValid, detailErrorCode = validation.CheckValidate(roomStruct); !isValid {
		errorMsg = "Validate fail!"
	}

	//check smart lock
	availableSmartLock := map[string]string{
		"SmartLockId": strconv.Itoa(roomStruct.SmartLockId),
	}
	roomNameConditions := map[string]string{
		"HotelId":  strconv.Itoa(roomStruct.HotelId),
		"RoomName": roomStruct.RoomName,
	}
	if isUpdate {
		availableSmartLock["RoomId"] = strconv.Itoa(roomStruct.Id)
		//validate exist room
		roomConditions := map[string]string{
			"Id":      strconv.Itoa(roomStruct.Id),
			"HotelId": strconv.Itoa(roomStruct.HotelId),
		}
		isExistRoom, oldRoom = models.IsRoomAvailable(roomConditions)
		if !isExistRoom {
			errorMsg = "Room " + roomConditions["Id"] + " not found!"
			detailErrorCode["Room"] = conf.RECORD_NOT_FOUND
		}
		//validate name
		roomNameConditions["Id.isnotin"] = strconv.Itoa(roomStruct.Id)

	}
	//validate name
	isExistRoomName, _ := models.IsRoomAvailable(roomNameConditions)
	if isExistRoomName {
		errorMsg = "Room name exists"
		detailErrorCode["RoomName"] = conf.RECORD_EXISTS
	}

	if isAvailable, err := models.CheckAvailableSmartLock(availableSmartLock); !isAvailable || err != nil {
		errorMsg = "Smart lock not available"
		detailErrorCode["SmartLock"] = conf.LIMIT_RECORD
	}

	if len(detailErrorCode) > 0 {
		jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		return
	}

	//save data
	if isUpdate {
		roomStruct.CreatedAt = oldRoom.CreatedAt
		if err := models.UpdateRoomsById(&roomStruct); err != nil {
			errorCode = strconv.Itoa(conf.SAVE_FAILURES)
			errorMsg = err.Error()
			detailErrorCode["Room"] = conf.SAVE_FAILURES
		}
	} else {
		if _, err := models.AddRooms(&roomStruct); err != nil {
			errorCode = strconv.Itoa(conf.SAVE_FAILURES)
			errorMsg = err.Error()
			detailErrorCode["Room"] = conf.SAVE_FAILURES
		}
	}

	jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
	//override out put
	if errorMsg == "" {
		jsonBody = libs.ResultJson(roomStruct, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
	}
	return
}

// Delete ...
// @Title Delete
// @Description delete the Rooms
// @Param Authorization header string true "Bearer token"
// @Param	Authorization	header	string	true	"token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty <br> 302 Record not found <br> 303 Save data failures
// @router /:id [delete]
func (c *RoomsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	userGlobal := libs.GlobalUser
	var user *models.AdminAccounts
	user.Id = userGlobal.Id
	// Check Room
	if !models.CheckRoomExists(user, id) {
		c.Data["json"] = libs.ResultJson(
			"",
			fmt.Sprint(conf.ERROR_STATUS),
			"Room is not exists",
			map[string]interface{}{"Room": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}
	// Destroy Room
	if err := models.DestroyRoom(id); err != nil {
		c.Data["json"] = libs.ResultJson(
			"",
			"403",
			"Delete room false",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(nil, "200", "Delete room success", nil)
	c.ServeJSON()
}

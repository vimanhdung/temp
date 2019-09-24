package hotel

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

// HotelsController operations for Hotels
type HotelsController struct {
	libs.Middleware
}

// URLMapping ...
func (c *HotelsController) URLMapping() {
	c.Mapping("CreateHotel", c.CreateHotel)
	c.Mapping("UpdateHotel", c.UpdateHotel)
	c.Mapping("GetAllHotels", c.GetAllHotels)
	c.Mapping("DetailHotel", c.DetailHotel)
	c.Mapping("Delete", c.Delete)

}

// CreateHotel ...
// @Title Create Hotel
// @Description create Hotels
// @Param	Authorization	header	string	true	"token"
// @Param	body	body 	models.InsertHotel	true	"body for Hotels content"
// @Success 201 {int} models.Hotels.Id
// @Failure 403 body is empty <br> 211 Variable is not json type <br> 217 Variable required <br> 303 Save data failures
//
// @router / [post]
func (c *HotelsController) CreateHotel() {
	var v models.Hotels
	jsonBody, _ := SaveHotelProcess(v, c.Ctx.Input.RequestBody, false)
	c.Data["json"] = jsonBody
	c.ServeJSON()
}

// UpdateHotel ...
// @Title Update Hotel
// @Description update the Hotels
// @Param	Authorization	header	string				true	"token"
// @Param	id				path 	string				true	"The id you want to update"
// @Param	body			body 	models.InsertHotel	true	"body for Hotels content"
// @Success 200 {object} models.SwaggerCreateHotels
// @Failure 403 :id is not int <br> 211 Variable is not json type <br> 217 Variable required <br> 303 Save data failures
// @router /:id [put]
func (c *HotelsController) UpdateHotel() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	var v = models.Hotels{Id: id}

	jsonBody, _ := SaveHotelProcess(v, c.Ctx.Input.RequestBody, true)
	c.Data["json"] = jsonBody
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Hotels
// @Param Authorization header string true "Bearer token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty <br> 302 Record not found <br> 303 Save data failures
// @router /:id [delete]
func (c *HotelsController) Delete() {
	// Check Permission
	/*if !c.PermissionDenied(conf.HOTEL_DELETE) {
		return
	}*/
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	// Check hotel exists
	queryParam := map[string]interface{}{
		"hotel_id":         idStr,
		"admin_account_id": libs.GlobalUser.Id,
	}
	if isExists, _ := models.FindHotel(queryParam); !isExists {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Hotel is not exists",
			map[string]interface{}{"Hotel": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}
	// Delete hotel
	if err := models.SoftDeleteHotel(id); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Delete false",
			map[string]interface{}{"Database": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), "Delete Success", nil)
	c.ServeJSON()
}

// GetAllHotels ...
// @Title Get All Hotels
// @Description get Hotels
// @Param	Authorization	header	string	true	"token"
// @Param	hotelName		query	string	false	"search like hotel name"
// @Param	createdAtTo		query	string	false	"hotel created at. Must be date eg. 2018-12-27"
// @Param	createdAtFrom	query	string	false	"hotel created at. Must be date eg. 2018-12-27"
// @Param	limit			query	string	false	"Limit the size of result set. Must be an integer"
// @Param	page			query	string	false	"Page of result set. Must be an integer"
// @Success 200 {object} models.SwaggerListHotels
// @Failure 403 <br> 211 Variable is not json type <br> 217 Variable required <br> 220 Variable is not date type
// @router /GetAllHotels/ [get]
func (c *HotelsController) GetAllHotels() {
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64
	var page int64 = 1

	var detailErrorCode = make(map[string]interface{})
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
	// hotelName
	if v := strings.Trim(c.GetString("hotelName"), " "); v != "" {
		query["HotelName.icontains"] = v
	}
	//created from date
	var validateDate = true
	if v := c.GetString("createdAtFrom"); v != "" {
		if validation.CheckDate(v) {
			query["CreatedAt.gte"] = v
		} else {
			validateDate = false
			detailErrorCode["CreatedAtFrom"] = fmt.Sprint(conf.VARIABLE_IS_NOT_DATE)
			errorMsg = "Wrong date format. Default date format: " + fmt.Sprint(conf.RFC3339)
		}
	}
	//created to date
	if v := c.GetString("createdAtTo"); v != "" {
		if validation.CheckDate(v) {
			query["CreatedAt.lte"] = v + " 23:59:59"
		} else {
			validateDate = false
			detailErrorCode["CreatedAtTo"] = fmt.Sprint(conf.VARIABLE_IS_NOT_DATE)
			errorMsg = "Wrong date format. Default date format: " + fmt.Sprint(conf.RFC3339)
		}
	}

	//validate created_at
	if !validateDate {
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	if query["CreatedAt.gte"] == "" && query["CreatedAt.lte"] != "" {
		errorMsg = "Missing param"
		detailErrorCode["CreatedAtFrom"] = fmt.Sprint(conf.MISSNG_PARAM)
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	totalRecord, err := models.CountHotels(query)
	if totalRecord == 0 {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	if limit <= 0 {
		limit = totalRecord
	}

	//set offset
	offset = (page - 1) * limit

	tmpListHotel, err := models.GetAllHotels(query, []string{}, sortby, order, offset, limit)

	if err != nil || tmpListHotel == nil {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultPagingJson(tmpListHotel, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, limit, page, totalRecord)
	c.ServeJSON()
}

func SaveHotelProcess(hotelStruct models.Hotels, requestBody []byte, isUpdate bool) (jsonBody interface{}, httpStatus int) {
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	httpStatus = 403
	var detailErrorCode = make(map[string]interface{})
	var oldHotelData models.Hotels
	if isUpdate {
		queryParam := map[string]interface{}{
			"hotel_id": hotelStruct.Id,
		}
		var isExists bool
		if isExists, oldHotelData = models.FindHotel(queryParam); !isExists {
			jsonBody = libs.ResultJson(nil, errorCode, "Hotel is not exists", detailErrorCode)
			return
		}
	}

	var hotelTmpInterface map[string]interface{}
	if errTmp := json.Unmarshal(requestBody, &hotelTmpInterface); errTmp != nil {
		errorMsg = "Wrong format"
		detailErrorCode["HotelData"] = conf.VARIABLE_IS_NOT_JSON
		jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		return
	}

	var files interface{}
	for index, value := range hotelTmpInterface {
		if index == "Files" && value != "" {
			files = value
			hotelTmpInterface[index] = ""
		}
	}

	if !isUpdate {
		delete(hotelTmpInterface, "Id")
	}

	//validate file
	var hotelFilesInsert = "[]"
	var rebuildValidate = make(map[string]interface{})
	if files != nil {
		hotelFilesInsert, rebuildValidate = models.ValidateHotelFiles(files, detailErrorCode)
		if hotelFilesInsert == "" {
			errorMsg = "Wrong format"
			detailErrorCode["Files"] = conf.VARIABLE_IS_NOT_JSON
			jsonBody = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
			return
		}
	}

	hotelId := hotelStruct.Id
	reMakeHotelInfo, _ := json.Marshal(hotelTmpInterface)
	if err := json.Unmarshal(reMakeHotelInfo, &hotelStruct); err == nil {
		if isUpdate {
			hotelStruct.Id = hotelId
			hotelStruct.CreatedAt = oldHotelData.CreatedAt
			hotelStruct.UpdatedAt = time.Now()
			hotelStruct.UpdatedUser = libs.GlobalUser.Id
		} else {
			hotelStruct.CreatedUser = libs.GlobalUser.Id
		}
		hotelStruct.Files = hotelFilesInsert

		var isValid bool
		var listErrorCode = make(map[string]interface{})
		isValid, listErrorCode = validation.CheckValidate(hotelStruct)

		if isValid {
			if isUpdate {
				if err := models.UpdateHotelsById(&hotelStruct); err != nil {
					rebuildValidate["UpdateHotel"] = conf.SAVE_FAILURES
					errorMsg = err.Error()
				}
			} else {
				if _, err := models.AddHotels(&hotelStruct); err != nil {
					rebuildValidate["InsertHotel"] = conf.SAVE_FAILURES
					errorMsg = err.Error()
				}
			}
		} else {
			for key, value := range listErrorCode {
				rebuildValidate[key] = value
			}
			errorMsg = "Wrong format"
		}
	} else {
		errorMsg = "Wrong format"
	}

	jsonBody = libs.ResultJson(nil, errorCode, errorMsg, rebuildValidate)

	//override out put
	if errorMsg == "" {
		if isUpdate {
			errorMsg = "Update hotel success!"
			httpStatus = 200
		} else {
			errorMsg = "Add new hotel success!"
			httpStatus = 201
		}
		jsonBody = libs.ResultJson(hotelStruct, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, rebuildValidate)
	}
	return
}

// DetailHotel ...
// @Title DetailHotel
// @Description get Hotels by id
// @Param	Authorization	header	string	true	"token"
// @Param	id				path 	string	true	"The key for staticblock"
// @Success 200 {object} models.SwaggerDetailHotels
// @Failure 403 :id is empty <br> 201 Missing param <br> 219 Variable is not positive integer
// @router /detailHotel/:id [get]
func (c *HotelsController) DetailHotel() {
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var detailErrorCode = make(map[string]interface{})

	idStr := c.Ctx.Input.Param(":id")

	if idStr == "" {
		errorMsg = "Missing hotel id"
		detailErrorCode["Id"] = conf.MISSNG_PARAM
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	id, idError := strconv.Atoi(idStr)
	if idError != nil || id < 0 {
		errorMsg = "Hotel id must be positive integer"
		detailErrorCode["Id"] = conf.VARIABLE_IS_NOT_POSITIVE_INTEGER
		c.Data["json"] = libs.ResultJson(nil, errorCode, errorMsg, detailErrorCode)
		c.ServeJSON()
		return
	}

	hotelDetail, err := models.GetHotelsById(id)
	if err != nil || hotelDetail == nil {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(hotelDetail, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, nil)
	c.ServeJSON()
}

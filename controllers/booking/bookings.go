package booking

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/migration"
	"indetail/conf"
	"indetail/libs"
	"indetail/libs/validation"
	"indetail/models"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// BookingsController operations for Bookings
type BookingsController struct {
	libs.Middleware
}

// URLMapping ...
func (c *BookingsController) URLMapping() {
	c.Mapping("CreateBooking", c.CreateBooking)
	c.Mapping("Detail", c.Detail)
	c.Mapping("GetListBooking", c.GetListBooking)
	c.Mapping("BookingVerifyPassport", c.BookingVerifyPassport)
	c.Mapping("Update", c.Update)
	c.Mapping("Reject", c.Reject)
}

// CreateBooking ...
// @Title CreateBooking
// @Description create Bookings
// @Param Authorization header string true "Bearer token"
// @Param	body		body 	models.AddBookingStruct	true		"body for Bookings content"
// @Success 201 {int} models.Booking.Id
// @Failure 403 101 Permission deny <br> 211 Variable is not json <br> 217 Variable required <br> 227 Date time invalid <br> 228 Param value not equal  <br> 302 Record not found <br> 303 Save data failure <br> 306 Create file error <br> 605 Room not available
// @router / [post]
func (c *BookingsController) CreateBooking() {
	var booking models.AddBookingStruct
	errStatus := fmt.Sprint(conf.ERROR_STATUS)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &booking); err != nil {
		c.Data["json"] = libs.ResultJson(nil, errStatus, "Booking data is not json type",
			map[string]interface{}{"Booking": conf.VARIABLE_IS_NOT_JSON})
		c.ServeJSON()
		return
	}

	var queryHotel = make(map[string]interface{})
	queryHotel["room_id"] = booking.RoomId
	exists, hotelInfo := models.CheckHotelAndRoom(queryHotel)

	if !exists {
		c.Data["json"] = libs.ResultJson(nil, errStatus, "Hotel permission deny",
			map[string]interface{}{"Hotel": conf.PERMISSION_DENY})
		c.ServeJSON()
		return
	}

	//validate booking
	var detailErrorCode = make(map[string]interface{})
	booking.HotelId = hotelInfo.Id
	if isValid, listErrorCode := validation.CheckValidate(booking); !isValid {
		for key, value := range listErrorCode {
			detailErrorCode[key] = value
		}
	}

	if booking.CheckinDate > booking.CheckoutDate {
		detailErrorCode["CheckinDate"] = conf.DATE_TIME_INVALID
	}

	if booking.GuestCount != (booking.AdultCount + booking.ChildrenCount) {
		detailErrorCode["GuestCount"] = conf.VARIABLE_NOT_EQUAL
	}
	if len(detailErrorCode) > 0 {
		c.Data["json"] = libs.ResultJson(nil, errStatus, "Booking wrong format",
			detailErrorCode)
		c.ServeJSON()
		return
	}

	booking.Status = conf.BOOKING_ACTIVE

	//check Available
	var queryRoom = make(map[string]string)
	queryRoom["RoomId"] = fmt.Sprint(booking.RoomId)
	queryRoom["HotelId"] = fmt.Sprint(booking.HotelId)
	queryRoom["checkinDate"] = fmt.Sprint(booking.CheckinDate)
	queryRoom["checkoutDate"] = fmt.Sprint(booking.CheckoutDate)
	total, errGetRoom := models.CountRoomsAvailable(queryRoom)
	if errGetRoom != nil || total < 1 {
		detailErrorCode["RoomId"] = conf.ROOM_NOT_AVAILABLE
		c.Data["json"] = libs.ResultJson(nil, errStatus, "Room not available",
			detailErrorCode)
		c.ServeJSON()
		return
	}

	//validate guest
	if len(booking.GuestInsert) < 1 {
		detailErrorCode["Guest"] = conf.VARIABLE_REQUIRED
		c.Data["json"] = libs.ResultJson(nil, errStatus, "Missing main guest",
			detailErrorCode)
		c.ServeJSON()
		return
	}
	var queryGuest = make(map[string]string)
	for _, guest := range booking.GuestInsert {
		queryGuest["GuestId"] += fmt.Sprint(guest.Id) + ","
	}
	queryGuest["GuestId"] = strings.TrimRight(queryGuest["GuestId"], ",")
	if totalRecord, err := models.CountTotalRecord(queryGuest, "guests.guest_id"); err != nil || totalRecord < int64(len(booking.GuestInsert)) {
		detailErrorCode["GuestId"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, errStatus, "Guest not found",
			detailErrorCode)
		c.ServeJSON()
		return
	}

	globalUser := libs.GlobalUser
	detailErrorCode, err, resultBooking := models.AddBookingWithTransaction(&booking, globalUser)
	if err != nil {
		c.Data["json"] = libs.ResultJson(nil, errStatus, err.Error(), detailErrorCode)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(resultBooking, fmt.Sprint(conf.SUCCESS_STATUS), "", detailErrorCode)
	c.ServeJSON()
}

// Detail ...
// @Title Detail
// @Description get Bookings by id
// @Param Authorization header string true "Bearer token"
// @Param	bookingId	path 	string	true		"Id of Booking"
// @Success 200 {object} models.SwaggerDetailBooking
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Booking not exists
// @router /:bookingId [get]
func (c *BookingsController) Detail() {
	bookingIdStr := c.Ctx.Input.Param(":bookingId")
	bookingId, _ := strconv.Atoi(bookingIdStr)
	params := map[string]string{
		"bookingId": bookingIdStr,
	}

	// Check Permission
	globalUser := libs.GlobalUser
	switch globalUser.Type {
	case conf.TYPE_GUEST:
		params["guestId"] = fmt.Sprint(globalUser.Id)
	case conf.TYPE_USER:
		params["userAppId"] = fmt.Sprint(globalUser.Id)
	default:
		params["adminAccountId"] = fmt.Sprint(globalUser.Id)
	}

	// Check Booking exists
	if _, exists := models.CheckBookingExist(params); !exists {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking is not exists",
			map[string]interface{}{"Booking": conf.RECORD_NOT_FOUND,
			})
		c.ServeJSON()
		return
	}
	// Get Detail Booking
	booking, _ := models.GetBookingsById(bookingId)
	var tmpBooking []models.BookingData
	tmpBooking = append(tmpBooking, *booking)
	query := map[string]string{
		"GuestId": "",
	}
	listBooking := models.GetExtraInfoForBooking(tmpBooking, query)
	c.Data["json"] = libs.ResultJson(listBooking[0], fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

// GetListBooking ...
// @Title Get List Booking
// @Description get Bookings
// @Param Authorization header string true "Bearer token
// @Param	hotelId				query	string		false	"hotelId"
// @Param	checkinFromDate		query	string	false	"check in from date. e.g 2018/12/11"
// @Param	checkinToDate		query	string	false	"check in to date. e.g 2018/12/13"
// @Param	checkoutFromDate	query	string	false	"check out from date. e.g 2018/12/11"
// @Param	checkoutToDate		query	string	false	"check out to date. e.g 2018/12/13"
// @Param	fields				query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	limit				query	string	false	"Limit the size of result set. Must be an integer"
// @Param	page				query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SwaggerListBooking
// @Failure 403 101 Permission deny <br> 201 Missing param <br> 211 Variable is not json <br> 219 Variable is not positive integer <br> 220 Variable is not date type <br> 302 Record not found
// @router / [get]
func (c *BookingsController) GetListBooking() {
	//authen
	loginUser := libs.GlobalUser
	//process booking
	var fields = "bookings.*"
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
	var errorMsg = ""
	var errorCode = fmt.Sprint(conf.ERROR_STATUS)
	var detailErrorCode = make(map[string]interface{})

	// offset: 0 (default is 0)
	query["HotelId"] = ""
	if v := c.GetString("hotelId"); v != "" {
		if intHotelId, hotelIdErr := strconv.Atoi(v); hotelIdErr == nil {
			query["HotelId"] = fmt.Sprint(intHotelId)
		} else {
			errorMsg = "HotelId must be positive integer"
			detailErrorCode["hotelId"] = fmt.Sprint(conf.VARIABLE_IS_NOT_POSITIVE_INTEGER)
			c.Data["json"] = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
			c.ServeJSON()
			return
		}
	}

	if loginUser.Type == conf.TYPE_USER {
		if val, ok := query["HotelId"]; !ok || val == "" {
			errorMsg = "missing param"
			detailErrorCode["hotelId"] = fmt.Sprint(conf.VARIABLE_REQUIRED)
			c.Data["json"] = libs.ResultJson("", errorCode, errorMsg, detailErrorCode)
			c.ServeJSON()
			return
		}
		query["UserAppId"] = fmt.Sprint(loginUser.Id)
	}

	query["GuestId"] = ""
	if loginUser.Type == conf.TYPE_GUEST {
		query["GuestId"] = fmt.Sprint(loginUser.Id)
	}

	if v := c.GetString("checkinFromDate"); v != "" {
		query["CheckinDate.gte"] = v + " 00:00:00"
	}
	if v := c.GetString("checkinToDate"); v != "" {
		query["CheckinDate.lte"] = v + " 23:59:59"
	}
	if v := c.GetString("checkoutFromDate"); v != "" {
		query["CheckoutDate.gte"] = v + " 00:00:00"
	}
	if v := c.GetString("checkoutToDate"); v != "" {
		query["CheckoutDate.lte"] = v + " 23:59:59"
	}

	//created from date
	if jsonResult, err := ValidateBookingDate(&query, &detailErrorCode); err == fmt.Sprint(conf.ERROR_STATUS) {
		c.Data["json"] = jsonResult
		c.ServeJSON()
		return
	}

	//count total
	totalRecord, err := models.CountBooking(fields, query)
	if totalRecord == 0 || err != nil {
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
	listBooking, msgErr := models.GetAllBookingsByCondition(query, fields, sortby, order, offset, limit)
	if msgErr != "" || listBooking == nil {
		errorMsg = "No record found"
		c.Data["json"] = libs.ResultJson("", fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, nil)
		c.ServeJSON()
		return
	}

	resultListBooking := models.GetExtraInfoForBooking(listBooking, query)

	c.Data["json"] = libs.ResultPagingJson(resultListBooking, fmt.Sprint(conf.SUCCESS_STATUS), errorMsg, limit, page, totalRecord)
	c.ServeJSON()
}

func ValidateBookingDate(query *map[string]string, detailErrorCode *map[string]interface{}) (jsonResult interface{}, errorCode string) {
	errorCode = fmt.Sprint(conf.ERROR_STATUS)
	//checkin date
	if (*query)["CheckinDate.gte"] == "" && (*query)["CheckinDate.lte"] != "" {
		(*detailErrorCode)["checkinFromDate"] = conf.MISSNG_PARAM
		jsonResult = libs.ResultJson("", errorCode, "Missing param checkinFromDate", *detailErrorCode)
		return
	}
	if _, err := time.Parse(conf.RFC_DATE_TIME, (*query)["CheckinDate.gte"]); err != nil && (*query)["CheckinDate.gte"] != "" {
		(*detailErrorCode)["checkinFromDate"] = conf.VARIABLE_IS_NOT_DATE
		jsonResult = libs.ResultJson("", errorCode, "Wrong checkinFromDate date format. Default date format: "+fmt.Sprint(conf.RFC3339), *detailErrorCode)
		return
	}
	if _, err := time.Parse(conf.RFC_DATE_TIME, (*query)["CheckinDate.lte"]); err != nil && (*query)["CheckinDate.lte"] != "" {
		(*detailErrorCode)["checkinToDate"] = conf.VARIABLE_IS_NOT_DATE
		jsonResult = libs.ResultJson("", errorCode, "Wrong checkinToDate date format. Default date format: "+fmt.Sprint(conf.RFC3339), *detailErrorCode)
		return
	}

	//checkout date
	if (*query)["CheckoutDate.gte"] == "" && (*query)["CheckoutDate.lte"] != "" {
		(*detailErrorCode)["checkoutFromDate"] = conf.MISSNG_PARAM
		jsonResult = libs.ResultJson("", errorCode, "Missing param checkoutFromDate", *detailErrorCode)
		return
	}
	if _, err := time.Parse(conf.RFC_DATE_TIME, (*query)["CheckoutDate.gte"]); err != nil && (*query)["CheckoutDate.gte"] != "" {
		(*detailErrorCode)["checkoutFromDate"] = conf.VARIABLE_IS_NOT_DATE
		jsonResult = libs.ResultJson("", errorCode, "Wrong checkoutFromDate date format. Default date format: "+fmt.Sprint(conf.RFC3339), *detailErrorCode)
		return
	}
	if _, err := time.Parse(conf.RFC_DATE_TIME, (*query)["CheckoutDate.lte"]); err != nil && (*query)["CheckoutDate.lte"] != "" {
		(*detailErrorCode)["checkoutToDate"] = conf.VARIABLE_IS_NOT_DATE
		jsonResult = libs.ResultJson("", errorCode, "Wrong checkoutToDate date format. Default date format: "+fmt.Sprint(conf.RFC3339), *detailErrorCode)
		return
	}

	errorCode = fmt.Sprint(conf.SUCCESS_STATUS)
	return
}

// Update ...
// @Title Update
// @Description update the Bookings
// @Param Authorization header string true "Bearer token
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.BookingUpdateFields	true		"body for Bookings content"
// @Success 200 {object} {object} Update Success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 217 : Fields required <br> 302 : Booking not exists <br> 303 : Update booking false <br> 221 : Parse json false
// @router /:id [put,post]
func (c *BookingsController) Update() {
	// Check Permission
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)

	var ob models.BookingUpdateFields
	userGlobal := libs.GlobalUser
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

	// Check Validate
	if b, errorCode := validation.CheckValidate(ob); !b {
		c.Data["json"] = libs.ResultJson("", fmt.Sprint(conf.ERROR_STATUS), "Validate false", errorCode)
		c.ServeJSON()
		return
	}
	// Check Booking exists
	par := map[string]string{
		"bookingId": fmt.Sprint(id),
	}

	// Check Permission
	globalUser := libs.GlobalUser
	switch globalUser.Type {
	case conf.TYPE_GUEST:
		par["guestId"] = fmt.Sprint(globalUser.Id)
	case conf.TYPE_USER:
		par["userAppId"] = fmt.Sprint(globalUser.Id)
	}

	// Check Booking exists
	var listBooking []models.Bookings
	var existsBooking bool
	if listBooking, existsBooking = models.CheckBookingExist(par); !existsBooking {
		c.Data["json"] = libs.ResultJson(
			"",
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking is not exists",
			map[string]interface{}{"Id": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}
	booking := listBooking[0]

	//validate booking
	detailErrorCode := make(map[string]interface{})
	listProcessGuest := make(map[string][]int)
	if detailErrorCode, listProcessGuest = ValidateUpdateBooking(booking, ob); len(detailErrorCode) > 0 {
		c.Data["json"] = libs.ResultJson("", fmt.Sprint(conf.ERROR_STATUS), "Validate fail", detailErrorCode)
		c.ServeJSON()
		return
	}

	// Update Booking
	if ob.BookingStatus == conf.BOOKING_REJECT {
		c.Data["json"] = libs.ResultJson(
			"",
			fmt.Sprint(conf.ERROR_STATUS),
			"Permission denied",
			map[string]interface{}{"Auth": conf.PERMISSION_DENY},
		)
		c.ServeJSON()
		return
	}

	//send fcm notification
	params := models.GetParamsUpdate(ob, userGlobal)

	//create QR code
	if booking.GuestCount == int8(len(ob.GuestUpdate)) && booking.QrImageUrl == "" {
		errGenQr, path := models.GenQrCodeForBooking(booking.Id)
		if errGenQr != nil {
			fmt.Println(errGenQr.Error())
			c.Data["json"] = libs.ResultJson(
				nil,
				fmt.Sprint(conf.ERROR_STATUS),
				"Gen QR code fail",
				map[string]interface{}{"QR": conf.SAVE_FAILURES},
			)
			c.ServeJSON()
			return
		}
		params["qr_image_url"] = beego.AppConfig.String("baseServer") + path
	}

	if errBooking := models.UpdateBookingTransaction(booking, params, listProcessGuest); errBooking != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			errBooking.Error(),
			map[string]interface{}{"Update": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}

	//get info
	resultBooking, _ := models.GetBookingsById(id)
	var tmpBookingData []models.BookingData
	tmpBookingData = append(tmpBookingData, *resultBooking)
	query := map[string]string{
		"GuestId": "",
	}
	listBookingData := models.GetExtraInfoForBooking(tmpBookingData, query)
	c.Data["json"] = libs.ResultJson(
		listBookingData[0],
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Update Success",
		nil,
	)
	c.ServeJSON()
}

func ValidateUpdateBooking(booking models.Bookings, ob models.BookingUpdateFields) (map[string]interface{}, map[string][]int) {
	var detailErrorCode = make(map[string]interface{})
	if booking.GuestCount < ob.ActualGuestCount {
		detailErrorCode["ActualGuestCount"] = conf.LIMIT_RECORD
	}

	if booking.GuestCount < int8(len(ob.GuestUpdate)) {
		detailErrorCode["GuestUpdate"] = conf.LIMIT_RECORD
	}

	if ob.CheckinDate != "" {
		tmpTimeValue, err := time.Parse(conf.RFC3339, ob.CheckinDate)
		if err != nil {
			detailErrorCode["CheckinDate"] = conf.VARIABLE_IS_NOT_DATE
		}
		ob.CheckinDate = tmpTimeValue.String()
		booking.CheckinDate = ob.CheckinDate
	}

	if ob.CheckoutDate != "" {
		tmpTimeValue, err := time.Parse(conf.RFC3339, ob.CheckoutDate)
		if err != nil {
			detailErrorCode["CheckoutDate"] = conf.VARIABLE_IS_NOT_DATE
		}
		ob.CheckoutDate = tmpTimeValue.String()
		booking.CheckoutDate = ob.CheckoutDate
	}

	if ob.GuestCount > 0 {
		booking.GuestCount = ob.GuestCount
	}
	if ob.AdultCount > 0 {
		booking.AdultCount = ob.AdultCount
	}
	if ob.ChildrenCount > 0 {
		booking.ChildrenCount = ob.ChildrenCount
	}

	if booking.CheckinDate > booking.CheckoutDate {
		detailErrorCode["CheckinDate"] = conf.DATE_TIME_INVALID
	}
	if booking.GuestCount != (booking.AdultCount + booking.ChildrenCount) {
		detailErrorCode["GuestCount"] = conf.VARIABLE_NOT_EQUAL
	}

	if ob.RoomId > 0 {
		var queryRoom = make(map[string]string)
		queryRoom["RoomId"] = fmt.Sprint(booking.RoomId)
		queryRoom["HotelId"] = fmt.Sprint(booking.HotelId)
		queryRoom["checkinDate"] = fmt.Sprint(booking.CheckinDate)
		queryRoom["checkoutDate"] = fmt.Sprint(booking.CheckoutDate)
		total, errGetRoom := models.CountRoomsAvailable(queryRoom)
		if errGetRoom != nil || total < 1 {
			detailErrorCode["RoomId"] = conf.ROOM_NOT_AVAILABLE
		}
	}

	var listProcessGuest = make(map[string][]int)
	if len(ob.GuestUpdate) < 1 {
		return detailErrorCode, listProcessGuest
	}

	var queryGuest = make(map[string]string)
	queryGuest["BookingId"] = fmt.Sprint(booking.Id)
	oldBookingGuest, _ := models.GetListBookingGuest(queryGuest)
	if oldBookingGuest == nil {
		detailErrorCode["BookingGuest"] = conf.RECORD_NOT_FOUND
	}
	//rebuild update guest id
	var rebuildGuestArray = make(map[int]int)
	stringGuestUpdateId := ""
	for _, guest := range ob.GuestUpdate {
		rebuildGuestArray[int(guest.Id)] = int(guest.Id)
		stringGuestUpdateId += fmt.Sprint(guest.Id) + ","
	}
	stringGuestUpdateId = strings.TrimRight(stringGuestUpdateId, ",")

	//check exists main guest
	var tmpListIdRemove = make(map[int]int)
	var listIdInsert []int
	var listIdRemove []int
	var listIdDoNothing = make(map[int]int)

	for _, oldGuest := range oldBookingGuest {
		if oldGuest.IsMainGuest != conf.GUEST_IS_MAIN {
			continue
		}
		//add main guest to list do nothing
		listIdDoNothing[oldGuest.GuestId] = oldGuest.GuestId
		if _, ok := rebuildGuestArray[oldGuest.GuestId]; !ok {
			detailErrorCode["MainGuest"] = conf.RECORD_NOT_FOUND
			break
		}
	}

	//get insert, delete booking guest
	for _, oldGuest := range oldBookingGuest {
		_, ok := rebuildGuestArray[oldGuest.GuestId]
		_, okNothing := listIdDoNothing[int(oldGuest.GuestId)]
		if !ok && !okNothing {
			tmpListIdRemove[oldGuest.GuestId] = oldGuest.GuestId
			listIdRemove = append(listIdRemove, oldGuest.GuestId)
		} else {
			listIdDoNothing[oldGuest.GuestId] = oldGuest.GuestId
		}
	}

	for _, guest := range ob.GuestUpdate {
		_, ok := tmpListIdRemove[int(guest.Id)]
		_, okNothing := listIdDoNothing[int(guest.Id)]
		if !ok && !okNothing {
			listIdInsert = append(listIdInsert, int(guest.Id))
		}
	}
	listProcessGuest["insert"] = listIdInsert
	listProcessGuest["remove"] = listIdRemove

	//check exists guest
	var queryCountGuest = make(map[string]string)
	queryCountGuest["GuestId"] = stringGuestUpdateId
	if totalRecord, err := models.CountTotalRecord(queryCountGuest, "guests.guest_id"); err != nil || int(totalRecord) != len(ob.GuestUpdate) {
		detailErrorCode["Guest"] = conf.RECORD_NOT_FOUND
	}

	return detailErrorCode, listProcessGuest
}

// Reject ...
// @Title Reject
// @Description update the Bookings
// @Param Authorization header string true "Bearer token
// @Param	id		path 	string	true		"The id you want to reject"
// @Success 200 {object} {object} Reject success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Booking not exists <br> 303 : Update booking false
// @router /reject/:id [put]
func (c *BookingsController) Reject() {
	// Check Permission
	/*if !c.PermissionDenied(conf.BOOKING_REJECT_STATUS) {
		return
	}*/
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	userGlobal := libs.GlobalUser
	var user *models.AdminAccounts
	user.Id = userGlobal.Id
	// Check Booking exists
	par := map[string]string{
		"bookingId": fmt.Sprint(id),
	}
	if _, exists := models.CheckBookingExist(par); !exists {
		c.Data["json"] = libs.ResultJson(
			"",
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking is not exists",
			map[string]interface{}{"Id": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}
	// Update Status Booking To Reject
	params := models.GetParamsReject()
	if err := models.UpdateBookingWithParams(id, params); err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Update booking false",
			map[string]interface{}{"Update": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), "Reject Success", nil)
	c.ServeJSON()
}

// BookingCharge ...
// @Title BookingCharge
// @Description update the Booking Charges
// @Param Authorization header string true "Bearer token
// @Param	id		path 	string	true		"The booking id you want to charge"
// @Success 200 {object} {object} Charge success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Booking not exists <br> 303 : Update booking false
// @router /reject/:id [put]
func (c *BookingsController) BookingCharge() {
	// build booking charge
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	log.Println(id)

	c.Data["json"] = libs.ResultJson(nil, "200", "Reject Success", nil)
	c.ServeJSON()
}

func CreateBookingChargeData(booking models.Bookings) int64 {
	chargeTime, err := time.Parse(conf.RFC_DATE_TIME, fmt.Sprint(booking.CheckinDate))
	if err != nil {
		chargeTime = time.Now()
	}

	var bookingChargeInput models.BookingCharges
	bookingChargeInput.ListBookingServiceId = ""
	bookingChargeInput.BookingId = booking.Id
	bookingChargeInput.Note = "Charge success"
	bookingChargeInput.Status = conf.BOOKING_CHARGE_STATUS_CHARGED
	bookingChargeInput.Amount = booking.TotalAmount
	bookingChargeInput.Discount = ""
	bookingChargeInput.TotalAmount = booking.TotalAmount
	bookingChargeInput.ChargeDatetime = chargeTime
	bookingChargeInput.CreatedAt = time.Now()

	id, err := models.AddBookingCharge(&bookingChargeInput)
	if err != nil {
		return 0
	}
	return id
}

func UpdateBookingChargeData(bookingChargeId int64, chargeInfo []byte) bool {
	bookingCharge, err := models.GetBookingChargeById(int(bookingChargeId))
	if err != nil {
		return false
	}

	var chargeStruct models.ChargeInfo
	if err := json.Unmarshal(chargeInfo, &chargeStruct); err != nil {
		return false
	}

	if chargeStruct.Status != conf.CHARGE_INFO_SUCCESS_STATUS {
		bookingCharge.Status = conf.BOOKING_CHARGE_STATUS_FAILURE
	} else {
		bookingCharge.Status = conf.BOOKING_CHARGE_STATUS_CHARGED
	}
	bookingCharge.Note = ""
	bookingCharge.UpdatedAt = time.Now()
	userGlobal := libs.GlobalUser
	var user *models.AdminAccounts
	user.Id = userGlobal.Id
	if err := models.UpdateBookingChargeById(bookingCharge, user); err != nil {
		return false
	}
	return true
}

// Confirm Checkin ...
// @Title Confirm Checkin
// @Description Confirm Checkin
// @Param Authorization header string true "Bearer token
// @Param	id		path 	string	true		"The booking id you want to checkin"
// @Success 200 {object} {object} Checkin success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Booking not exists <br> 303 : Checkin false <br> 601 : Room not vacant <br> 603 : Booking is checkin already <br> 604 : Not checkin day
// @router /confirmCheckin/:id [post]
func (c *BookingsController) ConfirmCheckin() {
	responseErrStatus := fmt.Sprint(conf.ERROR_STATUS)
	//check permission
	userGlobal := libs.GlobalUser
	if userGlobal.Type != conf.TYPE_ADMIN {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, "permission denied", map[string]interface{}{
			"Booking": conf.PERMISSION_DENY,
		})
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, idErr := strconv.Atoi(idStr)
	if idErr != nil {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, idErr.Error(), map[string]interface{}{
			"Id": conf.VARIABLE_IS_NOT_NUMERIC,
		})
		c.ServeJSON()
		return
	}

	existsParams := map[string]string{
		"bookingId":      strconv.Itoa(id),
		"adminAccountId": strconv.Itoa(userGlobal.Id),
	}
	// Check Booking exists
	listBooking, exists := models.CheckBookingExist(existsParams)
	if !exists {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, "Booking is not exists", map[string]interface{}{
			"Booking": conf.RECORD_NOT_FOUND,
		})
		c.ServeJSON()
		return
	}
	booking := listBooking[0]
	//check charge
	if listBookingCharge, err := models.IsPaidByListBookingId(strconv.Itoa(booking.Id)); (len(listBookingCharge) < 1 || err != nil) &&
		booking.PayType == conf.BOOKING_PAYMENT_ARRIVAL {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, "Booking not pay", map[string]interface{}{
			"BookingCharge": conf.BOOKING_NOT_PAID,
		})
		c.ServeJSON()
		return
	}

	if booking.Status != conf.BOOKING_ACTIVE {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, "Booking not available", map[string]interface{}{
			"BookingId": conf.BOOKING_CHECKINED,
		})
		c.ServeJSON()
		return
	}

	if booking.CheckinDate[:10] != time.Now().Format(migration.DBDateFormat)[:10] {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, "Today is not checkin day", map[string]interface{}{
			"Date": conf.NOT_CHECKIN_DAY,
		})
		c.ServeJSON()
		return
	}

	params := models.GetParamsCheckin(userGlobal.Id)
	if errUpdate := models.UpdateBookingCheckinWithParams(id, params, false); errUpdate != nil {
		c.Data["json"] = libs.ResultJson(nil, responseErrStatus, "Update booking false", map[string]interface{}{
			"Update": conf.SAVE_FAILURES,
		})
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), "Checkin Success", nil)
	c.ServeJSON()
}

// Confirm Checkout ...
// @Title Confirm Checkout
// @Description Confirm Checkout
// @Param Authorization header string true "Bearer token
// @Param	id		path 	string	true		"The booking id you want to checkin"
// @Success 200 {object} {object} Charge success
// @Failure 403 101 : Permission denied <br> 104 : Invalid Token <br> 302 : Booking not exists <br> 303 : Checkout false <br> 602 : Booking not checkin <br> 606 : Booking has check out <br> 607 : Booking has reject
// @router /confirmCheckout/:id [post]
func (c *BookingsController) ConfirmCheckout() {
	//check permission
	userGlobal := libs.GlobalUser
	if userGlobal.Type != conf.TYPE_ADMIN {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "permission denied", map[string]interface{}{
			"Booking": conf.PERMISSION_DENY,
		})
		c.ServeJSON()
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, idErr := strconv.Atoi(idStr)
	if idErr != nil {
		c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), idErr.Error(), map[string]interface{}{
			"Id": conf.VARIABLE_IS_NOT_NUMERIC,
		})
		c.ServeJSON()
		return
	}

	existsParams := map[string]string{
		"bookingId":      strconv.Itoa(id),
		"adminAccountId": strconv.Itoa(userGlobal.Id),
	}
	// Check Booking exists
	listBooking, exists := models.CheckBookingExist(existsParams)
	if !exists {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking is not exists",
			map[string]interface{}{"Booking": conf.RECORD_NOT_FOUND,
			})
		c.ServeJSON()
		return
	}

	booking := listBooking[0]
	if booking.Status < conf.BOOKING_CHECKIN {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking not checkin",
			map[string]interface{}{"Booking": conf.NOT_CHECKIN},
		)
		c.ServeJSON()
		return
	} else if booking.Status == conf.BOOKING_CHECKOUT {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking has check out",
			map[string]interface{}{"Booking": conf.BOOKING_HAS_CHECKOUT},
		)
		c.ServeJSON()
		return
	} else if booking.Status == conf.BOOKING_REJECT {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Booking has reject",
			map[string]interface{}{"Booking": conf.BOOKING_HAS_REJECT},
		)
		c.ServeJSON()
		return
	}

	params := models.GetParamsCheckout(userGlobal.Id)
	if errUpdate := models.UpdateBookingCheckinWithParams(id, params, true); errUpdate != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Update booking false",
			map[string]interface{}{"Update": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}

	c.Data["json"] = libs.ResultJson(nil, fmt.Sprint(conf.SUCCESS_STATUS), "Checkout Success", nil)
	c.ServeJSON()
}

// BookingVerifyPassport ...
// @Title BookingVerifyPassport
// @Description Verify Passport
// @Param Authorization header string true "Bearer token"
// @Param	body		body 	models.VerifyPassport	true		"body for verify passport"
// @Success 200 {string} Verify Success
// @Failure 403 104 : Invalid Token <br> 217 : Fields required <br> 302 : Customer not exists <br> 221 : Parse json false <br> 501 : Not match
// @router /checkin/verifyPassport/ [post]
func (c *BookingsController) BookingVerifyPassport() {
	detailErrorCode := make(map[string]interface{})
	responseErrCode := fmt.Sprint(conf.ERROR_STATUS)
	var ob models.VerifyPassport
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob); err != nil {
		detailErrorCode["Body"] = conf.PARSE_JSON_BODY_FALSE
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, err.Error(), detailErrorCode)
		c.ServeJSON()
		return
	}

	// Check Validate
	if b, detailErrorCode := validation.CheckValidate(ob); !b {
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "Validate false", detailErrorCode)
		c.ServeJSON()
		return
	}
	// Check Booking exists
	par := map[string]string{
		"bookingId": fmt.Sprint(ob.BookingId),
	}
	globalUser := libs.GlobalUser
	switch globalUser.Type {
	case conf.TYPE_GUEST:
		par["guestId"] = fmt.Sprint(globalUser.Id)
	case conf.TYPE_USER:
		par["userAppId"] = fmt.Sprint(globalUser.Id)
	default:
		par["adminAccountId"] = fmt.Sprint(globalUser.Id)
	}

	// Check Booking exists
	if _, exists := models.CheckBookingExist(par); !exists {
		detailErrorCode["Booking"] = conf.RECORD_NOT_FOUND
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "Booking is not exists", detailErrorCode)
		c.ServeJSON()
		return
	}

	//call API check passport
	request := httplib.Get(conf.BOOKING_VERIFY_PASSPORT_URL + "/passport")
	stringParam := "passport_url=" + ob.PassportImage
	checkFace := false
	if ob.GuestImage != "" {
		request = httplib.Get(conf.BOOKING_VERIFY_PASSPORT_URL + "/complete")
		stringParam += "<SoURL>portrait_url=" + ob.GuestImage
		checkFace = true
	}
	request.Body(stringParam)
	response, errAi := request.String()
	if errAi != nil {
		detailErrorCode["Passport"] = conf.NOT_MATCH
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "Passport not match", detailErrorCode)
		c.ServeJSON()
		return
	}
	var aiResponse models.AiResponse
	if errAi := json.Unmarshal([]byte(response), &aiResponse); errAi != nil {
		detailErrorCode["PassportInformation"] = conf.PARSE_JSON_BODY_FALSE
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, errAi.Error(), detailErrorCode)
		c.ServeJSON()
		return
	}
	if (aiResponse.Code != fmt.Sprint(conf.SUCCESS_STATUS) || !aiResponse.FaceMatching) && checkFace {
		detailErrorCode["Portrait"] = conf.NOT_MATCH
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "Portrait not match", detailErrorCode)
		c.ServeJSON()
		return
	}

	//Get guest info
	if aiResponse.InfoPassport == nil {
		detailErrorCode["PassportInformation"] = conf.VARIABLE_REQUIRED
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "Can't read passport information", detailErrorCode)
		c.ServeJSON()
		return
	}

	tmpReflectCusInfo := reflect.ValueOf(aiResponse.InfoPassport)
	var test map[string]interface{}
	if tmpReflectCusInfo.Type() != reflect.TypeOf(test) {
		detailErrorCode["PassportInformation"] = conf.FIELD_FORMAT_INVALID
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "Can't read passport information", detailErrorCode)
		c.ServeJSON()
		return
	}

	var responseGuest models.Guests
	for _, key := range tmpReflectCusInfo.MapKeys() {
		if key.String() == "full_name" {
			responseGuest.FullName = tmpReflectCusInfo.MapIndex(key).Interface().(string)
		} else if key.String() == "passport_number" {
			responseGuest.PassportNumber = tmpReflectCusInfo.MapIndex(key).Interface().(string)
		} else if key.String() == "nationality" {
			responseGuest.Nationality = tmpReflectCusInfo.MapIndex(key).Interface().(string)
		} else if key.String() == "sex" {
			responseGuest.Gender = conf.GUEST_FEMALE
			if tmpReflectCusInfo.MapIndex(key).Interface().(string) == "M" {
				responseGuest.Gender = conf.GUEST_MALE
			}
		} else if key.String() == "date_of_birth" {
			var convertBirthDayErr error
			responseGuest.BirthDay, convertBirthDayErr = ConvertDate(conf.RFC3339, fmt.Sprint(tmpReflectCusInfo.MapIndex(key).Interface().(string)))
			if convertBirthDayErr != nil {
				detailErrorCode["BirthDay"] = conf.VARIABLE_IS_NOT_DATE
			}
		} else if key.String() == "date_of_expiry" {
			var convertExpiredErr error
			responseGuest.PassportExpired, convertExpiredErr = ConvertDate(conf.RFC3339, fmt.Sprint(tmpReflectCusInfo.MapIndex(key).Interface().(string)))
			if convertExpiredErr != nil {
				detailErrorCode["PassportExpired"] = conf.VARIABLE_IS_NOT_DATE
			}
		}

	}

	c.Data["json"] = libs.ResultJson(responseGuest, fmt.Sprint(conf.SUCCESS_STATUS), "passport verify success", nil)
	if len(detailErrorCode) > 0 {
		c.Data["json"] = libs.ResultJson(nil, responseErrCode, "fail", detailErrorCode)
	}
	c.ServeJSON()
}

func ConvertDate(format string, stringDate string) (time.Time, error) {
	resultDate, err := time.Parse(format, stringDate)
	return resultDate, err
}

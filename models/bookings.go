package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"indetail/conf"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Bookings struct {
	Id               int       `orm:"column(booking_id);auto"`
	RoomId           int       `orm:"column(room_id)"`
	HotelId          int       `orm:"column(hotel_id)"`
	BookingOtaId     string    `orm:"column(booking_ota_id);size(30)"`
	PayType          int8      `orm:"column(pay_type)"`
	CheckinDate      string    `orm:"column(checkin_date);type(timestamp)"`
	CheckoutDate     string    `orm:"column(checkout_date);type(timestamp)"`
	GuestCount       int8      `orm:"column(guest_count);null"`
	ActualGuestCount int8      `orm:"column(actual_guest_count);null"`
	AdultCount       int8      `orm:"column(adult_count);null"`
	ChildrenCount    int8      `orm:"column(children_count);null"`
	Description      string    `orm:"column(description);size(255);null"`
	Note             string    `orm:"column(note)"`
	TotalAmount      string    `orm:"column(total_amount);size(255);null"`
	Status           int8      `orm:"column(status);null"`
	QrImageUrl       string    `orm:"column(qr_image_url);size(255);null"`
	DeletedAt        int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser      int       `orm:"column(created_user)"`
	UpdatedUser      int       `orm:"column(updated_user);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type AddBookingStruct struct {
	RoomId           int    `orm:"column(room_id)"`
	HotelId          int    `orm:"column(hotel_id)"`
	BookingOtaId     string `orm:"column(booking_ota_id);size(30)"`
	PayType          int8   `orm:"column(pay_type)"`
	CheckinDate      string `orm:"column(checkin_date);type(timestamp)"`
	CheckoutDate     string `orm:"column(checkout_date);type(timestamp)"`
	GuestCount       int8   `orm:"column(guest_count);null"`
	ActualGuestCount int8   `orm:"column(actual_guest_count);null"`
	AdultCount       int8   `orm:"column(adult_count);null"`
	ChildrenCount    int8   `orm:"column(children_count);null"`
	Description      string `orm:"column(description);size(255);null"`
	IsPaid           int
	Note             string `orm:"column(note)"`
	TotalAmount      string `orm:"column(total_amount);size(255);null"`
	Status           int8   `orm:"column(status);null"`
	GuestInsert      []BookingGuestInsert
}

type BookingReport struct {
	Id               int       `orm:"column(booking_id);auto"`
	RoomId           int       `orm:"column(room_id)"`
	HotelId          int       `orm:"column(hotel_id)"`
	BookingOtaId     string    `orm:"column(booking_ota_id);size(30)"`
	PayType          int8      `orm:"column(pay_type)"`
	CheckinDate      string    `orm:"column(checkin_date);type(timestamp)"`
	CheckoutDate     string    `orm:"column(checkout_date);type(timestamp)"`
	GuestCount       int8      `orm:"column(guest_count);null"`
	ActualGuestCount int8      `orm:"column(actual_guest_count);null"`
	AdultCount       int8      `orm:"column(adult_count);null"`
	ChildrenCount    int8      `orm:"column(children_count);null"`
	Description      string    `orm:"column(description);size(255);null"`
	Note             string    `orm:"column(note)"`
	TotalAmount      string    `orm:"column(total_amount);size(255);null"`
	Status           int8      `orm:"column(status);null"`
	DeletedAt        int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser      int       `orm:"column(created_user)"`
	UpdatedUser      int       `orm:"column(updated_user);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type BookingWithGuest struct {
	Id               int    `orm:"column(booking_id);auto"`
	RoomId           int    `orm:"column(room_id)"`
	HotelId          int    `orm:"column(hotel_id)"`
	BookingOtaId     string `orm:"column(booking_ota_id);size(30)"`
	PayType          int8   `orm:"column(pay_type)"`
	CheckinDate      string `orm:"column(checkin_date);type(timestamp)"`
	CheckoutDate     string `orm:"column(checkout_date);type(timestamp)"`
	GuestCount       int8   `orm:"column(guest_count);null"`
	ActualGuestCount int8   `orm:"column(actual_guest_count);null"`
	AdultCount       int8   `orm:"column(adult_count);null"`
	ChildrenCount    int8   `orm:"column(children_count);null"`
	Description      string `orm:"column(description);size(255);null"`
	Note             string `orm:"column(note)"`
	TotalAmount      string `orm:"column(total_amount);size(255);null"`
	Status           int8   `orm:"column(status);null"`
	Guest            []Guests
}

type BookingDetail struct {
	Id               int    `orm:"column(booking_id);auto"`
	RoomId           int    `orm:"column(room_id)"`
	HotelId          int    `orm:"column(hotel_id)"`
	BookingOtaId     string `orm:"column(booking_ota_id);size(30)"`
	PayType          int8   `orm:"column(pay_type)"`
	CheckinDate      string `orm:"column(checkin_date);type(timestamp)"`
	CheckoutDate     string `orm:"column(checkout_date);type(timestamp)"`
	GuestCount       int8   `orm:"column(guest_count);null"`
	ActualGuestCount int8   `orm:"column(actual_guest_count);null"`
	AdultCount       int8   `orm:"column(adult_count);null"`
	ChildrenCount    int8   `orm:"column(children_count);null"`
	Description      string `orm:"column(description);size(255);null"`
	Note             string `orm:"column(note)"`
	TotalAmount      string `orm:"column(total_amount);size(255);null"`
	Status           int8   `orm:"column(status);null"`
	QrImageUrl       string `orm:"column(qr_image_url);size(255);null"`
	Guest            []Guests
	Room             Rooms
}

type BookingUpdateFields struct {
	RoomId           int    `json:"room_id, omitempty"`
	BookingStatus    int    `json:"status" valid:"Max(4)"`
	ActualGuestCount int8   `json:"actual_guest_count, omitempty"`
	CheckinDate      string `json:"checkin_date, omitempty"`
	CheckoutDate     string `json:"checkout_date, omitempty"`
	Description      string `json:"description, omitempty"`
	Note             string `json:"note, omitempty"`
	PayType          int8   `json:"pay_type, omitempty"`
	GuestCount       int8   `json:"guest_count, omitempty"`
	AdultCount       int8   `json:"adult_count, omitempty"`
	ChildrenCount    int8   `json:"children_count, omitempty"`
	TotalAmount      string `json:"total_amount, omitempty"`
	GuestUpdate      []BookingGuestInsert
}

type BookingData struct {
	Id               int       `orm:"column(booking_id)"`
	RoomId           int       `orm:"column(room_id)"`
	HotelId          int       `orm:"column(hotel_id)"`
	BookingOtaId     string    `orm:"column(booking_ota_id);size(30)"`
	PayType          int8      `orm:"column(pay_type)"`
	CheckinDate      string    `orm:"column(checkin_date);type(timestamp)"`
	CheckoutDate     string    `orm:"column(checkout_date);type(timestamp)"`
	GuestCount       int8      `orm:"column(guest_count);null"`
	ActualGuestCount int8      `orm:"column(actual_guest_count);null"`
	AdultCount       int8      `orm:"column(adult_count);null"`
	ChildrenCount    int8      `orm:"column(children_count);null"`
	Description      string    `orm:"column(description);size(255);null"`
	Note             string    `orm:"column(note)"`
	TotalAmount      string    `orm:"column(total_amount);size(255);null"`
	QrImageUrl       string    `orm:"column(qr_image_url);size(255);null"`
	Status           int8      `orm:"column(status);null"`
	DeletedAt        int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser      int       `orm:"column(created_user)"`
	UpdatedUser      int       `orm:"column(updated_user);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
	ListGuest        []GuestWithIsMain
	Room             Rooms
	IsPaid           int
}

type VerifyPassport struct {
	BookingId     int    `json:"bookingId" valid:"Required"`
	PassportImage string `json:"passportImage" valid:"Required"`
	GuestImage    string `json:"guestImage"`
}

type AiResponse struct {
	Code         string      `json:"code"`
	FaceMatching bool        `json:"face_matching"`
	InfoPassport interface{} `json:"info_passport"`
	Message      string      `json:"message"`
}

func (t *Bookings) TableName() string {
	return "bookings"
}

func init() {
	orm.RegisterModel(new(Bookings))
}

// AddBookings insert a new Bookings into database and returns
// last inserted Id on success.
func AddBookings(m *Bookings) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOneBookings retrieves Bookings by Id. Returns error if
// Id doesn't exist
func GetOneBookings(id int) (booking *Bookings, err error) {
	o := orm.NewOrm()
	booking = &Bookings{Id: id, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(booking); err == nil {
		return booking, nil
	}
	return nil, err
}

// GetBookingsById retrieves Bookings by Id. Returns error if
// Id doesn't exist
func GetBookingsById(id int) (booking *BookingData, err error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("bookings.*").
		From("bookings").
		Where("booking_id = " + fmt.Sprint(id)).
		And("deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).String()

	if err := o.Raw(sql).QueryRow(&booking); err != nil {
		return nil, err
	}
	return booking, err
}

// GetBookingsByConditions retrieves Bookings by conditions. Returns error if
// Id doesn't exist
func GetBookingsByConditions(params map[string]string) ([]Bookings, bool) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("bookings.*").
		From("bookings").
		Where("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.booking_id = " + params["bookingId"]).
		And("bookings.status = " + fmt.Sprint(conf.BOOKING_NEW))

	query := sql.String()
	var bookings []Bookings
	num, err := o.Raw(query).QueryRows(&bookings)
	if num > 0 && err == nil {
		return bookings, true
	}
	return bookings, false
}

// GetAllBookings retrieves all Bookings matches certain condition. Returns empty list if
// no records exist
func GetAllBookings(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Bookings))
	//rebuild query
	qs = RebuildConditions(qs, query)
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Bookings
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// GetAllBookingsByCondition retrieves all Bookings matches certain condition. Returns empty list if
// no records exist
func GetAllBookingsByCondition(query map[string]string, fields string, sortby []string, order []string,
	offset int64, limit int64) (listBooking []BookingData, msgErr string) {
	o := orm.NewOrm()

	//build sql
	sql, sqlWhere := BuildBookingListSql(fields, query)

	// order by:
	sortFields, sortErr := MakeOrderForSqlQuery(sortby, order)
	if sortErr != nil {
		return nil, sortErr.Error()
	}
	sqlWhereString := sqlWhere.String()
	if sortFields != "" {
		sqlWhereString += " Order By " + sortFields
	}

	if limit > 0 {
		sqlWhereString += " LIMIT " + fmt.Sprint(limit)
	}

	if offset > 0 {
		sqlWhereString += " OFFSET " + fmt.Sprint(offset)
	}
	sqlQuery := sql.String() + " " + sqlWhereString
	_, err := o.Raw(sqlQuery).QueryRows(&listBooking)
	if err != nil {
		return nil, "Get list failure"
	}
	return listBooking, ""
}

func GetExtraInfoForBooking(listBooking []BookingData, query map[string]string) []BookingData {
	var listBoongId string = ""
	for _, booking := range listBooking {
		listBoongId += fmt.Sprint(booking.Id) + ","
	}
	listBoongId = strings.TrimRight(listBoongId, ",")

	var guestWithBookingId []GuestWithBookingId
	var guestId int = 0
	if query["GuestId"] != "" {
		guestId, _ = strconv.Atoi(query["GuestId"])
	}
	guestWithBookingId, _ = GetGuestByListBookingId(listBoongId, guestId)

	var tmpListGuest = make(map[int][]GuestWithIsMain)
	for _, guest := range guestWithBookingId {
		var tmpGuest GuestWithIsMain
		tmpGuest.Id = guest.Id
		tmpGuest.DeletedAt = guest.DeletedAt
		tmpGuest.Password = ""
		tmpGuest.IsMainGuest = guest.IsMainGuest
		tmpGuest.PassportNumber = guest.PassportNumber
		tmpGuest.Email = guest.Email
		tmpGuest.Address = guest.Address
		tmpGuest.BirthDay = guest.BirthDay
		tmpGuest.CreatedAt = guest.CreatedAt
		tmpGuest.CreatedUser = guest.CreatedUser
		tmpGuest.Files = guest.Files
		tmpGuest.FullName = guest.FullName
		tmpGuest.Gender = guest.Gender
		tmpGuest.Occupation = guest.Occupation
		tmpGuest.Phone = guest.Phone
		tmpGuest.Status = guest.Status
		tmpGuest.UpdatedAt = guest.UpdatedAt
		tmpGuest.UpdatedUser = guest.UpdatedUser
		tmpListGuest[guest.BookingId] = append(tmpListGuest[guest.BookingId], tmpGuest)
	}

	//get room info
	listRoom, _ := GetRoomByListBookingId(listBoongId)
	var tmpListRoom = make(map[int]Rooms)
	for _, room := range listRoom {
		var tmpRoom Rooms
		tmpRoom.Id = room.Id
		tmpRoom.DeletedAt = room.DeletedAt
		tmpRoom.HotelId = room.HotelId
		tmpRoom.RoomName = room.RoomName
		tmpRoom.RoomOtaId = room.RoomOtaId
		tmpRoom.SmartLockId = room.SmartLockId
		tmpRoom.CreatedAt = room.CreatedAt
		tmpRoom.CreatedUser = room.CreatedUser
		tmpRoom.UpdatedAt = room.UpdatedAt
		tmpRoom.UpdatedUser = room.UpdatedUser
		tmpListRoom[room.BookingId] = tmpRoom
	}

	//calculate booking charge
	listBookingCharge, _ := IsPaidByListBookingId(listBoongId)
	var tmpListBookingCharge = make(map[int]BookingCharges)
	for _, bookingCharge := range listBookingCharge {
		tmpListBookingCharge[bookingCharge.BookingId] = bookingCharge
	}

	for index, booking := range listBooking {
		//guest
		if val, ok := tmpListGuest[booking.Id]; ok {
			listBooking[index].ListGuest = val
		}
		if val, ok := tmpListGuest[0]; ok {
			listBooking[index].ListGuest = val
		}
		//room
		if val, ok := tmpListRoom[booking.Id]; ok {
			listBooking[index].Room = val
		}
		//charge
		if _, ok := tmpListBookingCharge[booking.Id]; ok {
			listBooking[index].IsPaid = 1
		} else {
			listBooking[index].IsPaid = 0
		}

	}

	return listBooking

}

// GetListBookings retrieves all Bookings matches certain condition. Returns empty list if
// no records exist
func GetListBookingReport(query map[string]string, bookingFields string) (listBooking []BookingReport, err error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select(bookingFields).From("bookings").
		InnerJoin("customers").On("bookings.customer_id = customers.customer_id").
		InnerJoin("rooms").On("bookings.room_id = rooms.room_id").String()
	qdWhere, _ := orm.NewQueryBuilder("mysql")
	sqlWhere := qdWhere.Where("customers.booking_count > 1").
		And("customers.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.booking_id != customers.first_booking_id")
	if query["HotelId"] != "" {
		qdWhere = qdWhere.And("hotel_id = " + query["HotelId"])
	}
	if query["Status"] != "" {
		qdWhere = qdWhere.And("bookings.status IN (" + query["Status"] + ")")
	}
	sqlWhere = sqlWhere.And("bookings.created_at >= '" + query["CreatedAt.gte"] + "'").
		And("bookings.created_at < '" + query["CreatedAt.lt"] + "'")

	if query["OrderBy"] != "" {
		sqlWhere.OrderBy(query["OrderBy"])
	}
	sql = sql + " " + sqlWhere.String()
	num, err := o.Raw(sql).QueryRows(&listBooking)
	if num > 0 && err == nil {
		return listBooking, err
	}
	return nil, err
}

// UpdateBookings updates Bookings by Id and returns error if
// the record to be updated doesn't exist
func UpdateBookingsById(o orm.Ormer, m *Bookings) (err error) {
	v := Bookings{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBookings deletes Bookings by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBookings(id int) (err error) {
	o := orm.NewOrm()
	v := Bookings{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Bookings{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// Check booking exists
// @param bookingId int true "id of booking"
// @return bool
func CheckBookingExist(params map[string]string) ([]Bookings, bool) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	qd.Select("bookings.*").
		From("bookings")
	qdWhere, _ := orm.NewQueryBuilder("mysql")
	qdWhere.Where("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.booking_id = " + params["bookingId"])

	if params["guestId"] != "" {
		qd.InnerJoin("booking_guests").On("bookings.booking_id = booking_guests.booking_id")
		qd.InnerJoin("guests").On("booking_guests.guest_id = guests.guest_id")
		qdWhere.And("guests.guest_id = " + params["guestId"]).And("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED))
	}
	if params["userAppId"] != "" {
		qd.InnerJoin("user_app_accounts").On("bookings.booking_id = user_app_accounts.booking_id")
		qdWhere.And("user_app_accounts.user_app_account_id = " + params["userAppId"]).And("user_app_accounts.deleted_at = " + fmt.Sprint(conf.NOT_DELETED))
	}
	if params["adminAccountId"] != "" {
		qd.InnerJoin("hotels").On("bookings.hotel_id = hotels.hotel_id")
		qd.InnerJoin("admin_accounts").On("hotels.hotel_id = admin_accounts.hotel_id")
		qdWhere.And("admin_accounts.admin_account_id = " + params["adminAccountId"]).And("admin_accounts.deleted_at = " + fmt.Sprint(conf.NOT_DELETED))
	}

	query := qd.String() + " " + qdWhere.String()
	var bookings []Bookings
	num, err := o.Raw(query).QueryRows(&bookings)
	if num > 0 && err == nil {
		return bookings, true
	}
	return nil, false
}

// Get booking exists
// @param bookingId int true "id of booking"
// @return bool
func GetBookingExist(params map[string]string) (*Bookings, bool) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("bookings.*").
		From("bookings").
		InnerJoin("hotels").On("bookings.hotel_id = hotels.hotel_id").
		Where("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("hotels.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.booking_id = " + params["bookingId"])

	query := sql.String()
	var bookings []Bookings
	num, err := o.Raw(query).QueryRows(&bookings)
	if num > 0 && err == nil {
		return &bookings[0], true
	}
	return nil, false
}

// Update Booking With Params
// @Param id int "id of booking"
// @Param param BookingUpdateFields "body of booking update"
// @return err error
func UpdateBookingWithParams(id int, param orm.Params) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	_, errUpdate := o.QueryTable(new(Bookings)).Filter("Id", id).Update(param)
	//end transaction
	if errUpdate != nil {
		err = o.Rollback()
		return errUpdate
	}
	err = o.Commit()
	return err
}

func UpdateBookingWithParamsTransaction(id int, param orm.Params, o orm.Ormer) (err error) {
	_, errUpdate := o.QueryTable(new(Bookings)).Filter("Id", id).Update(param)
	return errUpdate
}

// Update Booking Checkin With Params
// @Param id int "id of booking"
// @Param param BookingUpdateFields "body of booking update"
// @Param user *AdminAccounts
// @return err error
func UpdateBookingCheckinWithParams(id int, params orm.Params, isCheckout bool) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	_, errUpdate := o.QueryTable(new(Bookings)).Filter("Id", id).Update(params)
	//end transaction
	if errUpdate != nil {
		err = o.Rollback()
		return errUpdate
	}

	if !isCheckout {
		err = o.Commit()
		return err
	}

	userAppParams := orm.Params{
		"status": conf.USER_APP_STATUS_INACTIVE,
	}
	if userAppErr := UpdateUserAppAccountByTransaction(id, userAppParams, o); userAppErr != nil {
		err = o.Rollback()
		return userAppErr
	}

	err = o.Commit()
	return err
}

// Get Params Update Booking
func GetParamsUpdate(ob BookingUpdateFields, user GlobalUsers) orm.Params {
	var resultParams = make(map[string]interface{})

	if ob.RoomId > 0 {
		resultParams["room_id"] = ob.RoomId
	}
	if ob.BookingStatus > 0 {
		resultParams["status"] = ob.BookingStatus
	}
	if ob.CheckinDate != "" {
		resultParams["checkin_date"] = ob.CheckinDate
	}
	if ob.CheckoutDate != "" {
		resultParams["checkout_date"] = ob.CheckoutDate
	}
	if ob.ActualGuestCount > 0 {
		resultParams["actual_guest_count"] = ob.ActualGuestCount
	}
	if ob.Description != "" {
		resultParams["description"] = ob.Description
	}
	if ob.Note != "" {
		resultParams["note"] = ob.Note
	}
	if ob.PayType > 0 {
		resultParams["pay_type"] = ob.PayType
	}
	if ob.GuestCount > 0 {
		resultParams["guest_count"] = ob.GuestCount
	}
	if ob.AdultCount > 0 {
		resultParams["adult_count"] = ob.AdultCount
	}
	if ob.ChildrenCount > 0 {
		resultParams["children_count"] = ob.ChildrenCount
	}
	if ob.TotalAmount != "" {
		resultParams["total_amount"] = ob.TotalAmount
	}

	resultParams["updated_user"] = user.Id
	resultParams["updated_at"] = time.Now()
	return resultParams
}

// Get Params Checkin Booking
func GetParamsCheckin(userId int) orm.Params {
	return orm.Params{
		"status":       conf.BOOKING_CHECKIN,
		"updated_user": strconv.Itoa(userId),
	}
}

// Get Params Checkout Booking
func GetParamsCheckout(userId int) orm.Params {
	return orm.Params{
		"status":       conf.BOOKING_CHECKOUT,
		"updated_user": strconv.Itoa(userId),
	}
}

// Get Params Reject Booking
func GetParamsReject() orm.Params {
	return orm.Params{
		"status": conf.BOOKING_REJECT,
	}
}

// count booking
func CountBooking(fields string, query map[string]string) (total int64, err error) {
	o := orm.NewOrm()
	//build sql
	sql, sqlWhere := BuildBookingListSql(fields, query)
	sqlQuery := sql.String() + " " + sqlWhere.String()
	var bookings []Bookings
	num, err := o.Raw(sqlQuery).QueryRows(&bookings)
	if num > 0 && err == nil {
		return num, err
	}
	return -9, err

}

func BuildBookingListSql(fields string, query map[string]string) (sql orm.QueryBuilder, sqlWhere orm.QueryBuilder) {
	qd, _ := orm.NewQueryBuilder("mysql")
	sql = qd.Select(fields).From("bookings")

	qdWhere, _ := orm.NewQueryBuilder("mysql")
	sqlWhere = qdWhere.Where("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED))

	if query["HotelId"] != "" {
		sqlWhere = qdWhere.And("bookings.hotel_id = " + query["HotelId"])
	}

	if query["GuestId"] != "" {
		sql = qd.InnerJoin("booking_guests").On("bookings.booking_id = booking_guests.booking_id")
		sqlWhere = sqlWhere.And("booking_guests.guest_id = " + query["GuestId"])
	}

	if query["UserAppId"] != "" {
		sql = qd.InnerJoin("user_app_accounts").On("bookings.booking_id = user_app_accounts.booking_id")
		sqlWhere = sqlWhere.And("user_app_accounts.user_app_account_id = " + query["UserAppId"])
	}

	if query["CheckinDate.gte"] != "" {
		sqlWhere = qdWhere.And("bookings.checkin_date >= '" + query["CheckinDate.gte"] + "'")
	}
	if query["CheckinDate.lte"] != "" {
		sqlWhere = qdWhere.And("bookings.checkin_date <= '" + query["CheckinDate.lte"] + "'")
	}
	if query["CheckoutDate.gte"] != "" {
		sqlWhere = qdWhere.And("bookings.checkout_date >= '" + query["CheckoutDate.gte"] + "'")
	}
	if query["CheckoutDate.lte"] != "" {
		sqlWhere = qdWhere.And("bookings.checkout_date <= '" + query["CheckoutDate.lte"] + "'")
	}
	return sql, sqlWhere
}

// Booking Encode Base64
func BookingEncodeBase64(str string) string {
	str = str + conf.DOT + conf.SALT
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Booking Decode Base64
func BookingDecodeBase64(str string) string {
	decode, _ := base64.StdEncoding.DecodeString(str)
	arr := strings.Split(string(decode), conf.DOT)
	if len(arr) > 1 {
		return arr[0]
	}
	return ""
}

// Add Hotel insert a new Hotel into database and returns
// last inserted Id on success.
func AddBookingWithTransaction(addBooking *AddBookingStruct, globalUser GlobalUsers) (detailErrorCode map[string]interface{}, errorMsg error, booking Bookings) {
	o := orm.NewOrm()
	//begin transaction
	var err error
	err = o.Begin()
	detailErrorCode = make(map[string]interface{})
	//insert booking
	booking = RebuildBookingStruct(addBooking, globalUser)
	id, errInsert := o.Insert(&booking)
	if errInsert != nil {
		detailErrorCode["Booking"] = conf.SAVE_FAILURES
		err = o.Rollback()
		return detailErrorCode, errInsert, booking
	}

	//QR code
	params := make(map[string]interface{})
	errGenQr, path := GenQrCodeForBooking(booking.Id)
	if errGenQr != nil {
		detailErrorCode["QR"] = conf.CREATE_FILE_FAILURES
		err = o.Rollback()
		return detailErrorCode, errGenQr, booking
	}
	params["qr_image_url"] = beego.AppConfig.String("baseServer") + path
	if errInsertQr := UpdateBookingWithParamsTransaction(int(id), params, o); errInsertQr != nil {
		detailErrorCode["QR"] = conf.SAVE_FAILURES
		err = o.Rollback()
		return detailErrorCode, errInsertQr, booking
	}

	//guest
	bookingGuest := RebuildBookingGuest(*addBooking, int(id))
	fmt.Println(bookingGuest)
	_, errBookingGuest := AddBookingGuestTransaction(&bookingGuest, o)
	//end transaction
	if errBookingGuest != nil {
		err = o.Rollback()
		detailErrorCode["BookingGuests"] = conf.SAVE_FAILURES
		return detailErrorCode, errBookingGuest, booking
	}
	err = o.Commit()
	return detailErrorCode, err, booking
}

func RebuildBookingStruct(addBooking *AddBookingStruct, globalUser GlobalUsers) (booking Bookings) {

	booking.RoomId = addBooking.RoomId
	booking.HotelId = addBooking.HotelId
	booking.BookingOtaId = addBooking.BookingOtaId
	booking.PayType = addBooking.PayType
	booking.CheckinDate = addBooking.CheckinDate
	booking.CheckoutDate = addBooking.CheckoutDate
	booking.GuestCount = addBooking.GuestCount
	booking.AdultCount = addBooking.AdultCount
	booking.ChildrenCount = addBooking.ChildrenCount
	booking.Description = addBooking.Description
	booking.Note = addBooking.Note
	booking.TotalAmount = addBooking.TotalAmount
	booking.Status = addBooking.Status
	booking.DeletedAt = conf.NOT_DELETED
	booking.CreatedUser = globalUser.Id
	return booking
}

func RebuildBookingGuest(addBooking AddBookingStruct, bookingId int) (bookingGuest []BookingGuests) {
	for _, guest := range addBooking.GuestInsert {
		var tmp BookingGuests
		tmp.GuestId = int(guest.Id)
		tmp.BookingId = bookingId
		tmp.IsMainGuest = 1
		bookingGuest = append(bookingGuest, tmp)
	}
	return bookingGuest
}

func UpdateBookingTransaction(oldBooking Bookings, params orm.Params, listProcessGuest map[string][]int) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	_, errUpdate := o.QueryTable(new(Bookings)).Filter("Id", oldBooking.Id).Update(params)
	//end transaction
	if errUpdate != nil {
		err = o.Rollback()
		return errUpdate
	}

	if len(listProcessGuest["remove"]) > 0 {
		stringRemoveGuestId := ""
		for _, guestId := range listProcessGuest["remove"] {
			stringRemoveGuestId += fmt.Sprint(guestId) + ","
		}
		stringRemoveGuestId = strings.TrimRight(stringRemoveGuestId, ",")

		if errRemove := UpdateBookingGuestTransaction(oldBooking.Id, stringRemoveGuestId, o); errRemove != nil {
			err = o.Rollback()
			return errRemove
		}
	}

	if len(listProcessGuest["insert"]) < 1 {
		err = o.Commit()
		return err
	}

	var bookingGuest []BookingGuests
	for _, guestId := range listProcessGuest["insert"] {
		var tmp BookingGuests
		tmp.GuestId = guestId
		tmp.BookingId = oldBooking.Id
		tmp.IsMainGuest = conf.GUEST_IS_NOT_MAIN
		bookingGuest = append(bookingGuest, tmp)
	}
	if _, errBookingGuest := AddBookingGuestTransaction(&bookingGuest, o); errBookingGuest != nil {
		err = o.Rollback()
		return errBookingGuest
	}

	err = o.Commit()
	return err
}

// SaveToFile saves uploaded file to new path.
// it only operates the first one of mutil-upload form file field.
func GenQrCodeForBooking(bookingId int) (err error, path string) {
	imageString := fmt.Sprint(bookingId) + conf.DOT + conf.SALT

	qrCode, encodeErr := qr.Encode(imageString, qr.L, qr.Auto)
	if encodeErr != nil {
		return encodeErr, ""
	}
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	path = "/storage/qr_code/qr_" + fmt.Sprint(bookingId) + conf.QR_CODE_EXTENSION
	f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err, ""
	}
	if errSave := png.Encode(f, qrCode); errSave != nil {
		return errSave, ""
	}
	defer f.Close()
	return nil, path
}

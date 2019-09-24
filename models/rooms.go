package models

import (
	"fmt"
	"indetail/conf"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type Rooms struct {
	Id          int       `orm:"column(room_id);auto"`
	HotelId     int       `valid:"Required";orm:"column(hotel_id)"`
	SmartLockId int       `valid:"Required";orm:"column(smart_lock_id)"`
	RoomOtaId   int       `valid:"Required";orm:"column(room_ota_id)"`
	RoomName    string    `valid:"Required";orm:"column(room_name);size(255);null"`
	DeletedAt   int8      `orm:"column(deleted_at);null"`
	CreatedUser int       `orm:"column(created_user)"`
	UpdatedUser int       `orm:"column(updated_user);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type InsertRooms struct {
	HotelId     int    `valid:"Required";orm:"column(hotel_id)"`
	SmartLockId int    `valid:"Required";orm:"column(smart_lock_id)"`
	RoomOtaId   int    `valid:"Required";orm:"column(room_ota_id)"`
	RoomName    string `valid:"Required";orm:"column(room_name);size(255);null"`
}

type ResultRoomData struct {
	Id          int    `orm:"column(room_id)"`
	HotelId     int    `valid:"Required";orm:"column(hotel_id)"`
	SmartLockId int    `valid:"Required";orm:"column(smart_lock_id)"`
	RoomOtaId   int    `valid:"Required";orm:"column(room_ota_id)"`
	RoomName    string `valid:"Required";orm:"column(room_name);size(255);null"`
	DeletedAt   int8
	CreatedUser int
	UpdatedUser int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoomWithBookingId struct {
	Id          int    `orm:"column(room_id)"`
	BookingId   int    `orm:"column(booking_id);null" json:"booking_id"`
	HotelId     int    `valid:"Required";orm:"column(hotel_id)"`
	SmartLockId int    `valid:"Required";orm:"column(smart_lock_id)"`
	RoomOtaId   int    `valid:"Required";orm:"column(room_ota_id)"`
	RoomName    string `valid:"Required";orm:"column(room_name);size(255);null"`
	DeletedAt   int8
	CreatedUser int
	UpdatedUser int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MaxMinPrice struct {
	HotelId  int    `orm:"column(hotel_id)"`
	MaxPrice string `orm:"column(max)"`
	MinPrice string `orm:"column(min)"`
}

func (t *Rooms) TableName() string {
	return "rooms"
}

func init() {
	orm.RegisterModel(new(Rooms))
}

// AddRooms insert a new Rooms into database and returns
// last inserted Id on success.
func AddRooms(m *Rooms) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoomsById retrieves Rooms by Id. Returns error if
// Id doesn't exists
func GetRoomsById(id int) (v *Rooms, err error) {
	o := orm.NewOrm()
	v = &Rooms{Id: id, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "Id", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetRoomsCondition retrieves Rooms by condition. Returns error if
// no records exists
func GetRoomsCondition(query map[string]string) (result Rooms, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rooms))
	qs = RebuildConditions(qs, query)
	var room []Rooms
	var roomId int64
	if roomId, err = qs.All(&room); err == nil && roomId > 0 {
		result = room[0]
		return result, nil
	}
	return result, err
}

// GetRoomByListBookingId retrieves Guests by Id. Returns error if
// Id doesn't exists
func GetRoomByListBookingId(listBookingId string) (roomWithBookingId []RoomWithBookingId, err error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("bookings.booking_id, rooms.*").
		From("rooms").
		InnerJoin("bookings").On("rooms.room_id = bookings.room_id").
		Where("bookings.booking_id ").In(listBookingId).
		And("rooms.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		GroupBy("bookings.booking_id, rooms.room_id").String()
	if _, err := o.Raw(sql).QueryRows(&roomWithBookingId); err != nil {
		return nil, err
	}
	return roomWithBookingId, nil
}

// GetRoomDataById retrieves Rooms by Id. Returns error if
// Id doesn't exists
func GetRoomDataById(id int) (ResultRoomData, bool) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	sql := qb.Select("rooms.*").From("rooms").
		Where("room_id = " + fmt.Sprint(id)).
		And("deleted_at = " + strconv.Itoa(conf.NOT_DELETED))

	var tmpRooms []ResultRoomData
	num, tmpErr := o.Raw(sql.String()).QueryRows(&tmpRooms)
	if num > 0 && tmpErr == nil {
		return tmpRooms[0], false
	}
	return ResultRoomData{}, true
}

// GetAllRooms retrieves all Rooms matches certain condition. Returns empty list if
// no records exists
func GetAllRooms(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rooms))
	// query k=v
	qs = RebuildConditions(qs, query)
	// order by:
	qs, errorMq := MakeOrderForQuery(qs, sortby, order)

	if errorMq != nil {
		return ml, errorMq
	}

	var l []Rooms
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		var result1 []interface{}
		for _, v := range l {
			result1 = append(result1, v)
		}
		//filter data by using array field
		ml, err := FilterResultByField(fields, result1)
		return ml, err
	}
	return nil, err
}

// UpdateRooms updates Rooms by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoomsById(m *Rooms) (err error) {
	o := orm.NewOrm()
	v := Rooms{Id: m.Id, DeletedAt: conf.NOT_DELETED}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRooms deletes Rooms by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRooms(id int) (err error) {
	o := orm.NewOrm()
	v := Rooms{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Rooms{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// Check Room Exists
func CheckRoomExists(user *AdminAccounts, roomId int) bool {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("rooms.room_id").
		From("rooms").
		InnerJoin("hotels").On("hotels.hotel_id = rooms.hotel_id").
		InnerJoin("hotel_administrators").On("hotel_administrators.hotel_id = hotels.hotel_id").
		InnerJoin("admin_accounts").On("admin_accounts.admin_account_id = hotel_administrators.admin_account_id").
		Where("rooms.room_id = " + fmt.Sprint(roomId)).
		And("rooms.deleted_at = 0").
		And("admin_accounts.admin_account_id = " + fmt.Sprint(user.Id)).
		String()
	var room []Rooms
	num, err := o.Raw(sql).QueryRows(&room)
	if num > 0 && err == nil {
		return true
	}
	return false
}

// Destroy room with soft delete
func DestroyRoom(roomId int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(new(Rooms)).Filter("Id", roomId).Update(orm.Params{
		"deleted_at": 1,
	})
	return err
}

//check room available
func IsRoomAvailable(conditions map[string]string) (bool, Rooms) {
	var room = Rooms{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rooms))
	qs = RebuildConditions(qs, conditions)
	if qs.Exist() {
		qs.One(&room)
		return true, room
	}
	return false, room
}

// count hotel records
func CountRooms(query map[string]string) (totalRecord int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Rooms))
	// query k=v
	qs = RebuildConditions(qs, query)

	if totalRecord, err = qs.Count(); err == nil {
		return totalRecord, err
	} else {
		return 0, err
	}
}

// count rooms records
func CountRoomsAvailable(query map[string]string) (totalRecord int64, err error) {
	o := orm.NewOrm()
	sql := BuildRoomAvailableQuery(query, 0, -1)
	var rooms []Rooms
	num, err := o.Raw(sql).QueryRows(&rooms)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func BuildRoomAvailableQuery(query map[string]string, offset int, limit int) (sql string) {
	qbBooking, _ := orm.NewQueryBuilder("mysql")
	sqlBooking := qbBooking.Select("bookings.room_id").
		From("bookings").
		Where("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.hotel_id = " + query["HotelId"]).
		And("NOT bookings.status IN (" + fmt.Sprint(conf.BOOKING_REJECT) + "," + fmt.Sprint(conf.BOOKING_CHECKOUT) + ")").
		And("(bookings.checkin_date >= '" + query["checkinDate"] + "' AND bookings.checkin_date <= '" + query["checkoutDate"] + "')").
		Or("(bookings.checkout_date >= '" + query["checkinDate"] + "' AND bookings.checkout_date < '" + query["checkoutDate"] + "')").
		GroupBy("bookings.room_id").String()

	qb, _ := orm.NewQueryBuilder("mysql")
	sqlTmp := qb.Select("rooms.*").
		From("rooms").
		Where("rooms.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("rooms.hotel_id = " + query["HotelId"]).
		And("NOT rooms.room_id IN ( " + sqlBooking + " )")

	if query["RoomId"] != "" {
		sqlTmp = qb.And("rooms.room_id = " + query["RoomId"] + "")
	}

	if limit >= 0 {
		sqlTmp = qb.Limit(limit)
	}

	return sqlTmp.String()
}

// get rooms records
func GetRoomsAvailable(query map[string]string, offset int, limit int) (rooms []Rooms, err error) {
	o := orm.NewOrm()
	sql := BuildRoomAvailableQuery(query, offset, limit)
	_, err = o.Raw(sql).QueryRows(&rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
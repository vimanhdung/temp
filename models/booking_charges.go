package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"time"
)

type BookingCharges struct {
	Id                   int       `orm:"column(booking_charge_id);auto"`
	BookingId            int       `orm:"column(booking_id)"`
	TotalAmount          string    `orm:"column(total_amount)"`
	Discount             string    `orm:"column(discount)"`
	Amount               string    `orm:"column(amount)"`
	ChargeDatetime       time.Time `orm:"column(charge_datetime);type(timestamp)"`
	Type                 int       `orm:"column(type)"`
	Note                 string    `orm:"column(note)"`
	ListBookingServiceId string    `orm:"column(list_booking_service_id)"`
	Status               int8      `orm:"column(status);null"`
	DeletedAt            int8      `orm:"column(deleted_at);null"`
	CreatedAt            time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt            time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
	CreatedUser          int       `orm:"column(created_user)"`
	UpdatedUser          int       `orm:"column(updated_user)"`
}

type ChargeInfo struct {
	Amount float64
	Status string
	Detail string
}

func (t *BookingCharges) TableName() string {
	return "booking_charges"
}

func init() {
	orm.RegisterModel(new(BookingCharges))
}

// AddBookingCharge insert a new AddBookingCharge into database and returns
// last inserted Id on success.
func AddBookingCharge(m *BookingCharges) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBookingChargeById retrieves GetBookingCharge by Id. Returns error if
// Id doesn't exists
func GetBookingChargeById(id int) (v *BookingCharges, err error) {
	o := orm.NewOrm()
	v = &BookingCharges{Id: id, DeletedAt: 0, Status: conf.BOOKING_CHARGE_STATUS_NEW}
	if err = o.Read(v, "Id", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateBookingCharge updates BookingCharge by Id and returns error if
// the record to be updated doesn't exists
func UpdateBookingChargeById(m *BookingCharges, user *AdminAccounts) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	v := BookingCharges{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}

	//end transaction
	if err != nil {
		err = o.Rollback()
		return err
	}
	//update booking
	_, errBooking := GetOneBookings(m.BookingId)
	if errBooking != nil {
		err = o.Rollback()
		return errBooking
	}
	/*if booking.LastRoomCharge == "" {
		var param orm.Params
		param["last_room_charge"] = m.UpdatedAt.Format(conf.RFC_DATE_TIME)
		if updateBookingErr := UpdateBookingWithParams(m.BookingId, param, user); updateBookingErr != nil {
			err = o.Rollback()
			return errBooking
		}
	}*/

	err = o.Commit()
	return
}

func IsPaidByListBookingId(listBookingId string) (listBookingCharge []BookingCharges, err error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("booking_charges.*").
		From("booking_charges").
		Where("booking_guests.booking_id ").In(listBookingId).
		And("booking_guests.status = " + fmt.Sprint(conf.BOOKING_CHARGE_STATUS_CHARGED)).String()
	if _, err := o.Raw(sql).QueryRows(&listBookingCharge); err != nil {
		return nil, err
	}
	return listBookingCharge, nil
}

package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"time"
)

type BookingGuests struct {
	Id          int       `orm:"column(booking_guest_id);auto"`
	BookingId   int       `orm:"column(booking_id);"`
	GuestId     int       `orm:"column(guest_id);"`
	IsMainGuest int8      `orm:"column(is_main_guest);"`
	DeletedAt   int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type AddBookingGuests struct {
	BookingId   int       `orm:"column(booking_id);"`
	GuestId     int       `orm:"column(guest_id);"`
	IsMainGuest int8      `orm:"column(is_main_guest);"`
	DeletedAt   int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

func (t *BookingGuests) TableName() string {
	return "booking_guests"
}

func init() {
	orm.RegisterModel(new(BookingGuests))
}

// AddBookingGuest insert a new Guests into database and returns
// last inserted Id on success.
func AddBookingGuest(m *BookingGuests) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddBookingGuestTransaction insert a new Guests into database and returns
// last inserted Id on success.
func AddBookingGuestTransaction(m *[]BookingGuests, o orm.Ormer) (successNumber int64, err error) {
	successNumber, err = o.InsertMulti(100, m)
	return successNumber, err
}

// UpdateBookingGuestTransaction update Guests
// last inserted Id on success.
func UpdateBookingGuestTransaction(bookingId int, stringGuestId string, o orm.Ormer) error {
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Update("booking_guests").Set("deleted_at = " + fmt.Sprint(conf.IS_DELETED)).
		Where("booking_id = " + fmt.Sprint(bookingId)).And("guest_id").In(stringGuestId)

	query := sql.String()
	if _, err := o.Raw(query).Exec(); err != nil {
		return err
	}
	return nil
}

func GetListBookingGuest(query map[string]string) ([]BookingGuests, error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	qd.Select("booking_guests.*").
		From("booking_guests").
		Where("deleted_at = " + fmt.Sprint(conf.NOT_DELETED))

	if query["BookingId"] != "" {
		qd.And("booking_id").In(query["BookingId"])
	}
	if query["GuestId"] != "" {
		qd.And("guest_id").In(query["GuestId"])
	}

	var bookingGuests []BookingGuests
	if num, err := o.Raw(qd.String()).QueryRows(&bookingGuests); err != nil || num < 1 {
		return nil, err
	}
	return bookingGuests, nil
}

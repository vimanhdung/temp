package models

import (
	"errors"
	"fmt"
	"indetail/conf"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SmartLocks struct {
	Id          int       `orm:"column(smart_lock_id);auto"`
	Name        string    `orm:"column(name);size(255)"`
	DeviceId    string    `orm:"column(device_id);size(255)"`
	DeletedAt   int8      `orm:"column(deleted_at);null"`
	CreatedUser int       `orm:"column(created_user)"`
	UpdatedUser int       `orm:"column(updated_user);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type GetStatusSmartLock struct {
	Locked     bool `orm:"column(locked)" json:"locked"`
	Battery    int  `orm:"column(battery)" json:"battery"`
	Responsive bool `orm:"column(responsive)" json:"responsive"`
}

func (t *SmartLocks) TableName() string {
	return "smart_locks"
}

func init() {
	orm.RegisterModel(new(SmartLocks))
}

// AddSmartLocks insert a new SmartLocks into database and returns
// last inserted Id on success.
func AddSmartLocks(m *SmartLocks) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSmartLocksById retrieves SmartLocks by Id. Returns error if
// Id doesn't exist
func GetSmartLocksById(id int) (v *SmartLocks, err error) {
	o := orm.NewOrm()
	v = &SmartLocks{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSmartLocks retrieves all SmartLocks matches certain condition. Returns empty list if
// no records exist
func GetAllSmartLocks(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SmartLocks))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
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

	var l []SmartLocks
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

// UpdateSmartLocks updates SmartLocks by Id and returns error if
// the record to be updated doesn't exist
func UpdateSmartLocksById(m *SmartLocks) (err error) {
	o := orm.NewOrm()
	v := SmartLocks{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSmartLocks deletes SmartLocks by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSmartLocks(id int) (err error) {
	o := orm.NewOrm()
	v := SmartLocks{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SmartLocks{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func CheckAvailableSmartLock(query map[string]string) (isAvailable bool, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("smart_locks.smart_lock_id").
		From("smart_locks").
		InnerJoin("rooms").On("smart_locks.smart_lock_id = rooms.smart_lock_id").
		Where("rooms.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("smart_locks.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("smart_locks.smart_lock_id = " + query["SmartLockId"])

	if query["RoomId"] != "" {
		qb.And("rooms.room_id != " + query["RoomId"])
	}
	o := orm.NewOrm()
	var listSmartLock []SmartLocks
	if num, err := o.Raw(qb.String()).QueryRows(&listSmartLock); err != nil || num > 0 {
		return false, err
	}
	return true, nil
}

// count smart lock records
func CountSmartLocks(query map[string]string) (totalRecord int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SmartLocks))
	// query k=v
	qs = RebuildConditions(qs, query)

	if totalRecord, err = qs.Count(); err == nil {
		return totalRecord, err
	} else {
		return 0, err
	}
}

func GetSmartLockDeviceIdByBookingId(query map[string]string) (deviceId string, err error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("smart_locks.*").
		From("smart_locks").
		InnerJoin("rooms").On("smart_locks.smart_lock_id = rooms.smart_lock_id").
		InnerJoin("bookings").On("rooms.room_id = bookings.room_id")

	qbWhere, _ := orm.NewQueryBuilder("mysql")
	qbWhere.Where("rooms.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("smart_locks.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.booking_id = " + query["bookingId"]).
		And("bookings.status = " + fmt.Sprint(conf.BOOKING_CHECKIN))

	if query["guestId"] != "" {
		qb.InnerJoin("booking_guests").On("bookings.booking_id = booking_guests.booking_id").
			InnerJoin("guests").On("booking_guests.guest_id = guests.guest_id")

		qbWhere.And("booking_guests.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
			And("guests.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
			And("guests.guest_id = " + query["guestId"])
	}

	if query["userAppId"] != "" {
		qb.InnerJoin("user_app_accounts").On("bookings.booking_id = user_app_accounts.booking_id")

		qbWhere.And("user_app_accounts.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
			And("user_app_accounts.user_app_account_id = " + query["userAppId"])
	}

	o := orm.NewOrm()
	var listSmartLock []SmartLocks
	sql := qb.String() + " " + qbWhere.String()
	if num, err := o.Raw(sql).QueryRows(&listSmartLock); err != nil || num < 1 {
		return "", err
	}
	return listSmartLock[0].DeviceId, nil
}

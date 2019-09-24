package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type State struct {
	Id             int    `orm:"column(state_id);auto" description:"state_id"`
	LineId         string `orm:"column(line_id);size(63);null" description:"line_id"`
	Status         int    `orm:"column(status);null" description:"status"`
	Checkin        string `orm:"column(checkin);type(date);null" description:"checkin"`
	Checkout       string `orm:"column(checkout);type(date);null" description:"checkout"`
	Count          int    `orm:"column(count);null" description:"count"`
	RoomId         int    `orm:"column(room_id);null" description:"room_id"`
	BookingId      int    `orm:"column(booking_id);null" description:"booking_id"`
	HotelId        int    `orm:"column(hotel_id);null" description:"hotel_id"`
	FullName       string `orm:"column(full_name);size(60);null"`
	PassportNumber string `orm:"column(passport_number);size(60);null"`
	TotalPrice     string `orm:"column(total_price);size(60);null"`
}

func (t *State) TableName() string {
	return "state"
}

func init() {
	orm.RegisterModel(new(State))
}

// AddState insert a new State into database and returns
// last inserted Id on success.
func AddState(m *State) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetStateByLineId retrieves State by LineId. Returns error if
// Id doesn't exist
func GetStateByLineId(lineId string) (State, bool) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	sql := qb.Select("*").From("state").
		Where("line_id = '" + lineId + "'").String()

	var tmpState []State
	num, tmpErr := o.Raw(sql).QueryRows(&tmpState)
	if num > 0 && tmpErr == nil {
		return tmpState[0], true
	}
	return State{}, false
}

// UpdateState updates State by Id and returns error if
// the record to be updated doesn't exist
func UpdateStateById(m *State) (err error) {
	o := orm.NewOrm()
	v := State{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

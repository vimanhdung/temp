package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"time"
)

type FcmTokenManage struct {
	Id        int       `orm:"column(id);auto"`
	FcmToken  string    `orm:"column(fcm_token)"`
	DeviceId  string    `orm:"column(device_id)"`
	HotelId   int       `orm:"column(hotel_id)"`
	DeletedAt int8      `orm:"column(deleted_at);null"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type AddFcmTokenManage struct {
	FcmToken string
	DeviceId string
	HotelId  int
}

func (t *FcmTokenManage) TableName() string {
	return "fcm_token_manage"
}

func init() {
	orm.RegisterModel(new(FcmTokenManage))
}

/**
*Func create record Fcm Token
 */
func CreateFcmTokenManage(m *FcmTokenManage) (id int64, err error) {
	//create record Fcm Token
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetFcmTokenManageByFcmToken(fcmToken string) (token FcmTokenManage, bool bool) {
	o := orm.NewOrm()
	qs := o.QueryTable("fcm_token_manage")
	err := qs.Filter("fcm_token", fcmToken).Filter("deleted_at", conf.NOT_DELETED).One(&token)
	if err == nil {
		return token, true
	}
	return token, false
}

func GetFcmTokenManageByDeviceId(deviceId string) (token FcmTokenManage, bool bool) {
	o := orm.NewOrm()
	qs := o.QueryTable("fcm_token_manage")
	err := qs.Filter("device_id", deviceId).Filter("deleted_at", conf.NOT_DELETED).One(&token)
	if err == nil {
		return token, true
	}
	return token, false
}

func GetFcmTokenManageByHotelId(hotelId int) (token FcmTokenManage, bool bool) {
	o := orm.NewOrm()
	qs := o.QueryTable("fcm_token_manage")
	err := qs.Filter("hotel_id", hotelId).Filter("deleted_at", conf.NOT_DELETED).One(&token)
	if err == nil {
		return token, true
	}
	return token, false
}

func UpdateFcmTokenManage(m *FcmTokenManage) (err error) {
	o := orm.NewOrm()
	v := FcmTokenManage{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func DeleteFcmTokenManage(id int) (err error) {
	o := orm.NewOrm()
	if num, err := o.Delete(&FcmTokenManage{Id: id}); err == nil {
		fmt.Sprintf("Delete succesfull %s row", num)
	}
	return
}

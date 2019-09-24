package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableBookings_20190820_234707 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableBookings_20190820_234707{}
	m.Created = "20190820_234707"

	migration.Register("CreateTableBookings_20190820_234707", m)
}

// Run the migrations
func (m *CreateTableBookings_20190820_234707) Up() {
	m.SQL("DROP TABLE IF EXISTS `bookings`")
	m.SQL("SET sql_mode = ''")
	m.SQL("SET GLOBAL sql_mode = ''")

	m.CreateTable("bookings", "InnoDB", "utf8")
	m.PriCol("booking_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("room_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("hotel_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("booking_ota_id").SetDataType("VARCHAR(30)").SetNullable(false)
	m.NewCol("pay_type").SetDataType("TINYINT(2)").SetNullable(false)
	m.NewCol("guest_count").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("adult_count").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("children_count").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("checkin_date").SetDataType("TIMESTAMP").SetNullable(false)
	m.NewCol("checkout_date").SetDataType("TIMESTAMP").SetNullable(false)
	m.NewCol("description").SetDataType("VARCHAR(255)").SetNullable(true)
	m.NewCol("note").SetDataType("VARCHAR(255)").SetNullable(true)
	m.NewCol("total_amount").SetDataType("VARCHAR(30)").SetNullable(true)
	m.NewCol("status").SetDataType("TINYINT(2)").SetNullable(false).SetDefault("0")
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableBookings_20190820_234707) Down() {
	m.SQL("DROP TABLE bookings")

}

package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableBookingGuests_20190820_234708 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableBookingGuests_20190820_234708{}
	m.Created = "20190820_234708"

	migration.Register("CreateTableBookingGuests_20190820_234708", m)
}

// Run the migrations
func (m *CreateTableBookingGuests_20190820_234708) Up() {
	m.SQL("DROP TABLE IF EXISTS `booking_guests`")

	m.CreateTable("booking_guests", "InnoDB", "utf8")
	m.PriCol("booking_guest_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("booking_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("guest_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("is_main_guest").SetDataType("TINYINT(2)").SetNullable(false).SetDefault("0")
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableBookingGuests_20190820_234708) Down() {
	m.SQL("DROP TABLE booking_guests")

}

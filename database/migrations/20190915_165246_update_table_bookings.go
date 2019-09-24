package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UpdateTableBookings_20190915_165246 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UpdateTableBookings_20190915_165246{}
	m.Created = "20190915_165246"

	migration.Register("UpdateTableBookings_20190915_165246", m)
}

// Run the migrations
func (m *UpdateTableBookings_20190915_165246) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.AlterTable("bookings")
	m.NewCol("actual_guest_count").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("qr_image_url").SetDataType("VARCHAR(255)").SetNullable(true)
	sql := m.GetSQL()
	m.SQL(sql)

}

// Reverse the migrations
func (m *UpdateTableBookings_20190915_165246) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.AlterTable("guests")
	m.NewCol("actual_guest_count").Remove()
	m.NewCol("qr_image_url").Remove()
	sql := m.GetSQL()
	m.SQL(sql)
}

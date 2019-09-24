package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type UpdateTableGuests_20190915_155246 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &UpdateTableGuests_20190915_155246{}
	m.Created = "20190915_155246"

	migration.Register("UpdateTableGuests_20190915_155246", m)
}

// Run the migrations
func (m *UpdateTableGuests_20190915_155246) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.AlterTable("guests")
	m.NewCol("passport_expired").SetDataType("date").SetNullable(true)
	m.NewCol("nationality").SetDataType("VARCHAR(10)").SetNullable(true)
	m.NewCol("login_name").SetDataType("VARCHAR(255)").SetNullable(true)
	sql := m.GetSQL()
	m.SQL(sql)

}

// Reverse the migrations
func (m *UpdateTableGuests_20190915_155246) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.AlterTable("guests")
	m.NewCol("passport_expired").Remove()
	m.NewCol("nationality").Remove()
	m.NewCol("login_name").Remove()
	sql := m.GetSQL()
	m.SQL(sql)
}

package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableGuests_20190820_234715 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableGuests_20190820_234715{}
	m.Created = "20190820_234715"

	migration.Register("CreateTableGuests_20190820_234715", m)
}

// Run the migrations
func (m *CreateTableGuests_20190820_234715) Up() {
	m.SQL("DROP TABLE IF EXISTS `guests`")

	m.CreateTable("guests", "InnoDB", "utf8")
	m.PriCol("guest_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("full_name").SetDataType("VARCHAR(60)").SetNullable(true)
	m.NewCol("status").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("email").SetDataType("VARCHAR(255)").SetNullable(true)
	m.NewCol("password").SetDataType("VARCHAR(60)").SetNullable(true)
	m.NewCol("phone").SetDataType("VARCHAR(30)").SetNullable(true)
	m.NewCol("address").SetDataType("VARCHAR(255)").SetNullable(true)
	m.NewCol("passport_number").SetDataType("VARCHAR(30)").SetNullable(true)
	m.NewCol("birth_day").SetDataType("date").SetNullable(true)
	m.NewCol("gender").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("occupation").SetDataType("VARCHAR(255)").SetNullable(true)
	m.NewCol("files").SetDataType("JSON").SetNullable(true)
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableGuests_20190820_234715) Down() {
	m.SQL("DROP TABLE guests")

}

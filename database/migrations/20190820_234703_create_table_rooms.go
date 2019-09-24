package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableRooms_20190820_234703 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableRooms_20190820_234703{}
	m.Created = "20190820_234703"

	migration.Register("CreateTableRooms_20190820_234703", m)
}

// Run the migrations
func (m *CreateTableRooms_20190820_234703) Up() {
	m.SQL("DROP TABLE IF EXISTS `rooms`")

	m.CreateTable("rooms", "InnoDB", "utf8")
	m.PriCol("room_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("hotel_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("smart_lock_id").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("room_name").SetDataType("VARCHAR(255)").SetNullable(true)
	m.NewCol("room_ota_id").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableRooms_20190820_234703) Down() {
	m.SQL("DROP TABLE rooms")

}

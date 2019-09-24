package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableHotels_20190820_234701 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableHotels_20190820_234701{}
	m.Created = "20190820_234701"

	migration.Register("CreateTableHotels_20190820_234701", m)
}

// Run the migrations
func (m *CreateTableHotels_20190820_234701) Up() {
	m.SQL("DROP TABLE IF EXISTS `hotels`")

	m.CreateTable("hotels", "InnoDB", "utf8")
	m.PriCol("hotel_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("hotel_name").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("sesami_api_auth_key").SetDataType("VARCHAR(255)").SetNullable(false)
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
func (m *CreateTableHotels_20190820_234701) Down() {
	m.SQL("DROP TABLE hotels")
}

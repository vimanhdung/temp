package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableOtas_20190820_234718 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableOtas_20190820_234718{}
	m.Created = "20190820_234718"

	migration.Register("CreateTableOtas_20190820_234718", m)
}

// Run the migrations
func (m *CreateTableOtas_20190820_234718) Up() {
	m.SQL("DROP TABLE IF EXISTS `otas`")

	m.CreateTable("otas", "InnoDB", "utf8")
	m.PriCol("ota_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("hotel_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("ota_kind_code").SetDataType("TINYINT(2)").SetNullable(false).SetDefault("0")
	m.NewCol("ota_login_name").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("ota_password").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("status").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableOtas_20190820_234718) Down() {
	m.SQL("DROP TABLE otas")

}

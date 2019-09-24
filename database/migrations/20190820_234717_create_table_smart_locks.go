package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableSmartLocks_20190820_234717 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableSmartLocks_20190820_234717{}
	m.Created = "20190820_234717"

	migration.Register("CreateTableSmartLocks_20190820_234717", m)
}

// Run the migrations
func (m *CreateTableSmartLocks_20190820_234717) Up() {
	m.SQL("DROP TABLE IF EXISTS `smart_locks`")

	m.CreateTable("smart_locks", "InnoDB", "utf8")
	m.PriCol("smart_lock_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("name").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("device_id").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableSmartLocks_20190820_234717) Down() {
	m.SQL("DROP TABLE smart_locks")

}

package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableAdminAccounts_20190820_234711 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableAdminAccounts_20190820_234711{}
	m.Created = "20190820_234711"

	migration.Register("CreateTableAdminAccounts_20190820_234711", m)
}

// Run the migrations
func (m *CreateTableAdminAccounts_20190820_234711) Up() {
	m.SQL("DROP TABLE IF EXISTS `admin_accounts`")

	m.CreateTable("admin_accounts", "InnoDB", "utf8")
	m.PriCol("admin_account_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("hotel_id").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("email").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("password").SetDataType("VARCHAR(60)").SetNullable(false)
	m.NewCol("status").SetDataType("TINYINT(2)").SetNullable(false).SetDefault("0")
	m.NewCol("full_name").SetDataType("VARCHAR(60)").SetNullable(true)
	m.NewCol("deleted_at").SetDataType("TINYINT(2)").SetNullable(true).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableAdminAccounts_20190820_234711) Down() {
	m.SQL("DROP TABLE admin_accounts")

}

package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableUserAppAccount_20190820_234706 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableUserAppAccount_20190820_234706{}
	m.Created = "20190820_234706"

	migration.Register("CreateTableUserAppAccount_20190820_234706", m)
}

// Run the migrations
func (m *CreateTableUserAppAccount_20190820_234706) Up() {
	m.SQL("DROP TABLE IF EXISTS `user_app_accounts`")

	m.CreateTable("user_app_accounts", "InnoDB", "utf8")
	m.PriCol("user_app_account_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.PriCol("booking_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("login_name").SetDataType("VARCHAR(60)").SetNullable(false)
	m.NewCol("password").SetDataType("VARCHAR(60)").SetNullable(false)
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
func (m *CreateTableUserAppAccount_20190820_234706) Down() {
	m.SQL("DROP TABLE user_app_accounts")

}

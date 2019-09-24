package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableOtaEmails_20190820_234719 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableOtaEmails_20190820_234719{}
	m.Created = "20190820_234719"

	migration.Register("CreateTableOtaEmails_20190820_234719", m)
}

// Run the migrations
func (m *CreateTableOtaEmails_20190820_234719) Up() {
	m.SQL("DROP TABLE IF EXISTS `ota_emails`")

	m.CreateTable("ota_emails", "InnoDB", "utf8")
	m.PriCol("ota_email_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("hotel_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("subject").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("email_body").SetDataType("TEXT").SetNullable(false)
	m.NewCol("message_body").SetDataType("TEXT").SetNullable(true)
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
func (m *CreateTableOtaEmails_20190820_234719) Down() {
	m.SQL("DROP TABLE ota_emails")

}

package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTablePasswordResets_20190828_140409 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTablePasswordResets_20190828_140409{}
	m.Created = "20190828_140409"

	migration.Register("CreateTablePasswordResets_20190828_140409", m)
}

// Run the migrations
func (m *CreateTablePasswordResets_20190828_140409) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("DROP TABLE IF EXISTS `password_resets`")

	m.CreateTable("password_resets", "InnoDB", "utf8")
	m.PriCol("password_reset_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("account_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("email").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("code").SetDataType("VARCHAR(50)").SetNullable(false)
	m.NewCol("expire").SetDataType("VARCHAR(50)").SetNullable(false)
	m.NewCol("token").SetDataType("VARCHAR(50)").SetNullable(true)
	m.NewCol("type").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTablePasswordResets_20190828_140409) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE password_resets")
}

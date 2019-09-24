package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateJwtTokensTable_20190820_234722 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateJwtTokensTable_20190820_234722{}
	m.Created = "20190820_234722"

	migration.Register("CreateJwtTokensTable_20190820_234722", m)
}

// Run the migrations
func (m *CreateJwtTokensTable_20190820_234722) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("DROP TABLE IF EXISTS `jwt_tokens`")

	m.CreateTable("jwt_tokens", "InnoDB", "utf8")
	m.PriCol("jwt_token_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("account_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("jti").SetDataType("VARCHAR(255)").SetNullable(false)
	m.NewCol("type").SetDataType("TINYINT(2)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")

	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateJwtTokensTable_20190820_234722) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE jwt_tokens")
}

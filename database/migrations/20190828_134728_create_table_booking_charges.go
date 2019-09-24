package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTableBookingCharges_20190828_134728 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTableBookingCharges_20190828_134728{}
	m.Created = "20190828_134728"

	migration.Register("CreateTableBookingCharges_20190828_134728", m)
}

// Run the migrations
func (m *CreateTableBookingCharges_20190828_134728) Up() {
	m.SQL("DROP TABLE IF EXISTS `booking_charges`")

	m.CreateTable("booking_charges", "InnoDB", "utf8")
	m.PriCol("booking_charge_id").SetAuto(true).SetNullable(false).SetDataType("INT(10)").SetUnsigned(true)
	m.NewCol("booking_id").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("total_amount").SetDataType("VARCHAR(30)").SetNullable(false)
	m.NewCol("charge_datetime").SetDataType("TIMESTAMP").SetNullable(false)
	m.NewCol("transaction").SetDataType("TEXT").SetNullable(false)
	m.NewCol("status").SetDataType("TINYINT(2)").SetNullable(false).SetDefault("0")
	m.NewCol("created_user").SetDataType("INT(10)").SetNullable(false)
	m.NewCol("updated_user").SetDataType("INT(10)").SetNullable(true)
	m.NewCol("created_at").SetDataType("TIMESTAMP").SetNullable(false).SetDefault("CURRENT_TIMESTAMP")
	m.NewCol("updated_at").SetDataType("TIMESTAMP").SetNullable(true).SetDefault("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	sql := m.GetSQL()
	m.SQL(sql)
}

// Reverse the migrations
func (m *CreateTableBookingCharges_20190828_134728) Down() {
	m.SQL("DROP TABLE booking_charges")
}

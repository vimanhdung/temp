package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateSeederBase_20190911_155246 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateSeederBase_20190911_155246{}
	m.Created = "20190911_155246"

	migration.Register("CreateSeederBase_20190911_155246", m)
}

// Run the migrations
func (m *CreateSeederBase_20190911_155246) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("INSERT INTO guests (full_name, status, email, files, password, created_user, created_at, updated_at) VALUES ('QuanNguyen', 1, 'hoangquan.it.hcm@gmail.com', '{}', '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")
	m.SQL("INSERT INTO guests (full_name, status, email, files, password, created_user, created_at, updated_at) VALUES ('QuanNguyen', 1, 'hoangquan.it.hcm3@gmail.com', '{}', '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")
	m.SQL("INSERT INTO guests (full_name, status, email, files, password, created_user, created_at, updated_at) VALUES ('ThuyTran', 1, 'thuy.tran@mor.com.vn', '{}', '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")

	m.SQL("INSERT INTO user_app_accounts (booking_id, login_name, password, created_user, created_at, updated_at) VALUES (1, 'QuanNguyen',  '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")
	m.SQL("INSERT INTO user_app_accounts (booking_id, login_name, password, created_user, created_at, updated_at) VALUES (1, 'ThuyTran',  '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")
	m.SQL("INSERT INTO user_app_accounts (booking_id, login_name, password, created_user, created_at, updated_at) VALUES (1, 'ThanhTung',  '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")
	m.SQL("INSERT INTO user_app_accounts (booking_id, login_name, password, created_user, created_at, updated_at) VALUES (1, 'Indetail',  '$2a$10$rfatcwgpnt3cheVdnlQkseOt8l3ROhlCmW1w4KpPuXOp.LbbupEXq', 1, '2019-09-09 02:36:13', '2019-09-09 02:36:13')")

}

// Reverse the migrations
func (m *CreateSeederBase_20190911_155246) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}

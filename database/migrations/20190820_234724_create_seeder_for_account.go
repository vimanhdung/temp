package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateSeederForAccount_20190820_234724 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateSeederForAccount_20190820_234724{}
	m.Created = "20190820_234724"

	migration.Register("CreateSeederForAccount_20190820_234724", m)
}

// Run the migrations
func (m *CreateSeederForAccount_20190820_234724) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("INSERT INTO admin_accounts (email, password, created_user) VALUES ('admin@gmail.com',  '$2a$10$tLjipl9yWSaSUnSUXDSFGOsa2dYZ20zWLksaPzVIAGhwG4R10gXuq', 1)")
	m.SQL("INSERT INTO admin_accounts (email, password, created_user) VALUES ('hotel@gmail.com',  '$2a$10$tLjipl9yWSaSUnSUXDSFGOsa2dYZ20zWLksaPzVIAGhwG4R10gXuq', 1)")
}

// Reverse the migrations
func (m *CreateSeederForAccount_20190820_234724) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}

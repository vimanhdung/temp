package models

import (
	"fmt"
	"indetail/conf"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type PasswordResets struct {
	Id        int       `orm:"column(password_reset_id);auto"`
	AccountId int       `orm:"column(account_id);int(10)"`
	Email     string    `orm:"column(email);size(255)"`
	Code      string    `orm:"column(code);size(50)"`
	Expire    string    `orm:"column(expire);size(50)"`
	Token     string    `orm:"column(token);size(50)"`
	Type      int       `orm:"column(type);int(2)"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

func (t *PasswordResets) TableName() string {
	return "password_resets"
}

func init() {
	orm.RegisterModel(new(PasswordResets))
}

// AddPasswordResets insert a new PasswordResets into database and returns
// last inserted Id on success.
func AddPasswordResets(m *PasswordResets) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPasswordResetsById retrieves PasswordResets by Id. Returns error if
// Id doesn't exist
func GetPasswordResetsById(id int) (v *PasswordResets, err error) {
	o := orm.NewOrm()
	v = &PasswordResets{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdatePasswordResets updates PasswordResets by Id and returns error if
// the record to be updated doesn't exist
func UpdatePasswordResetsById(m *PasswordResets) (err error) {
	o := orm.NewOrm()
	v := PasswordResets{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePasswordResets deletes PasswordResets by Id and returns error if
// the record to be deleted doesn't exist
func DeletePasswordResets(o orm.Ormer, accountId int, typeReset int) (err error) {
	v := PasswordResets{AccountId: accountId, Type: typeReset}
	// ascertain id exists in the database
	if err = o.Read(&v, "AccountId", "Type"); err == nil {
		var num int64
		if num, err = o.Delete(&v, "AccountId", "Type"); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetStructInsert(m *AdminAccounts) (rp PasswordResets) {
	rp.AccountId = m.Id
	rp.Email = m.Email
	rp.Code = strconv.FormatInt(11111+rand.Int63n(50000), 10)
	rp.Expire = strconv.FormatInt(time.Now().Add(conf.EmailExpire).Unix(), 10)
	return rp
}

func GetStructGuestInsert(guest *Guests) (rp PasswordResets) {
	rp.AccountId = guest.Id
	rp.Email = guest.Email
	rp.Type = conf.TYPE_GUEST
	rp.Code = strconv.FormatInt(11111+rand.Int63n(50000), 10)
	rp.Expire = strconv.FormatInt(time.Now().Add(conf.EmailExpire).Unix(), 10)
	return rp
}

// GetPasswordResetsById retrieves PasswordResets by Id. Returns error if
// Id doesn't exist
func GetPasswordResetsByAccountId(accountId int) (v *PasswordResets, err error) {
	o := orm.NewOrm()
	v = &PasswordResets{AccountId: accountId}
	if err = o.Read(v, "AccountId"); err == nil {
		return v, nil
	}
	return nil, err
}

// Create Or Update Password Reset
// @Param m *PasswordResets
// @return err error
func CreateOrUpdatePasswordReset(m *PasswordResets) (err error, code string) {
	if rp, err := GetPasswordResetsByAccountId(m.AccountId); err == nil && rp != nil {
		rp.Code = strconv.FormatInt(11111+rand.Int63n(50000), 10)
		rp.Expire = strconv.FormatInt(time.Now().Add(conf.EmailExpire).Unix(), 10)
		return UpdatePasswordResetsById(rp), rp.Code
	} else {
		_, err := AddPasswordResets(m)
		return err, m.Code
	}
}

// Check Code Exists
// @Param ob ForgotPassword
// @return bool
func CheckCodeExists(ob ForgotConfirmCode) (bool, PasswordResets) {
	var rp PasswordResets
	o := orm.NewOrm()
	qs := o.QueryTable("password_resets").Filter("email", ob.Email).Filter("code", ob.Code).Filter("expire__gte", time.Now().Unix())
	qs.One(&rp)
	return qs.Exist(), rp
}

// Check Code Exists
// @Param ob ForgotPassword
// @return bool
func CheckTokenResetPassword(ob ForgotNewPassword) (bool, PasswordResets) {
	var rp PasswordResets
	o := orm.NewOrm()
	qs := o.QueryTable("password_resets").Filter("token", ob.TokenReset).Filter("expire__gte", time.Now().Unix())
	qs.One(&rp)
	return qs.Exist(), rp
}

// Update Account and delete Reset Password Row
// @Param user *AdminAccounts
// @return bool
func UpdateAccountAndDeleteResetRow(user *AdminAccounts) bool {
	o := orm.NewOrm()
	o.Begin()
	// Update Account
	update := UpdateAdminAccountsById(o, user)
	// Delete Password Reset
	rp, _ := GetPasswordResetsByAccountId(user.Id)
	destroy := DeletePasswordResets(o, rp.Id, conf.TYPE_ADMIN)
	// Destroy all token of this user
	destroyToken := DestroyAllTokenWithAccountId(o, user.Id, conf.TYPE_ADMIN)
	if update != nil || destroy != nil || destroyToken != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

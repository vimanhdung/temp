package models

import (
	"fmt"
	"indetail/conf"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdminAccounts struct {
	Id          int       `orm:"column(admin_account_id);auto"`
	HotelId     int       `orm:"column(hotel_id)"`
	Email       string    `orm:"column(email);size(255)"`
	Password    string    `orm:"column(password);size(60)"`
	Status      int8      `orm:"column(status);null"`
	FullName    string    `valid:"Required";orm:"column(full_name);size(120);null"`
	DeletedAt   int8      `orm:"column(deleted_at);null"`
	CreatedUser int       `orm:"column(created_user)"`
	UpdatedUser int       `orm:"column(updated_user);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type ListAdminAccount struct {
	Id          int       `orm:"column(admin_account_id);auto"`
	HotelId     int       `orm:"column(hotel_id)"`
	Email       string    `orm:"column(email);size(255)"`
	Password    string    `orm:"column(password);size(60)"`
	Status      int8      `orm:"column(status);null"`
	FullName    string    `valid:"Required";orm:"column(full_name);size(120);null"`
	DeletedAt   int8      `orm:"column(deleted_at);null"`
	CreatedUser int       `orm:"column(created_user)"`
	UpdatedUser int       `orm:"column(updated_user);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type AddAdminAccountStruct struct {
	RoleId   int    `orm:"column(role_id);null"`
	IsRoot   int8   `orm:"column(is_root)"`
	Email    string `orm:"column(email);size(255)"`
	Status   int8   `orm:"column(status);null"`
	FullName string `orm:"column(full_name);size(120);null"`
}

type LoginFields struct {
	Email    string `json:"email" valid:"Required;Email"`
	Password string `json:"password" valid:"Required"`
}

type ChangePassFields struct {
	OldPass string `json:"oldPass" valid:"Required"`
	NewPass string `json:"newPass" valid:"Required;Password"`
}

type ForgotSendCode struct {
	Email string `json:"email" valid:"Required;Email"`
}

type RessetPassword struct {
	Email string `json:"email" valid:"Required;Email"`
}

type ForgotConfirmCode struct {
	Email string `json:"email" valid:"Required;Email"`
	Code  string `json:"code" valid:"Required"`
}

type ForgotNewPassword struct {
	TokenReset  string `json:"tokenReset" valid:"Required"`
	NewPassword string `json:"newPassword" valid:"Required;Password"`
}

type KioskLogout struct {
	Password string `json:"password" valid:"Required"`
}

type HotelAdmin struct {
	HotelId        int    `orm:"column(hotel_id)"`
	AdminAccountId int    `orm:"column(admin_account_id)"`
	FullName       string `orm:"column(full_name);size(120);null"`
}

type UpdateMyAccount struct {
	FullName string `orm:"column(full_name);size(120);null"`
	Files    string `json:"files"`
}

type InsertAccountFields struct {
	RoleId   int    `orm:"column(role_id)" valid:"Required; Max(2)"`
	Email    string `orm:"column(email)" valid:"Required;Email"`
	Password string `orm:"column(password)" valid:"Required;Password"`
	Status   int8   `orm:"column(status)" valid:"Max(2)"`
	FullName string `valid:"Required";orm:"column(full_name);size(120);null"`
}

type UpdateAccountFields struct {
	RoleId   int    `orm:"column(role_id)" valid:"Required; Max(2)"`
	Email    string `orm:"column(email)" valid:"Required;Email"`
	Password string `orm:"column(password)" valid:"Password"`
	Status   int8   `orm:"column(status)" valid:"Max(2)"`
	FullName string `valid:"Required";orm:"column(full_name);size(120);null"`
}

type AdminFiles struct {
	Type string `json:"type" valid:"Required"`
	Url  string `json:"url" valid:"Required"`
}

type CountTotal struct {
	Total int `orm:"column(total)"`
}

type GlobalUsers struct {
	Id        int    `orm:"column(admin_account_id);auto"`
	Password  string `orm:"column(password);size(255)" json:"password,omitempty"`
	Email     string `orm:"column(email);size(255)" json:"email,omitempty"`
	LoginName string `orm:"column(login_name);size(255)" json:"login_name,omitempty"`
	FullName  string `orm:"column(full_name);size(120);null" json:"full_name,omitempty"`
	Type      int    `orm:"column(type);null" json:"type,omitempty"`
}

func (t *AdminAccounts) TableName() string {
	return "admin_accounts"
}

func init() {
	orm.RegisterModel(new(AdminAccounts))
}

// AddAdminAccounts insert a new AdminAccounts into database and returns
// last inserted Id on success.
func AddAdminAccounts(m *AdminAccounts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAdminAccountsById retrieves AdminAccounts by Id. Returns error if
// Id doesn't exist
func GetAdminAccountsById(id int) (v *AdminAccounts, err error) {
	o := orm.NewOrm()
	v = &AdminAccounts{Id: id, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "Id", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAdminAccounts retrieves all AdminAccounts matches certain condition. Returns empty list if
// no records exist
func GetAllAdminAccounts(query map[string]string, fields string, sortby []string, order []string,
	offset int64, limit int64) (adminAccount []AdminAccounts, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select(fields).
		From("admin_accounts").
		LeftJoin("hotel_administrators").On("admin_accounts.admin_account_id = hotel_administrators.admin_account_id").
		Where("admin_accounts.deleted_at = " + strconv.Itoa(conf.NOT_DELETED))
	for key, value := range query {
		if strings.Contains(key, "searchLike") {
			arrKey := strings.Split(key, "__searchLike")
			sql = sql.And(arrKey[0] + " LIKE '%" + fmt.Sprint(value) + "%'")
		} else {
			sql = sql.And(key + " = " + fmt.Sprint(value))
		}
	}

	sql = sql.GroupBy("admin_accounts.admin_account_id")

	// order by:
	sortFields, sortErr := MakeOrderForSqlQuery(sortby, order)
	if sortErr != nil {
		return nil, sortErr
	}

	stringSql := sql.String()
	if sortFields != "" {
		stringSql = sql.String() + " Order By " + sortFields
	}

	if limit > 0 {
		stringSql += " LIMIT " + fmt.Sprint(limit)
		stringSql += " OFFSET " + fmt.Sprint(offset)
	}

	_, err = o.Raw(stringSql).QueryRows(&adminAccount)
	return
}

func CountAdminAccounts(query map[string]string) (total int64, err error) {
	o := orm.NewOrm()
	subQb, _ := orm.NewQueryBuilder("mysql")

	subSql := subQb.Select("admin_accounts.admin_account_id").
		From("admin_accounts").
		LeftJoin("hotel_administrators").On("admin_accounts.admin_account_id = hotel_administrators.admin_account_id").
		Where("admin_accounts.deleted_at = " + strconv.Itoa(conf.NOT_DELETED))
	for key, value := range query {
		if strings.Contains(key, "__searchLike") {
			arrKey := strings.Split(key, "__searchLike")
			subSql = subSql.And(arrKey[0] + " LIKE '%" + fmt.Sprint(value) + "%'")
		} else if strings.Contains(key, "__in") {
			arrKey := strings.Split(key, "__in")
			subSql = subSql.And(arrKey[0] + " IN (" + fmt.Sprint(value) + ")")
		} else {
			subSql = subSql.And(key + " = " + fmt.Sprint(value))
		}
	}
	subSql = subSql.GroupBy("admin_accounts.admin_account_id")

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("count(*) as total").
		From("(" + subSql.String() + ") sl").String()
	var countTotal []CountTotal
	if count, err := o.Raw(sql).QueryRows(&countTotal); err != nil || count < 1 {
		return 0, err
	}
	return int64(countTotal[0].Total), nil
}

// UpdateAdminAccounts updates AdminAccounts by Id and returns error if
// the record to be updated doesn't exist
func UpdateAdminAccountsById(o orm.Ormer, m *AdminAccounts) (err error) {
	v := AdminAccounts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// Get User By Email
// Return user *AdminAccounts, error
func GetUserByEmail(email string) (user *AdminAccounts, err error) {
	o := orm.NewOrm()
	user = &AdminAccounts{Email: email, DeletedAt: 0}
	if err := o.Read(user, "Email", "DeletedAt"); err == nil {
		return user, nil
	}
	return nil, err
}

// Get Admin By Email
// Return user *libs.GlobalUsers, error
func PermissionGetAdminByEmail(email string) (GlobalUsers, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("admin_account_id, full_name, password, email").
		From("admin_accounts").
		Where("email = '" + email + "'").
		And("deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).String()

	var accounts []GlobalUsers
	_, err := o.Raw(sql).QueryRows(&accounts)
	return accounts[0], err
}

// Update Admin Account And Destroy Token
// destroy all token of this user
// @Param user *AdminAccounts
// @return bool
func UpdateAdminAccountsAndDestroyToken(user *AdminAccounts) bool {
	o := orm.NewOrm()
	o.Begin()
	// Update Admin Account
	update := UpdateAdminAccountsById(o, user)
	// Destroy token of this user
	destroy := DestroyAllTokenWithAccountId(o, user.Id, conf.TYPE_ADMIN)
	if update != nil || destroy != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

func GetListHotelIdByAdminName(qeryParams map[string]interface{}) (string) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("ha.hotel_id").
		From("admin_accounts as ad").
		InnerJoin("hotel_administrators as ha").On("ad.admin_account_id = ha.admin_account_id").
		Where("ad.full_name LIKE '%" + fmt.Sprint(qeryParams["hotelAdminName"]) + "%'").
		And("ha.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).
		And("ad.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).String()
	var hotelAdmin []HotelAdmin
	num, err := o.Raw(sql).QueryRows(&hotelAdmin)
	if num > 0 && err == nil {
		var resultHotelId = ""
		for _, value := range hotelAdmin {
			resultHotelId += fmt.Sprint(value.HotelId) + ","
		}
		return strings.TrimRight(resultHotelId, ",")
	}
	return ""
}

func GetListAdminNameByHotelId(qeryParams map[string]interface{}) ([]HotelAdmin) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("ha.hotel_id, ad.admin_account_id,	ad.full_name").
		From("admin_accounts as ad").
		InnerJoin("hotel_administrators as ha").On("ad.admin_account_id = ha.admin_account_id").
		InnerJoin("hotels as ho").On("ha.hotel_id = ho.hotel_id").
		Where("ho.hotel_id IN (" + fmt.Sprint(qeryParams["listHotelId"]) + ")").
		And("ho.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).
		And("ha.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).
		And("ad.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).String()
	var hotelAdmin []HotelAdmin
	num, err := o.Raw(sql).QueryRows(&hotelAdmin)
	if num > 0 && err == nil {
		return hotelAdmin
	}
	return nil
}

// Check Email Exists
// @param action int "action when check email 1: add new, 2: edit"
// @param user *AdminAccounts
// @param email string
// @return bool
func CheckEmailExists(action int, accountId int, email string) bool {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("admin_accounts.admin_account_id").
		From("admin_accounts").
		Where("admin_accounts.email = '" + email + "'").
		And("admin_accounts.deleted_at = 0")
	if action == conf.EDIT_ACCOUNT {
		qb = qb.And("admin_accounts.admin_account_id != " + fmt.Sprint(accountId))
	}
	sql := qb.String()
	var account []AdminAccounts
	num, err := o.Raw(sql).QueryRows(&account)
	if num > 0 && err == nil {
		return true
	}
	return false
}

// Soft Delete Account
// @Param id int "id of account"
// @return exists bool
// @return ok bool
func SoftDeleteAccount(id int) (exists bool, ok bool) {
	o := orm.NewOrm()
	v := AdminAccounts{Id: id, DeletedAt: conf.NOT_DELETED}
	// ascertain id exists in the database
	if err := o.Read(&v, "Id", "DeletedAt", "IsRoot"); err == nil {
		o.Begin()
		// Soft Delete Account
		v.DeletedAt = 1
		_, softDel := o.Update(&v)
		// Destroy Token
		destroy := DestroyAllTokenWithAccountId(o, id, conf.TYPE_ADMIN)
		if softDel != nil || destroy != nil {
			o.Rollback()
			return true, false
		}
		o.Commit()
		return true, true
	}
	return false, false
}

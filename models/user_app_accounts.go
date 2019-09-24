package models

import (
	"errors"
	"fmt"
	"indetail/conf"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserAppAccounts struct {
	Id          int       `orm:"column(user_app_account_id);auto"`
	BookingId   int    `orm:"column(booking_id);null"`
	LoginName   string    `orm:"column(login_name);size(255);null"`
	Password    string    `orm:"column(password);size(60)"`
	Status      string    `orm:"column(status);size(30);null"`
	DeletedAt   int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser int       `orm:"column(created_user)"`
	UpdatedUser int       `orm:"column(updated_user);null"`
	CreatedAt   time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type LoginUserAppFields struct {
	LoginName string `json:"login_name" valid:"Required"`
	Password  string `json:"password" valid:"Required"`
}

type InsertUserAppFields struct {
	LoginName string `orm:"column(login_name)" valid:"Required"`
	BookingId int    `orm:"column(booking_id)"`
	Password  string `orm:"column(password)" valid:"Required;Password"`
	Status    int8   `orm:"column(status)" valid:"Max(2)"`
}

func init() {
	orm.RegisterModel(new(UserAppAccounts))
}

func GetUserAppAccountByLoginName(loginName string) (v *UserAppAccounts, err error) {
	o := orm.NewOrm()
	v = &UserAppAccounts{LoginName: loginName, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "LoginName", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

func ChangePasswordUserAppAndLogout(userApp *UserAppAccounts) bool {
	o := orm.NewOrm()
	o.Begin()
	update := UpdateUserAppAccountsById(userApp)

	destroyToken := DestroyAllTokenWithAccountId(o, userApp.Id, conf.TYPE_USER)
	if update != nil || destroyToken != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

// AddUserAppAccounts insert a new AddUserAppAccounts into database and returns
// last inserted Id on success.
func AddUserAppAccounts(m *UserAppAccounts) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUser_app_accountsById retrieves User_app_accounts by Id. Returns error if
// Id doesn't exist
func GetUserAppAccountsById(id int) (v *UserAppAccounts, err error) {
	o := orm.NewOrm()
	v = &UserAppAccounts{Id: id, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "Id", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUser_app_accounts retrieves all User_app_accounts matches certain condition. Returns empty list if
// no records exist
func GetAllUserAppAccounts(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserAppAccounts))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []UserAppAccounts
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUserAppAccountByTransaction updates User_app_accounts by params and returns error if
// the record to be updated doesn't exist
func UpdateUserAppAccountByTransaction(bookingId int, params orm.Params, o orm.Ormer) (err error) {
	_, errUpdate := o.QueryTable(new(UserAppAccounts)).Filter("BookingId", bookingId).Update(params)
	return errUpdate
}

// UpdateUser_app_accounts updates User_app_accounts by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserAppAccountsById(m *UserAppAccounts) (err error) {
	o := orm.NewOrm()
	v := UserAppAccounts{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// Update Admin Account And Destroy Token
// destroy all token of this user
// @Param user *AdminAccounts
// @return bool
func UpdateUserAppAccountAndDestroyToken(userAppAccount *UserAppAccounts) bool {
	o := orm.NewOrm()
	o.Begin()
	// Update Admin Account
	update_err := UpdateUserAppAccountsById(userAppAccount)
	// Destroy token of this user
	//destroy_err := DeleteJwt_token_user_app_accounts_by_user_app_account_id(userAppAccount.Id)
	if update_err != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

// Permission Get User App By Login Name
// Return user *libs.GlobalUsers, error
func PermissionGetUserAppByLoginName(loginName string) (GlobalUsers, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("user_app_account_id as admin_account_id, password, login_name as email, login_name as full_name").
		From("user_app_accounts").
		Where("login_name = '" + loginName + "'").
		And("deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).String()

	var accounts []GlobalUsers
	_, err := o.Raw(sql).QueryRows(&accounts)
	return accounts[0], err
}

// CheckLoginNameExists
// @param action int "action when check login name 1: add new, 2: edit"
// @param accountId int "account id"
// @param email string
// @return bool
func CheckLoginNameExists(action int, accountId int, loginName string) bool {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("user_app_account_id").
		From("user_app_accounts").
		Where("login_name = '" + loginName + "'").
		And("deleted_at = " + strconv.Itoa(conf.NOT_DELETED))
	if action == conf.EDIT_ACCOUNT {
		qb = qb.And("user_app_account_id != " + fmt.Sprint(accountId))
	}
	sql := qb.String()
	var account []UserAppAccounts
	num, err := o.Raw(sql).QueryRows(&account)
	if num > 0 && err == nil {
		return true
	}
	return false
}

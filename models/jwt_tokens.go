package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type JwtTokens struct {
	Id        int       `orm:"column(jwt_token_id);auto"`
	AccountId int       `orm:"column(account_id)"`
	Jti       string    `orm:"column(jti);size(255)"`
	Type      int       `orm:"column(type)"`
	CreatedAt time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type JwtTokenLogin struct {
	Id    int
	Type  int
	Email string
}

func GetJwtTokenLogin(id int, typeLogin int, email string) (jwtTokenLogin JwtTokenLogin) {
	jwtTokenLogin.Id = id
	jwtTokenLogin.Type = typeLogin
	jwtTokenLogin.Email = email

	return jwtTokenLogin
}

func (t *JwtTokens) TableName() string {
	return "jwt_tokens"
}

func init() {
	orm.RegisterModel(new(JwtTokens))
}

// AddJwtTokens insert a new JwtTokens into database and returns
// last inserted Id on success.
func AddJwtTokens(o orm.Ormer, m *JwtTokens) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// DeleteJwtTokens deletes JwtTokens by Id and returns error if
// the record to be deleted doesn't exists
func DeleteJwtTokens(o orm.Ormer, jti string) (err error) {
	v := JwtTokens{Jti: jti}
	// ascertain id exists in the database
	if err = o.Read(&v, "jti"); err == nil {
		var num int64
		if num, err = o.Delete(&JwtTokens{Jti: jti}, "jti"); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// Check JWT Token exists
// @Param jwt_id string
// @return bool
func CheckJWTExists(jwtId string) bool {
	o := orm.NewOrm()
	qs := o.QueryTable("jwt_tokens").Filter("jti", jwtId)
	return qs.Exist()
}

// Get JWT Token exists
// @Param jwt_id string
// @return bool
func GetJWT(jwtId string) (isExists bool, jwtToken *JwtTokens, err error) {
	o := orm.NewOrm()

	qs := o.QueryTable("jwt_tokens").Filter("jti", jwtId)
	if !qs.Exist() {
		return false, nil, nil
	}

	jwtToken = &JwtTokens{Jti: jwtId}
	if err = o.Read(jwtToken, "Jti"); err == nil {
		return true, jwtToken, nil
	}
	return false, nil, err
}

// Destroy Token And Save New Token
// @Param oldJti string "old jti of old token"
// @Param newJti string "new jti of new token"
// @Param user *AdminAccounts "data of user"
// @return bool
func DestroyTokenAndSaveNewToken(oldJti string, newJti string, user *AdminAccounts) bool {
	o := orm.NewOrm()
	o.Begin()
	// Delete Old Token
	delete := DeleteJwtTokens(o, oldJti)
	// Save New Token
	jwtToken := JwtTokens{AccountId: user.Id, Jti: newJti}
	_, save := AddJwtTokens(o, &jwtToken)
	if delete != nil || save != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

// Destroy All Token With Account ID
// @Param accountId int "id of account"
// @return err error
func DestroyAllTokenWithAccountId(o orm.Ormer, accountId int, typeToken int) (err error) {
	v := JwtTokens{AccountId: accountId, Type: typeToken}
	// ascertain id exists in the database
	if err := o.Read(&v, "AccountId", "Type"); err == nil {
		var num int64
		if num, err = o.Delete(&v, "AccountId", "Type"); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

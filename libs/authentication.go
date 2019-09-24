package libs

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"indetail/conf"
	"indetail/models"
	"strings"
	"time"
)

// Check Auth
// @param email string
// @param password string
// @return bool
func CheckAuth(email string, password string) (bool, *models.AdminAccounts) {
	user, err := models.GetUserByEmail(email)
	if err != nil || user == nil || !CheckHash(password, user.Password) {
		return false, nil
	}
	return true, user
}

// Check Auth User App Account
// @param email string
// @param password string
// @return bool
func CheckAuthUserAppAccount(loginName string, password string) (bool, *models.UserAppAccounts) {
	userAppAccount, err := models.GetUserAppAccountByLoginName(loginName)
	
	if err != nil || userAppAccount == nil || !CheckHash(password, userAppAccount.Password) {
		return false, nil
	}
	return true, userAppAccount
}

// Check Auth Guest
// @param email string
// @param password string
// @return bool
func CheckAuthGuest(email string, password string) (bool, *models.Guests) {
	guest, err := models.GetGuestByEmail(email)
	if err != nil || guest == nil || !CheckHash(password, guest.Password) {
		return false, nil
	}
	return true, guest
}

// Check Hash Password
// @param inputPassword string
// @param dbPasword string
// @return bool
func CheckHash(inputPassword string, dbPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(inputPassword)); err != nil {
		return false
	}
	return true
}

// Convert string to Hased Password
// @param password string
// @return string
func GetHashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// Get Token with email
// @Param email string
// @Param expires time.Duration
// @return interface{}
//func GetToken(user *models.AdminAccounts, expires time.Duration) interface{} {
func GetToken(user *models.JwtTokenLogin, expires time.Duration) interface{} {
	// get jti for token
	jti := GenerateMd5String()
	et := EasyToken{
		Username: user.Email,
		Jti:      jti,
	}
	tokenStr, _ := et.GetToken(expires)
	jwtToken := models.JwtTokens{AccountId: user.Id, Jti: jti, Type: user.Type}
	if _, err := models.AddJwtTokens(orm.NewOrm(), &jwtToken); err != nil {
	}
	return ResultJson(
		map[string]interface{}{"token": tokenStr},
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Login Success!",
		nil,
	)
}

// Parse Token
// @Param token string
// @return token string without Bearer
func ParseToken(token string) (string) {
	splitToken := strings.Split(token, "Bearer ")
	if len(splitToken) != 1 {
		return splitToken[1]
	}
	return ""
}

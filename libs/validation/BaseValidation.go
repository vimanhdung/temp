package validation

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"indetail/conf"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func init() {
	validation.AddCustomFunc("Password", Password)
	validation.AddCustomFunc("PhoneHotel", PhoneHotel)
}

// Check Validate
func CheckValidate(ob interface{}) (bool, map[string]interface{}) {
	var message = make(map[string]interface{})
	//modify msg validate
	rebuildValidationMessage()
	valid := validation.Validation{}
	b, err := valid.Valid(ob)
	if err != nil {
	}
	if !b {
		for _, err := range valid.Errors {
			message[err.Field] = err.Message
			if err.Name == "Match" {
				message[err.Field] = strconv.Itoa(conf.FIELD_FORMAT_INVALID)
			}
			if err.Name == "Min" {
				message[err.Field] = strconv.Itoa(conf.FIELD_MIN_INVALID)
			}
			if err.Name == "Max" {
				message[err.Field] = strconv.Itoa(conf.FIELD_MAX_INVALID)
			}
		}
		return false, message
	}
	return true, message
}

func rebuildValidationMessage() {
	defaultMessage := map[string]string{
		"Required": strconv.Itoa(conf.VARIABLE_REQUIRED),
		"Range":    strconv.Itoa(conf.VARIABLE_OUT_OF_RANGE),
		"Length":   strconv.Itoa(conf.VARIABLE_IS_OVER_LENGTH),
		"Numeric":  strconv.Itoa(conf.VARIABLE_IS_NOT_NUMERIC),
		"Email":    strconv.Itoa(conf.VARIABLE_IS_NOT_EMAIL),
		"Phone":    strconv.Itoa(conf.VARIABLE_IS_NOT_PHONENUMBER),
		"ZipCode":  strconv.Itoa(conf.VARIABLE_IS_NOT_ZIPCODE),
	}
	validation.SetDefaultMessage(defaultMessage)
}

// Check Date Format
// @Param str string
// @return bool
func CheckDate(str string) bool {
	var valid = regexp.MustCompile(conf.DATE_REGEXP)
	return valid.MatchString(str)
}

// Function Password validate custom
func Password(v *validation.Validation, obj interface{}, key string) {
	password, _ := obj.(string)

	if password == "" {
		return
	}

	if spaceIndex := strings.Index(password, " "); spaceIndex != -1 ||
		!CheckUpperCase(password) ||
		!CheckLength(password) ||
		!ContainNumber(password) {
		v.SetError("Password", fmt.Sprint(conf.PASSWORD_FORMAT_INVALID))
	}
}

// Function Phone validate custom
func PhoneHotel(v *validation.Validation, obj interface{}, key string) {
	phone, _ := obj.(string)

	if phone == "" {
		return
	}

	BaseValidatePhone(v, phone, "Hotline")
}

func BaseValidatePhone(v *validation.Validation, inputPhone string, fieldName string) {
	var valid = regexp.MustCompile(conf.PHONE_REGEX)
	if !valid.MatchString(inputPhone){
		v.SetError(fieldName, fmt.Sprint(conf.VARIABLE_IS_NOT_PHONENUMBER))
	}
}

// Check Upper Case
// @Param str string to check
// @return bool
func CheckUpperCase(str string) bool {
	for _, c := range []rune(str) {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

// Check String Contain Number
// @Param str string to check
// @return bool
func ContainNumber(str string) bool {
	for _, c := range []rune(str) {
		if unicode.IsNumber(c) {
			return true
		}
	}
	return false
}

// Check Length
// @Param str string
// @return bool
func CheckLength(str string) bool {
	return len(str) >= conf.PASSWORD_MIN && len(str) <= conf.PASSWORD_MAX
}
package conf

import (
	"time"
)

/**
*** Define Token Expires
 */
const (
	TokenExpires        = time.Hour * 72
	TokenExpiresForLine = time.Hour * 87600
)

/**
** Define Status Response
 */
const (
	STATUS_SUCCESS = "success"
	STATUS_ERROR   = "error"
)

const (
	TYPE_ADMIN = 1
	TYPE_GUEST = 2
	TYPE_USER  = 3
)

/**
** Define Prefectures
 */
func GetPrefecture() (Prefecture []map[int]string) {
	prefecture := map[int]string{
		1: "Hokkaido",
		2: "Aomori",
		3: "Iwate",
		4: "Miyagi",
		5: "Akita",
		6: "Yamagata",
		7: "Fukushima",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		8:  "Ibaraki",
		9:  "Tochigi",
		10: "Gunma",
		11: "Saitama",
		12: "Chiba",
		13: "Tokyo",
		14: "Kanagawa",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		15: "Niigata",
		16: "Toyama",
		17: "Ishikawa",
		18: "Fukui",
		19: "Yamanashi",
		20: "Nagano",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		21: "Gifu",
		22: "Shizuoka",
		23: "Aichi",
		24: "Mie",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		25: "Shiga",
		26: "Kyoto",
		27: "Osaka",
		28: "Hyogo",
		29: "Nara",
		30: "Wakayama",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		31: "Tottori",
		32: "Shimane",
		33: "Okayama",
		34: "Hiroshima",
		35: "Yamaguchi",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		36: "Tokushima",
		37: "Kagawa",
		38: "Ehime",
		39: "Kochi",
	}
	Prefecture = append(Prefecture, prefecture)
	prefecture = map[int]string{
		40: "Fukuoka",
		41: "Saga",
		42: "Nagasaki",
		43: "Kumamoto",
		44: "Oita",
		45: "Miyazaki",
		46: "Kagoshima",
		47: "Okinawa",
	}
	Prefecture = append(Prefecture, prefecture)
	return Prefecture
}

/**
** Define Error Code
 */
const (
	//0 define for success
	SUCCESS = 0
	//1XX define for permission
	PERMISSION_DENY           = 101
	LOGIN_EXPIRED             = 102
	PASSWORD_CHANGED          = 103
	TOKEN_INVALID             = 104
	LOGOUT_FALSE              = 105
	BOX_AUTHENTICATION_FAILED = 110
	INCORRECT                 = 111

	//2XX define for validate
	MISSNG_PARAM   = 201
	TOO_MANY_PARAM = 202
	MISSNG_CONF    = 203

	VARIABLE_IS_NOT_NUMERIC          = 210
	VARIABLE_IS_NOT_JSON             = 211
	VARIABLE_IS_NOT_TIMESTAMP        = 212
	VARIABLE_IS_NOT_EMAIL            = 213
	VARIABLE_IS_NOT_PHONENUMBER      = 214
	VARIABLE_IS_NOT_ZIPCODE          = 215
	VARIABLE_OUT_OF_RANGE            = 216
	VARIABLE_REQUIRED                = 217
	VARIABLE_IS_OVER_LENGTH          = 218
	VARIABLE_IS_NOT_POSITIVE_INTEGER = 219
	VARIABLE_IS_NOT_DATE             = 220
	PARSE_JSON_BODY_FALSE            = 221
	FIELD_FORMAT_INVALID             = 222
	FIELD_MIN_INVALID                = 223
	FIELD_MAX_INVALID                = 224
	PASSWORD_FORMAT_INVALID          = 225
	LIMIT_RECORD                     = 226
	DATE_TIME_INVALID                = 227
	VARIABLE_NOT_EQUAL               = 228

	//3XX define for database action
	RECORD_EXISTS        = 301
	RECORD_NOT_FOUND     = 302
	SAVE_FAILURES        = 303
	SQL_ERROR            = 304
	LIMITED_RECORD       = 305
	CREATE_FILE_FAILURES = 306

	//4XX define for upload file
	FILE_TOO_LARGE  = 401
	FILE_TYPE_WRONG = 402

	//5XX define for AI calling
	NOT_MATCH         = 501
	CHANGE_STATE_FAIL = 502

	//6XX define for checkin, checkout api
	NOT_VACANT           = 601
	NOT_CHECKIN          = 602
	BOOKING_CHECKINED    = 603
	NOT_CHECKIN_DAY      = 604
	ROOM_NOT_AVAILABLE   = 605
	BOOKING_HAS_CHECKOUT = 606
	BOOKING_HAS_REJECT   = 607
	BOOKING_NOT_PAID     = 608

	//PROCESS FAILURES
	PROCESS_FAILURES = 1000
)

/**
** Set Status Header
 */
const (
	SUCCESS_STATUS = 200
	ERROR_STATUS   = 403
)

//Ota constant
const (
	OTA_NAME_LIST = `["Expedia", "Booking.com", "Airbnb", "Rakuten"]`
)

//Hotel constant
const (
	HOTEL_PUBLISH_STATUS = 1
	HOTEL_STATUS_LIST    = `["New", "Active", "InActive"]`
	HOTEL_TYPE_LIST      = `["Other", "Motels", "Hotels", "Resorts", "Floatels", "Transit Hotels", "Hostels", "Capsule Hotels", "Villas", "Home Stay", "Ecotels"]`
	HOTEL_SUB_TYPE_LIST  = `["Other", "1 Star", "2 Star", "3 Star", "4 Star", "5 Star", "6 Star", "7 Star"]`
)

//Room constant
const (
	ROOM_VACANT      = 0
	ROOM_OCCUPIED    = 1
	ROOM_DIRTY       = 2
	ROOM_INACTIVE    = 3
	ROOM_STATUS_LIST = `["Vacant", "Occupied", "Dirty", "InActive"]`
)

//General constant
const (
	NOT_DELETED      = 0
	IS_DELETED       = 1
	RFC3339          = "2006-01-02"
	RFC_DATE_TIME    = "2006-01-02 15:04:05"
	RFC_YEAR_MONTH   = "200601"
	DATE_REGEXP      = "^[0-9]{4}(-|/)(0[1-9]|1[0-2])(-|/)(0[1-9]|[1-2][0-9]|3[0-1])$"
	PASSWORD_REGEX   = `^(?=.*?[A-Z]).{8,60}$`
	PHONE_REGEX      = `\d{10,16}$`
	DEFAULT_PASSWORD = "Abc123@"
)

/**
** Configuration Email
 */
const (
	EMAIL_FROM       = "example@mor.vn"
	HOST_NAME        = "smtp.gmail.com"
	PORT             = ":587"
	EMAIL_ACCOUNT    = "mor.noreply@gmail.com"
	PASSWORD_ACCOUNT = "dsldfvrviqdymamj"
)
const EmailExpire = time.Hour * 24

/**
** Define Action for Check Email exists
 */
const (
	NEW_ACCOUNT  = 1
	EDIT_ACCOUNT = 2
)

/**
** Define Booking Status
 */
const (
	BOOKING_NEW                 = 0
	BOOKING_ACTIVE              = 1
	BOOKING_CHECKIN             = 2
	BOOKING_CHECKOUT            = 3
	BOOKING_REJECT              = 4
	BOOKING_PREPAYMENT          = 1
	BOOKING_PAYMENT_ARRIVAL     = 2
	BOOKING_STATUS              = `["New", "Active", "Check-in", "Check-out", "Reject"]`
	BOOKING_VERIFY_PASSPORT_URL = "http://27.72.88.246:5000"
	BOOKING_VERIFY_NOT_MATCH    = "un-match"
	BOOKING_VERIFY_MATCH        = "match"
)

/**
** Define User App
 */
const (
	USER_APP_STATUS_ACTIVE   = 1
	USER_APP_STATUS_INACTIVE = 2
)

/**
** Define Smart Lock
 */
const (
	CANDY_HOUSE_API_BASE_URL = "https://api.candyhouse.co/public/sesame"
	CANDY_HOUSE_API_TOKEN    = "pJisOp9cQQ6LWtcTyvRCWe9wW10ABAPv26t-v6UiCJ_DX4bF4T5c9aRUyhx9sTH6AgZbjcE-tGAG"
)

/**
** Define Booking Charge
 */
const (
	BOOKING_CHARGE_STATUS_NEW     = 0
	BOOKING_CHARGE_STATUS_CHARGED = 1
	BOOKING_CHARGE_STATUS_FAILURE = 2
	CHARGE_INFO_SUCCESS_STATUS    = "000"
)

/**
** Define Customer
 */
const (
	GUEST_STATUS_LIST   = `["Active", "InActive"]`
	GUEST_OCCUPATION    = `["Office Worker", "Proprietor", "None"]`
	GUEST_IS_MAIN       = 1
	GUEST_IS_NOT_MAIN   = 0
	GUEST_MALE          = 1
	GUEST_FEMALE        = 2
	GUEST_FILE_PASSPORT = "passport"
	GUEST_FILE_PORTRAIT = "portrait"
)

/**
*** Define Permission Code
 */
const (
	// AUTH
	AUTH_CHANGE_PASSWORD = "auth01"
	AUTH_LOGIN_HOTEL     = "auth02"
	AUTH_LOGIN_CMS       = "auth03"
	AUTH_ME              = "auth04"
	AUTH_LOGIN_LINE      = "auth05"

	ADMIN_CMS_ROLE   = 1
	ADMIN_HOTEL_ROLE = 2
	ADMIN_LINE_ROLE  = 3

	// COMMON
	COMMON_UPLOAD_FILE = "common01"

	// ACCOUNT
	ACCOUNT_LIST             = "account01"
	ACCOUNT_CREATE           = "account02"
	ACCOUNT_UPDATE_MYACCOUNT = "account03"
	ACCOUNT_DETAIL           = "account04"
	ACCOUNT_UPDATE           = "account05"
	ACCOUNT_DELETE           = "account06"

	// REPORT
	REPORT_BOOKING          = "report01"
	REPORT_CUSTOMER         = "report02"
	REPORT_HOTEL_DASHBOARD  = "report03"
	REPORT_HOTEL_STATISTIC  = "report04"
	REPORT_HOTEL            = "report05"
	REPORT_PROMOTION        = "report06"
	REPORT_ROOM             = "report07"
	REPORT_SYSTEM_DASHBOARD = "report08"
	REPORT_SYSTEM_REPORT    = "report09"

	// ROOM
	ROOM_LIST   = "room01"
	ROOM_CREATE = "room02"
	ROOM_DETAIL = "room03"
	ROOM_UPDATE = "room04"
	ROOM_DELETE = "room05"

	// BOOKING
	BOOKING_LIST           = "booking01"
	BOOKING_CREATE         = "booking02"
	BOOKING_LIST_SERVICE   = "booking03"
	BOOKING_ADD_SERVICE    = "booking04"
	BOOKING_REMOVE_SERVICE = "booking05"
	BOOKING_INVOICE_DETAIL = "booking06"
	BOOKING_REJECT_STATUS  = "booking07"
	BOOKING_DETAIL         = "booking08"
	BOOKING_UPDATE         = "booking09"

	// HOTEL
	HOTEL_CREATE       = "hotel01"
	HOTEL_UPDATE       = "hotel02"
	HOTEL_DELETE       = "hotel03"
	HOTEL_GET_ALL      = "hotel04"
	HOTEL_LIST_BY_USER = "hotel05"
)

// Define Base64 salt
const SALT = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
const DOT = "."
const QR_CODE_EXTENSION = ".png"

const IS_ROOT = 1
const MANUAL_ACCOUNT = 0

// Define length
const (
	PASSWORD_MIN = 8
	PASSWORD_MAX = 60
)

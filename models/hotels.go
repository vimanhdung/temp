package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Hotels struct {
	Id               int       `orm:"column(hotel_id);auto"`
	HotelName        string    `valid:"Required";orm:"column(hotel_name);null"`
	SesamiApiAuthKey string    `orm:"column(sesami_api_auth_key);null"`
	Files            string    `orm:"column(files);null"`
	DeletedAt        int8      `orm:"column(deleted_at);null"`
	CreatedUser      int       `orm:"column(created_user)"`
	UpdatedUser      int       `orm:"column(updated_user);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type HotelFullValidate struct {
	Id               int       `orm:"column(hotel_id);auto"`
	HotelName        string    `valid:"Required";orm:"column(hotel_name);null"`
	SesamiApiAuthKey string    `orm:"column(sesami_api_auth_key);null"`
	Files            string    `orm:"column(files);null"`
	DeletedAt        int8      `orm:"column(deleted_at);null"`
	CreatedUser      int       `orm:"column(created_user)"`
	UpdatedUser      int       `orm:"column(updated_user);null"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt        time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type InsertHotel struct {
	HotelName        string `valid:"Required";orm:"column(hotel_name);null"`
	SesamiApiAuthKey string `orm:"column(sesami_api_auth_key);null"`
	Files            string `orm:"column(files);null"`
}

func (t *Hotels) TableName() string {
	return "hotels"
}

func init() {
	orm.RegisterModel(new(Hotels))
}

// AddHotels insert a new Hotels into database and returns
// last inserted Id on success.
func AddHotels(m *Hotels) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHotelsById retrieves Hotels by Id. Returns error if
// Id doesn't exists
func GetHotelsById(id int) (v *Hotels, err error) {
	o := orm.NewOrm()
	v = &Hotels{Id: id, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHotels retrieves all Hotels matches certain condition. Returns empty list if
// no records exists
func GetAllHotels(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Hotels))
	// query k=v
	qs = RebuildConditions(qs, query)
	// order by:
	qs, errorMq := MakeOrderForQuery(qs, sortby, order)

	if errorMq != nil {
		return ml, errorMq
	}

	var l []Hotels
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		var result1 []interface{}
		for _, v := range l {
			result1 = append(result1, v)
		}
		//filter data by using array field
		ml, err := FilterResultByField(fields, result1)
		return ml, err
	}
	return nil, err
}

// GetListHotelsWithVacantRooms retrieves all Hotels matches certain condition. Returns empty list if
// no records exists
func GetListHotelsWithVacantRooms(query map[string]string, fields string, sortby []string, order []string,
	offset int64, limit int64) (ml []Hotels, err error) {

	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select(fields).
		From("hotels").
		Where("hotels.hotel_id IN (" + query["Id.in"] + ")").
		And("hotels.deleted_at = " + strconv.Itoa(conf.NOT_DELETED))
	if query["hotels.status"] != "" {
		sql = sql.And("hotels.status IN (" + query["hotels.status"] + ")")
	}

	//build order
	orderFields, errOrder := MakeOrderForSqlQuery(sortby, order)
	if errOrder != nil {
		return nil, errOrder
	}
	if orderFields != "" {
		sql.OrderBy(orderFields)
	}
	if limit > 0 {
		sql.Limit(int(limit))
	}
	//set offset
	sql.Offset(int(offset))
	_, err = o.Raw(sql.String()).QueryRows(&ml)

	return ml, err
}

// GetListHotels retrieves all Hotels matches certain condition. Returns empty list if
// no records exists
func GetListHotels(query map[string]string) (listHotel []Hotels, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Hotels))
	// query k=v
	qs = RebuildConditions(qs, query)
	// order by:
	qs = qs.OrderBy("CreatedAt")

	if _, err = qs.All(&listHotel); err != nil {
		return nil, err
	}
	return
}

// count hotel records
func CountHotels(query map[string]string) (totalRecord int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Hotels))
	// query k=v
	qs = RebuildConditions(qs, query)

	if totalRecord, err = qs.Count(); err == nil {
		return totalRecord, err
	} else {
		return 0, err
	}
}

// UpdateHotels updates Hotels by Id and returns error if
// the record to be updated doesn't exists
func UpdateHotelsById(m *Hotels) (err error) {
	o := orm.NewOrm()
	v := Hotels{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHotels deletes Hotels by Id and returns error if
// the record to be deleted doesn't exists
func DeleteHotels(id int) (err error) {
	o := orm.NewOrm()
	v := Hotels{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Hotels{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//Validate and rebuild hotel multi language
func ValidateHotelMultiLanguage(hotelMultiLanguage interface{}, hotelStatus int8, detailErrorCode map[string]interface{}) (string, map[string]interface{}) {
	var rebuildHotelMultiLanguage map[string]interface{}
	if reflect.TypeOf(hotelMultiLanguage) == reflect.TypeOf("string") {
		if error := json.Unmarshal([]byte(hotelMultiLanguage.(string)), &rebuildHotelMultiLanguage); error != nil {
			return "", nil
		} else {
			hotelMultiLanguage = rebuildHotelMultiLanguage
		}
	}

	tmpReflectLanguage := reflect.ValueOf(hotelMultiLanguage)

	var test map[string]interface{}
	if tmpReflectLanguage.Type() != reflect.TypeOf(test) {
		return "", nil
	}
	var tmpLanguageInterface interface{}
	var hotelMultiLanguageInsert = ""
	var tmpMultiLanguageInsert = "{"
	var isPass = true
	for _, language := range tmpReflectLanguage.MapKeys() {
		tmpMultiLanguageInsert += "\"" + language.String() + "\":{"
		tmpLanguageInterface = tmpReflectLanguage.MapIndex(language).Interface()
		hotelJaLanguage := reflect.ValueOf(tmpLanguageInterface)
		for subIndex, key := range hotelJaLanguage.MapKeys() {
			strKey := key.String()
			var value = ""
			if hotelJaLanguage.MapIndex(key).Interface() != nil {
				value = strings.Trim(hotelJaLanguage.MapIndex(key).Interface().(string), " ")
			}

			if strKey == "hotel_name" && value == "" {
				detailErrorCode[strKey] = conf.VARIABLE_REQUIRED
				isPass = false
			}

			if hotelStatus == conf.HOTEL_PUBLISH_STATUS {
				if strKey == "description" && value == "" {
					detailErrorCode[strKey] = conf.VARIABLE_REQUIRED
					isPass = false
				}
				if strKey == "cancellation_policy" && value == "" {
					detailErrorCode[strKey] = conf.VARIABLE_REQUIRED
					isPass = false
				}
				if strKey == "rules" && value == "" {
					detailErrorCode[strKey] = conf.VARIABLE_REQUIRED
					isPass = false
				}
			}

			tmpMultiLanguageInsert += "\"" + strKey + "\": \"" + value + "\","
			if subIndex == (len(hotelJaLanguage.MapKeys()) - 1) {
				tmpMultiLanguageInsert = strings.TrimRight(tmpMultiLanguageInsert, ",")
			}
		}
		tmpMultiLanguageInsert += "},"
	}

	if tmpMultiLanguageInsert != "" {
		tmpMultiLanguageInsert = strings.TrimRight(tmpMultiLanguageInsert, ",")
		hotelMultiLanguageInsert += tmpMultiLanguageInsert + "}"
	}
	if !isPass {
		hotelMultiLanguageInsert = ""
	}
	return hotelMultiLanguageInsert, detailErrorCode
}

//Validate and rebuild files
func ValidateHotelFiles(hotelFiles interface{}, detailErrorCode map[string]interface{}) (string, map[string]interface{}) {
	var rebuildHotelFiles []interface{}
	if reflect.TypeOf(hotelFiles) == reflect.TypeOf("string") {
		if error := json.Unmarshal([]byte(hotelFiles.(string)), &rebuildHotelFiles); error != nil {
			return "", detailErrorCode
		} else {
			hotelFiles = rebuildHotelFiles
		}
	}
	tmpReflectFile := reflect.ValueOf(hotelFiles)
	var test []interface{}
	if tmpReflectFile.Type() != reflect.TypeOf(test) {
		return "", detailErrorCode
	}
	var hotelFilesInsert = ""
	var tmpHotelFilesInsert = ""
	for _, tmpFile := range hotelFiles.([]interface{}) {
		reflectFile := reflect.ValueOf(tmpFile)
		tmpHotelFilesInsert += "{"
		for subIndex, key := range reflectFile.MapKeys() {
			strKey := key.String()
			var value = ""
			if reflectFile.MapIndex(key).Interface() != nil {
				var objectValue = reflectFile.MapIndex(key).Interface()
				if reflect.TypeOf(objectValue) == reflect.TypeOf(test) {
					detailErrorCode["Files"] = conf.VARIABLE_IS_NOT_JSON
					return "", detailErrorCode
				}
				value = strings.Trim(objectValue.(string), " ")
			}

			tmpHotelFilesInsert += "\"" + strKey + "\": \"" + value + "\","
			if subIndex == (len(reflectFile.MapKeys()) - 1) {
				tmpHotelFilesInsert = strings.TrimRight(tmpHotelFilesInsert, ",")
			}
		}
		tmpHotelFilesInsert += "},"
	}

	if tmpHotelFilesInsert != "" {
		tmpHotelFilesInsert = strings.TrimRight(tmpHotelFilesInsert, ",")
		hotelFilesInsert += "[" + tmpHotelFilesInsert + "]"
	}
	return hotelFilesInsert, detailErrorCode
}

// Check Hotel Exists
// @Param hotelId int "id of hotel"
// @Param user *AdminAccounts
// @return bool
func CheckHotelExists(queryParams map[string]interface{}) (bool, Hotels) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("hotels.*").
		From("hotels").
		InnerJoin("hotel_administrators").On("hotel_administrators.hotel_id = hotels.hotel_id").
		InnerJoin("admin_accounts").On("admin_accounts.admin_account_id = hotel_administrators.admin_account_id").
		Where("hotels.hotel_id = " + fmt.Sprint(queryParams["hotel_id"])).
		And("admin_accounts.admin_account_id = " + fmt.Sprint(queryParams["admin_account_id"])).
		And("hotels.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).
		And("admin_accounts.deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).
		Limit(1).
		String()
	var hotel []Hotels
	num, err := o.Raw(sql).QueryRows(&hotel)
	if num > 0 && err == nil {
		return true, hotel[0]
	}
	return false, Hotels{}
}

// Check Hotel Exists And Room Exists
func CheckHotelAndRoom(queryParams map[string]interface{}) (bool, Hotels) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("hotels.*").
		From("hotels").
		InnerJoin("rooms").On("hotels.hotel_id = rooms.hotel_id").
		Where("rooms.room_id = " + fmt.Sprint(queryParams["room_id"])).
		And("rooms.deleted_at = 0").
		And("hotels.deleted_at = 0").
		String()
	var hotel []Hotels
	num, err := o.Raw(sql).QueryRows(&hotel)
	if num > 0 && err == nil {
		return true, hotel[0]
	}
	return false, Hotels{}
}

// Check Hotel Exists And Get Hotel
func FindHotel(queryParams map[string]interface{}) (bool, Hotels) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Hotels)).Filter("Id", queryParams["hotel_id"]).Filter("DeletedAt", conf.NOT_DELETED)
	if qs.Exist() {
		var hotel Hotels
		qs.One(&hotel)
		return true, hotel
	}
	return false, Hotels{}
}

// Soft Delete Hotel
// @Param hotelId int "id of hotel"
// @return error
func SoftDeleteHotel(hotelId int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(new(Hotels)).Filter("Id", hotelId).Update(orm.Params{
		"deleted_at": 1,
	})
	return err
}

// Add Hotel insert a new Hotel into database and returns
// last inserted Id on success.
func AddHotelWithTransaction(m *Hotels, stringHotelAdministrator string) (id int64, errorMsg string) {
	o := orm.NewOrm()
	//begin transaction
	var err error
	err = o.Begin()
	//insert room
	id, err = o.Insert(m)

	errorMsgInsert, errInsert := InsertHotelAdministratorInTransaction(o, stringHotelAdministrator, *m)
	if errInsert != nil || errorMsgInsert != "" {
		errorMsg = errorMsgInsert
		if errInsert != nil {
			errorMsg = err.Error()
		}
		err = o.Rollback()
		return
	}
	//end transaction
	if err != nil {
		errorMsg = err.Error()
		err = o.Rollback()
		return
	}
	err = o.Commit()
	return
}

// UpdateRooms updates Hotel by Id and returns error if
// the record to be updated doesn't exists
func UpdateHotelAdministratorWithTransaction(m *Hotels, stringHotelAdministrator string) (errorMsg string) {
	o := orm.NewOrm()
	//begin transaction
	var err error
	err = o.Begin()

	v := Hotels{Id: m.Id}

	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}

	errorMsg, err = InsertHotelAdministratorInTransaction(o, stringHotelAdministrator, *m)

	//end transaction
	if err != nil || errorMsg != "" {
		if err != nil {
			errorMsg = err.Error()
		}
		err = o.Rollback()
		return
	}
	err = o.Commit()
	return
}

func InsertHotelAdministratorInTransaction(o orm.Ormer, stringHotelAdministrator string, hotel Hotels) (errorMsg string, err error) {
	//insert room service
	return
}

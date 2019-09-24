package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"indetail/conf"
	"reflect"
	"strconv"

	//"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type FileImage struct {
	Type string `json:"Type" valid:"Required"`
	Url  string `json:"Url" valid:"Required"`
}

type Guests struct {
	Id              int       `orm:"column(guest_id);auto"`
	FullName        string    `orm:"column(full_name);size(120);null" json:"FullName" valid:"Required"`
	Status          int8      `orm:"column(status);null" json:"Status"`
	Email           string    `orm:"column(email);size(255);null" json:"Email"`
	LoginName       string    `orm:"column(login_name);size(255);null" json:"LoginName"`
	Password        string    `orm:"column(password);size(60)"`
	Phone           string    `orm:"column(phone);size(30);null" json:"Phone"`
	Address         string    `orm:"column(address);size(255);null" json:"Address"`
	PassportNumber  string    `orm:"column(passport_number);size(30);null" json:"PassportNumber" valid:"Required"`
	PassportExpired time.Time `orm:"column(passport_expired);type(date);null" json:"PassportExpired"`
	Nationality     string    `orm:"column(nationality);size(10);null" json:"Nationality"`
	BirthDay        time.Time `orm:"column(birth_day);type(date);null" json:"BirthDay"`
	Gender          int       `orm:"column(gender);null" json:"Gender" valid:"Required"`
	Occupation      int       `orm:"column(occupation);null" json:"Occupation" valid:"Required"`
	Files           string    `orm:"column(files);null" json:"Files"`
	DeletedAt       int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser     int       `orm:"column(created_user)"`
	UpdatedUser     int       `orm:"column(updated_user);null"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type GuestWithIsMain struct {
	Id              int       `orm:"column(guest_id);auto"`
	FullName        string    `orm:"column(full_name);size(120);null" json:"FullName" valid:"Required"`
	Status          int8      `orm:"column(status);null" json:"Status"`
	Email           string    `orm:"column(email);size(255);null" json:"Email"`
	LoginName       string    `orm:"column(login_name);size(255);null" json:"LoginName"`
	Password        string    `orm:"column(password);size(60)"`
	Phone           string    `orm:"column(phone);size(30);null" json:"Phone"`
	Address         string    `orm:"column(address);size(255);null" json:"Address"`
	PassportNumber  string    `orm:"column(passport_number);size(30);null" json:"PassportNumber" valid:"Required"`
	PassportExpired time.Time `orm:"column(passport_expired);type(date);null" json:"PassportExpired"`
	Nationality     string    `orm:"column(nationality);size(10);null" json:"Nationality"`
	BirthDay        time.Time `orm:"column(birth_day);type(date);null" json:"BirthDay"`
	Gender          int       `orm:"column(gender);null" json:"Gender" valid:"Required"`
	Occupation      int       `orm:"column(occupation);null" json:"Occupation" valid:"Required"`
	IsMainGuest     int8      `orm:"column(is_main_guest);null" json:"IsMainGuest"`
	Files           string    `orm:"column(files);null" json:"Files"`
	DeletedAt       int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser     int       `orm:"column(created_user)"`
	UpdatedUser     int       `orm:"column(updated_user);null"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type GuestWithBookingId struct {
	Id              int       `orm:"column(guest_id);auto"`
	BookingId       int       `orm:"column(booking_id);null" json:"booking_id"`
	FullName        string    `orm:"column(full_name);size(120);null" json:"FullName"`
	Status          int8      `orm:"column(status);null" json:"Status"`
	Email           string    `orm:"column(email);size(255);null" json:"Email"`
	LoginName       string    `orm:"column(login_name);size(255);null" json:"LoginName"`
	Password        string    `orm:"column(password);size(60)"`
	Phone           string    `orm:"column(phone);size(30);null" json:"Phone"`
	Address         string    `orm:"column(address);size(255);null" json:"Address"`
	PassportNumber  string    `orm:"column(passport_number);size(30);null" json:"PassportNumber"`
	PassportExpired time.Time `orm:"column(passport_expired);type(date);null" json:"PassportExpired"`
	Nationality     string    `orm:"column(nationality);size(10);null" json:"Nationality"`
	BirthDay        time.Time `orm:"column(birth_day);type(date);null" json:"BirthDay"`
	Gender          int       `orm:"column(gender);null" json:"Gender"`
	Occupation      int       `orm:"column(occupation);null" json:"Occupation"`
	IsMainGuest     int8      `orm:"column(is_main_guest);null" json:"IsMainGuest"`
	Files           string    `orm:"column(files);null" json:"Files"`
	DeletedAt       int8      `orm:"column(deleted_at);null" json:"deleted_at,omitempty"`
	CreatedUser     int       `orm:"column(created_user)"`
	UpdatedUser     int       `orm:"column(updated_user);null"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt       time.Time `orm:"column(updated_at);type(timestamp);auto_now"`
}

type AddGuests struct {
	BookingId       string    `orm:"column(booking_id);null"`
	FullName        string    `orm:"column(full_name);size(120);null"`
	Status          int8      `orm:"column(status);null"`
	Email           string    `orm:"column(email);size(255);null"`
	Phone           string    `orm:"column(phone);size(30);null"`
	Address         string    `orm:"column(address);size(255);null"`
	PassportNumber  string    `orm:"column(passport_number);size(30);null"`
	PassportExpired time.Time `orm:"column(passport_expired);type(date);null" json:"PassportExpired"`
	Nationality     string    `orm:"column(nationality);size(10);null" json:"Nationality"`
	BirthDay        time.Time `orm:"column(birth_day);type(date);null"`
	Gender          int       `orm:"column(gender);null"`
	Occupation      int       `orm:"column(occupation);null"`
	Files           string    `orm:"column(files);null"`
}

type UpdateGuests struct {
	FullName        string    `orm:"column(full_name);size(120);null"`
	Status          int8      `orm:"column(status);null"`
	Email           string    `orm:"column(email);size(255);null"`
	Phone           string    `orm:"column(phone);size(30);null"`
	Address         string    `orm:"column(address);size(255);null"`
	PassportNumber  string    `orm:"column(passport_number);size(30);null"`
	PassportExpired time.Time `orm:"column(passport_expired);type(date);null" json:"PassportExpired"`
	Nationality     string    `orm:"column(nationality);size(10);null" json:"Nationality"`
	BirthDay        time.Time `orm:"column(birth_day);type(date);null"`
	Gender          int       `orm:"column(gender);null"`
	Occupation      int       `orm:"column(occupation);null"`
	Files           string    `orm:"column(files);null"`
}

type GuestTotal struct {
	Total      int64 `orm:"column(total);null"`
	CountGuest int64 `orm:"column(guest_count);null"`
}

type BookingGuestInsert struct {
	Id int64 `orm:"column(guest_id)"`
}

func (t *Guests) TableName() string {
	return "guests"
}

func init() {
	orm.RegisterModel(new(Guests))
}

// AddGuests insert a new Guests into database and returns
// last inserted Id on success.
func AddGuest(m *Guests) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// AddGuests insert a new Guests into database and returns
// last inserted Id on success.
func AddGuestTransaction(m *Guests, o orm.Ormer) (id int64, err error) {
	id, err = o.Insert(m)
	return
}

// GetGuestById retrieves Guests by Id. Returns error if
// Id doesn't exist
func GetGuestById(id int) (v *Guests, err error) {
	o := orm.NewOrm()
	v = &Guests{Id: id, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "Id", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetGuestByConditions retrieves Guests by Id. Returns error if
// Id doesn't exist
func GetGuestByListBookingId(listBookingId string, guestId int) (guestWithBookingId []GuestWithBookingId, err error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("booking_guests.booking_id, booking_guests.is_main_guest, guests.*").
		From("guests").
		InnerJoin("booking_guests").On("guests.guest_id = booking_guests.guest_id").
		Where("booking_guests.booking_id ").In(listBookingId).
		And("guests.deleted_at = " + fmt.Sprint(conf.NOT_DELETED))

	if guestId > 0 {
		sql = qd.And("guests.guest_id = " + fmt.Sprint(guestId))
	}
	sql = qd.GroupBy("booking_guests.booking_id, guests.guest_id")

	sqlQuery := sql.String()
	if _, err := o.Raw(sqlQuery).QueryRows(&guestWithBookingId); err != nil {
		return nil, err
	}
	return guestWithBookingId, nil
}

func GetGuestsByEmail(email string) (v *Guests, err error) {
	o := orm.NewOrm()
	v = &Guests{Email: email, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "Email", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGuests retrieves all Guests matches certain condition. Returns empty list if
// no records exist
func GetAllGuests(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Guests))
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

	var l []Guests
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

// UpdateGuests updates Guests by Id and returns error if
// the record to be updated doesn't exist
func UpdateGuestById(m *Guests) (err error) {
	o := orm.NewOrm()
	v := Guests{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func GetGuestByEmail(email string) (v *Guests, err error) {
	o := orm.NewOrm()
	v = &Guests{Email: email, DeletedAt: conf.NOT_DELETED}
	if err = o.Read(v, "Email", "DeletedAt"); err == nil {
		return v, nil
	}
	return nil, err
}

func UpdatePasswordGuestAndDeleteToken(guest *Guests) bool {
	o := orm.NewOrm()
	o.Begin()
	update := UpdateGuestById(guest)

	destroy := DeletePasswordResets(o, guest.Id, conf.TYPE_GUEST)

	destroyToken := DestroyAllTokenWithAccountId(o, guest.Id, conf.TYPE_GUEST)
	if update != nil || destroy != nil || destroyToken != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

func ChangePasswordGuestAndLogout(guest *Guests) bool {
	o := orm.NewOrm()
	o.Begin()
	update := UpdateGuestById(guest)

	destroyToken := DestroyAllTokenWithAccountId(o, guest.Id, conf.TYPE_GUEST)
	if update != nil || destroyToken != nil {
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

// DeleteGuests deletes Guests by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGuests(id int) (err error) {
	o := orm.NewOrm()
	v := Guests{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Guests{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//Validate and rebuild files
func ValidateGuestFiles(guestFiles interface{}, detailErrorCode map[string]interface{}) (string, map[string]interface{}) {
	var rebuildFiles []interface{}
	if reflect.TypeOf(guestFiles) == reflect.TypeOf("string") {
		if error := json.Unmarshal([]byte(guestFiles.(string)), &rebuildFiles); error != nil {
			return "", detailErrorCode
		} else {
			guestFiles = rebuildFiles
		}
	}

	tmpReflectFile := reflect.ValueOf(guestFiles)
	var test []interface{}
	if tmpReflectFile.Type() != reflect.TypeOf(test) {
		return "", detailErrorCode
	}
	var filesInsert = ""
	var tmpFilesInsert = ""
	for _, tmpFile := range guestFiles.([]interface{}) {
		reflectFile := reflect.ValueOf(tmpFile)
		tmpFilesInsert += "{"
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
			tmpFilesInsert += "\"" + strKey + "\": \"" + value + "\","
			if subIndex == (len(reflectFile.MapKeys()) - 1) {
				tmpFilesInsert = strings.TrimRight(tmpFilesInsert, ",")
			}
		}
		tmpFilesInsert += "},"
	}

	if tmpFilesInsert != "" {
		tmpFilesInsert = strings.TrimRight(tmpFilesInsert, ",")
		filesInsert += "[" + tmpFilesInsert + "]"
	}
	return filesInsert, detailErrorCode
}

// DeleteGuests deletes Guests by Id and returns error if
// the record to be deleted doesn't exist
func TotalGuestByBookingId(bookingId int) (totalGuest int64, isAvailable bool, err error) {
	o := orm.NewOrm()
	qd, _ := orm.NewQueryBuilder("mysql")
	sql := qd.Select("COUNT(guests.guest_id) as total, bookings.guest_count as guest_count").
		From("guests").
		InnerJoin("booking_guests").On("guests.guest_id = booking_guests.guest_id").
		InnerJoin("bookings").On("booking_guests.booking_id = bookings.booking_id").
		Where("guests.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("booking_guests.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).
		And("bookings.booking_id = " + fmt.Sprint(bookingId))

	query := sql.String()
	var totalGuestStruct []GuestTotal
	_, err = o.Raw(query).QueryRows(&totalGuestStruct)
	if err != nil {
		return -2, false, err
	}
	if totalGuestStruct[0].Total >= totalGuestStruct[0].CountGuest {
		return -1, false, err
	}

	return totalGuestStruct[0].Total, true, err
}

// GetAllGuestCondition retrieves all Guests matches certain condition. Returns empty list if
// no records exist
func GetAllGuestCondition(query map[string]string, fields string, sortby []string, order []string,
	offset int64, limit int64) (listGuest []Guests, msgErr string) {
	o := orm.NewOrm()
	sql, sqlWhere := RebuildQueryListGuest(query, fields)
	// order by:
	sortFields, sortErr := MakeOrderForSqlQuery(sortby, order)
	if sortErr != nil {
		return nil, sortErr.Error()
	}
	sqlWhereString := sqlWhere.String()
	if sortFields != "" {
		sqlWhereString += " Order By " + sortFields
	}

	if limit > 0 {
		sqlWhereString += " LIMIT " + fmt.Sprint(limit)
	}

	if offset > 0 {
		sqlWhereString += " OFFSET " + fmt.Sprint(offset)
	}
	sqlBaseString := sql.String()
	sqlBaseString = sqlBaseString + " " + sqlWhereString
	_, err := o.Raw(sqlBaseString).QueryRows(&listGuest)
	if err != nil {
		return nil, "Get list failure"
	}
	return listGuest, ""
}

func RebuildQueryListGuest(query map[string]string, fields string) (sql orm.QueryBuilder, sqlWhere orm.QueryBuilder) {
	qd, _ := orm.NewQueryBuilder("mysql")
	sql = qd.Select(fields).From("guests")

	qdWhere, _ := orm.NewQueryBuilder("mysql")
	sqlWhere = qdWhere.Where("guests.deleted_at = " + fmt.Sprint(conf.NOT_DELETED))

	if query["FullName"] != "" {
		sqlWhere = qdWhere.And("guests.full_name LIKE '%" + query["FullName"] + "%'")
	}
	if query["Status"] != "" {
		sqlWhere = qdWhere.And("guests.status = " + query["Status"])
	}
	if query["AgeGte"] != "" {
		sqlWhere = qdWhere.And("guests.birth_day >= '" + query["AgeGte"] + "'")
	}
	if query["AgeLte"] != "" {
		sqlWhere = qdWhere.And("guests.birth_day <= '" + query["AgeLte"] + "'")
	}
	if query["GuestId"] != "" {
		sqlWhere = qdWhere.And("guests.guest_id").In(query["GuestId"])
	}

	if query["BookingId"] != "" {
		sql = qd.InnerJoin("booking_guests").On("guests.guest_id = booking_guests.guest_id").
			InnerJoin("bookings").On("booking_guests.booking_id = bookings.booking_id")
		sqlWhere = qdWhere.And("bookings.deleted_at = " + fmt.Sprint(conf.NOT_DELETED)).And("bookings.booking_id = " + query["BookingId"])
	}
	return sql, sqlWhere
}

func CountTotalRecord(query map[string]string, fields string) (totalRecord int64, err error) {
	o := orm.NewOrm()
	sql, sqlWhere := RebuildQueryListGuest(query, fields)
	sqlBaseString := sql.String()
	sqlWhereString := sqlWhere.String()
	sqlBaseString = sqlBaseString + " " + sqlWhereString

	var listGuest []Guests
	totalRecord, err = o.Raw(sqlBaseString).QueryRows(&listGuest)
	if err != nil {
		return 0, err
	}
	return totalRecord, nil
}

// Add Guest With Transaction
// @Param user Guests
// @return err error
func AddGuestWithTransaction(guest *Guests, bookingId int) string {
	o := orm.NewOrm()
	o.Begin()

	//create guest
	guestId, errSaveGuest := AddGuestTransaction(guest, o)
	if errSaveGuest != nil {
		o.Rollback()
		return errSaveGuest.Error()
	}

	//if booking id let than 0 that mean create guest without booking guest
	if bookingId < 0 {
		o.Commit()
		return ""
	}

	//create booking guest
	var bookingGuestStruct []BookingGuests
	var bookingGuest BookingGuests
	bookingGuest.BookingId = bookingId
	bookingGuest.GuestId = int(guestId)
	bookingGuest.IsMainGuest = 0
	bookingGuestStruct = append(bookingGuestStruct, bookingGuest)

	//calculate guest
	total, isAvailable, totalErr := TotalGuestByBookingId(bookingId)
	if !isAvailable || totalErr != nil {
		o.Rollback()
		return "Can't add more guest to booking"
	}

	if total == 0 {
		bookingGuestStruct[0].IsMainGuest = 1
	}

	_, bookingGuestErr := AddBookingGuestTransaction(&bookingGuestStruct, o)
	if bookingGuestErr != nil {
		o.Rollback()
		return bookingGuestErr.Error()
	}

	o.Commit()
	return ""
}

// CheckExitPassport
func CheckExitPassport(passport string, guestId string) bool {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Guests))
	// query k=v
	var query = make(map[string]string)
	query["PassportNumber"] = passport
	if len(guestId) > 0 {
		query["Id.isnotin"] = guestId
	}
	qs = RebuildConditions(qs, query)
	if totalRecord, _ := qs.Count(); totalRecord > 0 {
		return true
	} else {
		return false
	}
}

// Permission Get User App By Login Name
// Return user *libs.GlobalUsers, error
func PermissionGetGuestByLoginName(loginName string) (GlobalUsers, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("guest_id as admin_account_id, full_name, password, email").
		From("guests").
		Where("email = '" + loginName + "'").
		And("deleted_at = " + strconv.Itoa(conf.NOT_DELETED)).String()

	var accounts []GlobalUsers
	_, err := o.Raw(sql).QueryRows(&accounts)
	return accounts[0], err
}

package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"indetail/conf"
	"reflect"
	"strings"
)

type Base struct {
}

func FilterResultByField(fields []string, container []interface{}) (result []interface{}, err error) {
	if len(fields) == 0 {
		for _, v := range container {
			result = append(result, v)
		}
	} else {
		// trim unused fields
		for _, v := range container {
			m := make(map[string]interface{})
			val := reflect.ValueOf(v)
			for _, fname := range fields {
				m[fname] = val.FieldByName(fname).Interface()
			}
			result = append(result, m)
		}
	}
	return result, nil
}

func MakeOrderForQuery(qs orm.QuerySeter, sortby []string, order []string) (qsResult orm.QuerySeter, err error) {
	var sortFields []string
	if len(sortby) != 0 && len(sortby) == len(order) {
		if len(sortby) != 0 {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if strings.ToLower(order[i]) == "desc" {
					orderby = "-" + v
				} else if strings.ToLower(order[i]) != "asc" {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				} else {
					orderby = v
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else {
			if len(order) != 0 {
				return nil, errors.New("Error: unused 'order' fields")
			}
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)

	return qs, nil
}

func MakeOrderForSqlQuery(sortby []string, order []string) (sortFields string, err error) {
	if len(sortby) != 0 && len(sortby) == len(order) {
		if len(sortby) != 0 {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				lowerString := strings.ToLower(order[i])
				if lowerString != "asc" && lowerString != "desc" {
					return "", errors.New("Error: Invalid order. Must be either [asc|desc]")
				} else {
					sortFields += v + " " + strings.ToLower(order[i]) + ","
				}
			}
		} else {
			if len(order) != 0 {
				return "", errors.New("Error: unused 'order' fields")
			}
		}
	} else {
		if len(order) != 0 {
			return "", errors.New("Error: unused 'order' fields")
		}
	}
	sortFields = strings.TrimRight(sortFields, ",")
	return sortFields, nil
}

func rebuildValueForInQuery(paramValue string) []interface{} {
	arrValue := strings.Split(paramValue, ",")
	var rebuildValue []interface{}
	for _, value := range arrValue {
		rebuildValue = append(rebuildValue, value)
	}

	return rebuildValue
}

//rebuild conditions
//if deleted_at not existed in query then set deleted_at equal 0
func RebuildConditions(seter orm.QuerySeter, query map[string]string) orm.QuerySeter {
	var fieldGroupBy string
	for key, value := range query {
		// rewrite dot-notation to Object__Attribute
		key = strings.Replace(key, ".", "__", -1)
		if strings.Contains(key, "isnull") {
			seter = seter.Filter(key, (value == "true" || value == "1"))
		} else if strings.Contains(key, "groupby") {
			fieldGroupBy = value
		} else if strings.Contains(key, "isnotin") {
			rebuildValue := rebuildValueForInQuery(value)
			arrParam := strings.Split(key, "__isnotin")
			seter = seter.Exclude(arrParam[0]+"__in", rebuildValue)
		} else if strings.Contains(key, "maxprice") {
			if value == "true" || value == "1" {
				seter = seter.FilterRaw("Price", "IN (SELECT MAX(price) FROM rooms GROUP BY hotel_id)")
			} else {
				seter = seter.FilterRaw("Price", "IN (SELECT MIN(price) FROM rooms GROUP BY hotel_id)")
			}
		} else if strings.Contains(key, "in") {
			rebuildValue := rebuildValueForInQuery(value)
			seter = seter.Filter(key, rebuildValue)
		} else {
			seter = seter.Filter(key, value)
		}
	}

	if query["deleted_at"] == "" && query["DeletedAt"] == "" {
		seter = seter.Filter("deleted_at", conf.NOT_DELETED)
	}

	//group by
	if fieldGroupBy != "" {
		seter = seter.GroupBy(fieldGroupBy)
	}

	return seter
}

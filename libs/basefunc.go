package libs

import (
	"encoding/json"
	"fmt"
	"indetail/conf"
)

type ResponseJson struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
	DetailErrorCode map[string]interface{} `json:"detailErrorCode"`
	Success string `json:"success"`
}

//make json format
func ResultJson(data interface{}, code string, message string, detailErrorCode map[string]interface{}) interface{} {
	var result ResponseJson
	result.Data = data
	result.Message = message
	result.DetailErrorCode = detailErrorCode
	result.Code = code
	if code != "200" {
		result.Success = "False"
	} else {
		result.Success = "True"
	}
	return result
}

type PagingResponseJson struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
	DetailErrorCode map[string]interface{} `json:"detailErrorCode"`
	Success string `json:"success"`
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
	Page int64 `json:"page"`
}

//make json format
func ResultPagingJson(data interface{}, code string, message string, limit int64, page int64, total int64) interface{} {
	var result PagingResponseJson
	result.Limit = limit
	result.Page = page
	result.Total = total
	result.Data = data
	result.Message = message
	result.Code = code
	if code != "200" {
		result.Success = "False"
	} else {
		result.Success = "True"
	}
	return result
}

func ValidateStatus(stringJsonStatus string, statusCode int)(jsonResult interface{}, isPass bool) {
	isPass = false
	var mapListStatus = make(map[int]string)
	jsonResult, mapListStatus = CoreTypeValidate(stringJsonStatus)
	if jsonResult != nil {
		return
	}
	if _, exists := mapListStatus[statusCode]; !exists {
		jsonResult = ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Status not found!", map[string]interface{}{"Status": conf.VARIABLE_OUT_OF_RANGE})
		return
	}

	return nil, true
}

func ValidateHotelType(stringJsonStatus string, statusCode int)(jsonResult interface{}, isPass bool) {
	isPass = false
	var mapListStatus = make(map[int]string)
	jsonResult, mapListStatus = CoreTypeValidate(stringJsonStatus)
	if jsonResult != nil {
		return
	}
	if _, exists := mapListStatus[statusCode]; !exists {
		jsonResult = ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Hotel type not found!", map[string]interface{}{"HotelType": conf.VARIABLE_OUT_OF_RANGE})
		return
	}

	return nil, true
}

func ValidateHotelSubType(stringJsonStatus string, statusCode int)(jsonResult interface{}, isPass bool) {
	isPass = false
	var mapListStatus = make(map[int]string)
	jsonResult, mapListStatus = CoreTypeValidate(stringJsonStatus)
	if jsonResult != nil {
		return
	}
	if _, exists := mapListStatus[statusCode]; !exists {
		jsonResult = ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Hotel sub type not found!", map[string]interface{}{"HotelSubType": conf.VARIABLE_OUT_OF_RANGE})
		return
	}

	return nil, true
}

func CoreTypeValidate(stringJsonStatus string)(interface{}, map[int]string) {
	var listStatus []string
	if error := json.Unmarshal([]byte(stringJsonStatus), &listStatus); error != nil {
		jsonResult := ResultJson(nil, fmt.Sprint(conf.ERROR_STATUS), "Missing config", nil)
		return jsonResult, nil
	}
	var mapListStatus = make(map[int]string)
	for key, value := range listStatus {
		mapListStatus[key] = value
	}
	return nil, mapListStatus
}

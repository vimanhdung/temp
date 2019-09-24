package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/libs"
)

// BaseController operations for Base function
type BaseController struct {
	beego.Controller
}

// URLMapping ...
func (c *BaseController) URLMapping() {
	c.Mapping("GetConfig", c.GetConfig)
}

// GetConfig ...
// @Title Get Config
// @Description get config information
// @Success 200 {object} models.SwaggerGetConfig Success
// @router /configs/ [get]
func (c *BaseController) GetConfig() {
	var resultData = make(map[string]interface{})
	resultData["prefecture"] = conf.GetPrefecture()
	//OTA
	resultData["otaName"] = conf.OTA_NAME_LIST
	//guest
	resultData["otaName"] = conf.GUEST_OCCUPATION

	c.Data["json"] = libs.ResultJson(resultData, fmt.Sprint(conf.SUCCESS_STATUS), "Success", nil)
	c.ServeJSON()
}

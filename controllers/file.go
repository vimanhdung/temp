package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"indetail/conf"
	"indetail/libs"
	"time"
)

// FileController operations for File
type FileController struct {
	libs.Middleware
}

var images = []string{"image/jpeg", "image/png", "image/gif", "image/jpg"}

/**
** Validate Image
 */
func ValidateImage(image string) bool {
	for _, a := range images {
		if a == image {
			return true
		}
	}
	return false
}

// Post ...
// @Title Upload File
// @Description Upload file
// @Param Authorization header string true "Bearer token"
// @Param image formData file true "The image to upload (jpg,png,gif)"
// @Success 201 {object} libs.ResponseJson Upload success
// @Failure 403 104 : Invalid Token <br> 302 : Image not found <br> 402 : Image is not properly formatted <br> 303 : Save file false
// @router /file/upload [post]
func (c *FileController) UploadImage() {
	// Check Permission
	/*if !c.PermissionDenied(conf.COMMON_UPLOAD_FILE) {
		return
	}*/
	file, header, _ := c.GetFile("image")
	if file == nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"File not found",
			map[string]interface{}{"Image": conf.RECORD_NOT_FOUND},
		)
		c.ServeJSON()
		return
	}

	if !ValidateImage(header.Header.Get("Content-Type")) {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Image is not properly formatted",
			map[string]interface{}{"Image": conf.FILE_TYPE_WRONG},
		)
		c.ServeJSON()
		return
	}
	now := time.Now().Format("20060102150405")
	filename := header.Filename
	if header.Header.Get("Filename") != "" {
		filename = header.Header.Get("Filename")
	}
	path := "/storage/" + now + "_" + filename
	err := c.SaveToFile("image", "."+path)
	if err != nil {
		c.Data["json"] = libs.ResultJson(
			nil,
			fmt.Sprint(conf.ERROR_STATUS),
			"Save file false",
			map[string]interface{}{"Storage": conf.SAVE_FAILURES},
		)
		c.ServeJSON()
		return
	}
	c.Data["json"] = libs.ResultJson(
		map[string]interface{}{"path": beego.AppConfig.String("baseServer") + path},
		fmt.Sprint(conf.SUCCESS_STATUS),
		"Upload success",
		nil,
	)
	c.ServeJSON()
	return
}

package main

import (
	"github.com/astaxie/beego/orm"
	_ "indetail/routers"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql" // import your required driver
)

func init()  {
	orm.RegisterDataBase(beego.AppConfig.String("aliasName"), beego.AppConfig.String("driverName"), beego.AppConfig.String("databaseSource"))
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

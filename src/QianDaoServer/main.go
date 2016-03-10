package main

import (
	"QianDaoServer/controllers"
	"QianDaoServer/models"
	// "QianDaoServer/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	o := orm.NewOrm()
	o.Using("default")
	orm.RunSyncdb("default", false, true)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/logs", &controllers.LogsController{})
	beego.Router("/add", &controllers.AddController{})
	beego.Router("/error", &controllers.ErrorController{})

	beego.Run()
}

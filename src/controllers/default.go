package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	c.Data["Website"] = "smartxupt.com"
	c.Data["Email"] = "admin@smartxupt.com"
	c.TplName = "index.tpl"
}

package controllers

import (
	"QianDaoServer/models"
	"github.com/astaxie/beego"
	// "strings"
)

type UserListController struct {
	beego.Controller
}

func (this *UserListController) Get() {

	op := this.Input().Get("op")
	if op == "del" {
		_, err := models.DeleteUser(this.Input().Get("Uid"), this.Input().Get("Mac"))

		if !err {
			this.Redirect("/static/nouser.html", 302)
		}
		this.Redirect("/userlist", 302)
		return
	}
	this.TplName = "user_view.html"
	this.Data["IsLogin"] = checkAcount(this.Ctx)
	var err error
	this.Data["Users"], err = models.GetAllUser()
	if err != nil {
		beego.Error(err)
	}
}

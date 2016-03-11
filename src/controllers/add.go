package controllers

import (
	"models"
	// "fmt"
	"github.com/astaxie/beego"
)

type AddController struct {
	beego.Controller
}

func (this *AddController) Get() {
	this.TplName = "add.html"
}

func (this *AddController) Post() {
	if models.GetUser(this.Input().Get("mac")) {
		models.AddUsers(this.Input().Get("mac"), this.Input().Get("uid"), this.Input().Get("name"))
		this.Redirect("/", 301)
	} else {
		this.Redirect("/error", 301)
	}
}

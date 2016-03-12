package controllers

import (
	"QianDaoServer/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

var Today string
var flag bool = false

type LogsController struct {
	beego.Controller
}

// var UserTime map[string]int64 = make(map[string]int64, 0)

func (this *LogsController) Get() {
	//var UserTime map[string]int64 = make(map[string]int64, 0)
	if !flag {
		Today = models.Today()
		flag = true
	}
	a, _ := strconv.ParseUint(time.Now().String()[14:16], 10, 64)
	if a == 0 {
		fmt.Println("a:", a)
		Today = models.Today()
	}
	this.TplName = "logs.html"
	mac := this.Input().Get("mac")
	models.UpdateUser(mac)
	models.UpdateDayLog(mac, Today)
	models.UpdateLogs(mac)
	//modles.UpdateLogs(mac)
	/*
		a, _ := strconv.ParseUint(time.Now().String()[17:19], 10, 64)
		if a%5 == 0 {
			fmt.Println("a:", a)
			UserTime[mac]++
		}
		if a%10 == 0 {
			fmt.Println(mac, UserTime[mac])
			models.UpdateUser(mac)
		}
	*/
}

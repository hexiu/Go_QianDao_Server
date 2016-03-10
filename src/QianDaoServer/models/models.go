package models

import (
	// "fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	_MYSQL_DRIVER            = "mysql"
	_LEAST_TIME      float64 = 300.0
	_DATABASE_UNAME          = "root"
//	_DATABASE_PASSWD         = "zypc2016"
		_DATABASE_PASSWD         = "axiu"
	_DATABASE_NAME = "QianDao"
)

type User struct {
	Id      int64
	Uid     string
	Name    string
	Mac     string
	AllTime int64
}

type Daylog struct {
	Id      int64
	Mac     string
	Date    string
	DayTime int64
}

type Logs struct {
	Id    int64
	Mac   string
	Date  string
	Start time.Time `orm:"index"`
	End   time.Time `orm:"index"`
}

func RegisterDB() {

	//注册模型
	orm.RegisterModel(new(User), new(Daylog), new(Logs))
	orm.RegisterDriver(_MYSQL_DRIVER, orm.DRMySQL)
	orm.RegisterDataBase("default", _MYSQL_DRIVER, _DATABASE_UNAME+":"+_DATABASE_PASSWD+"@"+"/"+_DATABASE_NAME+"?charset=utf8&loc=Asia%2FShanghai")
	//	orm.RegisterDataBase("default", _MYSQL_DRIVER, _DATABASE_UNAME+":"+_DATABASE_PASSWD+"@"+"/"+_DATABASE_NAME+"?charset=utf8&loc=Asia%2FShanghai")
}

func GetUser(mac string) bool {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("mac", mac).One(user)
	if err == nil {
		//用户存在
		return true
	} else {
		//用户不存在
		return false
	}
}

//添加用户信息
func AddUsers(mac, uid, name string) error {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("mac", mac).One(user)

	if err == nil {
		//存在,就是更新
		user.Uid = uid
		user.Name = name
		_, err = o.Update(user)
		if err != nil {
			return err
		}
	} else {
		//不存在，插入
		o := orm.NewOrm()
		user = &User{
			Uid:     uid,
			Name:    name,
			Mac:     mac,
			AllTime: 0,
		}

		_, err := o.Insert(user)
		if err != nil {
			return err
		}
	}
	return nil
}

//更新用户表信息
func UpdateUser(mac string) error {

	o := orm.NewOrm()

	user := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("mac", mac).One(user)
	if err == nil {
		user.AllTime++
		_, err = o.Update(user)
		if err != nil {
			return err
		}
	} else {
		err = AddUsers(mac, "", "")
		if err != nil {
			return err

		}

	}
	return nil

}

func AddDayLog(mac, today string) error {
	o := orm.NewOrm()
	daylog := &Daylog{
		Mac:     mac,
		Date:    today,
		DayTime: 0,
	}

	_, err := o.Insert(daylog)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDayLog(mac, today string) error {
	o := orm.NewOrm()
	daylog := new(Daylog)
	qs := o.QueryTable("daylog")
	err := qs.Filter("mac", mac).Filter("date", today).One(daylog)
	if err == nil {
		//存在表单
		daylog.DayTime++
		_, err = o.Update(daylog)
		if err != nil {
			return err
		}
	} else {
		//不存在表单
		err = AddDayLog(mac, today)
		if err != nil {
			return err

		}
	}

	return nil
}

func AddLogs(mac string) error {

	o := orm.NewOrm()
	logs := &Logs{
		Mac:   mac,
		Date:  Today(),
		Start: time.Now(),
		End:   time.Now(),
	}
	_, err := o.Insert(logs)
	if err != nil {
		return err
	}
	return nil

}

func UpdateLogs(mac string) error {
	o := orm.NewOrm()
	logs := new(Logs)
	qs := o.QueryTable("logs")

	jud, _ := qs.Filter("mac", mac).Count()
	//err := qs.Filter("uid", uid).Filter("today", today()).One(logs)
	if jud != 0 {
		//之前存在记录
		logss := make([]*Logs, 0)
		qs.Filter("mac", mac).OrderBy("-End").All(&logss)
		if len(logss) == 0 {
			return nil
		} else {
			lastEnd := logss[0].End
			// fmt.Println(lastEnd)
			// fmt.Println(time.Now().Sub(lastEnd).Seconds())
			jud := time.Now().Sub(lastEnd).Seconds() > _LEAST_TIME
			// fmt.Println(jud)
			if jud {
				AddLogs(mac)
			} else {
				err := qs.Filter("mac", mac).Filter("End", lastEnd).One(logs)
				if err == nil {
					logs.End = time.Now()
					_, err := o.Update(logs)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}

		return nil
	} else {
		//之前不存在，添加记录
		err := AddLogs(mac)
		if err != nil {
			return err
		}

	}
	return nil
}

func Today() string {
	today := time.Now().String()[0:10]
	return today
}

//操作用户数据的包
package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

//一些静态变量值，后期写入配置文件
const (
	_MYSQL_DRIVER           = "mysql"
	_LEAST_TIME     float64 = 300.0
	_DATABASE_UNAME         = "root"
	//	_DATABASE_PASSWD         = "zypc2016"
	_DATABASE_PASSWD = "axiu"
	_DATABASE_NAME   = "QianDao"
)

//用户信息表单
type User struct {
	Id      int64
	Uid     int64
	Name    string
	Mac     string
	AllTime int64
}

//用户日在线表单
type Daylog struct {
	Id      int64
	Mac     string
	Date    string
	DayTime int64
}

//用户在线时间段表单
type Logs struct {
	Id    int64
	Mac   string
	Date  string
	Start time.Time `orm:"index"`
	End   time.Time `orm:"index"`
}

//注册数据库
func RegisterDB() {

	//注册模型
	orm.RegisterModel(new(User), new(Daylog), new(Logs))
	orm.RegisterDriver(_MYSQL_DRIVER, orm.DRMySQL)
	orm.RegisterDataBase("default", _MYSQL_DRIVER, _DATABASE_UNAME+":"+_DATABASE_PASSWD+"@"+"/"+_DATABASE_NAME+"?charset=utf8&loc=Asia%2FShanghai")
	//	orm.RegisterDataBase("default", _MYSQL_DRIVER, _DATABASE_UNAME+":"+_DATABASE_PASSWD+"@"+"/"+_DATABASE_NAME+"?charset=utf8&loc=Asia%2FShanghai")
}

//获取User，判断用户是否存在
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

//获取所有用户的信息
func GetAllUser() ([]*User, error) {
	o := orm.NewOrm()
	users := make([]*User, 0)
	qs := o.QueryTable("user")
	_, err := qs.All(&users)
	return users, err
}

//添加用户信息到数据库（实则更新，只要收到客户端信息，服务器便会插入数据，防止遗漏数据）
// 只有安装客户端之后，才可以添加用户信息
func AddUsers(mac, cid, name string) error {
	uid, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	err = qs.Filter("mac", mac).One(user)

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
			Mac:     mac,
			AllTime: 0,
		}

		_, err := o.Insert(user)
		if err != nil {
			return err
		}

	}

	//下面的代码是允许用户注册。
	/*
		 else {
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

	*/
	return nil
}

//更新用户表信息（记录总的在线时长）
func UpdateUser(mac string) error {
	if TimeOut() {
		return nil
	}

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
		err = AddUsers(mac, "123", "")
		if err != nil {
			return err

		}

	}
	return nil

}

//添加Daylog表项，接收mac地址
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

//更新Daylog表单，接收mac地址和当天日期
func UpdateDayLog(mac, today string) error {

	if TimeOut() {
		return nil
	}
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

//添加Logs表项
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

//更新Logs 表单，记录用户在线的时间段
func UpdateLogs(mac string) error {

	if TimeOut() {
		return nil
	}
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

//获得今天的日期
func Today() string {
	today := time.Now().String()[0:10]
	return today
}

//删除用户，如果删除成功返回True，否则返回False
func DeleteUser(cid, mac string) (error, bool) {
	uid, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return err, false
	}
	o := orm.NewOrm()
	fmt.Println(uid, mac)
	qs := o.QueryTable("user")
	user := new(User)
	qs.Filter("mac", mac).One(user)
	fmt.Println(user.Id)
	users := &User{Id: user.Id}
	_, err = o.Delete(users)

	daylog := make([]*Daylog, 0)
	qs = o.QueryTable("daylog")
	_, err = qs.Filter("mac", mac).All(&daylog)
	for _, value := range daylog {
		dayLog := &Daylog{Id: value.Id}
		_, err = o.Delete(dayLog)
	}

	logs := make([]*Logs, 0)
	qs = o.QueryTable("logs")
	_, err = qs.Filter("mac", mac).All(&logs)
	for _, value := range logs {
		Log := &Logs{Id: value.Id}
		_, err = o.Delete(Log)
	}

	/*
		fmt.Println(err.Error())
		logs := &Logs{Mac: mac}
		_, err = o.Delete(logs)
		daylog := &Daylog{Mac: mac}
		_, err = o.Delete(daylog)
	*/
	return err, true
}

func TimeOut() bool {

	now, _ := strconv.ParseInt(time.Now().String()[11:13], 10, 64)
	if now >= 23 || now <= 6 {
		return false
	}
	return true
}

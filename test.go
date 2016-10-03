package main

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"rwlock"
	_ "rwlock/driver/mysql"
)

func LogLevel(level string) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q ", level, err, l)
	}

	return l
}
func createTable(dbname string) {
	force := false //不强制建数据库
	verbose := true

	err := orm.RunSyncdb("default", force, verbose) //建表
	if err != nil {
		log.Errorf("err: %s", err)
	}
}

func createdb(conn string, dbname string) error {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Errorf("Cannot open database(%s), err: %v", conn, err)
		return err
	}
	usedbSql := fmt.Sprintf("use %s", dbname)

	_, err = db.Exec(usedbSql)
	if err == nil {
		log.Infof("DB(%s) already exists, no need to create it.", dbname)
		return err
	}

	defer db.Close()

	dbCreateSql := fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci", dbname)

	_, err = db.Exec(dbCreateSql)
	if err == nil {
		log.Infof("Succeed to create db(%s)", dbname)
	}
	usedbSql = fmt.Sprintf("use %s", dbname)

	_, err = db.Exec(usedbSql)
	if err == nil {
		log.Infof("DB(%s) already exists, no need to create it.", dbname)
		return err
	}
	return nil

}
func init() {

	conn := "root:00010001@tcp(127.0.0.1:3306)/mysql?charset=utf8"
	createdb(conn, "lock")
}
func main() {
	// Initialize database
	orm.RegisterDriver("mysql", orm.DRMySQL)
	conn := "root:00010001@tcp(127.0.0.1:3306)/lock?charset=utf8"
	orm.RegisterDataBase("default", "mysql", conn)

	log.SetLevel(LogLevel("debug"))

	o := orm.NewOrm()
	if err := o.Using("default"); err != nil {
		log.Errorf("err: %s", err)
	}
	createTable("lock")
	// Register driver
	rwlock.InitDriver("mysql")

	lock, err := rwlock.GetRwlocker("aaa")
	if err != nil {
		log.Errorf("dddd: %s", err)
	}
	var ok bool
	if ok, err = lock.Wlock(); err != nil {
		log.Errorf("Get Rlocker error: %s", err)
	}
	if ok {
		log.Debugf("Get Rlock successfuly")
	} else {
		log.Debugf("Can not Get rlock")
	}

	//lock.RUnlock()
	//	if ok, err = lock.Rlock(); err != nil {
	//		log.Errorf("Get Wlocker error: %s", err)
	//	}
	//	if ok {
	//		log.Debugf("Get Wlock successfuly")
	//	}

	lock1, err1 := rwlock.GetRwlocker("bbb")
	if err1 != nil {
		log.Errorf("dddd: %s", err1)
	}
	if ok, err = lock1.Wlock(); err != nil {
		log.Errorf("Get Rlocker error: %s", err)
	}
	if ok {
		log.Debugf("Get Rlock successfuly")
	} else {
		log.Debugf("Can not Get rlock")
	}
}

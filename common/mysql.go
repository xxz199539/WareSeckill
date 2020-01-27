package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var Engine *xorm.Engine
func init(){
	engine, _ := NewMysqlConn()
	engine.SetMaxOpenConns(1000)
}


// 创建mysql连接
func NewMysqlConn() (db *xorm.Engine, err error) {
	p, err := ReadConf("conf.toml")
	if err != nil {
		log.Fatal("toml file err")
	}
	fmt.Println(p)
	Engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
		p.Mysql.MysqlUser, p.Mysql.MysqlPassword, p.Mysql.MysqlHost, p.Mysql.MysqlPort, p.Mysql.MysqlDbName))
	if err != nil {
		log.Fatal("get mysql db conn faild")
		return nil, err
	}
	return Engine, nil
}


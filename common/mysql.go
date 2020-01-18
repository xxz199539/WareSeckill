package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var engine *xorm.Engine

// 创建mysql连接
func NewMysqlConn() (db *xorm.Engine, err error) {
	//p, err := ReadConf("conf.toml")
	//if err != nil {
	//	log.Fatal("toml file err")
	//}
	engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", "root", 123456, "192.168.1.6:3306", "ware_seckill"))
	if err != nil {
		log.Fatal("get mysql db conn faild")
		return nil, err
	}
	return engine, nil
}


package common

import (
	"github.com/BurntSushi/toml"
	"sync"
)
var once sync.Once


var HostArray = []string{"127.0.0.1", "127.0.0.1"}

var LocalHost = ""

var Port = "8013"

type mysql struct {
	MysqlUser     string
	MysqlPassword int
	MysqlHost     string
	MysqlPort     int
	MysqlDbName   string


}
type secret struct {
	Key string
}

type config struct {
	Title   string
	Mysql   mysql
	Secret  secret
}

func ReadConf(fname string) (*config, error) {
	var config config
	if _, err := toml.DecodeFile(fname, &config);err != nil{
		return nil, err
	}
	return &config, nil
}

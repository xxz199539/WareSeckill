package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type mysql struct {
	mysqlUser string
	mysqlPassword int
	mysqlHost string
	mysqlPort int
	mysqlDbName string


}

type config struct {
	Mysql mysql
}

func ReadConf(fname string) (p *mysql, err error) {
	c := new(config)

	if _, err := toml.DecodeFile(fname, &c); err != nil {
		fmt.Println("toml.Unmarshal error ", err)
		return nil, err
	}
	return &c.Mysql, nil
}

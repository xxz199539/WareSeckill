package models

type User struct {
	Id       int    `xorm:"not null pk autoincr INT(8)"`
	NickName string `xorm:"default '' VARCHAR(32)"`
	UserName string `xorm:"not null VARCHAR(32)"`
	Password string `xorm:"not null VARCHAR(255)"`
}

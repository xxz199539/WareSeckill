package models

import (
	"time"
)

type Product struct {
	Id           int       `xorm:"not null pk autoincr INT(4)"`
	ProductName  string    `xorm:"not null VARCHAR(64)"`
	ProductNum   int       `xorm:"not null INT(4)"`
	ProductImage string    `xorm:"not null default '' VARCHAR(256)"`
	ProductUrl   string    `xorm:"VARCHAR(256)"`
	DeleteAt     time.Time `xorm:"comment('deleted') TIME(6)"`
}

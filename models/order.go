package models

import "time"

type Order struct {
	Id          int `xorm:"not null pk autoincr INT(8)"`
	UserId      int `xorm:"not null INT(8)"`
	ProductId   int `xorm:"not null index INT(8)"`
	OrderStatus int `xorm:"not null INT(8)"`
	CreateTime  time.Time `xorm:"DATETIME(4)"`
}

type OrderGroup struct {
	Order  `xorm:"extends"`
	productName string
}

const (
	SUCCESS = 1
	FAILED  = 0
)
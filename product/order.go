package product

import (
	"time"
)

type Order struct {
	CreateTime  time.Time `xorm:"DATETIME(4)"`
	Id          int       `xorm:"not null pk autoincr INT(8)"`
	OrderStatus int       `xorm:"not null INT(8)"`
	ProductId   int       `xorm:"not null index INT(8)"`
	UserId      int       `xorm:"not null INT(8)"`
}

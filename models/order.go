package models

type Order struct {
	Id          int `xorm:"not null pk autoincr INT(8)"`
	UserId      int `xorm:"not null INT(8)"`
	ProductId   int `xorm:"not null index INT(8)"`
	OrderStatus int `xorm:"INT(8)"`
}

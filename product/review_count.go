package product

type ReviewCount struct {
	Id        int    `xorm:"not null pk autoincr INT(8)"`
	ProductId int    `xorm:"not null index INT(8)"`
	UserId    int    `xorm:"not null index INT(8)"`
	UserIp    string `xorm:"VARCHAR(16)"`
}

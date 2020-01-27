package product

import (
	"time"
)

type Product struct {
	Country      string    `xorm:"default 'Made In China' comment('产地') VARCHAR(128)"`
	CurrentPrice float32   `xorm:"not null comment('当前价格') FLOAT(16,2)"`
	DeleteAt     time.Time `xorm:"comment('deleted') TIME(6)"`
	Id           int       `xorm:"not null pk INT(8)"`
	Material     string    `xorm:"default 'Cotton 100%' comment('材料') VARCHAR(128)"`
	OldPrice     float32   `xorm:"not null comment('原价') FLOAT(16,2)"`
	ProductImage string    `xorm:"default '' VARCHAR(256)"`
	ProductName  string    `xorm:"not null VARCHAR(64)"`
	ProductNum   int       `xorm:"not null INT(4)"`
	ProductUrl   string    `xorm:"VARCHAR(256)"`
}

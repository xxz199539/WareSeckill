package repositories

import (
	"WareSeckill/common"
	"WareSeckill/models"
	"github.com/go-xorm/xorm"
	"log"
)

// 先开发对应接口，再实现接口
type IProduct interface {
	Insert(*models.Product)(int64, error)
	Delete(int64) bool
	Update(*models.Product) error
	SelectByKey(int64) (*models.Product, error)
    SelectAll()([]*models.Product, error)
}

type ProductManager struct {

}

var engine *xorm.Engine


func init(){
	engine, _ = common.NewMysqlConn()
}

func NewProductManager() IProduct {
	return &ProductManager{}
}


func (p *ProductManager)Insert(product *models.Product) (productID int64,err error)  {
	pro := new(models.Product)
	pro.ProductName = product.ProductName
	pro.ProductImage = product.ProductImage
	pro.ProductNum = product.ProductNum
	pro.ProductUrl = product.ProductUrl
	affected, err := engine.Insert(pro)
    if err != nil {
    	log.Fatalf("insert product failed: %v", err)
    	return 0, nil
	}
	return affected, nil
}

func (p *ProductManager)Update(product *models.Product) error {
	pro := new(models.Product)
	pro.ProductName = product.ProductName
	pro.ProductImage = product.ProductImage
	pro.ProductNum = product.ProductNum
	pro.ProductUrl = product.ProductUrl
	_, err := engine.Id(product.Id).Update(pro)
	if err != nil {
		log.Fatalf("update product failed: %v", err)
		return err
	}
	return nil
}

func(p *ProductManager)Delete(Id int64) bool {
	var product models.Product
	engine.Id(Id).Get(&product)
	_, err := engine.Id(Id).Delete(&product)
	if err != nil {
		log.Fatalf("delete product failed: %v", err)
		return false
	}
	return true
}

func(p *ProductManager)SelectByKey(Id int64)(*models.Product, error){
	product := new(models.Product)
	_, err := engine.Where("id = ?", Id).Get(&product)
	if err != nil {
		log.Fatalf("get product by id failed: %v", err)
		return &models.Product{}, err
	}
	return product, nil
}

func(p *ProductManager)SelectAll()([]*models.Product, error) {
	products := make([]*models.Product, 0)
	err := engine.Find(&products)
	if err != nil {
		log.Fatalf("get all product failed: %v", err)
		return nil, err
	}
   return products, nil
}
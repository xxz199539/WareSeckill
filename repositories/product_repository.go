package repositories

import (
	"WareSeckill/common"
	"WareSeckill/models"
	"log"
	"time"
)

// 先开发对应接口，再实现接口
type IProduct interface {
	Insert(*models.Product) (int64, error)
	Delete(int64) bool
	Update(*models.Product) error
	SelectByKey(int64) (*models.Product, error)
	SelectAll() ([]*models.Product, error)
	SubProductNumberOne(int64) (int64, error)
}

type ProductManager struct {
}

func NewProductManager() IProduct {
	return &ProductManager{}
}

func (p *ProductManager) Insert(product *models.Product) (productID int64, err error) {
	pro := new(models.Product)
	pro.ProductName = product.ProductName
	pro.ProductImage = product.ProductImage
	pro.ProductNum = product.ProductNum
	pro.ProductUrl = product.ProductUrl
	affected, err := common.Engine.Insert(pro)
	if err != nil {
		log.Fatalf("insert product failed: %v", err)
		return 0, nil
	}
	return affected, nil
}

func (p *ProductManager) Update(product *models.Product) error {
	pro := new(models.Product)
	pro.ProductName = product.ProductName
	pro.ProductImage = product.ProductImage
	pro.ProductNum = product.ProductNum
	pro.ProductUrl = product.ProductUrl
	_, err := common.Engine.Id(product.Id).Update(pro)
	if err != nil {
		log.Fatalf("update product failed: %v", err)
		return err
	}
	return nil
}

func (p *ProductManager) Delete(Id int64) bool {
	var product models.Product
	_, _ = common.Engine.Id(Id).Get(&product)
	_, err := common.Engine.Id(Id).Delete(&product)
	if err != nil {
		log.Fatalf("delete product failed: %v", err)
		return false
	}
	return true
}

func (p *ProductManager) SelectByKey(Id int64) (*models.Product, error) {
	product := new(models.Product)
	_, err := common.Engine.Where("id = ?", Id).Get(product)
	if err != nil {
		log.Fatalf("get product by id failed: %v", err)
		return &models.Product{}, err
	}
	return product, nil
}

func (p *ProductManager) SelectAll() ([]*models.Product, error) {
	products := make([]*models.Product, 0)
	err := common.Engine.Find(products)
	if err != nil {
		log.Fatalf("get all product failed: %v", err)
		return nil, err
	}
	return products, nil
}

func (p *ProductManager) SubProductNumberOne(id int64) (int64,error) {
	product, err := p.SelectByKey(id)
	if err != nil {
		return 0, err
	}
	if product.ProductNum >= 1 {
		newProduct := &models.Product{
			Country:      product.Country,
			CurrentPrice: product.CurrentPrice,
			DeleteAt:     time.Time{},
			Material:     product.Material,
			OldPrice:     product.OldPrice,
			ProductImage: product.ProductImage,
			ProductName:  product.ProductName,
			ProductNum:   product.ProductNum - 1,
			ProductUrl:   product.ProductUrl,
		}
		_, err = common.Engine.Id(product.Id).Cols("product_num").Update(newProduct)
		if err != nil {
			return 0, err
		}
		return int64(product.ProductNum), nil
	}
	return 0, nil
}

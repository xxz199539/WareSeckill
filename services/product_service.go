package services

import (
	"WareSeckill/models"
	"WareSeckill/repositories"
)

type IProductService interface {
	GetProductById(Id int64) (*models.Product, error)
	GetAllProduct()([]*models.Product, error)
	DeleteProduct(Id int64) bool
	InsertProduct(product *models.Product)(int64, error)
	UpdateProduct(product *models.Product) error
}

type ProductService struct {
    productRepository repositories.IProduct
}

func NewProductService(repositoriy repositories.IProduct) IProductService {
	return &ProductService{productRepository:repositoriy}
}

func(p *ProductService)GetProductById(Id int64)(*models.Product, error){
	return p.productRepository.SelectByKey(Id)
}

func(p *ProductService)GetAllProduct()([]*models.Product, error){
	return p.productRepository.SelectAll()
}

func(p *ProductService)DeleteProduct(Id int64) bool{
	return p.productRepository.Delete(Id)
}

func(p *ProductService)InsertProduct(product *models.Product)(int64, error){
	return p.productRepository.Insert(product)
}

func (p *ProductService)UpdateProduct(product *models.Product) error {
	return p.productRepository.Update(product)
}



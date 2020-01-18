package controllers

import (
	"WareSeckill/models"
	"WareSeckill/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
}

func(p *ProductController)GetAll() mvc.View{
	productArray, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name:"product/view.html",
		Data:iris.Map{
			"productArray": productArray,
		},
	}
}

func(p *ProductController)PostUpdate() {
	product := &models.Product{}
	err := p.Ctx.ReadForm(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func(p *ProductController) GetAdd() mvc.View{
	return mvc.View{
		Name: "product/add.html",
	}
}

func(p *ProductController)PostAdd(){
	product := &models.Product{}
	err := p.Ctx.ReadForm(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err = p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func(p *ProductController)GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductById(id)
	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func(p *ProductController)GetDelete(){
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	isOk := p.ProductService.DeleteProduct(id)
	if isOk {
		p.Ctx.Application().Logger().Debug("Delete product success")
	}else {

		p.Ctx.Application().Logger().Debug("Delete product failed")
	}
	p.Ctx.Redirect("/product/all")
}
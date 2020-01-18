package controllers

import (
	"WareSeckill/models"
	"WareSeckill/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
)

type OrderController struct {
	Ctx iris.Context
    OrderService  services.IOrderService
}

func (o *OrderController)GetAll() mvc.View{
	orderArray,_ := o.OrderService.SelectAll()
	return mvc.View{
		Name:"order/view.html",
		Data:iris.Map{
			"order": orderArray,
		},
	}
}

func (o *OrderController)PostUpdate() {
	order := &models.Order{}
	err := o.Ctx.ReadForm(order)
	if err != nil{
		o.Ctx.Application().Logger().Debug(err)
	}
	err = o.OrderService.Update(order)
	if err != nil{
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}

func (o *OrderController)GetDelete()  {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	hasDelete := o.OrderService.Delete(id)
	if hasDelete {
		o.Ctx.Application().Logger().Debug("Delete product success")
	}else {
		o.Ctx.Application().Logger().Debug("Delete product failed")
	}
	o.Ctx.Redirect("/order/all")
}

func (o *OrderController)PostAdd() {
	order := &models.Order{}
	if err := o.Ctx.ReadForm(order);err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	_, err := o.OrderService.Insert(order)
	if err != nil{
		o.Ctx.Application().Logger().Debug(err)
	}
	o.Ctx.Redirect("/order/all")
}



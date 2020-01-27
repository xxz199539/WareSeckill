package controllers

import (
	"WareSeckill/RabbitMQ"
	"WareSeckill/common"
	"WareSeckill/common/encrypt"
	"WareSeckill/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strconv"
)

type ProductController struct {
	Ctx           iris.Context
	ProductServer services.IProductService
	OrderService  services.IOrderService
	UserService   services.UserService
	RabbitMq      *RabbitMQ.RabbitMQ
}

func (p *ProductController) GetDetail() mvc.View {
	common.AuthValidate(p.Ctx)
	productId := p.Ctx.URLParam("productId")
	cookieUserId := p.Ctx.GetCookie("userId")
	dePwdCookieUserId, err := encrypt.DePwdCode(cookieUserId)
	userId := int(dePwdCookieUserId[0])
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return mvc.View{
			Name: "shared/error.html",
			Data: iris.Map{
				"Message": err.Error(),
			},
		}
	}
	isOk := common.AccessControlGlobal.GetDistributedRight(cookieUserId)
	if !isOk {
		p.Ctx.Redirect("/user/login")
	}else {
		// 新增商品浏览记录
		ipString := p.Ctx.Request().RemoteAddr
		productIdInt, _ := strconv.Atoi(productId)
		product, err := p.ProductServer.GetProductById(int64(productIdInt))
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
			return mvc.View{
				Name: "shared/error.html",
				Data: iris.Map{
					"Message": err.Error(),
				},
			}
		}
		_, err = common.ProductReviewCount(userId, productIdInt, ipString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)// 这里浏览记录添加出错不影响商品浏览
		}
		user, err := p.UserService.GetUserById(int64(userId))
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
			return mvc.View{
				Name: "shared/error.html",
				Data: iris.Map{
					"Message": err.Error(),
				},
			}
		}
		return mvc.View{
			Layout: "shared/productLayout.html",
			Name:   "product/view.html",
			Data: iris.Map{
				"product": product,
				"user":    user,
			},
		}
	}
	return mvc.View{}
}

func (p *ProductController) GetOrder() []byte {
	productIdString := p.Ctx.URLParam("productID")
	userId := p.Ctx.URLParam("userId")
	getOneHost, err := common.HashConsistent.Get(userId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return []byte("false")
	}
	getOneUrl := "http://" + getOneHost + ":" + common.Port + "/getOne?productId=" + productIdString

	res, body, err := common.CurlUrl(getOneUrl, p.Ctx.Request())
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return []byte("false")
	}
	if res.StatusCode == 200 {
		return body
	}
	return []byte("false")
}
package main

import (
	"WareSeckill/backend/web/controllers"
	"WareSeckill/common"
	"WareSeckill/repositories"
	"WareSeckill/services"
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

    // 注册模板
    template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
    // 设置模板目标
    app.HandleDir("/assets", "./backend/web/assets")
    // 出现异常跳转到制定页面
    app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问页面出错"))
		ctx.ViewLayout("")
		_ = ctx.View("./backend/web/views/shared/error.html")
	})
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
	accessControl := common.NewAccessControl()
	hashConsistent := common.NewConsistent()
	for _, v := range common.HostArray {
		hashConsistent.Add(v)
	}
    // 注册控制器
    productManager := repositories.NewProductManager(accessControl, hashConsistent)
    productService := services.NewProductService(productManager)
    productParty := app.Party("/product")
    product := mvc.New(productParty)
    product.Register(ctx, productService)
    product.Handle(new(controllers.ProductController))

	orderManager := repositories.NewOrderRepository()
	orderService := services.NewOrderService(orderManager)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
    err := app.Run(
    	iris.Addr("0.0.0.0:8080"),
    	iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
    if err != nil {
    	log.Fatalf("start server failed: %v", err)
	}
}

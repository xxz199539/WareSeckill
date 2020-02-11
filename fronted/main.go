package main

import (
	"WareSeckill/RabbitMQ"
	"WareSeckill/common"
	"WareSeckill/fronted/web/controllers"
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
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// 设置模板目标
	app.StaticWeb("/public", "./fronted/web/public")
	// 出现异常跳转到制定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问页面出错"))
		ctx.ViewLayout("")
		_ = ctx.View("./fronted/web/views/shared/error.html")
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitMq := RabbitMQ.NewRabbitMQSimple("BoomShakala")
	// 注册控制器
	userManager := repositories.NewUserRepository()
	userService := services.NewUserService(userManager)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	productManager := repositories.NewProductManager()
	productService := services.NewProductService(productManager)
	orderManager := repositories.NewOrderRepository()
	orderService := services.NewOrderService(orderManager)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	productParty.Use(common.AuthValidate)
	product.Register(productService, orderService, userService, rabbitMq)
	product.Handle(new(controllers.ProductController))

	captChaParty := app.Party("/captcha/:captChaId")
	captCha := mvc.New(captChaParty)
	captChaService := services.NewCaptChaService()
	captCha.Register(ctx, captChaService)
	captCha.Handle(new(controllers.CaptChaController))

	// 负载均衡设置
	// 采用一致性哈希算法
	hashConsistent := common.NewConsistent()
	// 采用一致性hash算法，添加节点
	for _, v := range common.HostArray {
		hashConsistent.Add(v)
	}
	err := app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	if err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}

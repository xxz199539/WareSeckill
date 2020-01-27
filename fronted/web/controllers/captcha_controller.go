package controllers

import (
	"WareSeckill/services"
	"github.com/kataras/iris"
)

type CaptChaController struct {
	Ctx iris.Context
	CaptChaServer services.CaptChaService
}

func (c *CaptChaController)Get(){
	c.CaptChaServer.GenerateCaptCha(c.Ctx)
}

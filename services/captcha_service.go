package services

import (
	"github.com/dchest/captcha"
	"github.com/kataras/iris"
)

type CaptChaService interface {
	ProcessFormHandler(captchaId, captchaSolution string)bool
	GenerateCaptCha(ctx iris.Context)
	GenerateCapId()string
}

type CaptChaManager struct {
}


func NewCaptChaService()CaptChaService{
	return &CaptChaManager{}
}
func(c *CaptChaManager)ProcessFormHandler(captchaId, captchaSolution string)bool {
	if captchaId == "" || captchaSolution=="" {
		return false
	}
	if !captcha.VerifyString(captchaId, captchaSolution) {
		return false
	}
	return true
}


func (c *CaptChaManager)GenerateCaptCha(ctx iris.Context) {
	captcha.Server(340, captcha.StdHeight).ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}

func (c *CaptChaManager)GenerateCapId()string{
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	return d.CaptchaId
}
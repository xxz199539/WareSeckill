package controllers

import (
	"WareSeckill/common"
	"WareSeckill/common/encrypt"
	"WareSeckill/models"
	"WareSeckill/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

type UserController struct {
	Ctx        iris.Context
	UserServer services.UserService
	Session    *sessions.Session
}

// 返回注册页面
func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "/user/register.html",
	}
}

func (u *UserController) GetLogin() mvc.View {
	id := services.NewCaptChaService().GenerateCapId()
	return mvc.View{
		Name: "/user/login.html",
		Data:iris.Map{
			"CaptchaId": id,
		},
	}
}


// 注册提交表单
func (u *UserController) PostRegister(){
	user := &models.User{}
	if err := u.Ctx.ReadForm(user); err != nil {
		u.Ctx.Application().Logger().Debug(err)
	}
	id, err := u.UserServer.AddUser(user)
	if id ==0 || err != nil {
		u.Ctx.Application().Logger().Debug(err)
		u.Ctx.Redirect("/user/login")
		return
	}
	u.Ctx.Redirect("/user/login")
	return
}

func (u *UserController)PostLogin(){
	userName := u.Ctx.FormValue("userName")
	passWord := u.Ctx.FormValue("password")
	captchaSolution := u.Ctx.FormValue("captchaSolution")
	captchaId := u.Ctx.FormValue("captchaId")
	isOk := services.NewCaptChaService().ProcessFormHandler(captchaId, captchaSolution)
	if !isOk{
		u.Ctx.Application().Logger().Debug("Validate CaptCha error")
		u.Ctx.Redirect("/user/login")
	}
	user, err := u.UserServer.CheckUserPasswordByName(userName, passWord)
	if err != nil{
		u.Ctx.Application().Logger().Debug(err)
		u.Ctx.Redirect("/user/login")
	}
	if user.Id !=0{
		hashedId, err := encrypt.EnPwdCode([]byte(string(user.Id)))
		if err != nil {
			u.Ctx.Application().Logger().Debug(err)
			u.Ctx.Redirect("/user/login")
		}
	    common.GlobalCookie(u.Ctx, "userId", hashedId)
		err = common.AccessControlGlobal.GenerateDistributedRight(user.Id, hashedId)
		if err != nil {
			u.Ctx.Redirect("/user/login")
		}else {
			u.Ctx.Redirect("/product/all")
		}
	}else {
		u.Ctx.Redirect("/user/login")
	}
	return
}
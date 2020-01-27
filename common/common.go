package common

import (
	"WareSeckill/models"
	"errors"
	"github.com/kataras/iris"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

var AccessControlGlobal =  &AccessControl{
	SourcesArray: make(map[int]interface{}),
	RWMutex:      sync.RWMutex{},
}

var HashConsistent = &Consistent{
	circle:         make(map[uint32]string),
	VirtualNodeNum: 20,
}

func init(){
	once.Do(func() {
		for _, v := range HostArray {
			HashConsistent.Add(v)
		}
	})
}

func ValidatePassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

//设置全局cookie
func GlobalCookie(ctx iris.Context, name string, value string) {
	ctx.SetCookie(&http.Cookie{Name: name, Value: value, Path: "/", MaxAge:3600})
}

func AuthValidate(ctx iris.Context) {
	userId := ctx.GetCookie("userId")
	if userId == "" {
		ctx.Redirect("/user/login")
	}
	ctx.Next()
}

func ProductReviewCount(userId, productId int, ip string) (int64, error) {
	data := &models.ReviewCount{
		ProductId: productId,
		UserId:    userId,
		UserIp:    ip,
	}
	affected, err := Engine.Insert(data)
	if err != nil {
		return 0, err
	}
	return affected, nil
}

// 获取本机IP地址
func GetIntranceIp()(string ,error)  {
	addrs,err:=net.InterfaceAddrs()
	if err !=nil {
		return "",err
	}

	for _,address:= range addrs{
		//检查Ip地址判断是否回环地址
		if ipnet,ok:=address.(*net.IPNet);ok&&!ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(),nil
			}
		}
	}

	return "",errors.New("获取地址异常")
}

func CurlUrl(uri string, req *http.Request) (response *http.Response,body []byte,err error) {
	hashedUId, err := req.Cookie("userId")
	if err != nil {
			return
		}
	client := &http.Client{}
	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}
	req.AddCookie(&http.Cookie{
		Name:       "userId",
		Value:      hashedUId.Value,
		Path:       "/",
		MaxAge:     3600,
	})
	response, err = client.Do(req)
	defer func() {
		_ = response.Body.Close()
	}()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return response, body, nil
}
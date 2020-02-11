package main

import (
	"WareSeckill/RabbitMQ"
	"WareSeckill/common"
	"WareSeckill/common/encrypt"
	"WareSeckill/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

// 提供cookie验证，提供rabbitMQ生产者服务
var sum int64
var mutex sync.Mutex

var accessControl = &common.AccessControl{
	SourcesArray: make(map[int]interface{}),
	RWMutex:      sync.RWMutex{},
}

type AccessControl struct {
	// 用来存放用户想要的信息
	sourcesArray map[int]interface{}
	sync.RWMutex
}

var RabbitMq *RabbitMQ.RabbitMQ

// 根据uid获取数据
func (m *AccessControl) GetNewRecord(uid int, keyName string) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

// 设置记录
func (m *AccessControl) SetNewRecord(uid int, value interface{}) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourcesArray[uid] = value
}

func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Exec Validate")
	cookieFromRequest, err := r.Cookie("userId")
	if err != nil {
		return err
	}
	decryptString, err := encrypt.DePwdCode(cookieFromRequest.Value)
	if err != nil {
		return errors.New("check cookie:userId failed")
	}
	cookieInt := int(decryptString[0])
	// 获取缓存中保存的cookie
	cookieCache := accessControl.GetNewRecord(cookieInt)
	if cookieCache == nil {
		return errors.New("validate fail")
	}
	if cookieCache != cookieFromRequest.Value {
		return errors.New("validate fail")
	}
	return nil
}

func Check(w http.ResponseWriter, r *http.Request) {
	byteSlice := make([]byte, 0)
	byteSlice = strconv.AppendBool(byteSlice, true)
	_, _ = w.Write(byteSlice)
}

func Generate(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("generating user cookie")
	hashedUserId, err := r.Cookie("userId")
	if err != nil {
		return err
	}
	userIdInt, err := encrypt.DePwdCode(hashedUserId.Value)
	if err != nil {
		return err
	}
	accessControl.SetNewRecord(int(userIdInt[0]), hashedUserId.Value)
	//accessControl.SetNewRecord(userIdInt, "expireTime", time.Now().AddDate(0,0,1).Format("2006-01-02 15:04:05"))
	return nil
}

func OrderProduct(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		_, _ = w.Write([]byte("false"))
		return
	}
	if len(queryForm["productId"]) <= 0 {
		_, _ = w.Write([]byte("false"))
		return
	}
	if sum < 10000 {
		productId := queryForm["productId"]
		userId, err := r.Cookie("userId")
		if err != nil {
			_, _ = w.Write([]byte("false"))
			return
		}
		productIdInt64, err := strconv.ParseInt(productId[0], 10, 64)
		if err != nil {
			_, _ = w.Write([]byte("false"))
			return
		}
		dePwdUserId, err := encrypt.DePwdCode(userId.Value)
		userIdInt64 := int64(dePwdUserId[0])
		newMessage := &models.Message{
			UserId:    userIdInt64,
			ProductId: productIdInt64,
		}
		byteMessage, err := json.Marshal(newMessage)
		if err != nil {
			_, _ = w.Write([]byte("false"))
			return
		}
		err = RabbitMq.PublishSimple(string(byteMessage))
		if err != nil {
			_, _ = w.Write([]byte("false"))
			return
		}
		_, _ = w.Write([]byte("true"))
		sum++
		return
	}
	_, _ = w.Write([]byte("false"))
	return

}

func main() {
	RabbitMq = RabbitMQ.NewRabbitMQSimple("BoomShakalaka")
	defer RabbitMq.Destroy()
	// 1.过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/generate", Generate)
	// 2.启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/generate", filter.Handle(Check))
	http.HandleFunc("/getOne", filter.Handle(OrderProduct))

	_ = http.ListenAndServe(":8013", nil)
}

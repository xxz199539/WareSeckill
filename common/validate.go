package common

import (
	"WareSeckill/common/encrypt"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type AccessControl struct {
	// 用来存放用户想要的信息
	SourcesArray map[int]interface{}
	sync.RWMutex
}

func NewAccessControl() *AccessControl {
	return &AccessControl{
		SourcesArray: make(map[int]interface{}),
		RWMutex:      sync.RWMutex{},
	}
}
// 根据uid获取数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.SourcesArray[uid]
	return data
}

// 设置记录
func (m *AccessControl) SetNewRecord(uid int,keyValue interface{}) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.SourcesArray[uid] = keyValue
}

func (m *AccessControl) GetDistributedRight(hashId string) bool {

	DePwdUId, err := encrypt.DePwdCode(hashId)
	if err != nil {
		return false
	}
	// 采用一致性算法，根据用户id，判断获取具体机器
	hostRequest, err := HashConsistent.Get(string(DePwdUId))
	if err != nil {
		return false
	}
	//判断是否为本机
	if hostRequest == LocalHost {
		// 执行本机数据读取和校验
		return m.GetDataFromLocal(hashId)
	}
	return m.GetDataFromReMote(hostRequest, hashId)
}

func (m *AccessControl) GetDataFromLocal(uid string) bool {
	// 拿到加密后的id
	DePwdUserId, err := encrypt.DePwdCode(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(int(DePwdUserId[0]))
	if data != nil {
		return true
	}
	return false
}

func (m *AccessControl) GetDataFromReMote(host string, userId string) bool {
	// 模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%v:%v/check", host, Port), nil)
	if err != nil {
		return false
	}

	// 将uid添加到cookie中
	cookieUid := &http.Cookie{Name: "userId", Value: userId, Path: "/"}
	req.AddCookie(cookieUid)

	// 获取返回结果
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		}
	}
	return false
}

func (m *AccessControl)GenerateDataRemote(host string, hashId string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%v:%v/generate", host, Port), nil)
	if err != nil {
		return err
	}
	cookieUid := &http.Cookie{Name: "userId", Value: hashId, Path: "/"}
	req.AddCookie(cookieUid)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (m *AccessControl)GenerateDataLocal(uid int, value interface{}) {
	m.SetNewRecord(uid, value)
}

func (m *AccessControl)GenerateDistributedRight(uId int, hashId string) error {
	host, err := HashConsistent.Get(string(uId))
	if err != nil {
		return err
	}
	localHost, err := GetIntranceIp()
	fmt.Println(localHost)
	if err != nil {
		return err
	}
	if host != localHost {
		err := AccessControlGlobal.GenerateDataRemote(host,hashId)
		if err != nil {
			return err
		}
	}else {
		AccessControlGlobal.GenerateDataLocal(uId, hashId)
		fmt.Println(AccessControlGlobal)
	}
	return nil
}


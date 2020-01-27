package repositories

import (
	"WareSeckill/common"
	"WareSeckill/models"
	"errors"
	"fmt"
	"log"
)

type IUserRepository interface {
	Insert(*models.User)(int64, error)
	SelectByName(string)(*models.User, error)
	SelectById(int64)(*models.User, error)
}

type UserRepository struct {

}

func NewUserRepository() IUserRepository{
	return &UserRepository{}
}

func (u *UserRepository)Insert(user *models.User)(int64, error) {
	newUser := &models.User{
		NickName: user.NickName,
		UserName: user.UserName,
		Password: user.Password,
	}
	existUser, err := u.SelectByName(user.UserName)
	if err != nil {
		return 0, errors.New("get user by name failed")
	}
	if existUser.UserName != "" {
		return 0, errors.New(fmt.Sprintf("name :%v exists", existUser.UserName))
	}
	id, err := common.Engine.Insert(newUser)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserRepository)SelectByName(name string)(*models.User, error){
	user := new(models.User)
	_, err := common.Engine.Where("user_name = ?", name).Get(user)
	if err != nil {
		log.Fatalf("get user by name failed: %v", err)
		return &models.User{},  nil
	}
	return user, nil
}


func (u *UserRepository)SelectById(id int64)(*models.User, error){
	user := new(models.User)
	_, err := common.Engine.Where("id = ?", id).Get(user)
	if err != nil {
		log.Fatalf("get user by id failed: %v", err)
		return &models.User{},  nil
	}
	return user, nil
}


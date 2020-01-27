package services

import (
	"WareSeckill/common"
	"WareSeckill/models"
	"WareSeckill/repositories"
	"log"
)

type UserService interface {
	AddUser(*models.User)(int64, error)
	CheckUserPasswordByName(name, password string)(user *models.User,err error)
	GetUserById(int64)(*models.User, error)
}

type UserServiceManager struct {
	repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository)UserService{
	return &UserServiceManager{repository}
}

func (u *UserServiceManager)AddUser(user *models.User)(int64, error){
	hashedPassword, err := common.GenerateHashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	userId, err := u.IUserRepository.Insert(user)
	if userId == 0 || err != nil{
		return 0, err
	}
	return userId, nil
}


func (u *UserServiceManager)CheckUserPasswordByName(name, password string)(*models.User, error){
    user, err := u.IUserRepository.SelectByName(name)
    if err != nil{
    	log.Fatalf("get user failed: %v", err)
    	return &models.User{}, err
	}
	_, err = common.ValidatePassword(user.Password, password)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func (u *UserServiceManager)GetUserById(id int64)(*models.User, error){
	return u.IUserRepository.SelectById(id)
}

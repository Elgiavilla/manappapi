package service

import (
	"time"

	"github.com/elgiavilla/manapp/models"
	"github.com/elgiavilla/manapp/user"
)

type userService struct {
	userRepo   user.Repository
	ctxTimeout time.Duration
}

func NewServiceUser(u user.Repository, timout time.Duration) user.Service {
	return &userService{
		userRepo:   u,
		ctxTimeout: timout,
	}
}

func (u *userService) Insert(userModel models.User) (*models.User, error) {
	d, err := u.userRepo.Insert(userModel)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (u *userService) Login(userModel models.User) (*models.User, error) {
	res, err := u.userRepo.Login(userModel)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userService) GetFindById(user_id uint) (*models.User, error) {
	res, err := u.userRepo.GetFindById(user_id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userService) Update(userModel models.User) (*models.User, error) {
	d, err := u.userRepo.Update(userModel)
	if err != nil {
		return nil, err
	}
	return d, nil
}

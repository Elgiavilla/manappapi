package user

import (
	"github.com/elgiavilla/manapp/models"
)

type Service interface {
	Insert(userModel models.User) (*models.User, error)
	Login(userModel models.User) (*models.User, error)
	GetFindById(user_id uint) (*models.User, error)
	Update(userModel models.User) (*models.User, error)
}

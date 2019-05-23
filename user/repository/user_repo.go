package repository

import (
	"github.com/elgiavilla/manapp/models"
	"github.com/elgiavilla/manapp/user"
	"github.com/jinzhu/gorm"
)

type dbUserRepo struct {
	db *gorm.DB
}

func NewUserRepo(Conn *gorm.DB) user.Repository {
	return &dbUserRepo{Conn}
}

func (u *dbUserRepo) Insert(userModel models.User) (*models.User, error) {
	row := new(models.User)
	d := u.db.Debug().Create(&userModel).Scan(&row)
	if d.Error != nil {
		return nil, d.Error
	}
	return row, nil
}

func (u *dbUserRepo) Login(userModel models.User) (*models.User, error) {
	list := u.db.Where("email = ?", &userModel.Email).Find(&userModel)
	return &userModel, list.Error
}

func (u *dbUserRepo) GetFindById(user_id uint) (*models.User, error) {
	var user models.User
	list := u.db.Where("id = ?", user_id).Find(&user)
	return &user, list.Error
}

func (u *dbUserRepo) Update(userModel models.User) (*models.User, error) {
	row := new(models.User)
	up := u.db.Model(&userModel).Update(&userModel).Scan(&row)
	return row, up.Error
}

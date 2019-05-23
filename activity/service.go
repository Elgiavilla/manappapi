package activity

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/elgiavilla/manapp/models"
)

type Service interface {
	Insert(activityModel models.Activity) (*models.Activity, error)
	GetByUserID(user_id uint, page int, limit int) (*pagination.Paginator, error)
	GetByID(activity_id uint) (*models.Activity, error)
	UpdateStatus(activityModel models.Activity, status int64) (*models.Activity, error)
	GetByUserToday(user_id uint, page int, limit int) (*pagination.Paginator, error)
	GetCountSuccessByUser(user_id uint) (*models.CountActivity, error)
}

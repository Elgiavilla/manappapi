package repository

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/elgiavilla/manapp/activity"
	"github.com/elgiavilla/manapp/models"
	"github.com/jinzhu/gorm"
)

type dbActivityRepo struct {
	db *gorm.DB
}

func NewActivityRepo(Conn *gorm.DB) activity.Repository {
	return &dbActivityRepo{Conn}
}

func (a *dbActivityRepo) Insert(activityModel models.Activity) (*models.Activity, error) {
	row := new(models.Activity)
	d := a.db.Debug().Create(&activityModel).Scan(&row)
	if d.Error != nil {
		return nil, d.Error
	}
	return row, nil
}

func (a *dbActivityRepo) GetByUserID(user_id uint, page int, limit int) (*pagination.Paginator, error) {
	var activity []*models.Activity
	list := a.db.Debug().Where("user_id = ? AND complete_activity != 0", user_id).Find(&activity)
	if list.Error != nil {
		return nil, list.Error
	}
	paginator := pagination.Paging(&pagination.Param{
		DB:      list,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &activity)
	return paginator, nil
}

func (a *dbActivityRepo) GetByID(activity_id uint) (*models.Activity, error) {
	row := new(models.Activity)
	d := a.db.Debug().Where("id = ?", activity_id).Find(&row)
	if d.Error != nil {
		return nil, d.Error
	}
	return row, nil
}

func (a *dbActivityRepo) UpdateStatus(activityModel models.Activity, status int64) (*models.Activity, error) {
	row := new(models.Activity)
	up := a.db.Debug().Model(&activityModel).Update("complete_activity", status).Scan(&row)
	if up.Error != nil {
		return nil, up.Error
	}
	return row, nil
}

func (a *dbActivityRepo) GetByUserToday(user_id uint, page int, limit int) (*pagination.Paginator, error) {
	var activity []*models.Activity
	list := a.db.Debug().Where("user_id = ? AND DATE(created_at) =  CURRENT_DATE AND complete_activity = 0", user_id).Find(&activity)
	if list.Error != nil {
		return nil, list.Error
	}
	paginator := pagination.Paging(&pagination.Param{
		DB:      list,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &activity)
	return paginator, nil
}

func (a *dbActivityRepo) GetCountSuccessByUser(user_id uint) (*models.CountActivity, error) {
	row := new(models.CountActivity)
	list := a.db.Debug().Table("activities").Select("users.created_at as created_at, SUM(CASE WHEN activities.complete_activity = 1 THEN 1 ELSE 0 END) as success").Joins("LEFT JOIN users ON users.id = activities.user_id").Group("users.created_at").Where("users.id = ? AND activities.created_at = activities.created_at AND activities.created_at <= NOW()", user_id).Scan(&row)
	if list.Error != nil {
		return nil, list.Error
	}
	return row, nil
}

package service

import (
	"time"

	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/elgiavilla/manapp/activity"
	"github.com/elgiavilla/manapp/models"
)

type activityService struct {
	activityRepo activity.Repository
	ctxTimeout   time.Duration
}

func NewActivityService(a activity.Repository, timot time.Duration) activity.Service {
	return &activityService{
		activityRepo: a,
		ctxTimeout:   timot,
	}
}

func (a *activityService) Insert(activityModel models.Activity) (*models.Activity, error) {
	d, err := a.activityRepo.Insert(activityModel)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *activityService) GetByUserID(user_id uint, page int, limit int) (*pagination.Paginator, error) {
	d, err := a.activityRepo.GetByUserID(user_id, page, limit)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *activityService) GetByID(activity_id uint) (*models.Activity, error) {
	d, err := a.activityRepo.GetByID(activity_id)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *activityService) UpdateStatus(activityModel models.Activity, status int64) (*models.Activity, error) {
	d, err := a.activityRepo.UpdateStatus(activityModel, status)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *activityService) GetByUserToday(user_id uint, page int, limit int) (*pagination.Paginator, error) {
	d, err := a.activityRepo.GetByUserToday(user_id, page, limit)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (a *activityService) GetCountSuccessByUser(user_id uint) (*models.CountActivity, error) {
	d, err := a.activityRepo.GetCountSuccessByUser(user_id)
	if err != nil {
		return nil, err
	}
	return d, nil
}

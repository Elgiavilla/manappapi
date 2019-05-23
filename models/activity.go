package models

import (
	"github.com/jinzhu/gorm"
)

type Activity struct {
	gorm.Model
	UserID            uint
	Title             string `json:"title"`
	Description       string `json:"description"`
	Complete_activity int64  `json:"complete_activity"`
	Time_end          string `json:"time_end"`
}

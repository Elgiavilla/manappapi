package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email        string `json:"email"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Notification *bool  `json:"notification"`
	Device_id    string `json:"device_id"`
	Activity     []Activity
}

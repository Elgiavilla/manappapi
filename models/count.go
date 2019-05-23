package models

import "time"

type CountActivity struct {
	Success   int       `json:"success"`
	CreatedAt time.Time `json:"CreatedAt"`
}

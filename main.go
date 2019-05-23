package main

import (
	"net/url"
	"time"

	"github.com/elgiavilla/manapp/middleware"
	"github.com/elgiavilla/manapp/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"

	_userHttp "github.com/elgiavilla/manapp/user/http"
	_userRepo "github.com/elgiavilla/manapp/user/repository"
	_userSvc "github.com/elgiavilla/manapp/user/service"

	_actHttp "github.com/elgiavilla/manapp/activity/http"
	_actRepo "github.com/elgiavilla/manapp/activity/repository"
	_actSvc "github.com/elgiavilla/manapp/activity/service"
)

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=manapp password=123 dbname=manapp sslmode=disable")
	//fmt.Println(db)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&models.User{}, &models.Activity{})
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	e := echo.New()
	middl := middleware.InitMiddleware()
	timeoutContext := time.Duration(2) * time.Second

	/*
		USER CONFIG REPO AND SERVICE
	*/
	userRepo := _userRepo.NewUserRepo(db)
	userSvc := _userSvc.NewServiceUser(userRepo, timeoutContext)

	/*
		ACTIVITY CONFIG REPO AND SERVICE
	*/
	actRepo := _actRepo.NewActivityRepo(db)
	actSvc := _actSvc.NewActivityService(actRepo, timeoutContext)

	_userHttp.NewUserHandler(e, userSvc)
	_actHttp.NewActivityHandler(e, actSvc)
	e.Use(middl.CORS)
	e.Start(":8090")
}

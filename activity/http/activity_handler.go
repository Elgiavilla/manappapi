package http

import (
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/elgiavilla/manapp/activity"
	"github.com/elgiavilla/manapp/models"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type Token struct {
	ID int64
	jwt.StandardClaims
}

type SuccessMessage struct {
	Token string `json:"token"`
}

type ReponseError struct {
	Message string `json:"message"`
}

type HttpActivityHandler struct {
	ActivityService activity.Service
}

func NewActivityHandler(e *echo.Echo, activityService activity.Service) {
	handler := &HttpActivityHandler{
		ActivityService: activityService,
	}
	e.GET("/activity/user", handler.GetByUserID)
	e.POST("/activity", handler.Insert)
	e.GET("/activity/:id/:status", handler.UpdateStatus)
	e.GET("/activity/user/today", handler.GetByUserToday)
	e.GET("/activity/:id", handler.GetByID)
	e.GET("/activity/user/count", handler.GetCountSuccessByUser)
}

func (a *HttpActivityHandler) Insert(c echo.Context) error {
	var activity models.Activity
	err := c.Bind(&activity)
	idP, err := a.Auth(c.Request().Header.Get("Auth"))
	idInter := idP["ID"]
	str := uint(idInter.(float64))
	activity.UserID = str
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	d, err := a.ActivityService.Insert(activity)
	if err != nil {
		return c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (a *HttpActivityHandler) GetByUserID(c echo.Context) error {
	idP, err := a.Auth(c.Request().Header.Get("Auth"))
	idInter := idP["ID"]
	str := uint(idInter.(float64))
	pageS := c.QueryParam("page")
	page, _ := strconv.Atoi(pageS)
	limitS := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitS)

	d, err := a.ActivityService.GetByUserID(str, page, limit)
	if err != nil {
		return c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (a *HttpActivityHandler) GetByID(c echo.Context) error {
	idP, _ := strconv.Atoi(c.Param("id"))
	id := uint(idP)

	d, err := a.ActivityService.GetByID(id)
	if err != nil {
		return c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (a *HttpActivityHandler) UpdateStatus(c echo.Context) error {
	var activity models.Activity

	idP, _ := strconv.Atoi(c.Param("id"))
	id := uint(idP)

	statusP, _ := strconv.Atoi(c.Param("status"))
	status := int64(statusP)
	activity.ID = id

	d, err := a.ActivityService.UpdateStatus(activity, status)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (a *HttpActivityHandler) GetByUserToday(c echo.Context) error {
	idP, _ := a.Auth(c.Request().Header.Get("Auth"))
	idInter := idP["ID"]
	str := uint(idInter.(float64))

	pageS := c.QueryParam("page")
	page, _ := strconv.Atoi(pageS)
	limitS := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitS)

	d, err := a.ActivityService.GetByUserToday(str, page, limit)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, d)
}

func (a *HttpActivityHandler) GetCountSuccessByUser(c echo.Context) error {
	idP, _ := a.Auth(c.Request().Header.Get("Auth"))
	idInter := idP["ID"]
	str := uint(idInter.(float64))

	d, err := a.ActivityService.GetCountSuccessByUser(str)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, d)
}

/*
	HELPER
*/

func (a *HttpActivityHandler) Auth(header string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("password"), nil
	})

	if token != nil && err == nil {
		return token.Claims.(jwt.MapClaims), nil
	} else {
		return nil, err
	}
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

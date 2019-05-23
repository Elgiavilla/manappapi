package http

import (
	"fmt"
	"net/http"

	"github.com/elgiavilla/manapp/models"
	"github.com/elgiavilla/manapp/user"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	jwt "github.com/dgrijalva/jwt-go"
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

type HttpUserHandler struct {
	UserService user.Service
}

func NewUserHandler(e *echo.Echo, userService user.Service) {
	handler := &HttpUserHandler{
		UserService: userService,
	}
	e.POST("/user/regist", handler.Insert)
	e.POST("/user/login", handler.Login)
	e.PUT("/user", handler.Update)
	e.GET("/user/id", handler.GetByID)
}

func (u *HttpUserHandler) Insert(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	d, err := u.UserService.Insert(user)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	tkID := int64(d.ID)
	tk := &Token{ID: tkID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("password"))
	return c.JSON(http.StatusCreated, SuccessMessage{Token: tokenString})
}

func (u *HttpUserHandler) Login(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	d, err := u.UserService.Login(user)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	tkID := int64(d.ID)
	tk := &Token{ID: tkID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("password"))

	user.ID = d.ID
	_, err = u.UserService.Update(user)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, SuccessMessage{Token: tokenString})
}

func (u *HttpUserHandler) Update(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	idP, err := u.Auth(c.Request().Header.Get("Auth"))
	idInter := idP["ID"]
	str := uint(idInter.(float64))
	user.ID = str
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	d, err := u.UserService.Update(user)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, d)
}

func (u *HttpUserHandler) GetByID(c echo.Context) error {
	idP, err := u.Auth(c.Request().Header.Get("Auth"))
	idInter := idP["ID"]
	str := uint(idInter.(float64))

	list, err := u.UserService.GetFindById(str)
	if err != nil {
		c.JSON(getStatusCode(err), ReponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (u *HttpUserHandler) Auth(header string) (jwt.MapClaims, error) {
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

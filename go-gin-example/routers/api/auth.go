package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/EGGYC/go-gin-example/pkg/app"
	"github.com/EGGYC/go-gin-example/pkg/e"
	"github.com/EGGYC/go-gin-example/pkg/util"
	"github.com/EGGYC/go-gin-example/service/auth_service"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

//
//type auth struct {
//	Username string `valid:"Required; MaxSize(50)"`
//	Password string `valid:"Required; MaxSize(50)"`
//}
//
//// @Summary Get Auth
//// @Produce  json
//// @Param username query string true "userName"
//// @Param password query string true "password"
//// @Success 200 {object} app.Response
//// @Failure 500 {object} app.Response
//// @Router /auth [get]
//func GetAuth(c *gin.Context) {
//	appG := app.Gin{C: c}
//	valid := validation.Validation{}
//
//	username := c.PostForm("username")
//	password := c.PostForm("password")
//
//	a := auth{Username: username, Password: password}
//	ok, _ := valid.Valid(&a)
//
//	if !ok {
//		app.MarkErrors(valid.Errors)
//		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
//		return
//	}
//
//	authService := auth_service.Auth{Username: username, Password: password}
//	isExist, err := authService.Check()
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
//		return
//	}
//
//	if !isExist {
//		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
//		return
//	}
//
//	token, err := util.GenerateToken(username, password)
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
//		return
//	}
//
//	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
//		"token": token,
//	})
//}

//
//type auth struct {
//	Username string `valid:"Required; MaxSize(50)"`
//	Password string `valid:"Required; MaxSize(50)"`
//}
//
//func GetAuth(c *gin.Context) {
//	username := c.Query("username")
//	password := c.Query("password")
//
//	valid := validation.Validation{}
//	a := auth{Username: username, Password: password}
//	ok, _ := valid.Valid(&a)
//
//	data := make(map[string]interface{})
//	code := e.INVALID_PARAMS
//	if ok {
//		isExist := models.CheckAuth(username, password) // 转向models数据库的auth.go查询数据库中有无此用户
//		if isExist {
//			token, err := util.GenerateToken(username, password) // pkg/util/jwt.go中的GenerateToken函数
//			if err != nil {
//				code = e.ERROR_AUTH_TOKEN
//			} else {
//				data["token"] = token
//
//				code = e.SUCCESS
//			}
//
//		} else {
//			code = e.ERROR_AUTH
//		}
//	} else {
//		for _, err := range valid.Errors {
//			logging.Info(err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": data,
//	})
//}

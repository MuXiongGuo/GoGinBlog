// Package app 要让错误提前返回，c.JSON 的侵入是不可避免的，但是可以让其更具可变性，指不定哪天就变 XML 了呢？
package app

import (
	"github.com/gin-gonic/gin"

	"github.com/EGGYC/go-gin-example/pkg/e"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}

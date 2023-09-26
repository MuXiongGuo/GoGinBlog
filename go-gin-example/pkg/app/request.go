// Package app 要让错误提前返回，c.JSON 的侵入是不可避免的，但是可以让其更具可变性，指不定哪天就变 XML 了呢？
package app

import (
	"github.com/astaxie/beego/validation"

	"github.com/EGGYC/go-gin-example/pkg/logging"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}

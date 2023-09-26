// Package v1 编写路由空壳
package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/EGGYC/go-gin-example/pkg/app"
	"github.com/EGGYC/go-gin-example/pkg/e"
	"github.com/EGGYC/go-gin-example/pkg/export"
	"github.com/EGGYC/go-gin-example/pkg/logging"
	"github.com/EGGYC/go-gin-example/pkg/setting"
	"github.com/EGGYC/go-gin-example/pkg/util"
	"github.com/EGGYC/go-gin-example/service/tag_service"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add article tag
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddTagForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Update article tag
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Export article tag
// @Produce  json
// @Param name body string false "Name"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/export [post]
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// @Summary Import article tag
// @Produce  json
// @Param file body file true "Excel File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//
//// GetTags 编写标签列表的路由逻辑
//// 获取多个文章标签
//func GetTags(c *gin.Context) {
//	name := c.Query("name")
//
//	maps := make(map[string]interface{})
//	data := make(map[string]interface{})
//
//	if name != "" {
//		maps["name"] = name
//	}
//
//	var state int = -1
//	if arg := c.Query("state"); arg != "" {
//		state = com.StrTo(arg).MustInt() // 工具包 str先转化同等的 Strto 再调用已经实现的函数来转INT
//		maps["state"] = state
//	}
//
//	code := e.SUCCESS
//
//	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
//	data["total"] = models.GetTagTotal(maps)
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": data,
//	})
//}
//
////c.Query可用于获取?name=test&state=1这类URL参数，而c.DefaultQuery则支持设置一个默认值
////code变量使用了e模块的错误编码，这正是先前规划好的错误码，方便排错和识别记录
////util.GetPage保证了各接口的page处理是一致的
////c *gin.Context是Gin很重要的组成部分，可以理解为上下文，它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应
//
//// @Summary 新增文章标签
//// @Produce  json
//// @Param name query string true "Name"
//// @Param state query int false "State"
//// @Param created_by query int false "CreatedBy"
//// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
//// @Router /api/v1/tags [post]
//func AddTag(c *gin.Context) { // AddTag 新增文章标签
//	name := c.Query("name")
//	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
//	createdBy := c.Query("created_by")
//
//	valid := validation.Validation{}
//	valid.Required(name, "name").Message("名称不能为空")
//	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
//	valid.Required(createdBy, "created_by").Message("创建人不能为空")
//	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
//	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		if !models.ExistTagByName(name) {
//			code = e.SUCCESS
//			models.AddTag(name, state, createdBy)
//		} else {
//			code = e.ERROR_EXIST_TAG
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": make(map[string]string),
//	})
//}
//
//// @Summary 修改文章标签
//// @Produce  json
//// @Param id path int true "ID"
//// @Param name query string true "ID"
//// @Param state query int false "State"
//// @Param modified_by query string true "ModifiedBy"
//// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
//// @Router /api/v1/tags/{id} [put]
//func EditTag(c *gin.Context) { // EditTag 修改文章标签
//	id := com.StrTo(c.Param("id")).MustInt()
//	name := c.Query("name")
//	modifiedBy := c.Query("modified_by")
//
//	valid := validation.Validation{}
//
//	var state int = -1
//	if arg := c.Query("state"); arg != "" {
//		state = com.StrTo(arg).MustInt()
//		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
//	}
//
//	valid.Required(id, "id").Message("ID不能为空")
//	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
//	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
//	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		code = e.SUCCESS
//		if models.ExistTagByID(id) {
//			data := make(map[string]interface{})
//			data["modified_by"] = modifiedBy
//			if name != "" {
//				data["name"] = name
//			}
//			if state != -1 {
//				data["state"] = state
//			}
//
//			models.EditTag(id, data)
//		} else {
//			code = e.ERROR_NOT_EXIST_TAG
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": make(map[string]string),
//	})
//}
//
//// @Summary Delete article tag
//// @Produce  json
//// @Param id path int true "ID"
//// @Router /api/v1/tags/{id} [delete]
//func DeleteTag(c *gin.Context) { // DeleteTag 删除文章标签
//	id := com.StrTo(c.Param("id")).MustInt()
//
//	valid := validation.Validation{}
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		code = e.SUCCESS
//		if models.ExistTagByID(id) {
//			models.DeleteTag(id)
//		} else {
//			code = e.ERROR_NOT_EXIST_TAG
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": make(map[string]string),
//	})
//}

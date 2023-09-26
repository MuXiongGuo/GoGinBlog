package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/EGGYC/go-gin-example/pkg/app"
	"github.com/EGGYC/go-gin-example/pkg/e"
	"github.com/EGGYC/go-gin-example/pkg/qrcode"
	"github.com/EGGYC/go-gin-example/pkg/setting"
	"github.com/EGGYC/go-gin-example/pkg/util"
	"github.com/EGGYC/go-gin-example/service/article_service"
	"github.com/EGGYC/go-gin-example/service/tag_service"
)

// @Summary Get a single article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary Get multiple articles
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add article
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		CreatedBy:     form.CreatedBy,
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Update article
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

const (
	QRCODE_URL = "https://github.com/EDDYCJY/blog#gin%E7%B3%BB%E5%88%97%E7%9B%AE%E5%BD%95"
)

func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{C: c}
	article := &article_service.Article{}
	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}

//// GetArticle 获取单个文章
//func GetArticle(c *gin.Context) {
//	appG := app.Gin{c}
//	id := com.StrTo(c.Param("id")).MustInt()
//	valid := validation.Validation{}
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//
//	if valid.HasErrors() {
//		app.MarkErrors(valid.Errors)
//		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
//		return
//	}
//
//	articleService := article_service.Article{ID: id}
//	exists, err := articleService.ExistByID()
//	if err != nil {
//		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
//		return
//	}
//	if !exists {
//		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
//		return
//	}
//
//	article, err := articleService.Get()
//	if err != nil {
//		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
//		return
//	}
//
//	appG.Response(http.StatusOK, e.SUCCESS, article)
//}

//func GetArticle(c *gin.Context) {
//	id := com.StrTo(c.Param("id")).MustInt()
//
//	valid := validation.Validation{}
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//
//	code := e.INVALID_PARAMS
//	var data interface{}
//	if !valid.HasErrors() {
//		if models.ExistArticleByID(id) {
//			data = models.GetArticle(id)
//			code = e.SUCCESS
//		} else {
//			code = e.ERROR_NOT_EXIST_ARTICLE
//		}
//	} else {
//		for _, err := range valid.Errors {
//			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": data,
//	})
//}

//// 获取多个文章
//func GetArticles(c *gin.Context) {
//	data := make(map[string]interface{})
//	maps := make(map[string]interface{})
//	valid := validation.Validation{}
//
//	var state int = -1
//	if arg := c.Query("state"); arg != "" {
//		state = com.StrTo(arg).MustInt()
//		maps["state"] = state
//
//		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
//	}
//
//	var tagId int = -1
//	if arg := c.Query("tag_id"); arg != "" {
//		tagId = com.StrTo(arg).MustInt()
//		maps["tag_id"] = tagId
//
//		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
//	}
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		code = e.SUCCESS
//
//		data["lists"], _ = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
//		data["total"], _ = models.GetArticleTotal(maps)
//
//	} else {
//		for _, err := range valid.Errors {
//			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": data,
//	})
//}
//
//// 新增文章
//func AddArticle(c *gin.Context) {
//	tagId := com.StrTo(c.Query("tag_id")).MustInt()
//	title := c.Query("title")
//	desc := c.Query("desc")
//	content := c.Query("content")
//	createdBy := c.Query("created_by")
//	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
//
//	valid := validation.Validation{}
//	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
//	valid.Required(title, "title").Message("标题不能为空")
//	valid.Required(desc, "desc").Message("简述不能为空")
//	valid.Required(content, "content").Message("内容不能为空")
//	valid.Required(createdBy, "created_by").Message("创建人不能为空")
//	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		if models.ExistTagByID(tagId) {
//			data := make(map[string]interface{})
//			data["tag_id"] = tagId
//			data["title"] = title
//			data["desc"] = desc
//			data["content"] = content
//			data["created_by"] = createdBy
//			data["state"] = state
//
//			models.AddArticle(data)
//			code = e.SUCCESS
//		} else {
//			code = e.ERROR_NOT_EXIST_TAG
//		}
//	} else {
//		for _, err := range valid.Errors {
//			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": make(map[string]interface{}),
//	})
//}
//
//// 修改文章
//func EditArticle(c *gin.Context) {
//	valid := validation.Validation{}
//
//	id := com.StrTo(c.Param("id")).MustInt()
//	tagId := com.StrTo(c.Query("tag_id")).MustInt()
//	title := c.Query("title")
//	desc := c.Query("desc")
//	content := c.Query("content")
//	modifiedBy := c.Query("modified_by")
//
//	var state int = -1
//	if arg := c.Query("state"); arg != "" {
//		state = com.StrTo(arg).MustInt()
//		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
//	}
//
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
//	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
//	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
//	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
//	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		if ok, _ := models.ExistArticleByID(id); ok {
//			if models.ExistTagByID(tagId) {
//				data := make(map[string]interface{})
//				if tagId > 0 {
//					data["tag_id"] = tagId
//				}
//				if title != "" {
//					data["title"] = title
//				}
//				if desc != "" {
//					data["desc"] = desc
//				}
//				if content != "" {
//					data["content"] = content
//				}
//
//				data["modified_by"] = modifiedBy
//
//				models.EditArticle(id, data)
//				code = e.SUCCESS
//			} else {
//				code = e.ERROR_NOT_EXIST_TAG
//			}
//		} else {
//			code = e.ERROR_NOT_EXIST_ARTICLE
//		}
//	} else {
//		for _, err := range valid.Errors {
//			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
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
//// 删除文章
//func DeleteArticle(c *gin.Context) {
//	id := com.StrTo(c.Param("id")).MustInt()
//
//	valid := validation.Validation{}
//	valid.Min(id, 1, "id").Message("ID必须大于0")
//
//	code := e.INVALID_PARAMS
//	if !valid.HasErrors() {
//		if ok, _ := models.ExistArticleByID(id); ok {
//			models.DeleteArticle(id)
//			code = e.SUCCESS
//		} else {
//			code = e.ERROR_NOT_EXIST_ARTICLE
//		}
//	} else {
//		for _, err := range valid.Errors {
//			logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": code,
//		"msg":  e.GetMsg(code),
//		"data": make(map[string]string),
//	})
//}

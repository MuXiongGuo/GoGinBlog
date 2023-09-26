package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EGGYC/go-gin-example/pkg/e"
	"github.com/EGGYC/go-gin-example/pkg/logging"
	"github.com/EGGYC/go-gin-example/pkg/upload"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image") // c.Request.FormFile：获取上传的图片（返回提供的表单键的第一个文件）
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) { // CheckImageExt、CheckImageSize检查图片大小，检查图片后缀
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath) // CheckImage：检查上传图片所需（权限、文件夹）
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil { // SaveUploadedFile：保存图片
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	// 延长时间5s  测试 上传过程中 关闭会不会影响它  测试成功 不会影响!
	time.Sleep(5 * time.Second)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// Package upload 由原来基础的file工具包 封装一个专用于图片上传的包，包中的方法都是适用于图片
// 并且隐藏底层细节
package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/EGGYC/go-gin-example/pkg/file"
	"github.com/EGGYC/go-gin-example/pkg/logging"
	"github.com/EGGYC/go-gin-example/pkg/setting"
	"github.com/EGGYC/go-gin-example/pkg/util"
)

// GetImageFullUrl 获取图片完整访问URL
func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

// GetImageName 获取图片名称
func GetImageName(name string) string {
	ext := path.Ext(name)                     // 获取文件扩展名 如.txt
	fileName := strings.TrimSuffix(name, ext) // 去除后缀
	fileName = util.EncodeMD5(fileName)       // md5加密重新命名

	return fileName + ext
}

// GetImagePath 返回图片路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// GetImageFullPath 获取图片完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// CheckImageExt 检查图片后缀
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckImageSize 检查图片大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// CheckImage 检查图片
func CheckImage(src string) error {
	dir, err := os.Getwd() // 取当前工作目录
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	// 相对绝对都可以 相对格式为 runtime/logs/images
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src) // 为啥这个不加dir呢 为什么呀！
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

// 路径问题，工作路径，绝对相对路径

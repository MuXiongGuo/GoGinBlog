// Package logging 自定义log文件日志
// 解决目前错误信息只能在控制台输出的问题
package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/EGGYC/go-gin-example/pkg/file"
	"github.com/EGGYC/go-gin-example/pkg/setting"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd() // 相对路径和绝对路径都可以哦 用相对就是工作目录为根 绝对就是多个前缀
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

//func getLogFilePath() string {
//	return fmt.Sprintf("%s", setting.AppSetting.LogSavePath)
//}
//
//func getLogFileFullPath() string {
//	prefixPath := getLogFilePath()
//	suffixPath := fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
//
//	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
//}
//
//func openLogFile(filePath string) *os.File {
//	_, err := os.Stat(filePath) // 返回文件信息结构描述文件
//	switch {                    // 询问每个case 是 true的执行，false的跳过继续询问
//	case os.IsNotExist(err):
//		mkDir()
//	case os.IsPermission(err):
//		log.Fatalf("Permission :%v", err)
//	}
//
//	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		log.Fatalf("Fail to OpenFile :%v", err)
//	}
//
//	return handle
//}
//
//func mkDir() {
//	dir, _ := os.Getwd()
//	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
//	if err != nil {
//		panic(err)
//	}
//}

//os.Stat：返回文件信息结构描述文件。如果出现错误，会返回*PathError
//type PathError struct {
//	Op   string
//	Path string
//	Err  error
//}
//os.IsNotExist：能够接受ErrNotExist、syscall的一些错误，它会返回一个布尔值，能够得知文件不存在或目录不存在
//os.IsPermission：能够接受ErrPermission、syscall的一些错误，它会返回一个布尔值，能够得知权限是否满足
//os.OpenFile：调用文件，支持传入文件名称、指定的模式调用文件、文件权限，返回的文件的方法可以用于I/O。如果出现错误，则为*PathError。
//const (
//	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
//	O_RDONLY int = syscall.O_RDONLY // 以只读模式打开文件
//	O_WRONLY int = syscall.O_WRONLY // 以只写模式打开文件
//	O_RDWR   int = syscall.O_RDWR   // 以读写模式打开文件
//	// The remaining values may be or'ed in to control behavior.
//	O_APPEND int = syscall.O_APPEND // 在写入时将数据追加到文件中
//	O_CREATE int = syscall.O_CREAT  // 如果不存在，则创建一个新文件
//	O_EXCL   int = syscall.O_EXCL   // 使用O_CREATE时，文件必须不存在
//	O_SYNC   int = syscall.O_SYNC   // 同步IO
//	O_TRUNC  int = syscall.O_TRUNC  // 如果可以，打开时
//)
//os.Getwd：返回与当前目录对应的根路径名
//os.MkdirAll：创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error
//os.ModePerm：const定义ModePerm FileMode = 0777

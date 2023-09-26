package main

// 利用endless实现优雅重启
//import (
//	"fmt"
//	"log"
//	"syscall"
//
//	"github.com/fvbock/endless"
//
//	"github.com/EGGYC/go-gin-example/pkg/setting"
//	"github.com/EGGYC/go-gin-example/routers"
//)
//
//func main() {
//	// endless为了能优雅重启 代替了 原来的http.Serve 应该是内部已经封装了这个功能
//	endless.DefaultReadTimeOut = setting.ReadTimeout
//	endless.DefaultWriteTimeOut = setting.WriteTimeout
//	endless.DefaultMaxHeaderBytes = 1 << 20
//	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)
//
//	server := endless.NewServer(endPoint, routers.InitRouter()) // endless.NewServer 返回一个初始化的 endlessServer 对象
//	// 在 BeforeBegin 时输出当前进程的 pid，调用 ListenAndServe 将实际“启动”服务
//	server.BeforeBegin = func(add string) {
//		log.Printf("Actual pid is %d", syscall.Getpid()) // syscall.Getpid() 返回当前进程的 pid
//	}
//
//	err := server.ListenAndServe()
//	if err != nil {
//		log.Printf("Server err: %v", err)
//	}
//}

//package main
//
//import (
//	"fmt"
//	"net/http"
//
//	"github.com/EGGYC/go-gin-example/pkg/setting"
//	"github.com/EGGYC/go-gin-example/routers"
//)
//
//func main() {
//	router := routers.InitRouter()
//
//	s := &http.Server{
//		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
//		Handler:        router,
//		ReadTimeout:    setting.ReadTimeout,
//		WriteTimeout:   setting.WriteTimeout,
//		MaxHeaderBytes: 1 << 20,
//	}
//
//	s.ListenAndServe()
//}

//gin.Default()：返回Gin的type Engine struct{...}，里面包含RouterGroup，相当于创建一个路由Handlers，可以后期绑定各类的路由规则和函数、中间件等
//router.GET(...){...}：创建不同的HTTP方法绑定到Handlers中，也支持POST、PUT、DELETE、PATCH、OPTIONS、HEAD 等常用的Restful方法
//gin.H{...}：就是一个map[string]interface{}
//gin.Context：Context是gin中的上下文，它允许我们在中间件之间传递变量、管理流、验证JSON请求、响应JSON请求等，在gin中包含大量Context的方法，例如我们常用的DefaultQuery、Query、DefaultPostForm、PostForm等等

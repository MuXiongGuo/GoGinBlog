package main

//import (
//	"log"
//	"time"
//
//	"github.com/robfig/cron"
//
//	"github.com/EGGYC/go-gin-example/models"
//)
//
//func main() {
//	log.Println("Starting...")
//
//	c := cron.New()                   // 会根据本地时间创建一个新（空白）的 Cron job runner
//	c.AddFunc("5 * * * * *", func() { // AddFunc 会向 Cron job runner 添加一个 func ，以按给定的时间表运行
//		log.Println("Run models.CleanAllTag...")
//		models.CleanAllTag()
//	})
//	c.AddFunc("5 * * * * *", func() {
//		log.Println("Run models.CleanAllArticle...")
//		models.CleanAllArticle()
//	})
//
//	c.Start() // 在当前执行的程序中启动 Cron 调度程序。其实这里的主体是 goroutine + for + select + timer 的调度控制哦
//
//	t1 := time.NewTimer(time.Second * 10)
//	for {
//		select {
//		case <-t1.C:
//			t1.Reset(time.Second * 10)
//		}
//	}
//}

// time.NewTimer + for + select + t1.Reset
// 如果你是初学者，大概会有疑问，这是干嘛用的？
//（1）time.NewTimer
// 会创建一个新的定时器，持续你设定的时间 d 后发送一个 channel 消息
//（2）for + select
// 阻塞 select 等待 channel
//（3）t1.Reset
// 会重置定时器，让它重新开始计时
// 注：本文适用于 “t.C已经取走，可直接使用 Reset”。
// 总的来说，这段程序是为了阻塞主程序而编写的，希望你带着疑问来想，有没有别的办法呢？
// 有的，你直接 select{} 也可以完成这个需求 :)

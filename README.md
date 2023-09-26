# go-gin-blog
项目描述：项目核心为使用gin框架进行路由分配、gorm操控数据库搭建的博客系统，用户获取token后可进行图片上传、查看文章及tag、上传文章数据、删除文章及tag功能。
项目内容：项目中使用JWT进行身份验证、自定义文件日志、使用http.Server-Shutdown()进行优雅重启服务、配备Swagger API文档、部署到Docker中运行(构建Scratch超小镜像、Mysql挂载数据卷)、定制GORM Callbacks并搭配Cron定时任务实现软硬删除功能、实现了图片上传功能、实现了Redis缓存等功能。
涉及技术：gin、gorm、net/http、Docker、MySQL、Redis、Swagger、JWT、logger等。

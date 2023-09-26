# go-gin-blog
项目描述：项目核心为使用gin框架进行路由分配、gorm操控数据库搭建的博客系统，用户获取token后可进行图片上传、查看文章及tag、上传文章数据、删除文章及tag功能。
项目内容：项目中使用JWT鉴权包括生成token，利用gin接入中间件对token验证、自定义文件日志，对发生的错误操根据设置的错误编码将错误信息记录在日志中，并跳过此次错误操作继续保持程序运行、
使用http.Server-Shutdown()进行优雅重启服务、配备Swagger API文档、部署到Docker中运行(构建Scratch超小镜像、Mysql挂载数据卷)、
定制GORM Callbacks并搭配Cron定时任务实现软硬删除功能、实现了图片上传功能，使用封装的file包将图片保存在服务端，通过gin框架提供的StaticFS使图片可见、实现了Redis缓存，
对获取数据类的接口增加缓存设置，数据更新时删除相应缓存，以及给缓存设置过期时间来实现MySQL与Redis保持一致性等功能。。
涉及技术：gin、gorm、net/http、Docker、MySQL、Redis、Swagger、JWT、logger等。

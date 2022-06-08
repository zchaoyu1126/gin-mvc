package main

import (
	_ "gin-mvc/common/logger"
	"gin-mvc/router"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	// 配置路由中间件以及路由信息
	router.InitRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

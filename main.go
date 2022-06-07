package main

import (
	"fmt"
	"gin-mvc/common/logger"
	"gin-mvc/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

// 全局初始化
func init() {
	fmt.Println("hi")
}

func main() {
	fmt.Println("hi")
	fmt.Println("hi")
	fmt.Println("hi")
	fmt.Println("hi")
	fmt.Println("hi")
	fmt.Println("hi")
	logger.InitLog(zapcore.DebugLevel, logger.LOGFORMAT_CONSOLE, "all.log")
	r := gin.Default()

	// 配置路由中间件以及路由信息
	router.InitRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

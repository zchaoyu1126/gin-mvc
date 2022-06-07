package router

import (
	"gin-mvc/app/controller"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(r *gin.Engine) {
	// 配置全局的中间件信息
	r.Use(LoggerMiddlerWare(zap.L(), time.RFC3339, true))
	r.Use(RecoveryMiddleWare(zap.L(), true))

	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	apiRouter.GET("/user/", UserAuthMiddleWare("Query"), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

}

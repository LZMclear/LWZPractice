package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func Router(r *gin.Engine) {
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
}

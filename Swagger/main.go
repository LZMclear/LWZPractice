package main

import (
	"Swagger/router"
	"github.com/gin-gonic/gin"
)

// @title Go-site Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8889
// @BasePath

func main() {
	r := gin.New()
	router.Router(r)
	r.Run(":8080")
}

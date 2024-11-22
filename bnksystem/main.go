package main

import (
	"bnksystem/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.POST("/transfer", controller.GetTransfer)
	r.Run(":8080")
}

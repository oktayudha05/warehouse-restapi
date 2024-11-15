package main

import (
	"warehouse-restapi/controller"
	"warehouse-restapi/middleware"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	router.GET("/", controller.Home)
	router.POST("/registerKaryawan", controller.RegisterKaryawan)
	router.POST("/loginKaryawan", controller.LoginKaryawan)
	router.POST("/registerPengunjung", controller.RegisterPengunjung)
	router.POST("/barang", controller.PostBarang)
	router.GET("/barang", middleware.MiddlewareJwt(), controller.GetBarang)

	router.Run(":3000")
}
package main

import (
	"warehouse-restapi/controller"

	"github.com/gin-gonic/gin"
)


func main(){
	router := gin.Default()

	router.GET("/", controller.Home)
	router.POST("/barang", controller.PostBarang)
	router.GET("/barang", controller.GetBarang)

	router.Run(":3000")
}
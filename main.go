package main

import (
	"warehouse-restapi/controller"
	"warehouse-restapi/middleware"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	router.GET("/", controller.Home)

	karyawan := router.Group("/karyawan")
	{
		karyawan.POST("/register", controller.RegisterKaryawan)
		karyawan.POST("/login", controller.LoginKaryawan)
	}
	pengunjung := router.Group("/pengunjung")
	{
		pengunjung.POST("/register", controller.RegisterPengunjung)
		pengunjung.POST("/login", controller.LoginPengunjung)
	}
	barang := router.Group("/")
	{
		barang.POST("/barang",middleware.JwtAndAuthorization("karyawan"), controller.PostBarang)
		barang.GET("/barang", middleware.JwtAndAuthorization(), controller.GetBarang)
	}

	router.Run(":3000")
}
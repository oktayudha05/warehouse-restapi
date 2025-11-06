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
	barang := router.Group("/barang")
	{
		barang.GET("/", middleware.JwtAndAuthorization(), controller.GetBarang)
		barang.PUT("/", middleware.JwtAndAuthorization("karyawan"), controller.UpdateBarang)
		barang.POST("/", middleware.JwtAndAuthorization("karyawan"), controller.PostBarang)
		barang.DELETE("/", middleware.JwtAndAuthorization("karyawan"), controller.DeleteBarang)
	}

	router.Run(":3000")
}
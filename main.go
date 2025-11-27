package main

import (
	"warehouse-restapi/controller"
	"warehouse-restapi/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})
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

	router.Run("0.0.0.0:3333")
}

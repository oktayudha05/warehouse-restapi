package main

import (
	"warehouse-restapi/controller"
	"warehouse-restapi/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://192.168.56.1:5173", "http://localhost:5173", "https://warehouse-api.oyudha.me"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

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

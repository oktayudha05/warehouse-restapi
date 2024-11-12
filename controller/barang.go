package controller

import (
	"net/http"
	"warehouse-restapi/database"
	"warehouse-restapi/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var collBarang = database.Db.Collection("barang")
var validate = validator.New()

func Home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func PostBarang(c *gin.Context){
	ctx := c.Request.Context()
	var barang model.Barang
	err := c.BindJSON(&barang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal bind data"})
		return
	}
	err = validate.Struct(barang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "format data salah"})
		return
	}
	_, err = collBarang.InsertOne(ctx, barang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memasukan barang ke database"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "berhasil menambahkan barang ke database"})
}
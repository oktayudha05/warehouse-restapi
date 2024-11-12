package controller

import (
	"net/http"
	"warehouse-restapi/database"
	"warehouse-restapi/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
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

func GetBarang(c *gin.Context){
	ctx := c.Request.Context()
	var storeBarang []model.Barang
	cur, err := collBarang.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan barang dari database"})
		return
	}
	for cur.Next(ctx){
		var barang model.Barang
		cur.Decode(&barang)
		storeBarang = append(storeBarang, barang)
	}
	if len(storeBarang) == 0{
		c.JSON(http.StatusNoContent, gin.H{"message": "data belum ada"})
		return
	}
	c.IndentedJSON(http.StatusOK, storeBarang)
}
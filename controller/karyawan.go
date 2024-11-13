package controller

import (
	"net/http"
	"warehouse-restapi/database"
	"warehouse-restapi/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var collKaryawan = database.Db.Collection("karyawan")

func RegisterKaryawan(c *gin.Context){
	ctx := c.Request.Context()
	var postKaryawan model.Karyawan
	err := c.BindJSON(&postKaryawan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal bind data"})
		return
	}
	err = validate.Struct(postKaryawan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "format data salah"})
		return
	}
	count, err := collKaryawan.CountDocuments(ctx, bson.M{"username": postKaryawan.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mencari username"})
		return
	}
	if count > 0{
		c.JSON(http.StatusConflict, gin.H{"message": "username sudah ada"})
		return
	}
	_, err = collKaryawan.InsertOne(ctx, postKaryawan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal post data ke database"})
		return
	}
	c.IndentedJSON(http.StatusOK,gin.H{"message": "berhasil menambah karyawan", "data": postKaryawan})
}

func LoginKaryawan(c *gin.Context){

}
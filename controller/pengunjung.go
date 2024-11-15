package controller

import (
	"net/http"
	"warehouse-restapi/database"
	"warehouse-restapi/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collPengunjung = database.Db.Collection("pengunjung")

func RegisterPengunjung(c *gin.Context){
	ctx := c.Request.Context()
	var postData model.Pnegunjung
	err := c.BindJSON(&postData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal bind data"})
		return
	}
	err = validate.Struct(postData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data tidak sesuai format"})
		return
	}
	count, err := collPengunjung.CountDocuments(ctx, bson.M{"username": postData.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error ketika mencari data"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusFound, gin.H{"message": "username sudah digunakan"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "berhasil melakukan register akun", "data": postData})
}

func LoginPengunjung(c *gin.Context){
	ctx := c.Request.Context()
	var postData model.Pnegunjung
	c.BindJSON(&postData)
	err := validate.Struct(postData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data tidak sesuai format"})
		return
	}
	var pengunjung model.Pnegunjung
	err = collPengunjung.FindOne(ctx, bson.M{"username": postData.Username, "password": postData.Password}).Decode(&pengunjung)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			c.JSON(http.StatusBadRequest, gin.H{"message": "username atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan akun"})
		return
	}
}
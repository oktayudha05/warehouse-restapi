package controller

import (
	"net/http"
	"warehouse-restapi/database"
	"warehouse-restapi/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var collPengunjung = database.Db.Collection("pengunjung")

func RegisterPengunjung(c *gin.Context){
	ctx := c.Request.Context()
	var postPengunjung model.Pnegunjung
	err := c.BindJSON(&postPengunjung)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal bind data"})
		return
	}
	err = validate.Struct(postPengunjung)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data tidak sesuai format"})
		return
	}
	count, err := collPengunjung.CountDocuments(ctx, bson.M{"username": postPengunjung.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error ketika mencari data"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusFound, gin.H{"message": "username sudah digunakan"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "berhasil melakukan register akun", "data": postPengunjung})
}
package controller

import (
	"fmt"
	"net/http"
	"warehouse-restapi/database"
	"warehouse-restapi/middleware"
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

	validateCh := make(chan error)
	checkUsernameCh := make(chan error)
	go func(){
		err := validate.Struct(postData)
		validateCh <- err
	}()
	go func(){
		count, err := collPengunjung.CountDocuments(ctx, bson.M{"username": postData.Username})
		if err != nil {
			checkUsernameCh <- err
			return
		}
		if count > 0 {
			checkUsernameCh <- fmt.Errorf("username sudah ada")
			return
		}
		checkUsernameCh <- nil
	}()
	validateErr := <- validateCh
	checkUsernameErr := <- checkUsernameCh
	if validateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "data tidak sesuai format"})
		return
	}
	if checkUsernameErr != nil {
		if checkUsernameErr.Error() == "username sudah ada"{
			c.JSON(http.StatusConflict, gin.H{"message": "username sudah digunakan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mencari username"})
		return
	}

	_, err = collPengunjung.InsertOne(ctx, postData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendaftarkan akun"})
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
	token, err := middleware.GenerateJwt(pengunjung.Username, "pengunjung")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan token"})
		return
	}
	pengunjungRes := model.PengunjungRes{
		Username: pengunjung.Username,
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "berhasil login sebagai " + pengunjung.Username, "token": token, "data": pengunjungRes})
}
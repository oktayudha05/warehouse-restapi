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

var collKaryawan = database.Db.Collection("karyawan")

func RegisterKaryawan(c *gin.Context){
	ctx := c.Request.Context()
	var postKaryawan model.Karyawan
	err := c.BindJSON(&postKaryawan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal bind data"})
		return
	}

	validateCh := make(chan error)
	checkUsernameCh := make(chan error)
	go func(){
		err := validate.Struct(postKaryawan)
		validateCh <- err
	}()
	go func(){
		count, err := collKaryawan.CountDocuments(ctx, bson.M{"username": postKaryawan.Username})
		if err != nil{
			checkUsernameCh <- err
			return
		} 
		if count > 0{
			checkUsernameCh <- fmt.Errorf("username sudah ada")
			return
		}
		checkUsernameCh <- nil
	}()

	validateErr := <- validateCh
	checkUsernameErr := <- checkUsernameCh
	if validateErr != nil{
		c.JSON(http.StatusBadRequest, gin.H{"message": "format data salah"})
		return
	}
	if checkUsernameErr != nil {
		if checkUsernameErr.Error() == "username sudah ada"{
			c.JSON(http.StatusConflict, gin.H{"message": "username sudah ada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mencari username"})
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
	ctx := c.Request.Context()
	var reqKaryawan model.Karyawan
	err := c.BindJSON(&reqKaryawan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind reques data"})
		return
	}
	var karyawan model.Karyawan
	err = collKaryawan.FindOne(ctx, bson.M{"username": reqKaryawan.Username, "password": reqKaryawan.Password}).Decode(&karyawan)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			c.JSON(http.StatusBadRequest, gin.H{"message": "akun tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan user"})
		return
	}
	token, err := middleware.GenerateJwt(karyawan.Username, "karyawan")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan token"})
		return
	}
	karyawanRes := model.KaryawanRes{
		NamaKaryawan: karyawan.NamaKaryawan,
		Username: karyawan.Username,
		Jabatan: karyawan.Jabatan,
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "berhasil login sebagai " + karyawan.NamaKaryawan, "token": token, "data": karyawanRes})
}
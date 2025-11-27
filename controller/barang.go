package controller

import (
	"net/http"
	"strconv"
	"time"
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
	barang.TanggalMasukBarang = time.Now()
	err = validate.Struct(barang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "format data salah"})
		return
	}
	fillCariBarang := bson.M{"namabarang": barang.NamaBarang, "jenisbarang": barang.JenisBarang}
	result := collBarang.FindOne(ctx, fillCariBarang)
	if result.Err() == nil {
		barang.TanggalMasukBarang = time.Now()
		_, err = collBarang.UpdateOne(ctx, fillCariBarang, bson.M{"$inc": bson.M{"jumlah": barang.Jumlah}, "$set": bson.M{"tanggalmasukbarang": barang.TanggalMasukBarang},})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menambahkan jumlah"})
			return
		}
		c.IndentedJSON(http.StatusCreated, gin.H{"data": barang, "message": "berhasil menambahkan " + strconv.FormatInt(int64(barang.Jumlah), 10) + " " + barang.NamaBarang + " ke database"})
		return
	}
	_, err = collBarang.InsertOne(ctx, barang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memasukan barang ke database"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"data": barang, "message": "berhasil menambahkan barang ke database"})
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

func DeleteBarang(c *gin.Context){
	ctx := c.Request.Context()
	var reqBarang model.Barang
	err := c.BindJSON(&reqBarang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind data"})
		return
	}
	filter := bson.M{"namabarang": reqBarang.NamaBarang}
	_, err = collBarang.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menghapus data"})
		return
	}
}

func UpdateBarang(c *gin.Context){
	ctx := c.Request.Context()
	var reqBarang model.Barang
	err := c.BindJSON(&reqBarang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind data"})
		return
	}
	filter := bson.M{"namabarang": reqBarang.NamaBarang}
	update := bson.M{"$set": bson.M{
		"jenisbarang":       reqBarang.JenisBarang,
		"hargabarang":             reqBarang.HargaBarang,
		"jumlah":            reqBarang.Jumlah,
		"tanggalmasukbarang": time.Now(),
	}}
	_, err = collBarang.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengupdate data"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": reqBarang, "message": "berhasil mengupdate data barang"})
}
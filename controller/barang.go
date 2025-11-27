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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	defer cur.Close(ctx)
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
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	
	filter := bson.M{"_id": objectID}
	_, err = collBarang.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "berhasil menghapus barang"})
}

func UpdateBarang(c *gin.Context){
	ctx := c.Request.Context()
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	
	var reqBarang model.Barang
	err = c.BindJSON(&reqBarang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind data"})
		return
	}
	
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"namabarang": reqBarang.NamaBarang,
		"jenisbarang": reqBarang.JenisBarang,
		"hargabarang": reqBarang.HargaBarang,
		"jumlah": reqBarang.Jumlah,
		"tanggalmasukbarang": time.Now(),
	}}
	
	_, err = collBarang.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengupdate data"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": reqBarang, "message": "berhasil mengupdate data barang"})
}

// Endpoint untuk menambah stok (gunakan field jumlah dari Barang)
func TambahStokBarang(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	
	// Ambil jumlah yang ingin ditambahkan dari body
	var reqData struct {
		Jumlah int `json:"jumlah" validate:"required"`
	}
	err = c.BindJSON(&reqData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind data"})
		return
	}
	
	if reqData.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "jumlah harus lebih besar dari 0"})
		return
	}
	
	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"jumlah": reqData.Jumlah}}
	
	_, err = collBarang.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menambah stok"})
		return
	}
	
	// Ambil data barang terbaru
	var updatedBarang model.Barang
	err = collBarang.FindOne(ctx, filter).Decode(&updatedBarang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data barang terbaru"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": updatedBarang,
		"message": "berhasil menambah stok barang",
	})
}

// Endpoint untuk mengurangi stok (gunakan field jumlah dari Barang)
func KurangiStokBarang(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	
	// Ambil jumlah yang ingin dikurangi dari body
	var reqData struct {
		Jumlah int `json:"jumlah" validate:"required"`
	}
	err = c.BindJSON(&reqData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "gagal bind data"})
		return
	}
	
	if reqData.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "jumlah harus lebih besar dari 0"})
		return
	}
	
	// Cek jumlah stok saat ini
	var currentBarang model.Barang
	err = collBarang.FindOne(ctx, bson.M{"_id": objectID}).Decode(&currentBarang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "barang tidak ditemukan"})
		return
	}
	
	if currentBarang.Jumlah < reqData.Jumlah {
		c.JSON(http.StatusBadRequest, gin.H{"message": "stok tidak mencukupi"})
		return
	}
	
	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"jumlah": -reqData.Jumlah}}
	
	_, err = collBarang.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengurangi stok"})
		return
	}
	
	// Ambil data barang terbaru
	var updatedBarang model.Barang
	err = collBarang.FindOne(ctx, filter).Decode(&updatedBarang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data barang terbaru"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": updatedBarang,
		"message": "berhasil mengurangi stok barang",
	})
}
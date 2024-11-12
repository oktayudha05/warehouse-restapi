# Petunjuk penggunaan
***Untuk menjalankan server***
- Server akan berjalan pada port 3000, pastikan port 3000 tidak sedang digunakan. Anda dapat menggantinya pada file `main.go`
- Masuk folder server, buat .env file isikan dengan `MONGO_URI="mongodb://linkmongodb"`, dan `JWT_KEY="bebas_sebagai_secret_key"`
- Arahkan terminal ke folder server, lalu jalankan `go run main.go`
- Test API pada Postman
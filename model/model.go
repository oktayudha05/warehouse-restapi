package model

import "time"

type Barang struct {
	NamaBarang string `json:"nama" bson:"namabarang" validate:"required"`
	JenisBarang string `json:"jenis" bson:"jenisbarang" validate:"required"`
	HargaBarang int `json:"harga" bson:"hargabarang"`
	Jumlah int `json:"jumlah" bson:"jumlah"`
	TanggalMasukBarang time.Time `json:"tanggal_masuk" bson:"tanggalmasukbarang" validate:"required"`
}

type UpdateBarang struct {
	NamaBarang *string `json:"nama"`
	JenisBarang *string `json:"jenis"`
	HargaBarang *int `json:"harga"`
	Jumlah *int `json:"jumlah"`
}

type Karyawan struct {
	NamaKaryawan string `json:"nama" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Jabatan string `json:"jabatan" validate:"required"`
}
type KaryawanRes struct {
	NamaKaryawan string `json:"nama"`
	Username string `json:"username"`
	Jabatan string `json:"jabatan"`
}

type Pnegunjung struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type PengunjungRes struct {
	Username string `json:"username"`
}
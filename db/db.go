// db/db.go
package db

import (
	"log"
	"mesin-atm/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Conn *gorm.DB // Inisialisasi variabel Conn

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/atm_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	// Membuka koneksi ke MySQL dan menyimpan hasilnya ke dalam Conn
	Conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Jika berhasil, lakukan migrasi tabel
	err = Conn.AutoMigrate(&models.Account{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Migrasi gagal:", err)
	}

	log.Println("Database berhasil dimigrasi!")
}

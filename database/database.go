package database

import (
	"fmt"
	"log"
	"os"

	"learn_project/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB

func Connect() {

		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
    // Create the DSN (Data Source Name)
    // Create DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
	dbHost, dbUser, dbName, dbPassword, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database:", err)
	}

	log.Println("✅ Berhasil koneksi ke database!")

	// Jalankan migrasi otomatis
	autoMigrate()

}

// Fungsi untuk migrasi otomatis
func autoMigrate() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Gagal melakukan migrasi:", err)
	}
	log.Println("✅ Migrasi berhasil!")
}
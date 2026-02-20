package postgres

import (
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func NewDatabaseConnection() {
	var err error

	dsn := "host=localhost user=postgres password=123456 dbname=mydb port=5432 sslmode=disable"
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = Db.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Printf("Error during migration: %v", err)
		return
	}
	log.Println("Successfully connected to database")
}

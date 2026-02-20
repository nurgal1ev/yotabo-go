package postgres

import (
	"github.com/nurgal1ev/yotabo-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func NewDatabaseConnection(cfg Config) {
	var err error

	Db, err = gorm.Open(postgres.Open(cfg.URI()), &gorm.Config{})
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

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
	dsn := "postgres://postgres:password@localhost:5432/mydb?sslmode=disable"
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = Db.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		return
	}
	log.Println("Successfully connected to database")
}

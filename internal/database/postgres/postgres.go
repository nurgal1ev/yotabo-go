package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewDatabaseConnection() {
	var err error

	dsn := "postgres://postgres:password@localhost:5432/mydb?sslmode=disable"
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = Db.AutoMigrate()
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Successfully connected to database")
}

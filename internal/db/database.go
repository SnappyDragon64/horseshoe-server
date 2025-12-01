package db

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID           int    `gorm:"primary_key"`
	Username     string `gorm:"uniqueIndex;not null'"`
	PasswordHash string `gorm:"not null'"`
	CreatedAt    time.Time
}

var DB *gorm.DB

func Init(dbPath string) {
	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized successfully at:", dbPath)
}

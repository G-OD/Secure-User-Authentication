package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupDB() *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}

	return db
}

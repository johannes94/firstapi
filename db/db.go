package db

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	Id    uint    `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Name  *string `json:"name"`
	Price *uint   `json:"price"`
}

var DB *gorm.DB

func InitDB() *gorm.DB {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		println(err)
		panic("Failed to open the database connection")
	}

	database.AutoMigrate(&Product{})

	DB = database
	return database
}

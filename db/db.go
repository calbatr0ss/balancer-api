package db

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	environment = os.Getenv("GO_ENV")
	dbURL       = os.Getenv("DATABASE_URL")
)

func Setup() {
	var err error
	// Init gorm
	if environment == "dev" {
		DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("gorm failed to connect to the sqlite database")
		}
	} else {
		DB, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{})
		if err != nil {
			panic("gorm failed to connect to the mysql database")
		}
	}
}

// func GetAllRecords() (*[]models.Record, error) {
// 	var records []models.Record
//
// 	if result := dbc.DB.Find(&records); result.Error != nil {
// 		return nil, result.Error
// 	}
//
// 	return &records, nil
// }

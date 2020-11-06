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
	dbUser      = os.Getenv("DB_USER")
	dbPass      = os.Getenv("DB_PASS")
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
		dsn := dbUser + ":" + dbPass + "@/dbname..."
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("gorm failed to connect to the mysql database")
		}
	}
}

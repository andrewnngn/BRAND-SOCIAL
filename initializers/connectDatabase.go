package initializers

import (
	//"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := os.Getenv("DB_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//development only
	//database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

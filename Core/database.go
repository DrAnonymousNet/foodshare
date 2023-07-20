package core

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error
var DB *gorm.DB // Assume you have a GORM database connection


func SetupDatabase() (*gorm.DB, error){
	godotenv.Load(".env")

	dsn :=fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=UTC",
		 os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("SSL_MODE"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		return DB, err
	}
	return DB, nil

}
package db

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func Connect()(*gorm.DB, error) {
    port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

	portInt, err := strconv.Atoi(port)
    if err != nil {
        fmt.Println("Error converting port to integer:", err)
        return nil, err
    }

	/**
	Constructing connection string for connectiong with PostgreSQL instance
	*/

	connection_string := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=require",
    host, portInt, user, password, dbname)

	/**
	If there is error while connecting to database
	Return nil and error
	
	If there is no error while connecting to database
	Return db and nil 
	*/

	db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//db.AutoMigrate(&models.Categories{}, &models.Products{})
	
	return db, nil
	
}
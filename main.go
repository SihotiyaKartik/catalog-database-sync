package main

import (
	"ecommerce_store/catalogsync"
	"ecommerce_store/db"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error occured while loading .env file")
	}
}

func main(){
	/**
	Intializing the Gin router
	*/
    r := gin.Default()

	db, err := db.Connect()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgreSQL database successfully")
	
	catalogsync.FetchAndStore(db)

	r.Run(":8080")
}
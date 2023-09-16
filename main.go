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

	catalogsync.FetchAndStore()

	/**
	Intializing the Gin router
	*/
    r := gin.Default()

	_, err := db.Connect()
	if err != nil{
		log.Fatal(err)
	}

	// defer db.Close()

	fmt.Println("Connected to PostgreSQL database successfully")

	r.Run(":8080")
}
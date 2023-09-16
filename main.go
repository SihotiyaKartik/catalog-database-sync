package main

import (
	"ecommerce_store/catalogsync"
	"ecommerce_store/db"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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

	c := cron.New()
	_, e := c.AddFunc("0 2 * * *", func() {catalogsync.FetchAndStore(db)})

	if e != nil{
		fmt.Printf("Error while adding FetchAndStore function for cron job: %v", e)
	}

	c.Start()

	r.Run(":8080")
}
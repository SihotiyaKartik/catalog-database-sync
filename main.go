package main

import (
	"ecommerce_store/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ecommerce_store/catalog"
	"ecommerce_store/catalogsync"

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

	base_url := os.Getenv("EXTERNAL_BASE_URL")

	/**
	Creating an http client for listening request
	*/

	client := &http.Client {
		Timeout: 15 * time.Second,
	}

	c := cron.New()
	_, e := c.AddFunc("0 2 * * *", func() {catalogsync.FetchAndStore(db, base_url, client)})

	if e != nil{
		fmt.Printf("Error while adding FetchAndStore function for cron job: %v", e)
	}

	c.Start()

	
	r.GET("/shop/categories", func(c *gin.Context) {catalog.GetCategories(c, db)})
	r.GET("/shop/products", func(c *gin.Context) {catalog.GetProducts(c, db)})

	r.Run(":8080")
}
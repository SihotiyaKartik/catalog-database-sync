package catalogsync

import (
	"ecommerce_store/consumer"
	"ecommerce_store/models"
	"ecommerce_store/producer"
	"ecommerce_store/utils"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

/**
Configuring the color of logs printed in terminal
For better readability
*/

var green = color.New(color.FgGreen).PrintlnFunc()
var red = color.New(color.FgRed).PrintlnFunc()

func FetchAndStoreProducts(db *gorm.DB, base_url string, client *http.Client){

	pro, err := producer.GetProducer()
	if err != nil{
		red(err)
	}

	con, err := consumer.GetConsumer()
	if err != nil{
		red(err)
	}

	/**
	Handling our products message queue
	checking if there is old request present (failed previously)
	processing that request
	*/

	handleProductConsumer(con)

	/**
	Fetching category_id in batches from our database
	Fetching products associated with this category_id from admin database
	*/

	var count int64

	db.Table("categories").Count(&count)

	limit := 100
	offset := 0

	for offset < int(count) {
		var categories []models.Categories

		d := db.Table("categories").Offset(offset).Limit(limit).Find(&categories)
		if d.Error != nil{
			fmt.Println(d.Error)
		}

		for _, categ := range categories {
			
			products_page := 1

			for {
				productsResponse, err := getProductsFromAdminData(base_url, products_page, categ.SuperId, *client, pro)
				if err != nil{
					red(err)
				}

				if productsResponse.Products == nil {
					break
				}

				/**
				Saving prooducts in our database
				*/

				for _, product := range productsResponse.Products {
					utils.AddOrUpdateProductsData(utils.Product(product), categ.ID, db)
				}


			}

		}

		

	}

}


/**
part of Cron function scheduled every day for fetching categories
and storing it in our database
*/

func FetchAndStoreCategories(db *gorm.DB, base_url string, client *http.Client){

	/**
	Creating instacne of producer and consumer
	*/

	pro, err := producer.GetProducer()
	if err != nil{
		red(err)
	}

	con, err := consumer.GetConsumer()
	if err != nil{
		red(err)
	}

	/**
	First checking if there is some request present in category queue or not
	Which failed in previous process
	*/
	handleCategoryConsumer(con)

	categories_page := 1

	for {

		categoriesResponse, err := getCategoriesFromAdminData(base_url, categories_page, *client, pro)
		if err != nil{
			red(err)
		}

		if categoriesResponse.Categories == nil{
			break
		}

		for _, category := range categoriesResponse.Categories {

			/**
			Checking if the category which we want to add
			is already present in our database or not

			If it is not present, then only add the category and
			its coressponding products in our database
			*/
			_, err = utils.FindCategoriesBySuperId(db, category.Id, category.Name)
			
			 
		}

		categories_page++

	}
	
	
}

func FetchAndStore(db *gorm.DB, base_url string, client *http.Client){

	green("Cron functtion Started")

	FetchAndStoreCategories(db, base_url, client)

	FetchAndStoreProducts(db, base_url, client)

	green("Cron functtion completed")

}
	


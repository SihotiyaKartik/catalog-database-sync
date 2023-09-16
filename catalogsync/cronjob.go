package catalogsync

import (
	"net/http"
	"os"
	"time"

	"ecommerce_store/utils"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

/**
Cron function scheduled every daya for fetching categories and its products
and storing it in our database
*/

func FetchAndStore(db *gorm.DB){

	/**
	Configuring the color of logs printed in terminal
	For better readability
	*/
	green := color.New(color.FgGreen).PrintlnFunc()
	red := color.New(color.FgRed).PrintlnFunc()

	green("Cron function started")

	base_url := os.Getenv("EXTERNAL_BASE_URL")

	/**
	Creating an http client for listening request
	*/

	client := &http.Client {
		Timeout: 15 * time.Second,
	}

	/**
	For paginating the API in order to get all categories
	*/

	categories_page := 1

	for {

		categoriesResponse, err := getCategoriesFromAdminData(base_url, categories_page, *client)
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
			_, err = utils.FindCategoriesBySuperId(db, category.Id)
			if err != nil {
				
				/**
				Saving category record in our database
				*/
				categoryData, e := utils.AddCategoriesData(utils.Category(category), db)
				if e != nil {
					red(e)
				}

				/**
				For paginating the API in order to get all products
				*/

				products_page := 1

				for {
					productsResponse, err := getProductsFromAdminData(base_url, products_page, category.Id, *client)
					if err != nil{
						red(err)
					}

					if productsResponse.Products == nil{
						break
					}

					for _, product := range productsResponse.Products{
						/**
						Saving product record in our database
						*/
						if err := utils.AddProductsData(categoryData.ID, utils.Product(product), db); err != nil{
							red(err)
						}
					}

					products_page++
				}
				
			}
		}

		categories_page++

	}
	
	green("Cron functtion completed")
}


	


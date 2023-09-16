package catalogsync

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"ecommerce_store/utils"

	"github.com/fatih/color"
	"gorm.io/gorm"
)


func FetchAndStore(db *gorm.DB){

	/**
	Configuring the color of logs printed in terminal
	For better readability
	*/
	green := color.New(color.FgGreen).PrintlnFunc()
	red := color.New(color.FgRed).PrintlnFunc()

	green("Cron function started")

	base_url := os.Getenv("EXTERNAL_BASE_URL")

	client := &http.Client {
		Timeout: 15 * time.Second,
	}

	categories_page := 1

	for {

		categoriesResponse, err := getCategoriesFromAdminData(base_url, categories_page, *client)
		if err != nil{
			red(err)
		}

		if categoriesResponse.Categories == nil || categories_page > 3{
			break
		}

		for _, category := range categoriesResponse.Categories {
			_, err = utils.FindCategoriesBySuperId(db, category.Id)
			if err != nil {
				
				categoryData, e := utils.AddCategoriesData(utils.Category(category), db)
				if e != nil {
					red(e)
				}

				fmt.Println(categoryData)

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
						fmt.Println(product.Images)
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
	

}


	


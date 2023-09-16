package catalogsync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Category struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type CategoriesResponse struct {
	Page int `json:"page"`
	Categories []Category `json:"categories"`
}

func FetchAndStore(){

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

		categories_url := fmt.Sprintf("%s/task/categories?limit=100&page=%d", base_url, categories_page) 
		
		req, err := http.NewRequest("GET", categories_url, nil)
		if err != nil {
			red("Error occurred while creating GET categories API request:", err)
			return
		}

		req.Header.Set("x-api-key", os.Getenv("EXTERNAL_API_KEY"))
		
		resp, err := client.Do(req)
		if err != nil{
			red("Error occureed while making GET categories API request:", err)
			return		
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			red("Error occurrde while reading response body:", err)
			return
		}

		var categoriesResponse CategoriesResponse
		if err := json.Unmarshal(body, &categoriesResponse); err != nil{
			red("Error occurred while unmarshing categories response, err")
			return
		}

		if categoriesResponse.Categories == nil{
			break
		}

		categories_page++

	}
	

}


	


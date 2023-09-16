package catalogsync

import (
	"ecommerce_store/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Category struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type CategoriesResponse struct {
	Page int `json:"page"`
	Categories []Category `json:"categories"`
}

type Image struct {
	Href string `json:"href"`
}

type Product struct {
	CategoryId uint `json:"category_id"`
    Sku int `json:"sku"`
	Name string `json:"name"`
    SalePrice float64 `json:"salePrice"`
    Images []models.Image `json:"images"`
    Digital bool `json:"digital"`
    ShippingCost interface{} `json:"shippingCost"`
    Description string `json:"description"`
    CustomerReviewCount int `json:"customerReviewCount"`
}

type ProductsResponse struct {
	Page int `json:"page"`
	Products []Product `json:"products"`
}

func getCategoriesFromAdminData(base_url string, page int, client http.Client) (CategoriesResponse, error){

	var categoriesResponse CategoriesResponse

	categories_url := fmt.Sprintf("%s/task/categories?limit=1&page=%d", base_url, page)

	req, err := http.NewRequest("GET", categories_url, nil)
	if err != nil {
		return CategoriesResponse{}, fmt.Errorf("Error occurred while creating GET categories API request: %v", err)
	}

	req.Header.Set("x-api-key", os.Getenv("EXTERNAL_API_KEY"))
	
	resp, err := client.Do(req)
	if err != nil{
		return CategoriesResponse{}, fmt.Errorf("Error occureed while making GET categories API request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CategoriesResponse{}, fmt.Errorf("Error occurred while reading response body: %v", err)
	}

	if err := json.Unmarshal(body, &categoriesResponse); err != nil{
		return CategoriesResponse{}, fmt.Errorf("Error occurred while unmarshing categories response, err")
	}

	return categoriesResponse, nil
}

func getProductsFromAdminData(base_url string, page int, categoryId string, client http.Client) (ProductsResponse, error){

	var productsResponse ProductsResponse

	products_url := fmt.Sprintf("%s/task/products?categoryID=%s&page=%d&limit=100", base_url, categoryId, page)

	req, err := http.NewRequest("GET", products_url, nil)
	if err != nil {
		return ProductsResponse{}, fmt.Errorf("Error occurred while creating GET products API request: %v", err)
	}

	req.Header.Set("x-api-key", os.Getenv("EXTERNAL_API_KEY"))
	
	resp, err := client.Do(req)
	if err != nil{
		return ProductsResponse{}, fmt.Errorf("Error occureed while making GET products API request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProductsResponse{}, fmt.Errorf("Error occurred while reading response body: %v", err)
	}

	if err := json.Unmarshal(body, &productsResponse); err != nil{
		return ProductsResponse{}, fmt.Errorf("Error occurred while unmarshing products response, err")
	}

	return productsResponse, nil
}
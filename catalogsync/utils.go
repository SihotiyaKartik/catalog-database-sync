package catalogsync

import (
	"ecommerce_store/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type categoryKafka struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	Status string `json:"status"`
	Type string `json:"type"`
}

type productKafka struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	Status string `json:"status"`
	Type string `json:"type"`
}

type Category struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type CategoriesResponse struct {
	Page int `json:"page"`
	Categories []Category `json:"categories"`
}

/**
Making shippingCost type as interface{} for adjusting values between float and string
*/

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

var topic_category string = "cctupsea-categories"
var topic_product string = "cctupsea-products"

/**
Function for retrieving categories data from monk-commerce database
*/

func getCategoriesFromAdminData(base_url string, page int, client http.Client, pro *kafka.Producer) (CategoriesResponse, error){

	var categoriesResponse CategoriesResponse
	deliveryChan := make(chan kafka.Event)

	value := categoryKafka{
		Page: page,
		Type: "Categories",
		Limit: 100,
	}

	categories_url := fmt.Sprintf("%s/task/categories?limit=100&page=%d", base_url, page)

	req, err := http.NewRequest("GET", categories_url, nil)
	if err != nil {
		return CategoriesResponse{}, fmt.Errorf("Error occurred while creating GET categories API request: %v", err)
	}

	req.Header.Set("x-api-key", os.Getenv("EXTERNAL_API_KEY"))
	
	resp, err := client.Do(req)
	if err != nil{
		
		/**
		If any failure while making GET request, storing that request in my kafka category topic
		having 3 days rentention
		*/
		value.Status = "Failed"
		jsonVal, _ := json.Marshal(value)
		pro.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic_category, Partition: int32(kafka.PartitionAny)}, Value: jsonVal}, deliveryChan)
		
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

/**
Function for retrieving products data from monk-commerce database
*/
func getProductsFromAdminData(base_url string, page int, categoryId string, client http.Client, pro *kafka.Producer) (ProductsResponse, error){

	var productsResponse ProductsResponse
	deliveryChan := make(chan kafka.Event)

	value := productKafka{
		Page: page,
		Type: "Products",
		Limit: 100,
	}

	products_url := fmt.Sprintf("%s/task/products?categoryID=%s&page=%d&limit=100", base_url, categoryId, page)

	req, err := http.NewRequest("GET", products_url, nil)
	if err != nil {
		return ProductsResponse{}, fmt.Errorf("Error occurred while creating GET products API request: %v", err)
	}

	req.Header.Set("x-api-key", os.Getenv("EXTERNAL_API_KEY"))
	
	resp, err := client.Do(req)
	if err != nil{

		/**
		If any failure while making GET request, storing that request in my kafka products topic
		having 3 days rentention
		*/
		value.Status = "Failed"
		jsonVal, _ := json.Marshal(value)
		pro.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic_product, Partition: int32(kafka.PartitionAny)}, Value: jsonVal}, deliveryChan)

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

func handleCategoryConsumer(con *kafka.Consumer){

	con.Subscribe(topic_category, nil)

	// Create a channel for handling OS signals (e.g., SIGINT, SIGTERM)
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	run := true
	
	timeout := 2 * time.Second

	for run == true{
		select {
		case sig := <-sigchan:
			fmt.Println("Stopping categories consumer", sig)
			run = false
			
		default:
			msg, err := con.ReadMessage(timeout)
			if err == nil {
				d := string(msg.Value)
				fmt.Println(d)
				/**
				Adding logic for fetching and storing data for these categories request
				*/
			} else {
				//Stopping consumer
				run = false
				return
				
			}
		}
		
	}
	
	<-sigchan // Wait for OS signals to stop the consumer

}

func handleProductConsumer(con *kafka.Consumer){

	con.Subscribe(topic_product, nil)

	// Create a channel for handling OS signals (e.g., SIGINT, SIGTERM)
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	run := true
	
	timeout := 2 * time.Second

	for run == true{
		select {
		case sig := <-sigchan:
			fmt.Println("Stopping products consumer", sig)
			run = false
			
		default:
			msg, err := con.ReadMessage(timeout)
			if err == nil {
				d := string(msg.Value)
				fmt.Println(d)
				/**
				Adding logic for fetching and storing data for these products request
				*/
			} else {
				//Stopping consumer
				run = false
				return
				
			}
		}
		
	}
	
	<-sigchan // Wait for OS signals to stop the consumer

}
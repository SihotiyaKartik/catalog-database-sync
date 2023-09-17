package catalog

import (
	"ecommerce_store/models"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductsResponse struct {
	ID uint `json:"id"`
    CategoryId uint `json:"category_id"`
    Sku int `json:"sku"`
	Name string `json:"name"`
    SalePrice float64 `json:"salePrice"`
    Images json.RawMessage `json:"images"`
    Digital bool `json:"digital"`
    ShippingCost string `json:"shippingCost"`
    Description string `json:"description"`
    CustomerReviewCount int `json:"customerReviewCount"`
}

type Image struct {
    Href string `json:"href"`
}

func GetProducts(c *gin.Context, db *gorm.DB){

	/**
	Checking if x-api-key is present or not in headers of API
	For authorization purpose
	*/

	xApiKey := c.GetHeader("x-api-key")
	if xApiKey != os.Getenv("EXTERNAL_API_KEY"){
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	/**
	Getting page and limit from API url
	Setting their default valur if they are not present
	*/

	pageVal := c.DefaultQuery("page", "1")
	limitVal := c.DefaultQuery("limit", "10")
	categoryId := c.DefaultQuery("categoryId", "")

	/**
	Checking for any potential errors in categoryId, page and limit
	*/

	if categoryId == ""{
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide categoryId"})
		return
	}

	page, err := strconv.Atoi(pageVal)
	if err != nil || page < 1{
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page parameter"})
		return 
	}

	limit, err := strconv.Atoi(limitVal)
	if err != nil || limit < 1 || limit > 100{
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid limit parameter, should be between 1 and 100 inclusive"})
		return
	}

	/**
	Calculating total pages using total record and limit for a particular categoryId
	*/

	var totalRecord int64
	db.Model(&models.Products{}).Where("deleted_at IS NULL AND category_id = ?", categoryId).Count(&totalRecord)

	totalPages := int(math.Ceil(float64(totalRecord)/float64(limit)))

	/**
	Getting products for a particular categoryId sorted in descending order by reviews
	Configuring resultant data in ProductResponse data type
	*/

	var products []ProductsResponse
	offset := (page - 1)*limit
	db.Model(&models.Products{}).Where("category_id = ?", categoryId).Order("customer_review_count DESC").Offset(offset).Limit(limit).Find(&products)


	response := gin.H{
		"page": page,
		"products": products,
		"totalPages": totalPages,
	}

	c.JSON(http.StatusOK, response)
}
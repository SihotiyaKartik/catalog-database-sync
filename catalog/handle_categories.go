package catalog

import (
	"ecommerce_store/models"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryResult struct{
	ID uint `json:"id"`
	SuperId string `json:"superId"`
	Name string `json:"name"`
	TotalProducts int `json:"totalProducts"`
}

func GetCategories(c *gin.Context, db *gorm.DB){

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

	/**
	Checking for any potential errors in page and limit
	*/

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
	Calculating total pages using total record and limit
	*/

	var totalRecord int64
	db.Model(&models.Categories{}).Where("deleted_at IS NULL").Count(&totalRecord)

	totalPages := int(math.Ceil(float64(totalRecord)/float64(limit)))

	/**
	Getting categories, sorted in descending order by total products
	Configuring resultant data in CategoryResult data type
	*/

	var categories []CategoryResult
	offset := (page - 1)*limit
	db.Model(&models.Categories{}).Order("total_products DESC").Offset(offset).Limit(limit).Find(&categories)
	
	response := gin.H{
		"page": page,
		"categories":categories,
		"totalPages":totalPages,
	}

	c.JSON(http.StatusOK, response)

}
package utils

import (
	"ecommerce_store/models"
	"encoding/json"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Category struct {
	Id string `json:"id"`
	Name string `json:"name"`
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



func FindCategoriesBySuperId(db *gorm.DB, super_id string) (models.Categories, error){

	var category models.Categories

	if err := db.Where("super_id = ?", super_id).First(&category).Error; err != nil {
		return category, err
	} else {
		return category, nil
	}
}

func shippingCostToString(shippingCost interface{})(string){
	var shippingCostString string

	if val, ok := shippingCost.(float64); ok{
		shippingCostString = strconv.FormatFloat(val, 'f', -1, 64)
	} else if val, ok := shippingCost.(string); ok {
		shippingCostString = val
	} else {
		shippingCostString = ""
	}

	return shippingCostString
}

func AddCategoriesData(category Category, db *gorm.DB)(models.Categories, error){
	newCategory := models.Categories{
		Name: category.Name,
		SuperId: category.Id,
	}

	result := db.Create(&newCategory)
	if result.Error != nil{
		return models.Categories{}, fmt.Errorf("Error occurred while saving category record: %v", result.Error)
	}
	return newCategory, nil

}

func AddProductsData(categoryId uint, product Product, db *gorm.DB)(error){

	shippingCost := shippingCostToString(product.ShippingCost)
	imagesJSON, err := json.Marshal(product.Images)
    if err != nil {
        return fmt.Errorf("Error occurred while marshaling product images: %v", err)
    }

	newProduct := models.Products{
		Sku: product.Sku,
		Name: product.Name,
		SalePrice: product.SalePrice,
		CategoryId: categoryId,
		Images: imagesJSON,
		Digital: product.Digital,
		ShippingCost: shippingCost,
		Description: product.Description,
		CustomerReviewCount: product.CustomerReviewCount,
	}

	result := db.Create(&newProduct)
	if result.Error != nil{
		return fmt.Errorf("Error occurred while saving product record: %v", result.Error)
	}

	return nil

}
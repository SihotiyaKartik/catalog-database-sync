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

/**
function for checking if a particular category exist in our database through its super_id
If exist then update, otherwise add
*/

func FindCategoriesBySuperId(db *gorm.DB, super_id string, name string) (models.Categories, error){

	var category models.Categories

	if err := db.Where("super_id = ?", super_id).First(&category).Error; err != nil {
		/**
		We didnt found a record for this particular super_id
		so creating a new record
		*/
		newCategory := models.Categories{
            SuperId: super_id,
            Name: name,
            TotalProducts: 0,
        }

		if err := db.Create(&newCategory).Error; err != nil{
			return newCategory, err
		}
		return newCategory, nil
	}

	category.Name = name
	if err := db.Save(&category).Error; err != nil{
		return category, err
	}
	return category, nil

}

/**
Converting shippingCost into string
as shippingCost coming have multiple type
*/

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

/**
function for adding category record in our database
*/

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

/**
function for adding product record in our database
*/

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

func AddOrUpdateProductsData(product Product, id uint, db *gorm.DB){

	var prod models.Products

	shippingCost := shippingCostToString(product.ShippingCost)
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		fmt.Println("Error occurred while marshaling product images", err)
        return
    }

	if err := db.Where("sku = ?", product.Sku).First(&prod).Error; err != nil {

		/**
		We didnt found a record for this particular sku
		so creating a new record
		*/

		newProduct := models.Products{
			Sku: product.Sku,
			Name: product.Name,
			SalePrice: product.SalePrice,
			CategoryId: id,
			Images: imagesJSON,
			Digital: product.Digital,
			ShippingCost: shippingCost,
			Description: product.Description,
			CustomerReviewCount: product.CustomerReviewCount,
		}

		db.Create(&newProduct)
	}

	/**
	Updating the existing product record 
	*/

	prod.Name = product.Name
	prod.SalePrice = product.SalePrice
	prod.Description = product.Description
	prod.CustomerReviewCount = product.CustomerReviewCount
	prod.Digital = product.Digital
	prod.ShippingCost = shippingCost

	if err := db.Save(&prod).Error; err != nil {
		fmt.Println("Error occurred while updating products data", err)
		return
	}

}
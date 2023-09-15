package models

import (
	"github.com/jinzhu/gorm"
)

type Products struct {
	gorm.Model
    Category Categories `gorm:"foreignKey:CategoryId"`
    CategoryId uint `json:"category_id"`
    Sku int `json:"sku"`
	Name string `json:"name"`
    SalePrice float64 `json:"salePrice"`
    Images []map[string] string `gorm:"type:jsonb" json:"images"`
    Digital bool `json:"digital"`
    ShippingCost float64 `json:"shippingCost"`
    Description string `json:"description"`
    CustomerReviewCount int `json:"customerReviewCount"`
}

package models

import (
	"gorm.io/gorm"
)

type Image struct{
    Href string `json:"href"`
}

type Products struct {
	gorm.Model
    Category Categories `gorm:"foreignKey:CategoryId"`
    CategoryId uint `json:"category_id"`
    Sku int `json:"sku"`
	Name string `json:"name"`
    SalePrice float64 `json:"salePrice"`
    Images []byte `gorm:"type:jsonb" json:"images"`
    Digital bool `json:"digital"`
    ShippingCost string `json:"shippingCost"`
    Description string `json:"description"`
    CustomerReviewCount int `json:"customerReviewCount"`
}

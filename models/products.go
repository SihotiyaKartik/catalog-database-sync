package models

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
    Sku int `json:"sku"`
	Name string `json:"name"`
    SalePrice float64 `json:"salePrice"`
    Images []struct {
        Href string `json:"href"`
    } `json:"images"`
    Digital bool `json:"digital"`
    ShippingCost float64 `json:"shippingCost"`
    Description string `json:"description"`
    CustomerReviewCount int `json:"customerReviewCount"`
}

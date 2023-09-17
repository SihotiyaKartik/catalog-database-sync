package models

import "gorm.io/gorm"

type Categories struct {
	gorm.Model `json:"-"`
	SuperId string `json:"super_id" gorm:"unique"`
	Name string `json:"name"`
	TotalProducts int `json:"totalProducts"`
}


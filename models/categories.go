package models

import "github.com/jinzhu/gorm"

type Categories struct {
	gorm.Model
	SuperId string `json:"super_id"`
	Name string `json:"name"`
}


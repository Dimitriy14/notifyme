package models

import "github.com/jinzhu/gorm"

type ProductFiler struct {
	gorm.Model
	UserEmail   string `json:"user_email"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Count       string `json:"count"`
	PayedSum    string `json:"payed_sum"`
}

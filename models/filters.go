package models

type ProductFiler struct {
	UserEmail   string `json:"user_email" gorm:"primary_key"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Count       string `json:"count"`
	PayedSum    string `json:"payed_sum"`
}

package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string  `json:"name" gorm:"size:100;not null"`
	Price         float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	StockQuantity int     `json:"stock_quantity" gorm:"default:0"`
	Orders        []Order `gorm:"foreignKey:ProductID"`
}

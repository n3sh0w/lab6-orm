package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

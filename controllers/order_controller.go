package controllers

import (
	"go-postgres-orm/database"
	"go-postgres-orm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("User").Preload("Product").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func CreateOrder(c *gin.Context) {
	var input struct {
		UserID    uint `json:"user_id" binding:"required"`
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	if err := tx.First(&user, input.UserID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var product models.Product
	if err := tx.First(&product, input.ProductID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	if product.StockQuantity < input.Quantity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock"})
		return
	}

	order := models.Order{
		UserID:    input.UserID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	if err := tx.Model(&models.Product{}).Where("id = ?", product.ID).
		Update("stock_quantity", product.StockQuantity-input.Quantity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
		return
	}

	var orderWithDetails models.Order
	if err := tx.Preload("User").Preload("Product").First(&orderWithDetails, order.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load order details"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orderWithDetails})
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, order.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	quantityDiff := input.Quantity - order.Quantity

	if product.StockQuantity < quantityDiff {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock"})
		return
	}

	database.DB.Model(&order).Update("quantity", input.Quantity)

	database.DB.Model(&product).Update("stock_quantity", product.StockQuantity-quantityDiff)

	var updatedOrder models.Order
	database.DB.Preload("User").Preload("Product").First(&updatedOrder, order.ID)

	c.JSON(http.StatusOK, gin.H{"data": updatedOrder})
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, order.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	database.DB.Delete(&order)

	database.DB.Model(&product).Update("stock_quantity", product.StockQuantity+order.Quantity)

	c.JSON(http.StatusOK, gin.H{"data": "Order deleted successfully"})
}

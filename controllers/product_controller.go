package controllers

import (
	"go-postgres-orm/database"
	"go-postgres-orm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	database.DB.Preload("Orders").Preload("Orders.User").Find(&products)

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func CreateProduct(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:          input.Name,
		Price:         input.Price,
		StockQuantity: input.StockQuantity,
	}
	database.DB.Create(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		StockQuantity int     `json:"stock_quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&product).Updates(models.Product{
		Name:          input.Name,
		Price:         input.Price,
		StockQuantity: input.StockQuantity,
	})

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	database.DB.Delete(&product)

	c.JSON(http.StatusOK, gin.H{"data": "Product deleted successfully"})
}

package main

import (
	"go-postgres-orm/database"
	"go-postgres-orm/models"
	"log"
)

func main() {
	database.Connect()

	database.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	log.Println("migrate the schemas finished")
}

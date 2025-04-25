package main

import (
	"go-postgres-orm/database"
	"go-postgres-orm/routes"
)

func main() {
	database.Connect()

	r := routes.SetupRoutes()

	r.Run(":8080")
}

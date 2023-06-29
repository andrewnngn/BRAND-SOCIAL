package main

import (
	"github.com/Cedar-81/swype/initializers"
	"github.com/Cedar-81/swype/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Post{})
	initializers.DB.AutoMigrate(&models.Comment{})
	initializers.DB.AutoMigrate(&models.Notification{})
}

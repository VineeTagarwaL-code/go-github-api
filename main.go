package main

import (
	"go-github-api/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env")
	}
	routes.SetUpHandler(r)
	r.Run(":" + os.Getenv("PORT"))
}

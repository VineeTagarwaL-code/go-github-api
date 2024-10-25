package main

import (
	"go-github-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.SetUpHandler(r)
	r.Run(":3000")
}

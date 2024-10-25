package routes

import (
	"go-github-api/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpHandler(r *gin.Engine) {
	r.GET("/api/:name", handlers.GithubHandler)
}

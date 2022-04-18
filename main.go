package main

import (
	"github.com/eaaslan/go-jwt/routes"
	"github.com/gin-gonic/gin"

	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	router := gin.New()
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.Use(gin.Logger())

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "access granted for api-2"})
	})

	err := router.Run(port)
	if err != nil {
		return
	}

}

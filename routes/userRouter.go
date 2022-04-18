package routes

import (
	"github.com/eaaslan/go-jwt/controllers"
	"github.com/eaaslan/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/user/:user_id", controllers.GetUser())

}

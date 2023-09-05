package routes

import (
	"github.com/gin-gonic/gin"
)

// Regsiter user-store routes
func RegisterUserStoreRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/users")

	// public routes
	routes.GET("/")
	routes.GET("/:user_id", userController.GetUserById)

	// auth routes
	authRoutes := routes.Group("/auth")
	authRoutes.POST("/sign-up", userController.SignUp)
	authRoutes.POST("/login", userController.Login)
	authRoutes.GET("/logout", authMiddleware.Authenticate, userController.Logout)
	authRoutes.GET("/me", userController.GetMe)
	authRoutes.GET("/verify-email/:verification_code", userController.VerifyEmail)

	// private routes
	privateRoutes := routes.Group("/", authMiddleware.Authenticate)
	privateRoutes.POST("/")
	privateRoutes.PUT("/:user_id")
	privateRoutes.DELETE("/:user_id")
}

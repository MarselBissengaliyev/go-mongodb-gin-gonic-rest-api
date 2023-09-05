package routes

import (
	"github.com/gin-gonic/gin"
)

// Register post-store routes
func RegisterPostStoreRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/posts")

	// public routes
	routes.GET("/", postController.GetPosts)
	routes.GET("/:post_id", postController.GetPostById)

	// private routes
	privateRoutes := routes.Group("/", authMiddleware.Authenticate)
	privateRoutes.POST("/", postController.CreatePost)
	privateRoutes.PUT("/:post_id", postController.UpdatePostById)
	privateRoutes.DELETE("/:post_id", postController.DeletePostById)
}

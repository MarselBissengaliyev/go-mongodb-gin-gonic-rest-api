package routes

import "github.com/gin-gonic/gin"

// Register like-store routes
func RegisterLikeStoreRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/posts/:post_id/likes")

	// public routes
	routes.GET("/", likeController.GetLikesOfPost)

	// private routes
	privateRoutes := routes.Group("/", authMiddleware.Authenticate)
	privateRoutes.POST("/", likeController.CreateLike)
	privateRoutes.DELETE("/:like_id", likeController.DeleteLike)
}
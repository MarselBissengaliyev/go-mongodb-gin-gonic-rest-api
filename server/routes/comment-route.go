package routes

import "github.com/gin-gonic/gin"

// Register comment-store routes
func RegisterCommentStoreRoute(rg *gin.RouterGroup) {
	routes := rg.Group("/posts/:post_id/comments")

	// public routes
	routes.GET("/", commentController.GetComments)

	// private routes
	privateRoutes := routes.Group("/", authMiddleware.Authenticate)
	privateRoutes.POST("/", commentController.CreateComment)
	privateRoutes.PUT("/:comment_id", commentController.UpdateCommentById)
	privateRoutes.DELETE("/:comment_id", commentController.DeleteCommentById)
}

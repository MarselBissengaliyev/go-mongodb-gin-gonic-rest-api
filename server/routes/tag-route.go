package routes

import "github.com/gin-gonic/gin"

// Register tag-store routes
func RegisterTagStoreRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/posts/:post_id/tags")
	routes.GET("/", tagController.GetTagsOfPost)
}

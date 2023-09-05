package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagController struct{}

// Get tags handler
func (tc *TagController) GetTagsOfPost(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var tags []bson.M

	postId := c.Params.ByName("post_id")
	postObjectId, _ := primitive.ObjectIDFromHex(postId)

	if err := postCollection.FindOne(ctx, bson.M{"_id": postObjectId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error() + "\n post not found",
		})
		return
	}

	cursor, err := tagCollection.Find(ctx, bson.M{"post_id": postId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if err := cursor.All(ctx, &tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tags,
	})
}

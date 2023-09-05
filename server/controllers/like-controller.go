package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeController struct{}

// Get like of posts handler
func (lk *LikeController) GetLikesOfPost(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var post models.Post
	var likes []bson.M

	postId, err := primitive.ObjectIDFromHex(c.Params.ByName("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Decode(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cursor, err := likeCollection.Find(ctx, bson.M{
		"post_id": post.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := cursor.All(ctx, &likes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, likes)
}

// Create like handler
func (lk *LikeController) CreateLike(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var like models.Like
	var uid = fmt.Sprint(c.Keys["uid"])

	validate.Struct(like)

	postId, err := primitive.ObjectIDFromHex(c.Params.ByName("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.BindJSON(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	like.UserId = &uid

	insertResult, err := likeCollection.InsertOne(ctx, like)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, insertResult)
}

// Delete like handler
func (lk *LikeController) DeleteLike(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var uid = fmt.Sprint(c.Keys["uid"])
	var userType = fmt.Sprint(c.Keys["user_type"])
	var like models.Like

	likeId, err := primitive.ObjectIDFromHex(c.Params.ByName("like_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := likeCollection.FindOne(ctx, bson.M{"_id": likeId}).Decode(&like); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if *like.UserId != uid || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You don't have rights for this",
		})
		return
	}

	deleteResult, err := likeCollection.DeleteOne(ctx, bson.M{"_id": likeId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"error": deleteResult.DeletedCount,
	})
}

package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController struct{}

// Get comments handler
func (cc *CommentController) GetComments(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	postId := c.Query("post_id")

	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	var comments []bson.M

	cursor, err := commentCollection.Find(
		ctx,
		bson.M{"post_id": postId},
		mongoHelper.NewMongoPaginate(limit, page).GetPaginatedOpts(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := cursor.All(ctx, &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Create comment handler
func (cc *CommentController) CreateComment(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	postId := c.Query("post_id")
	userId := fmt.Sprint(c.Keys["uid"])

	var comment models.Comment

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	comment.PostId = &postId
	comment.UserId = &userId

	insertResult, insertErr := commentCollection.InsertOne(ctx, comment)
	if insertErr != nil {
		log.Panic(insertErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": insertErr.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, insertResult)
}

// Update comment by id handler
func (cc *CommentController) UpdateCommentById(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	postId, _ := primitive.ObjectIDFromHex(c.Params.ByName("post_id"))
	userId := fmt.Sprint(c.Keys["uid"])
	userType := fmt.Sprint(c.Keys["user_type"])
	commentId, _ := primitive.ObjectIDFromHex(c.Params.ByName("comment_id"))

	var updateCommentObj models.Comment
	var comment models.Comment

	if err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error() + "\n post not found",
		})
		return
	}

	if err := commentCollection.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error() + "\n comment not found",
		})
		return
	}

	if comment.UserId != &userId || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you don't have the right to do that",
		})
		return
	}

	updateResult, updateErr := commentCollection.UpdateByID(ctx, commentId, updateCommentObj)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": updateErr.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, updateResult.ModifiedCount)
}

// delete comment by id
func (cc *CommentController) DeleteCommentById(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var comment models.Comment

	postId, _ := strconv.Atoi(c.Params.ByName("post_id"))
	commentId, _ := primitive.ObjectIDFromHex(c.Params.ByName("comment_id"))
	userId := fmt.Sprint(c.Query("uid"))
	userType := fmt.Sprint(c.Keys["user_type"])

	if err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error() + "\n post not found",
		})
		return
	}

	if err := commentCollection.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error() + "\n comment not found",
		})
		return
	}

	if comment.UserId != &userId || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you don't have the right to do that",
		})
		return
	}

	deleteResult, deleteErr := commentCollection.DeleteOne(ctx, bson.M{"_id": commentId})
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": deleteErr.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, deleteResult.DeletedCount)
}

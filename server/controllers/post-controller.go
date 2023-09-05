package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostController struct{}

// Get posts
func (pc *PostController) GetPosts(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	page, _ := strconv.Atoi(c.Query("page"))
	if c.Query("page") == "" {
		page = 1
	}

	var posts []bson.M

	// Find posts
	cursor, err := postCollection.Find(
		ctx,
		bson.M{},
		mongoHelper.NewMongoPaginate(20, page).GetPaginatedOpts(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error() + "\n suka",
		})
		fmt.Println(err)
		return
	}

	if err := cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error() + "\n blyad",
		})
		fmt.Println(err)
		return
	}

	fmt.Println(posts)
	c.JSON(http.StatusOK, posts)
}

// Get post by id
func (pc *PostController) GetPostById(c *gin.Context) {
	postId := c.Params.ByName("post_id")
	postObjectId, _ := primitive.ObjectIDFromHex(postId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var post models.Post

	if err := postCollection.FindOne(ctx, bson.M{"_id": postObjectId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, bson.M{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	if err := postCollection.FindOne(ctx, bson.M{"_id": postObjectId}).Decode(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	views := *post.Views + 1

	post.Views = &views

	_, updateErr := postCollection.UpdateOne(ctx, bson.M{"_id": postId}, post)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  updateErr.Error(),
		})
		fmt.Println(updateErr)
		return
	}

	fmt.Println(post)
	c.JSON(http.StatusOK, post)
}

// Create post
func (pc *PostController) CreatePost(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var uid = fmt.Sprint(c.Keys["uid"])
	defer cancel()

	var post models.Post
	var tags []models.Tag

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err)
		return
	}

	tags = post.Tags
	views := 0

	createdAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	post = models.Post{
		UserId:    &uid,
		Content:   post.Content,
		Title:     post.Title,
		Views:     &views,
		Preview:   post.Preview,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	validationErr := validate.Struct(post)

	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  validationErr.Error(),
		})
		fmt.Println(validationErr)
		return
	}

	validationTagErr := tagHelper.ValidateTags(tags)
	if validationTagErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  validationTagErr.Error(),
		})
		fmt.Println(validationTagErr)
		return
	}

	insertedPost, insertErr := postCollection.InsertOne(ctx, post)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  insertErr.Error(),
		})
		fmt.Println(insertErr)
		return
	}

	_, insertTagsErr := tagCollection.InsertMany(
		ctx,
		tagHelper.MakeSliceOfInterfacesFromTags(tags),
	)

	if insertTagsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  insertErr.Error(),
		})
		fmt.Println(insertErr)
		return
	}

	c.JSON(http.StatusCreated, insertedPost)
}

// Update post by id
func (pc *PostController) UpdatePostById(c *gin.Context) {
	postId := c.Params.ByName("post_id")
	postObjectId, _ := primitive.ObjectIDFromHex(postId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var post models.Post
	var tags []models.Tag

	uid := fmt.Sprint(c.Keys["uid"])
	userType := fmt.Sprint(c.Keys["user_type"])

	if uid != postId || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error":  "You don't have rights for do this",
		})
		return
	}

	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	tags = post.Tags

	if err := tagHelper.ValidateTags(tags); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if err := validate.Struct(post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if err := postCollection.FindOne(ctx, bson.M{"_id": postObjectId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, bson.M{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	updatePostResult, err := postCollection.UpdateByID(ctx, postId, bson.M{
		"$set": bson.M{
			"author":  post.UserId,
			"content": post.Content,
			"title":   post.Title,
			"preview": post.Preview,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	_, err = tagCollection.UpdateMany(
		ctx,
		bson.M{"post_id": updatePostResult.UpsertedID},
		tagHelper.MakeSliceOfInterfacesFromTags(tags),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "Post succefuly updated",
	})
}

// Delete post by id
func (pc *PostController) DeletePostById(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	postId, err := primitive.ObjectIDFromHex(c.Params.ByName("post_id"))
	uid := fmt.Sprint(c.Keys["uid"])
	userType := fmt.Sprint(c.Keys["user_type"])

	if uid != postId.Hex() || userType != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error":  "You don't have rights for do this",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error() + "\n not valid post id",
		})
		fmt.Println(err)
		return
	}

	if err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Err(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	_, err = postCollection.DeleteOne(ctx, bson.M{"_id": postId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	_, deleteTagErr := tagCollection.DeleteMany(ctx, bson.M{"post_id": postId})
	if deleteTagErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "post succefully deleted",
	})
}

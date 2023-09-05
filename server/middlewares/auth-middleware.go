package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthMiddleware struct{}

// authenticate middleware
func (m *AuthMiddleware) Authenticate(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User

	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No authorization header provided",
		})
		c.Abort()
		return
	}

	claims, err := tokenHelper.ValidateToken(clientToken)
	if err != "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}

	uid, _ := primitive.ObjectIDFromHex(claims.Uid)

	if err := userCollection.FindOne(ctx, bson.M{
		"_id": uid,
	}).Decode(&user); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "the user belonging to this token no logger exists",
		})
		c.Abort()
		return
	}

	if !user.IsEmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"status": "failed", "message": "Please verify your email"})
		return
	}

	c.Set("email", claims.Email)
	c.Set("first_name", claims.FirstName)
	c.Set("last_name", claims.LastName)
	c.Set("uid", claims.Uid)
	c.Set("user_type", claims.UserType)
	c.Next()
}
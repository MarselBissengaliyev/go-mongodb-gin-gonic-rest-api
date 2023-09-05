package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/config"
	"github.com/MarselBisengaliev/go-react-blog/helpers"
	"github.com/MarselBisengaliev/go-react-blog/models"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct{}

// sign up handler
func (uc *UserController) SignUp(c *gin.Context) {
	conf, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal("could not load config + \n", err)
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		fmt.Println(err)
		return
	}

	if validationErr := validate.Struct(user); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  validationErr.Error(),
		})
		fmt.Println(validationErr)
		return
	}

	password := authHelper.HashPassword(*user.Password)
	user.Password = &password

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "\n error occured while checking email \n" + err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"erorr":  "this email already exists",
		})
		fmt.Println(err)
		return
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()

	userType := "USER"
	user.UserType = &userType

	token, refreshToken, generateTokenErr := tokenHelper.GenerateAllTokens(
		*user.Email,
		*user.FirstName,
		*user.LastName,
		*user.UserType,
		user.ID.Hex(),
	)

	if generateTokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "error occured while generate tokens \n" + generateTokenErr.Error(),
		})
		return
	}

	user.Token = &token
	user.RefreshToken = &refreshToken
	user.IsEmailVerified = false
	code := randstr.String(20)
	verificationCode := encodeHelper.Encode(code)
	user.VerificationCode = &verificationCode

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "User item was not created\n" + insertErr.Error(),
		})
		fmt.Println(insertErr)
		return
	}

	emailData := helpers.EmailData{
		URL:       conf.ClientOrigin + "/users/auth/verify-email/" + code,
		FirstName: *user.FirstName,
		Subject:   "Your account verification code",
	}

	emailHelper.SendEmail(&user, &emailData)

	message := "We sent an email with a verification code to " + *user.Email

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": message,
	})
}

// login handler
func (uc *UserController) Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User


	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if err := userCollection.FindOne(ctx, bson.M{
		"email": user.Email,
	}).Decode(&foundUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"erorr":  "email or password is incorrect",
		})
		return
	}

	passwordIsValid, msg := authHelper.VerifyPassword(*user.Password, *foundUser.Password)

	if !passwordIsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  msg,
		})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  "user not found",
		})
	}

	if !foundUser.IsEmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"status": "failed", "message": "Please verify your email"})
		return
	}

	token, refreshToken, generateTokenErr := tokenHelper.GenerateAllTokens(
		*foundUser.Email,
		*foundUser.FirstName,
		*foundUser.LastName,
		*foundUser.UserType,
		foundUser.ID.Hex(),
	)

	if generateTokenErr != nil {
		log.Panic(generateTokenErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "error occured while generate tokens",
		})
		return
	}

	tokenHelper.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  foundUser.Token,
	})
}

// logout handler
func (uc *UserController) Logout(c *gin.Context) {
	var uidFromKeys = fmt.Sprint(c.Keys["uid"])

	tokenHelper.UpdateAllTokens("", "", uidFromKeys)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "you succefully logout",
	})
}

// get me handler
func (uc *UserController) GetMe(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	uid, err := primitive.ObjectIDFromHex(fmt.Sprint(c.Keys["uid"]))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if err := userCollection.FindOne(ctx, bson.M{"_id": uid}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

// verify email handler
func (uc *UserController) VerifyEmail(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User

	code := c.Params.ByName("verification_code")
	verificationCode := encodeHelper.Encode(code)

	if err := userCollection.FindOneAndUpdate(ctx, bson.M{
		"verification_code": verificationCode,
	}, bson.M{
		"$set": bson.M{
			"verification_code": "",
			"is_email_verified": true,
		},
	}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Email verified succefully",
		"token":   user.Token,
	})
}

// get user by id handler
func (uc *UserController) GetUserById(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	userId, _ := primitive.ObjectIDFromHex(c.Params.ByName("user_id"))
	defer cancel()

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}
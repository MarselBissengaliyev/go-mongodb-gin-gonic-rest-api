package routes

import (
	"github.com/MarselBisengaliev/go-react-blog/controllers"
	"github.com/MarselBisengaliev/go-react-blog/middlewares"
)

var authMiddleware = new(middlewares.AuthMiddleware)
var userController = new(controllers.UserController)
var postController = new(controllers.PostController)
var commentController = new(controllers.CommentController)
var tagController = new(controllers.TagController)
var likeController = new(controllers.LikeController)

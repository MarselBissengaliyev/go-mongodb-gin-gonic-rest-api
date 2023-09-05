package main

import (
	"log"

	"github.com/MarselBisengaliev/go-react-blog/config"
	"github.com/MarselBisengaliev/go-react-blog/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig("./")

	if err != nil {
		log.Fatal("could not load config + \n", err)
	}

	
	port := conf.Port

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	routes.RegisterPostStoreRoutes(v1)
	routes.RegisterUserStoreRoutes(v1)
	routes.RegisterCommentStoreRoute(v1)
	routes.RegisterLikeStoreRoutes(v1)

	router.Run(":" + port)
}

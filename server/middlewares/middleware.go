package middlewares

import (
	"github.com/MarselBisengaliev/go-react-blog/database"
	"github.com/MarselBisengaliev/go-react-blog/helpers"
)

var tokenHelper = new(helpers.TokenHelper)
var userCollection = database.OpenCollection(database.Client, "users")
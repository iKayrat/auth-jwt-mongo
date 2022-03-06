package api

import (
	"github.com/gin-gonic/gin"
	"github.com/iKayrat/auth-jwt-mongo/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Store  db.Store
	Router *gin.Engine
}

func NewServer(client *mongo.Client) *Server {
	router := gin.Default()

	// new db
	store := db.NewStore(client)

	server := &Server{
		Store:  store,
		Router: router,
	}

	// routers
	router.GET("/rest/token/", server.NewToken)
	router.POST("/rest/token/refresh", server.RefreshToken)

	return server
}

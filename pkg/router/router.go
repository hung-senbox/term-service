package router

import (
	"term-service/internal/term/handler"
	"term-service/internal/term/repository"
	"term-service/internal/term/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(mongoCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// Setup dependency injection
	repo := repository.NewTermRepository(mongoCollection)
	svc := service.NewTermService(repo)
	h := handler.NewHandler(svc)

	// Register routes
	h.RegisterRoutes(r)

	return r
}

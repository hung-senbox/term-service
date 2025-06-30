package router

import (
	"term-info-service/internal/term/handler"
	"term-info-service/internal/term/repository"
	"term-info-service/internal/term/service"

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

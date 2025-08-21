package router

import (
	"term-service/internal/gateway"
	"term-service/internal/term/handler"
	"term-service/internal/term/repository"
	"term-service/internal/term/service"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(mongoCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()
	// consul
	consulClient, _ := api.NewClient(api.DefaultConfig())

	// Tạo UserGateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)

	// Setup dependency injection
	repo := repository.NewTermRepository(mongoCollection)
	svc := service.NewTermService(repo, userGateway)
	h := handler.NewHandler(svc)

	// Register routes
	h.RegisterRoutes(r)

	return r
}

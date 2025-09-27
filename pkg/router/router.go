package router

import (
	"term-service/internal/gateway"
	holiday_handler "term-service/internal/holiday/handler"
	holiday_repo "term-service/internal/holiday/repository"
	holiday_route "term-service/internal/holiday/route"
	holiday_service "term-service/internal/holiday/service"
	"term-service/internal/term/handler"
	"term-service/internal/term/repository"
	"term-service/internal/term/route"
	"term-service/internal/term/service"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(termCollection *mongo.Collection, holidayCollection *mongo.Collection, consulClient *api.Client) *gin.Engine {
	r := gin.Default()
	// consul
	//consulClient, _ := api.NewClient(api.DefaultConfig())

	// Gateway setup
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	orgGateway := gateway.NewOrganizationGateway("go-main-service", consulClient)
	messageLanguageGW := gateway.NewMessageLanguageGateway("go-main-service", consulClient)

	// Term
	termRepo := repository.NewTermRepository(termCollection)
	termSvc := service.NewTermService(termRepo, userGateway, orgGateway)
	termHandler := handler.NewHandler(termSvc)

	// Holiday
	holidayRepo := holiday_repo.NewHolidayRepository(holidayCollection)
	holidaySvc := holiday_service.NewHolidayService(holidayRepo, userGateway, orgGateway, messageLanguageGW)
	holidayHandler := holiday_handler.NewHandler(holidaySvc)

	// Register routes
	route.RegisterTermRoutes(r, termHandler)
	holiday_route.RegisterHolidayRoutes(r, holidayHandler)

	return r
}

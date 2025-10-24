package router

import (
	"term-service/internal/cache"
	"term-service/internal/gateway"
	holiday_handler "term-service/internal/holiday/handler"
	holiday_repo "term-service/internal/holiday/repository"
	holiday_route "term-service/internal/holiday/route"
	holiday_service "term-service/internal/holiday/service"
	"term-service/internal/term/handler"
	"term-service/internal/term/repository"
	"term-service/internal/term/route"
	"term-service/internal/term/service"
	"term-service/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"

	cached_service "term-service/internal/cache/service"

	goredis "github.com/redis/go-redis/v9"
)

func SetupRouter(consulClient *api.Client, cacheClient *goredis.Client, termCollection, holidayCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// Gateway setup
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	orgGateway := gateway.NewOrganizationGateway("go-main-service", consulClient)
	messageLanguageGW := gateway.NewMessageLanguageGateway("go-main-service", consulClient)

	// cache setup
	appCache := cache.NewRedisCache(cacheClient)
	cachedUserGateway := cached_service.NewCachedUserGateway(userGateway, appCache, config.AppConfig.Database.RedisCache.TTLSeconds)

	// Term
	termRepo := repository.NewTermRepository(termCollection)
	termSvc := service.NewTermService(termRepo, cachedUserGateway, orgGateway, messageLanguageGW)
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

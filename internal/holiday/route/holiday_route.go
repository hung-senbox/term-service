package route

import (
	"term-service/internal/holiday/handler"
	"term-service/internal/term/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterHolidayRoutes(r *gin.Engine, h *handler.HolidayHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured())
	{
		holidaysAdmin := adminGroup.Group("/holidays")
		{
			holidaysAdmin.POST("", h.UploadHolidays)
			holidaysAdmin.GET("", h.GetHolidays4Web)
		}
	}
}

package route

import (
	"term-service/internal/term/handler"
	"term-service/internal/term/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTermRoutes(r *gin.Engine, h *handler.TermHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured())
	{
		termsAdmin := adminGroup.Group("/terms")
		{
			termsAdmin.POST("", h.UploadTerm)
			termsAdmin.GET("", h.GetTerms4Web)
			termsAdmin.GET("/student/:student_id", h.GetTermsByStudent)
		}
	}

	// Organization routes
	orgGroup := r.Group("/api/v1/organization")
	orgGroup.Use(middleware.Secured())
	{
		orgGroup.GET("/:organization_id/terms", h.GetTermsByOrgID)
	}

	// User routes
	userGroup := r.Group("/api/v1")
	userGroup.Use(middleware.Secured())
	{
		termsUser := userGroup.Group("/terms")
		{
			termsUser.GET("", h.GetTerms4App)
			termsUser.GET("/current", h.GetCurrentTerm)
			termsUser.GET("/student/:student_id", h.GetTermsByStudent)
			termsUser.GET("/organization/:organization_id", h.GetTermsByOrg4App)
		}
	}

	// gw routes
	gatewayGroup := r.Group("/api/v1/gateway")
	gatewayGroup.Use(middleware.Secured())
	{
		termsGateway := gatewayGroup.Group("/terms")
		{
			termsGateway.GET("/:term_id", h.GetTerm4Gw)
		}
	}

}

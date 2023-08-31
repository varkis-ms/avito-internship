package v1

import (
	"avito-internship/docs"
	"avito-internship/internal/service"
	"avito-internship/pkg/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *gin.Engine, l *logging.Logger, services *service.Services) {
	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routers
	h := handler.Group("/api/v1")
	{
		newSegmentRoutes(h.Group("/segment"), services.Segment, l)
		newUserRoutes(h.Group("/user"), services.User, l)
		newReportRoutes(h.Group("/report"), services.Report, l)
	}

}

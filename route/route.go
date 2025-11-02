package route

import (
	"github.com/Prasang-money/searchSvc/handler"
	"github.com/gin-gonic/gin"
)

func GetRoute() *gin.Engine {

	router := gin.Default()
	handler := handler.NewHandler()
	router.GET("/health", handler.HealthCheck())
	router.POST("/createUrl", handler.SearchHandler())

	return router
}

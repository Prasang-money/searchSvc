package route

import (
	"github.com/Prasang-money/searchSvc/cache"
	"github.com/Prasang-money/searchSvc/handler"
	"github.com/Prasang-money/searchSvc/service"
	"github.com/gin-gonic/gin"
)

func GetRoute() *gin.Engine {

	router := gin.Default()
	cache := cache.NewCache(1000)
	service := service.NewService(cache)
	handler := handler.NewHandler(service)

	router.GET("/health", handler.HealthCheck())
	router.GET("/api/countries/search", handler.SearchHandler())

	return router
}

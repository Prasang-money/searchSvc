package handler

import (
	"net/http"

	"github.com/Prasang-money/searchSvc/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(svc service.ServiceInterface) *Handler {
	return &Handler{
		service: svc,
	}
}
func (handler Handler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	}
}

func (handler Handler) SearchHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		countryName := c.Query("name")
		resp, err := handler.service.SearchCountries(countryName)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, *resp)

	}

}

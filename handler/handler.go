package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
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
		// Implementation here
	}

}

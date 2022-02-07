package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	// Routes
	r.GET("/healthcheck", HealthCheckHandler)

	return r
}

func HealthCheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

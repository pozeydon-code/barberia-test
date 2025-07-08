package http

import (
	"barberia-api-class/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(apptSvc *service.AppointmentService, prodSvc *service.ProductService) *gin.Engine {
	r := gin.Default()

	// Middleware de CORS basico
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "barberia-api",
		})
	})

	v1 := r.Group("/api/v1")
	{
		// Appointments Routes
		appts := v1.Group("/appointments")
		{
			apptHandler := NewAppointmentHandler(apptSvc)
			appts.POST("", apptHandler.Create)
			appts.GET("", apptHandler.List)
			appts.GET("/:id", apptHandler.Get)
			appts.PUT("/:id", apptHandler.Update)
			appts.DELETE("/:id", apptHandler.Delete)
			appts.GET("/:id/total", apptHandler.GetTotal)
		}

		// Products Routes
		products := v1.Group("/products")
		{
			prodHandler := NewProductHandler(prodSvc)
			products.POST("", prodHandler.Create)
			products.GET("", prodHandler.List)
			products.GET("/:id", prodHandler.Get)
			products.PUT("/:id", prodHandler.Update)
			products.DELETE("/:id", prodHandler.Delete)
		}
	}

	return r
}
